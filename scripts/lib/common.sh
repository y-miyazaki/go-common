#!/bin/bash
#######################################
# Description: Common utility functions for shell scripts
# Usage: source /path/to/scripts/lib/common.sh
#
# This library provides standard functions used across multiple scripts:
# - Logging and output formatting
# - Error handling and exit functions
# - Help/usage display functions
# - Dependency validation
# - Section headers for organized output
#######################################

#######################################
# Function to display section headers with consistent formatting
# Arguments:
#   $1 - Section title
# Outputs:
#   Formatted section header to stdout
#######################################
function echo_section {
    echo "#--------------------------------------------------------------"
    echo "# $1"
    echo "#--------------------------------------------------------------"
}
#######################################
# Function to display section headers and capture start time
# Arguments:
#   $1 - Section title
# Outputs:
#   Formatted section header to stdout
# Usage:
#   start_time=$(get_start_time)
#   start_echo_section "Title"
#   # ... work ...
#   end_echo_section "Title" "$start_time"
#######################################
function start_echo_section {
    local title="$1"
    echo "#--------------------------------------------------------------"
    echo "# $title"
    echo "#--------------------------------------------------------------"
}

#######################################
# Function to get current time for timing measurements
# Outputs:
#   Current epoch time to stdout
#######################################
function get_start_time {
    date +%s
}
#######################################
# Function to display section completion with elapsed time
# Arguments:
#   $1 - Section title
#   $2 - Start time (epoch seconds from start_echo_section)
# Outputs:
#   Formatted section footer with elapsed time to stdout
#######################################
function end_echo_section {
    local title="$1"
    local start_time="$2"
    local end_time
    end_time=$(date +%s)
    local elapsed=$((end_time - start_time))

    echo "#--------------------------------------------------------------"
    echo "# $title completed in ${elapsed} seconds"
    echo "#--------------------------------------------------------------"
}

#######################################
# Function to display error message and exit with error code
# Arguments:
#   $1 - Error message
#   $2 - Exit code (optional, defaults to 1)
# Outputs:
#   Error message to stderr
# Returns:
#   Exits with specified code
#######################################
function error_exit {
    local message="$1"
    local exit_code="${2:-1}"
    echo "ERROR: $message" >&2
    exit "$exit_code"
}

#######################################
# Function to log messages with timestamp and level
# Arguments:
#   $1 - Log level (INFO, WARN, ERROR, DEBUG)
#   $2 - Log message
# Globals:
#   VERBOSE - If true, shows all levels; otherwise only ERROR/WARN
# Outputs:
#   Formatted log message to stdout
#######################################
function log {
    local level="$1"
    local message="$2"

    if [[ "$level" == "ERROR" ]] || [[ "$level" == "WARN" ]] || [[ "${VERBOSE:-false}" == "true" ]]; then
        echo "[$(date +'%Y-%m-%d %H:%M:%S')] [$level] $message"
    fi
}

#######################################
# Function to validate required command line tools
# Arguments:
#   $@ - List of required tools/commands
# Outputs:
#   Error message if tool is missing
# Returns:
#   0 if all tools are available, exits with error if any missing
#######################################
function validate_dependencies {
    local required_tools=("$@")
    local missing_tools=()

    for tool in "${required_tools[@]}"; do
        if ! command -v "$tool" &> /dev/null; then
            missing_tools+=("$tool")
        fi
    done

    if [[ ${#missing_tools[@]} -gt 0 ]]; then
        error_exit "Missing required tools: ${missing_tools[*]}. Please install them and ensure they are in PATH."
    fi

    log "INFO" "All required dependencies are available: ${required_tools[*]}"
}

#######################################
# Function to validate required environment variables
# Arguments:
#   $@ - List of required environment variable names
# Outputs:
#   Error message if variable is missing
# Returns:
#   0 if all variables are set, exits with error if any missing
#######################################
function validate_env_vars {
    local required_vars=("$@")
    local missing_vars=()

    for var in "${required_vars[@]}"; do
        if [[ -z "${!var:-}" ]]; then
            missing_vars+=("$var")
        fi
    done

    if [[ ${#missing_vars[@]} -gt 0 ]]; then
        error_exit "Missing required environment variables: ${missing_vars[*]}"
    fi

    log "INFO" "All required environment variables are set: ${required_vars[*]}"
}

#######################################
# Function to display standardized help message header
# Arguments:
#   $1 - Script name (usually $(basename "$0"))
#   $2 - Brief description
#   $3 - Usage pattern
# Outputs:
#   Formatted help header to stdout
#######################################
function show_help_header {
    local script_name="$1"
    local description="$2"
    local usage_pattern="$3"

    echo "Usage: $script_name $usage_pattern"
    echo ""
    echo "Description: $description"
    echo ""
}

#######################################
# Function to display standardized help footer
# Outputs:
#   Common help options to stdout
#######################################
function show_help_footer {
    echo ""
    echo "Common Options:"
    echo "  -h, --help     Display this help message"
    echo "  -v, --verbose  Enable verbose output"
    echo "  -d, --dry-run  Run in dry-run mode (no changes made)"
    echo ""
}

#######################################
# Function to check if running in dry-run mode
# Globals:
#   DRY_RUN - Boolean flag
# Returns:
#   0 if in dry-run mode, 1 otherwise
#######################################
function is_dry_run {
    [[ "${DRY_RUN:-false}" == "true" ]]
}

#######################################
# Function to execute command with dry-run support
# Arguments:
#   $@ - Command to execute
# Globals:
#   DRY_RUN - If true, only shows what would be executed
# Outputs:
#   Command being executed (if verbose or dry-run)
# Returns:
#   Command exit code (or 0 in dry-run mode)
#######################################
function execute_command {
    # Execute a command safely without eval. Accepts arguments and runs them
    # Example: execute_command aws s3 cp "src" "dest"
    if is_dry_run; then
        # Always print dry-run info regardless of VERBOSE so users see planned actions
        echo "DRY-RUN: Would execute: $*"
        return 0
    fi

    log "DEBUG" "Executing: $*"

    # Execute the command. Support two calling styles:
    # 1) execute_command cmd arg1 arg2 ...  --> safe, runs the command directly
    # 2) execute_command "cmd arg1 arg2"   --> common in older scripts; run via bash -lc
    if [[ $# -eq 1 ]]; then
        # Single-string command (may contain spaces/options) — run under bash -lc so
        # shell parsing behaves as the caller expects.
        bash -lc "$1"
    else
        # Multi-argument safe execution
        "${@}"
    fi
}
