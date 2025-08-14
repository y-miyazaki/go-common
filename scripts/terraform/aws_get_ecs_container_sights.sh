#!/bin/bash
#######################################
# Description: Retrieves ECS task definition families and outputs them in Terraform configuration format
# Usage: ./aws_ecs_container_sights.sh [options]
#   options:
#     -h, --help    Display this help message
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
    echo "Usage: $(basename "$0") [options]"
    echo ""
    echo "This script retrieves ECS task definition families and outputs them in a format"
    echo "suitable for Terraform configuration."
    echo ""
    echo "Options:"
    echo "  -h, --help    Display this help message"
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
# Retrieve ECS task definitions from AWS
#######################################
function get_task_definitions {
    echo_section "Retrieving ECS task definitions"

    if ! task_definitions=$(aws ecs list-task-definitions --output json | jq -r '.taskDefinitionArns[]'); then
        error_exit "Failed to retrieve ECS task definitions"
    fi

    echo "$task_definitions"
}

#######################################
# Extract unique task definition families
#######################################
function extract_unique_families {
    local task_definitions="$1"
    local unique_families=()

    echo_section "Processing task definition families"

    for task_definition in $task_definitions; do
        # Extract task definition family
        local family
        family=$(echo "$task_definition" | awk -F'/' '{print $NF}' | awk -F':' '{print $1}')

        # Add family only if it doesn't already exist in the array
        local family_exists=false
        for existing_family in "${unique_families[@]}"; do
            if [[ "$existing_family" == "$family" ]]; then
                family_exists=true
                break
            fi
        done

        if [[ "$family_exists" == "false" ]]; then
            unique_families+=("$family")
        fi
    done

    # Return the unique families as a space-separated string
    printf '%s\n' "${unique_families[@]}"
}

#######################################
# Generate Terraform configuration format output
#######################################
function generate_terraform_output {
    local families=("$@")
    local items=""

    echo_section "Generating Terraform configuration format"

    for family in "${families[@]}"; do
        items+="  {\n    ClusterName = \"\"\n    TaskDefinitionFamily = \"$family\"\n  },\n"
    done

    # Remove trailing comma if items is not empty
    if [ -n "$items" ]; then
        items=${items%,*}
    fi

    echo "$items"
}

#######################################
# Output the final result
#######################################
function output_result {
    local formatted_items="$1"

    echo_section "Generating output"
    echo "["
    echo -e "$formatted_items"
    echo "]"
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

    # Log script start
    echo_section "Starting ECS task definition family extraction"
    log "INFO" "Retrieving ECS task definitions from AWS"

    # Get task definitions from AWS
    local task_definitions
    task_definitions=$(get_task_definitions)

    # Extract unique families
    local unique_families
    readarray -t unique_families < <(extract_unique_families "$task_definitions")

    # Generate Terraform configuration format
    local terraform_output
    terraform_output=$(generate_terraform_output "${unique_families[@]}")

    # Output the final result
    output_result "$terraform_output"

    echo_section "Process completed successfully"
    log "INFO" "Successfully processed ${#unique_families[@]} unique task definition families"
}

# Only call main function if script is executed directly, not sourced
if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi
