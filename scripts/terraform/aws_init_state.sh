#!/bin/bash
#######################################
# Description: Create an S3 bucket for Terraform remote state with secure defaults
# Usage: ./aws_init_state.sh -r {region} -b {bucket name} [-s]
#   options:
#     -b {bucket name}   S3 bucket name (required)
#     -r {region}        AWS region for S3 bucket (required)
#     -s                 Add random hash suffix to bucket name for uniqueness
#     -h                 Show help
# Design Rules:
#   - Idempotent: safe to re-run if bucket already exists
#   - Secure defaults: block public access, versioning, encryption, lifecycle
#   - Clear logging and error handling
#######################################

# Error handling: exit on error, unset variable, or failed pipeline
set -euo pipefail # Error handling: exit on error, unset variable, or failed pipeline

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
IS_BUCKET_AUTO_HASH=0
BUCKET=""
REGION=""
AWS_ID=""

#######################################
# Display usage information
#######################################
function show_usage {
    local error_msg="${1:-}"

    if [ -n "$error_msg" ]; then
        error_exit "$error_msg"
    fi

    show_help_header "$(basename "$0")" "Create S3 bucket for Terraform state management" "-r {region} -b {bucket name} [options]"
    echo "This script creates an S3 bucket for Terraform state management with proper security settings."
    echo "You can optionally add a random hash suffix to the bucket name for uniqueness."
    echo ""
    echo "Options:"
    echo "  -b {bucket name}          S3 bucket name"
    echo "  -r {region}               AWS region for S3 bucket"
    echo "  -s                        Add random hash suffix to bucket name"
    echo "  -h                        Show this help message"
    show_help_footer
    echo "Examples:"
    echo "  $(basename "$0") -r us-east-1 -b my-terraform-state"
    echo "  $(basename "$0") -r ap-northeast-1 -b terraform-state -s"
    exit 0
}

#######################################
# Parse command line arguments
#######################################
function parse_arguments {
    while getopts b:r:sh opt; do
        case $opt in
            b) BUCKET=$OPTARG ;;
            r) REGION=$OPTARG ;;
            s) IS_BUCKET_AUTO_HASH=1 ;;
            h) show_usage ;;
            \?) show_usage "Invalid option: -$OPTARG" ;;
        esac
    done

    # Validate required parameters
    if [ -z "${BUCKET}" ]; then
        show_usage "Bucket name (-b) is required"
    fi

    if [ -z "${REGION}" ]; then
        show_usage "Region (-r) is required"
    fi
}

#######################################
# Generate unique bucket name if requested
#######################################
function generate_bucket_name {
    if [ $IS_BUCKET_AUTO_HASH -eq 1 ]; then
        echo_section "Generating unique bucket name"
        local random_hash
        random_hash=$(awk -v min=10000000 -v max=100000000 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
        BUCKET="${BUCKET}-${random_hash}"
        log "INFO" "Generated bucket name: $BUCKET"
    fi
}

#######################################
# Get AWS account information
#######################################
function get_aws_account_info {
    echo_section "Retrieving AWS Account information"
    if ! AWS_ID=$(get_aws_account_id); then
        error_exit "Failed to retrieve AWS Account ID. Check your AWS credentials."
    fi
    log "INFO" "AWS Account ID: ${AWS_ID}"
}

#######################################
# Create S3 bucket (handles us-east-1 special case)
#######################################
function create_s3_bucket {
    echo_section "Creating S3 bucket: $BUCKET"
    local create_args=(--bucket "${BUCKET}")

    # us-east-1 does not accept explicit LocationConstraint
    if [ "${REGION}" != "us-east-1" ]; then
        create_args+=(--create-bucket-configuration "LocationConstraint=${REGION}")
    fi

    if ! aws s3api create-bucket "${create_args[@]}" 2>/dev/null; then
        log "WARN" "Bucket creation may have failed or bucket already exists. Proceeding with configuration."
    else
        log "INFO" "S3 bucket created successfully"
    fi
}

#######################################
# Configure bucket security settings
#######################################
function configure_bucket_security {
    echo_section "Configuring bucket security settings"

    # Block public access
    log "INFO" "Setting public access block"
    if ! aws s3api put-public-access-block --bucket "${BUCKET}" \
        --public-access-block-configuration "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true"; then
        error_exit "Failed to set public access block"
    fi

    # Enable versioning
    log "INFO" "Enabling bucket versioning"
    if ! aws s3api put-bucket-versioning --bucket "${BUCKET}" \
        --versioning-configuration Status=Enabled; then
        error_exit "Failed to enable versioning"
    fi

    # Configure default encryption
    log "INFO" "Setting default encryption"
    if ! aws s3api put-bucket-encryption --bucket "${BUCKET}" \
        --server-side-encryption-configuration '{
      "Rules": [
        {
          "ApplyServerSideEncryptionByDefault": {
            "SSEAlgorithm": "AES256"
          }
        }
      ]
    }'; then
        error_exit "Failed to set bucket encryption"
    fi

    # Configure lifecycle management
    log "INFO" "Setting lifecycle configuration"
    if ! aws s3api put-bucket-lifecycle-configuration --bucket "${BUCKET}" \
        --lifecycle-configuration '{
        "Rules": [
            {
                "ID": "default",
                "Status": "Enabled",
                "Filter": {
                    "Prefix": ""
                },
                "AbortIncompleteMultipartUpload": {
                    "DaysAfterInitiation": 7
                }
            }
        ]
    }'; then
        error_exit "Failed to set lifecycle configuration"
    fi
}

#######################################
# Verify bucket configuration
#######################################
function verify_bucket_config {
    echo_section "Verifying bucket configuration"

    log "INFO" "Bucket location:"
    aws s3api get-bucket-location --bucket "${BUCKET}" || log "WARN" "Could not retrieve bucket location"

    log "INFO" "Bucket versioning:"
    aws s3api get-bucket-versioning --bucket "${BUCKET}" || log "WARN" "Could not retrieve versioning"

    log "INFO" "Bucket encryption:"
    aws s3api get-bucket-encryption --bucket "${BUCKET}" || log "WARN" "Could not retrieve encryption"
}

#######################################
# Apply bucket policy for terraform state management
#######################################
function apply_bucket_policy {
    echo_section "Applying bucket policy"

    local policy_template="${SCRIPT_DIR}/files/aws/terraform_state_policy.template.json"

    if [ -f "${policy_template}" ]; then
        log "INFO" "Creating bucket policy from template"

        # Create temporary policy file with substituted values
        local temp_policy_file="${SCRIPT_DIR}/files/aws/terraform_state_policy.json"

        # Substitute placeholders in template
        if ! sed -e "s/##AWS_ID##/${AWS_ID}/g" -e "s/##BUCKET##/${BUCKET}/g" \
            "${policy_template}" >"${temp_policy_file}"; then
            error_exit "Failed to create policy file from template"
        fi

        # Apply bucket policy
        if ! aws s3api put-bucket-policy --bucket "${BUCKET}" \
            --policy "file://${temp_policy_file}"; then
            rm -f "${temp_policy_file}"
            error_exit "Failed to apply bucket policy"
        fi

        # Clean up temporary file
        rm -f "${temp_policy_file}"
        log "INFO" "Bucket policy applied successfully"
    else
        log "WARN" "Policy template file not found at ${policy_template}"
        log "WARN" "Skipping bucket policy configuration"
    fi
}

#######################################
# Display completion summary and examples
#######################################
function display_completion_summary {
    echo_section "S3 Bucket Creation Completed"
    echo "Bucket name: ${BUCKET}"
    echo "Region: ${REGION}"
    echo "AWS Account: ${AWS_ID}"
    echo ""
    echo "Your S3 bucket for Terraform state is ready to use!"
    echo ""
    echo "Example Terraform backend configuration:"
    echo "terraform {"
    echo "  backend \"s3\" {"
    echo "    bucket = \"${BUCKET}\""
    echo "    key    = \"path/to/your/terraform.tfstate\""
    echo "    region = \"${REGION}\""
    echo "  }"
    echo "}"
}

#######################################
# Main execution function
#######################################
function main {
    # Parse command line arguments
    parse_arguments "$@"

    # Validate dependencies
    validate_dependencies "aws" "jq"

    # Check AWS credentials before any AWS CLI usage
    if ! check_aws_credentials; then
        error_exit "AWS credentials are not set or invalid."
    fi

    # Generate unique bucket name if requested
    generate_bucket_name

    # Get AWS account information
    get_aws_account_info

    # Create S3 bucket
    create_s3_bucket

    # Configure bucket security settings
    configure_bucket_security

    # Verify bucket configuration
    verify_bucket_config

    # Apply bucket policy
    apply_bucket_policy

    # Display completion summary
    display_completion_summary
}

# Only call main function if script is executed directly, not sourced
if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi
