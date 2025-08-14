#!/bin/bash
#######################################
# Description: Delete all QuickSight namespaces for the current AWS account
# Usage: ./aws_delete_quicksight_namespace.sh
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
# Global variables and default values
#######################################
AWS_ACCOUNT_ID=""

#######################################
# Display usage information
#######################################
function show_usage {
    echo "Usage: $(basename "$0")"
    echo ""
    echo "Delete all QuickSight namespaces for the current AWS account"
    echo ""
    echo "Options:"
    echo "  -h, --help        Display this help message"
    echo ""
    echo "Example: $(basename "$0")"
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
            -*)
                error_exit "Unknown option: $1"
                ;;
            *)
                error_exit "Unexpected argument: $1"
                ;;
        esac
    done
}

#######################################
# Delete all QuickSight namespaces
#######################################
function delete_quicksight_namespaces {
    echo_section "Deleting QuickSight Namespaces"
    log "INFO" "Target AWS account: $AWS_ACCOUNT_ID"

    # Get all namespaces and delete each one
    aws quicksight list-namespaces --aws-account-id "$AWS_ACCOUNT_ID" --output json | jq -r '.Namespaces[].Name' | while read -r namespace; do
        log "INFO" "Deleting namespace: $namespace"
        aws quicksight delete-namespace --aws-account-id "$AWS_ACCOUNT_ID" --namespace "$namespace"
    done

    log "INFO" "All namespaces have been deleted"
}

#######################################
# Main execution function
#######################################
function main {
    # Parse arguments
    parse_arguments "$@"

    # Validate required dependencies
    validate_dependencies "aws" "jq"

    # Check AWS credentials before any AWS CLI usage
    if ! check_aws_credentials; then
        error_exit "AWS credentials are not set or invalid."
    fi

    # Get AWS account ID after credential validation
    AWS_ACCOUNT_ID=$(get_aws_account_id)
    if [[ -z "$AWS_ACCOUNT_ID" ]]; then
        error_exit "Failed to get AWS account ID"
    fi

    # Log script start
    echo_section "Starting QuickSight namespace deletion"
    log "INFO" "Target AWS account: $AWS_ACCOUNT_ID"

    # Delete all QuickSight namespaces
    delete_quicksight_namespaces

    echo_section "Process completed successfully"
    log "INFO" "QuickSight namespace deletion completed"
}

# Only call main function if script is executed directly, not sourced
if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi
