#!/bin/bash
#######################################
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
    show_help_header "$(basename "$0")" "Deploy Terraform configuration" "[directory]"
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

    # Validate environment variables and dependencies
    validate_env_vars "ENV" "TF_PLUGIN_CACHE_DIR"
    validate_dependencies "terraform" "tfenv"

    # Change to target directory
    if ! cd "$dir"; then
        error_exit "Failed to change to directory: $dir"
    fi

    log "INFO" "Starting Terraform deployment in directory: $(pwd)"
}

#######################################
# Run terraform deployment
#######################################
function run_terraform_deployment {
    # Run complete Terraform workflow
    terraform_workflow "$ENV" "apply" "auto-approve"
}

#######################################
# Main execution function
#######################################
function main {
    # Set target directory
    local dir=${1:-./}

    # Validate environment and prepare
    validate_and_prepare "$dir"

    # Run terraform deployment
    run_terraform_deployment
}

# Only call main function if script is executed directly, not sourced
if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi
