#!/bin/bash
#######################################
# Description: Builds Go Lambda functions for deployment
# This script compiles Go code from the specified directory and packages
# the binaries for AWS Lambda deployment.
#
# Usage: ./build.sh <source_dir> <output_dir> [architecture] [--parallel] [--verbose]
#   <source_dir>    Directory containing main.go files (e.g., cmd/api)
#   <output_dir>    Output directory name (e.g., api)
#   [architecture]  Target architecture (default: arm64)
#   [--parallel]    Enable parallel builds (experimental)
#   [--verbose]     Show more detailed output
#######################################

# Error handling: exit on error, unset variable, or failed pipeline
set -euo pipefail

# Get script directory for library loading
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export SCRIPT_DIR

# Load common libraries - ALWAYS use this pattern
# shellcheck source=../lib/all.sh
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/../lib/all.sh"

#######################################
# Global variables and default values
#######################################
DIR=""
BINDIR=""
ARCH="arm64" # Default to arm64
PARALLEL="false"
VERBOSE="false"
MIN_VERSION="1.21"

#######################################
# Display usage information
#######################################
function show_usage {
    echo "Usage: $(basename "$0") [options]"
    echo ""
    echo "Description: Builds Go Lambda functions for deployment"
    echo ""
    echo "Options:"
    echo "  -h, --help      Display this help message"
    echo "  --parallel      Enable parallel builds (experimental)"
    echo "  --verbose       Show more detailed output"
    echo ""
    echo "Arguments:"
    echo "  source_dir      Directory containing main.go files (e.g., cmd/api)"
    echo "  output_dir      Output directory name (e.g., api)"
    echo "  architecture    Target architecture (default: arm64)"
    echo ""
    echo "Examples:"
    echo "  $(basename "$0") cmd/api api"
    echo "  $(basename "$0") cmd/api api arm64 --parallel --verbose"
    exit 0
}

#######################################
# Parse command line arguments
#######################################
function parse_arguments {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h | --help)
                show_usage
                ;;
            --parallel)
                PARALLEL="true"
                shift
                ;;
            --verbose)
                VERBOSE="true"
                shift
                ;;
            -*)
                error_exit "Unknown option: $1"
                ;;
            *)
                # Process positional arguments
                if [[ -z "$DIR" ]]; then
                    DIR="$1"
                elif [[ -z "$BINDIR" ]]; then
                    BINDIR="$1"
                elif [[ "$ARCH" == "arm64" ]]; then
                    ARCH="$1"
                else
                    error_exit "Unexpected argument: $1"
                fi
                shift
                ;;
        esac
    done

    # Validate required arguments
    if [[ -z "$DIR" ]] || [[ -z "$BINDIR" ]]; then
        echo "Error: Missing required arguments" >&2
        show_usage
    fi

    # Validate directory exists
    if [[ ! -d "$DIR" ]]; then
        error_exit "Source directory does not exist: $DIR"
    fi
}

#######################################
# Validate Go environment
#######################################
function validate_go_environment {
    # Check Go version
    local go_version
    go_version=$(go version | awk '{print $3}' | sed 's/go//')
    log "INFO" "Using Go version $go_version"

    # Simple version check
    if [[ "$(echo -e "$go_version\n$MIN_VERSION" | sort -V | head -n1)" != "$MIN_VERSION" ]]; then
        error_exit "Go version $go_version is less than minimum required version $MIN_VERSION"
    fi

    # Set security environment variables for Go 1.20+
    export GODEBUG=tarinsecurepath=0,zipinsecurepath=0
    log "INFO" "Go environment validated successfully"
}

#######################################
# Prepare build environment
#######################################
function prepare_build_environment {
    echo_section "Preparing build environment"

    # Download and update dependencies
    log "INFO" "Updating Go modules..."
    go mod download || error_exit "Failed to download Go modules"

    # Clean output directories
    log "INFO" "Cleaning output directories..."
    rm -rf outputs/"${BINDIR}"/* 2>/dev/null || true
    rm -rf bin/"${BINDIR}"/* 2>/dev/null || true

    # Create output directories
    log "INFO" "Creating output directories..."
    mkdir -p bin/"${BINDIR}" || error_exit "Failed to create bin directory"
    mkdir -p outputs/"${BINDIR}" || error_exit "Failed to create outputs directory"

    log "INFO" "Build environment prepared successfully"
}

#######################################
# Find Lambda functions to build
#######################################
function find_lambda_functions {
    # Keep visual separators using echo_section (print to stdout for clarity)
    echo_section "Finding Lambda functions"

    local files
    files=$(find "$DIR" -type f -name 'main.go')
    local file_count
    # Count only non-empty lines to avoid counting empty output
    file_count=$(printf "%s\n" "$files" | sed '/^$/d' | wc -l)

    if [[ $file_count -eq 0 ]]; then
        error_exit "No main.go files found in $DIR"
    fi

    log "INFO" "Found $file_count Lambda functions to build"
    # Print file paths only
    printf "%s\n" "$files"
}

#######################################
# Build individual Lambda function
#######################################
function build_function {
    local file=$1
    local dir
    local function

    dir=$(dirname "$file")
    function=$(echo "$dir" | sed -e "s/.*\///g")

    log "INFO" "  Building $function..."
    env GOOS=linux GOARCH="${ARCH}" go build -ldflags="-s -w" -o bootstrap "$file" || error_exit "Failed to build $function"

    log "INFO" "  Packaging $function..."
    zip outputs/"${BINDIR}"/go_"${function}".zip bootstrap || error_exit "Failed to create zip for $function"

    # Get binary size
    if [[ "$VERBOSE" == "true" ]]; then
        local binary_size
        binary_size=$(du -h bootstrap | cut -f1)
        log "INFO" "  Binary size: $binary_size"
    fi

    log "INFO" "  Moving binary for $function..."
    mkdir -p bin/"${BINDIR}"/"${function}"
    mv bootstrap bin/"${BINDIR}"/"${function}"/ || error_exit "Failed to move binary for $function"
    # Informational message should go to stderr so stdout remains clean for capturing the final count
    echo "$function built successfully" >&2
}

#######################################
# Build all Lambda functions
#######################################
function build_lambda_functions {
    local files="$1"
    local file_count
    file_count=$(echo "$files" | wc -l)

    echo_section "Building Lambda functions"
    log "INFO" "Target: $DIR -> $BINDIR ($ARCH)"
    if [[ "$PARALLEL" == "true" ]]; then
        log "INFO" "Parallel build enabled"
    fi

    local count=0

    if [[ "$PARALLEL" == "true" && $(command -v parallel) ]]; then
        # Use GNU parallel if available and parallel build is enabled
        export -f build_function log error_exit echo_section
        export BINDIR ARCH VERBOSE
        # Redirect informational output to stderr; keep stdout reserved for the final count
        echo "$files" | parallel build_function {} 1>/dev/null || error_exit "Failed in parallel build"
        count=$file_count
    else
        # Sequential build
        for file in $files; do
            # Ensure per-function informational output goes to stderr
            build_function "$file" 1>/dev/null
            count=$((count + 1))
        done
    fi
}

#######################################
# Main execution function
#######################################
function main {
    local start_time
    start_time=$(date +%s)

    # Parse arguments
    parse_arguments "$@"

    # Validate dependencies
    validate_dependencies "go"

    # Validate Go environment
    validate_go_environment

    # Prepare build environment
    prepare_build_environment

    # Find Lambda functions (filter out header/separator lines starting with '#')
    local files
    files=$(find_lambda_functions | sed '/^#/d')

    # Build all Lambda functions
    local count
    build_lambda_functions "$files"

    # Calculate elapsed time
    local end_time elapsed minutes seconds
    end_time=$(date +%s)
    elapsed=$((end_time - start_time))
    minutes=$((elapsed / 60))
    seconds=$((elapsed % 60))

    echo_section "Build completed successfully! Built Lambda functions in ${minutes}m ${seconds}s."
    log "INFO" "All Lambda functions built successfully"
}

# Only call main function if script is executed directly, not sourced
if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi
