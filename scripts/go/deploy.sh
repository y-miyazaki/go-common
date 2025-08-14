#!/bin/bash
#######################################
# Description: Deploys Go Lambda functions to AWS using Serverless Framework
# This script handles the deployment process of Go Lambda functions by
# installing dependencies, building the project, and deploying to AWS.
#
# Usage: ./deploy.sh <stage>
#   <stage>    Deployment stage/environment (e.g., dev, staging, prod)
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
# Validate input arguments
#######################################
function validate_arguments {
    local stage="$1"

    if [ -z "$stage" ]; then
        show_usage "Stage argument is required"
    fi
}

#######################################
# Install dependencies
#######################################
function install_dependencies {
    log "INFO" "Installing dependencies..."
    if ! execute_command "npm ci"; then
        error_exit "Failed to install dependencies"
    fi
}

#######################################
# Build project
#######################################
function build_project {
    log "INFO" "Building project..."
    if ! execute_command "make build"; then
        error_exit "Failed to build project"
    fi
}

#######################################
# Deploy to AWS
#######################################
function deploy_to_aws {
    local stage="$1"

    log "INFO" "Deploying to AWS ($stage)..."
    if ! execute_command "make deploy STAGE=${stage}"; then
        error_exit "Failed to deploy to $stage environment"
    fi
}

#######################################
# Display usage information
#######################################
function show_usage {
    local error_msg="$1"

    if [ -n "$error_msg" ]; then
        echo "Error: $error_msg" >&2
        echo ""
    fi

    show_help_header "$(basename "$0")" "Deploy Go Lambda functions to AWS using Serverless Framework" "<stage>"
    echo "This script handles the deployment process of Go Lambda functions by"
    echo "installing dependencies, building the project, and deploying to AWS."
    echo ""
    echo "Arguments:"
    echo "  stage       Deployment stage/environment (e.g., dev, staging, prod)"
    echo ""
    show_help_footer
    echo "Examples:"
    echo "  $0 dev"
    echo "  $0 staging"
    echo "  $0 prod"
    exit 1
}

#######################################
# Main function
#######################################
function main {
    local stage="$1"

    # Show usage if -h or --help is provided
    if [ "$stage" = "-h" ] || [ "$stage" = "--help" ]; then
        show_usage
    fi

    # Validate input arguments
    validate_arguments "$stage"

    echo_section "Deploying Lambda functions to $stage environment"

    # Install dependencies
    install_dependencies

    # Build project
    build_project

    # Deploy to AWS
    deploy_to_aws "$stage"

    echo_section "Deployment completed successfully!"
}

# Only call main function if script is executed directly, not sourced
if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi
