#!/bin/bash
#######################################
# Description: Collect AWS Lambda, Glue, and RDS available runtime/engine versions and output as CSV
# Usage: ./aws_get_version.sh [-h] [-v] [-d] [-r REGION] [-c CATEGORIES] [-o OUTPUT_FILE]
#   options:
#     -h, --help       Display this help message
#     -v, --verbose    Enable verbose output
#     -d, --dry-run    Run in dry-run mode (no changes made)
#     -r, --region     AWS region to use (default: $AWS_DEFAULT_REGION or ap-northeast-1)
#     -c, --categories Comma-separated list of categories to collect (optional)
#     -o, --output     Output CSV file (default: aws_runtime_versions.csv)
#
# Output:
# - Generates CSV file containing available runtime/engine versions for Lambda, Glue, and RDS
# - Marks latest versions available for each runtime/engine type
# - CSV columns: Category,Subcategory,Subsubcategory,Region,Version,Is_Latest,Status,Deprecation_Date,EOL_Date
#
# CSV Format Notes:
# - Compatible with Excel, Numbers, Google Sheets
# - All values with commas or newlines are properly quoted
# - Is_Latest column shows "Yes" for latest versions, empty for older versions
# - Status column shows availability status (available, deprecated, deprecation_scheduled)
# - Deprecation_Date column shows when runtime was/will be deprecated (YYYY-MM-DD format)
# - EOL_Date column shows when runtime will reach end-of-life (YYYY-MM-DD format)
#
# Design Rules:
# - Modular design with separate collection functions per runtime/engine type
# - Extensible to add new runtime/engine types easily
# - Structured data approach for Lambda and Glue with static definitions
# - AWS API calls for RDS to get real-time engine version information
# - Category-based output with consistent column alignment
# - Runtime categories processed in defined order (lambda, glue, rds)
# - No empty category headers - headers only appear when versions exist
# - All functions follow consistent naming pattern: collect_<type>_versions
# - Latest version marking for each runtime family/engine type
# - Status tracking with deprecation and end-of-life date information
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
VERBOSE=false
export VERBOSE
DRY_RUN=false
OUTPUT_FILE="aws_runtime_versions.csv"
AWS_REGION="${AWS_DEFAULT_REGION:-ap-northeast-1}"
CATEGORIES=""

# Categories of services to collect runtime/engine versions
RUNTIME_CATEGORIES=(
    "lambda"
    "glue"
    "rds"
)

# Categories that need to maintain grouping structure (no sorting)
NO_SORT_CATEGORIES=(
    "rds" # Maintain engine grouping structure
)

#######################################
# Lambda runtime data as structured text
# Format: family|runtime|status|deprecation_date|block_create_date|block_update_date|is_latest_in_family
# Based on AWS documentation: https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html
# Additional reference: https://docs.aws.amazon.com/lambda/latest/dg/runtimes-update.html
# Status values: deprecated, deprecation_scheduled, available
# Dates format: YYYY-MM-DD
# Update this variable to add/modify Lambda runtime versions
#######################################
LAMBDA_DATA='nodejs|nodejs16.x|deprecated|2024-06-12|2025-10-01|2025-11-01|false
nodejs|nodejs18.x|deprecation_scheduled|2025-09-01|2025-10-01|2025-11-01|false
nodejs|nodejs20.x|available|2026-04-30|2026-06-01|2026-07-01|false
nodejs|nodejs22.x|available|2027-04-30|2027-06-01|2027-07-01|true
python|python3.8|deprecated|2024-10-14|2025-10-01|2025-11-01|false
python|python3.9|deprecation_scheduled|2025-12-15|2026-01-15|2026-02-15|false
python|python3.10|available|2026-06-30|2026-07-31|2026-08-31|false
python|python3.11|available|2026-06-30|2026-07-31|2026-08-31|false
python|python3.12|available|2028-10-31|2028-11-30|2029-01-10|false
python|python3.13|available|2029-06-30|2029-07-31|2029-08-31|true
java|java8|deprecated|2024-01-08|2024-02-08|2025-11-01|false
java|java8.al2|deprecated|2024-01-08|2024-02-08|2025-11-01|false
java|java11|available|2026-06-30|2026-07-31|2026-08-31|false
java|java17|available|2026-06-30|2026-07-31|2026-08-31|false
java|java21|available|2029-06-30|2029-07-31|2029-08-31|true
dotnet|dotnet6|deprecated|2024-12-20|2025-10-01|2025-11-01|false
dotnet|dotnet8|available|2026-11-10|2026-12-10|2027-01-11|true
dotnet|dotnet9|available|2027-11-10|2027-12-10|2028-01-11|false
ruby|ruby2.7|deprecated|2023-12-07|2024-01-09|2025-11-01|false
ruby|ruby3.2|available|2026-03-31|2026-04-30|2026-05-31|false
ruby|ruby3.3|available|2027-03-31|2027-04-30|2027-05-31|true
ruby|ruby3.4|available|2028-03-31|2028-04-30|2028-05-31|false
provided|provided|available|2026-06-30|2026-07-31|2026-08-31|false
provided|provided.al2|available|2026-06-30|2026-07-31|2026-08-31|false
provided|provided.al2023|available|2029-06-30|2029-07-31|2029-08-31|true
go|provided.al2023|available|2029-06-30|2029-07-31|2029-08-31|true'

#######################################
# Glue version data as structured text
# Format: version|status|end_of_support_date|end_of_life_date|is_latest
# Based on AWS documentation: https://docs.aws.amazon.com/glue/latest/dg/glue-version-support-policy.html
# Migration guide: https://docs.aws.amazon.com/glue/latest/dg/migrating-version-50.html
# Release notes: https://docs.aws.amazon.com/glue/latest/dg/release-notes.html
# Status values: deprecated, available
# Dates format: YYYY-MM-DD (empty for current supported versions)
# Update this variable to add/modify Glue versions
#######################################
GLUE_DATA='1.0|deprecated|2022-09-30|2026-04-01|false
2.0|deprecated|2024-01-31|2026-04-01|false
3.0|available|||false
4.0|available|||false
5.0|available|||true'

#######################################
# Display usage information
#######################################
function show_usage {
    cat <<EOF
Usage: $(basename "$0") [options]

Description: Collect AWS Lambda, Glue, and RDS available runtime/engine versions and mark latest versions.

Options:
  -h, --help       Display this help message
  -v, --verbose    Enable verbose output
  -d, --dry-run    Run in dry-run mode (no changes made)
  -r, --region     AWS region to query (default: \$AWS_DEFAULT_REGION or ap-northeast-1)
  -c, --categories Comma-separated list of categories to collect (optional)
  -o, --output     Output CSV file (default: aws_runtime_versions.csv)

Available categories:
  lambda, glue, rds

Examples:
  $(basename "$0") -v -o my_versions.csv
  $(basename "$0") -r us-east-1 -o us_versions.csv
  $(basename "$0") -c lambda,glue -o lambda_glue_only.csv

EOF
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
            -v | --verbose)
                VERBOSE=true
                export VERBOSE
                shift
                ;;
            -d | --dry-run)
                DRY_RUN=true
                shift
                ;;
            -o | --output)
                OUTPUT_FILE="$2"
                shift 2
                ;;
            -r | --region)
                AWS_REGION="$2"
                shift 2
                ;;
            -c | --categories)
                CATEGORIES="$2"
                shift 2
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
# Function to collect Lambda runtime versions
#######################################
function collect_lambda_versions {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Version,Is_Latest,Status,Deprecation_Date,Block_Function_Create,Block_Function_Update"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    # Process each runtime entry from global LAMBDA_DATA
    while IFS='|' read -r family runtime status deprecation_date block_create_date block_update_date is_latest_in_family; do
        # Skip empty lines
        [[ -z "$family" ]] && continue

        local is_latest=""
        if [[ "$is_latest_in_family" == "true" ]]; then
            is_latest="Yes"
        fi

        buffer+="Lambda,Runtime,$family,$region,$runtime,$is_latest,$status,$deprecation_date,$block_create_date,$block_update_date\n"
    done <<<"$LAMBDA_DATA"

    echo "$buffer"
}

#######################################
# Function to collect Glue versions
#######################################
function collect_glue_versions {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Version,Is_Latest,Status,End_of_Support,End_of_Life"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    # Process each version entry from global GLUE_DATA
    while IFS='|' read -r version status end_of_support_date end_of_life_date is_latest_version; do
        # Skip empty lines
        [[ -z "$version" ]] && continue

        local is_latest=""
        if [[ "$is_latest_version" == "true" ]]; then
            is_latest="Yes"
        fi

        buffer+="Glue,Version,glue,$region,$version,$is_latest,$status,$end_of_support_date,$end_of_life_date\n"
    done <<<"$GLUE_DATA"

    echo "$buffer"
}

#######################################
# Function to collect RDS engine versions
#######################################
function collect_rds_versions {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Version,Is_Latest,Status,Deprecation_Date,EOL_Date"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    # Get available engine versions for major engines using AWS API
    local engines=("mysql" "postgres" "aurora-mysql" "aurora-postgresql" "mariadb")

    for engine in "${engines[@]}"; do
        local engine_versions_json
        engine_versions_json=$(aws rds describe-db-engine-versions --engine "$engine" --region "$region" 2>/dev/null | jq '.DBEngineVersions' || echo "[]")

        # Get latest version by sorting versions properly
        local latest_version
        latest_version=$(echo "$engine_versions_json" | jq -r 'sort_by(.EngineVersion) | reverse | .[0].EngineVersion // ""')

        # Process each version with AWS API data
        while IFS= read -r version_data; do
            [[ -z "$version_data" ]] && continue

            local version status
            version=$(extract_jq_value "$version_data" '.EngineVersion')
            status=$(extract_jq_value "$version_data" '.Status')

            if [[ -n "$version" && "$version" != "N/A" ]]; then
                local is_latest=""
                if [[ "$version" == "$latest_version" ]]; then
                    is_latest="Yes"
                fi

                # Use engine name as Subsubcategory, include status from AWS API
                # RDS versions don't have standardized deprecation dates from AWS API
                local deprecation_date=""
                local eol_date=""
                buffer+="RDS,Engine,$engine,$region,$version,$is_latest,$status,$deprecation_date,$eol_date\n"
            fi
        done < <(echo "$engine_versions_json" | jq -c '.[]?')
    done

    echo "$buffer"
}

#######################################
# Common utility functions for AWS runtime version collection
#######################################

# Generic function to collect AWS runtime versions across regions
function collect_runtime_versions {
    local category=$1

    log "INFO" "Collecting $category runtime/engine versions from AWS..."

    # Get header from the first call
    local collect_function="collect_${category}_versions"
    if ! declare -f "$collect_function" >/dev/null; then
        log "WARN" "Collection function $collect_function not found for category $category"
        return 1
    fi

    # Get header
    local csv_header
    csv_header=$($collect_function "header")

    local buffer=""
    log "INFO" "Checking $category versions in region: $AWS_REGION"
    buffer+=$($collect_function "$AWS_REGION")

    # Check if this category should maintain grouping structure (no sorting)
    local sort_output="true"
    for no_sort_category in "${NO_SORT_CATEGORIES[@]}"; do
        if [[ "$category" == "$no_sort_category" ]]; then
            sort_output="false"
            break
        fi
    done

    output_csv_data "$category" "$csv_header" "$buffer" "$sort_output"
}

# Function to output CSV data with standard formatting
function output_csv_data {
    local category=$1
    local header=$2
    local buffer=$3
    local sort_output=${4:-"true"} # Use explicit parameter or default to true

    if [[ -n "$buffer" ]]; then
        {
            echo "$header"
            if [[ "$sort_output" == "true" ]]; then
                printf "%b" "$buffer" | sort
            else
                printf "%b" "$buffer"
            fi
            echo ""
        } >>"$OUTPUT_FILE"
    fi
    log "INFO" "$category runtime/engine versions written to $OUTPUT_FILE"
}

#######################################
# Main execution function
#######################################
function main {
    parse_arguments "$@"

    # Record start time
    start_time=$(date +%s)

    # Validate dependencies
    validate_dependencies "aws" "jq"

    # Check AWS credentials before any AWS CLI usage
    if ! check_aws_credentials; then
        error_exit "AWS credentials are not set or invalid."
    fi

    log "INFO" "Starting AWS runtime version collection"
    log "INFO" "Output file: $OUTPUT_FILE"
    log "INFO" "AWS region: $AWS_REGION"

    if [[ "$DRY_RUN" == "true" ]]; then
        log "INFO" "Running in dry-run mode, no changes will be made"
        return 0
    fi

    # Initialize output file
    true >"$OUTPUT_FILE"

    # Determine which categories to process
    local categories_to_process=()
    if [[ -n "$CATEGORIES" ]]; then
        # Split comma-separated categories into array
        IFS=',' read -ra categories_to_process <<<"$CATEGORIES"
        log "INFO" "Processing specified categories: ${categories_to_process[*]}"

        # Validate specified categories
        for category in "${categories_to_process[@]}"; do
            local valid_category=false
            for valid in "${RUNTIME_CATEGORIES[@]}"; do
                if [[ "$category" == "$valid" ]]; then
                    valid_category=true
                    break
                fi
            done
            if [[ "$valid_category" == "false" ]]; then
                error_exit "Invalid category: $category. Valid categories are: ${RUNTIME_CATEGORIES[*]}"
            fi
        done
    else
        # Use all categories if none specified
        categories_to_process=("${RUNTIME_CATEGORIES[@]}")
        log "INFO" "Processing all categories: ${categories_to_process[*]}"
    fi

    # Collect runtime versions for specified categories
    for runtime_category in "${categories_to_process[@]}"; do
        collect_runtime_versions "$runtime_category"
    done

    # Record end time and calculate elapsed time
    end_time=$(date +%s)
    elapsed=$((end_time - start_time))
    log "INFO" "AWS runtime version collection completed in ${elapsed} seconds"

    # Display summary
    local total_entries
    total_entries=$(($(wc -l <"$OUTPUT_FILE") - $(echo "${RUNTIME_CATEGORIES[@]}" | wc -w))) # Subtract header lines
    log "INFO" "Total runtime/engine versions collected: $total_entries"

    echo "Results written to: $OUTPUT_FILE"
}

# Script entry point
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
