#!/bin/bash
#######################################
# Description: Terraform integration testing script (init, validate, plan, lint, security)
# Usage: ./integration.sh [directory]
#   directory: Target directory (optional, defaults to current)
#######################################

# Error handling: exit on error, unset variable, or failed pipeline
set -euo pipefail

# Get script directory for library loading
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export SCRIPT_DIR

# Load all-in-one library
# shellcheck source=../lib/all.sh
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/../lib/all.sh"

#######################################
# Display usage information
#######################################
function show_usage {
    show_help_header "$(basename "$0")" "Terraform integration testing (init, validate, plan)" "[directory]"
    echo "This script performs Terraform initialization, validation and planning for a directory."
    echo ""
    echo "Arguments:"
    echo "  directory       Target directory (optional, defaults to current)"
    echo ""
    show_help_footer
    echo "Required environment variables:"
    echo "  TF_PLUGIN_CACHE_DIR     - Terraform plugin cache directory"
    echo "  ENV                     - Environment (e.g., dev, staging, prod)"
    echo ""
    echo "Examples:"
    echo "  ENV=dev $0 ./terraform/base"
    echo "  ENV=prod $0 ./terraform/application"
    exit 1
}

# Show usage if -h or --help is provided
if [ "${1:-}" = "-h" ] || [ "${1:-}" = "--help" ]; then
    show_usage
fi

#######################################
# Validate environment and prepare
#######################################
function validate_and_prepare {
    local dir="$1"

    # Validate environment and dependencies
    validate_terraform_env
    validate_dependencies "terraform" "tflint" "trivy"

    # Change to target directory
    if ! cd "$dir"; then
        error_exit "Failed to change to directory: $dir"
    fi

    log "INFO" "Starting Terraform integration testing in directory: $(pwd)"
}

#######################################
# Run terraform workflow
#######################################
function run_terraform_workflow {
    # Run Terraform workflow (plan only)
    terraform_workflow "$ENV" "plan"
}

#######################################
# Run additional linting and security checks
#######################################
function run_additional_checks {
    echo_section "Additional Linting and Security Checks"

    # Run tflint if available
    if command -v tflint &>/dev/null; then
        log "INFO" "Running tflint with module support..."
        if ! execute_command "tflint --module"; then
            log "WARN" "tflint found issues that should be addressed."
        fi
    else
        log "WARN" "tflint not available, skipping linting"
    fi

    # Run trivy security scan
    if command -v trivy &>/dev/null; then
        log "INFO" "Running trivy security scan..."
        if ! execute_command "trivy fs . --format table"; then
            log "WARN" "trivy found security issues that should be addressed."
        fi
    else
        log "WARN" "trivy not available, skipping security scan"
    fi
}

#######################################
# Main execution function
#######################################
function main {
    # Set target directory
    local dir=${1:-./}

    # Validate environment and prepare
    validate_and_prepare "$dir"

    # Run terraform workflow
    run_terraform_workflow

    # Run additional checks
    run_additional_checks

    echo_section "Integration testing completed"
}

# Only call main function if script is executed directly, not sourced
if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi
