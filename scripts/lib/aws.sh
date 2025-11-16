#!/bin/bash
#######################################
# Description: AWS-specific utility functions for shell scripts
# Usage: source /path/to/scripts/lib/aws.sh
#
# This library provides AWS-related functions:
# - JSON parsing with jq
# - AWS CLI wrapper functions
# - Resource ARN parsing utilities
# - AWS region and account utilities
# - AWS credentials validation (check_aws_credentials)
#######################################

# Ensure common.sh is loaded for logging functions
if ! declare -f log > /dev/null 2>&1; then
    # Try to source common.sh from the same directory
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    # shellcheck source=./common.sh
    # shellcheck disable=SC1091
    source "$SCRIPT_DIR/common.sh"
fi

#######################################
# Function to safely execute AWS CLI command with error handling
# Arguments:
#   $@ - AWS CLI command and arguments
# Outputs:
#   Command output on success, error message on failure
# Returns:
#   0 on success, 1 on error
#######################################
function aws_safe_exec {
    local cmd="$*"
    local output
    local exit_code

    log "DEBUG" "Executing AWS CLI: $cmd"

    if output=$(eval "$cmd" 2>&1); then
        echo "$output"
        return 0
    else
        exit_code=$?
        log "ERROR" "AWS CLI command failed (exit code: $exit_code): $cmd"
        log "ERROR" "Output: $output"
        return $exit_code
    fi
}

#######################################
# Function to check AWS CLI credentials and identity
# Returns:
#   0 if credentials are valid
#   1 if credentials are missing or invalid
# Outputs:
#   Logs AWS identity on success, error on failure
#######################################
function check_aws_credentials {
    local identity
    identity=$(aws sts get-caller-identity --query 'Arn' --output text 2> /dev/null)
    if [[ -z "$identity" || "$identity" == "null" ]]; then
        log "ERROR" "AWS credentials are not set or invalid."
        return 1
    fi
    log "INFO" "AWS identity: $identity"
    return 0
}

#######################################
# Function to extract array values using jq and join with comma
# Arguments:
#   $1 - JSON data
#   $2 - jq query for array
#   $3 - default value (optional, defaults to "N/A")
#   $4 - separator (optional, defaults to ",")
# Outputs:
#   Comma-separated values (wrapped in quotes for CSV) or custom separator-separated values or default
# Returns:
#   0 on success
#######################################
function extract_jq_array {
    local json_data="$1"
    local jq_query="$2"
    local default_value="${3:-N/A}"
    local separator="${4:-,}"

    # Handle empty JSON data
    if [[ -z "$json_data" || "$json_data" == "null" ]]; then
        echo "$default_value"
        return 0
    fi

    # Execute jq query to get array and join with separator
    local result
    if result=$(echo "$json_data" | jq -r "${jq_query} | if type == \"array\" then join(\"${separator}\") else . end" 2> /dev/null); then
        # Handle null, empty, or array results
        if [[ -z "$result" || "$result" == "null" || "$result" == "[]" ]]; then
            echo "$default_value"
        else
            # For comma separator, wrap result in double quotes for CSV compatibility
            if [[ "$separator" == "," ]]; then
                echo "\"$result\""
            else
                echo "$result"
            fi
        fi
    else
        log "DEBUG" "jq array query failed: $jq_query"
        echo "$default_value"
    fi
}

#######################################
# Function to extract value using jq with default value support
# Arguments:
#   $1 - JSON data
#   $2 - jq query
#   $3 - default value (optional, defaults to "N/A")
# Outputs:
#   Extracted value or default
# Returns:
#   0 on success
#######################################
function extract_jq_value {
    local json_data="$1"
    local jq_query="$2"
    local default_value="${3:-N/A}"

    # Handle empty JSON data
    if [[ -z "$json_data" || "$json_data" == "null" ]]; then
        echo "$default_value"
        return 0
    fi

    # Execute jq query and handle errors
    local result
    if result=$(echo "$json_data" | jq -r "$jq_query" 2> /dev/null); then
        # Handle null or empty results
        if [[ -z "$result" || "$result" == "null" ]]; then
            echo "$default_value"
        else
            echo "$result"
        fi
    else
        log "DEBUG" "jq query failed: $jq_query"
        echo "$default_value"
    fi
}

#######################################
# Function to convert Unix timestamp to readable date
# Arguments:
#   $1 - Unix timestamp (seconds or milliseconds)
# Outputs:
#   Formatted date string (YYYY-MM-DD HH:MM:SS)
# Returns:
#   0 on success, 1 on error
#######################################
function format_aws_timestamp {
    local timestamp="$1"

    # Handle empty or non-numeric input
    if [[ ! "$timestamp" =~ ^[0-9]+$ ]]; then
        echo "N/A"
        return 0
    fi

    # Convert milliseconds to seconds if needed
    if [[ ${#timestamp} -gt 10 ]]; then
        timestamp=$((timestamp / 1000))
    fi

    # Format timestamp
    if date -d "@$timestamp" '+%Y-%m-%d %H:%M:%S' 2> /dev/null; then
        return 0
    else
        echo "N/A"
        return 1
    fi
}

#######################################
# Function to get AWS account ID
# Outputs:
#   AWS account ID
# Returns:
#   0 on success, 1 on error (instead of exit)
#######################################
function get_aws_account_id {
    local account_id
    if ! account_id=$(aws sts get-caller-identity --query Account --output text 2> /dev/null); then
        log "ERROR" "Failed to get AWS account ID. Check AWS credentials."
        return 1
    fi
    echo "$account_id"
    return 0
}

#######################################
# Function to get current AWS region
# Outputs:
#   Current AWS region
# Returns:
#   0 on success, exits on error
#######################################
function get_aws_region {
    local region
    local region_from_aws
    # aws configure get region returns empty string (exit 0) when not set, so capture value and test
    region_from_aws=$(aws configure get region 2> /dev/null || true)
    if [[ -n "$region_from_aws" ]]; then
        echo "$region_from_aws"
        return 0
    fi

    # Try to get from instance metadata if running on EC2
    local region_from_meta
    region_from_meta=$(curl -s --max-time 2 http://169.254.169.254/latest/meta-data/placement/region 2> /dev/null || true)
    if [[ -n "$region_from_meta" ]]; then
        echo "$region_from_meta"
        return 0
    fi

    error_exit "Failed to determine AWS region. Set AWS_DEFAULT_REGION or configure AWS CLI."
}

#######################################
# Function to get resource name from ARN
# Arguments:
#   $1 - AWS ARN
# Outputs:
#   Resource name (last part of ARN)
# Returns:
#   0 on success, 1 on error
#######################################
function get_resource_name_from_arn {
    local arn="$1"
    local resource_part

    if ! resource_part=$(parse_arn "$arn" | jq -r '.resource'); then
        return 1
    fi

    # Handle different resource formats
    if [[ "$resource_part" == */* ]]; then
        # Format: resource-type/resource-name
        echo "${resource_part##*/}"
    elif [[ "$resource_part" == *:* ]]; then
        # Format: resource-type:resource-name
        echo "${resource_part##*:}"
    else
        # Simple resource name
        echo "$resource_part"
    fi
}

#######################################
# Function to resolve a Security Group ID to a human-friendly name
# Arguments:
#   $1 - Security Group ID (e.g. sg-0123456789abcdef0)
#   $2 - AWS region (optional, defaults to current region)
# Outputs:
#   Security Group name (Tag Name or GroupName) or the original ID if not resolvable
# Returns:
#   0 on success, 1 on error
#######################################
function get_security_group_name {
    local sg_id="$1"
    local region="${2:-$(get_aws_region)}"

    if [[ -z "$sg_id" ]]; then
        echo "N/A"
        return 1
    fi

    # Only attempt to resolve well-formed SG IDs
    if [[ ! "$sg_id" =~ ^sg-[0-9a-fA-F]+$ ]]; then
        # Not an SG id; return as-is
        echo "$sg_id"
        return 0
    fi

    local cmd out
    cmd=(aws ec2 describe-security-groups --group-ids "$sg_id" --region "$region" --output json)
    if out=$(aws_safe_exec "${cmd[*]}" 2> /dev/null); then
        # Prefer Tag 'Name' if present, otherwise GroupName
        local sg_name
        sg_name=$(echo "$out" | jq -r '.SecurityGroups[0] | (.Tags[]? | select(.Key=="Name") | .Value) // .GroupName // ""' 2> /dev/null || true)
        if [[ -n "$sg_name" ]]; then
            echo "$sg_name"
            return 0
        fi
        # Fallback to returning the ID
        echo "$sg_id"
        return 0
    else
        # If AWS call failed, return the ID to avoid breaking consumers
        echo "$sg_id"
        return 1
    fi
}

#######################################
# Function to resolve a VPC ID to a human-friendly name
# Arguments:
#   $1 - VPC ID (e.g. vpc-0123456789abcdef0)
#   $2 - AWS region (optional, defaults to current region)
# Outputs:
#   VPC Name (Tag 'Name') or the original VPC ID if not resolvable
# Returns:
#   0 on success, 1 on error
#######################################
function get_vpc_name {
    local vpc_id="$1"
    local region="${2:-$(get_aws_region)}"

    if [[ -z "$vpc_id" ]]; then
        echo "N/A"
        return 1
    fi

    # Only attempt to resolve well-formed VPC IDs
    if [[ ! "$vpc_id" =~ ^vpc-[0-9a-fA-F]+$ ]]; then
        # Not a VPC id; return as-is
        echo "$vpc_id"
        return 0
    fi

    local cmd out
    cmd=(aws ec2 describe-vpcs --vpc-ids "$vpc_id" --region "$region" --output json)
    if out=$(aws_safe_exec "${cmd[*]}" 2> /dev/null); then
        # Prefer Tag 'Name' if present, otherwise fall back to VpcId
        local vpc_name
        vpc_name=$(echo "$out" | jq -r '.Vpcs[0] | (.Tags[]? | select(.Key=="Name") | .Value) // .VpcId // ""' 2> /dev/null || true)
        if [[ -n "$vpc_name" ]]; then
            echo "$vpc_name"
            return 0
        fi
        # Fallback to returning the ID
        echo "$vpc_id"
        return 0
    else
        # If AWS call failed, return the ID to avoid breaking consumers
        echo "$vpc_id"
        return 1
    fi
}

#######################################
# Function to resolve a Subnet ID to a human-friendly name
# Arguments:
#   $1 - Subnet ID (e.g. subnet-0123456789abcdef0)
#   $2 - AWS region (optional, defaults to current region)
# Outputs:
#   Subnet Name (Tag 'Name') or the original Subnet ID if not resolvable
# Returns:
#   0 on success, 1 on error
#######################################
function get_subnet_name {
    local subnet_id="$1"
    local region="${2:-$(get_aws_region)}"

    if [[ -z "$subnet_id" ]]; then
        echo "N/A"
        return 1
    fi

    # Only attempt to resolve well-formed Subnet IDs
    if [[ ! "$subnet_id" =~ ^subnet-[0-9a-fA-F]+$ ]]; then
        # Not a Subnet id; return as-is
        echo "$subnet_id"
        return 0
    fi

    local cmd out
    cmd=(aws ec2 describe-subnets --subnet-ids "$subnet_id" --region "$region" --output json)
    if out=$(aws_safe_exec "${cmd[*]}" 2> /dev/null); then
        # Prefer Tag 'Name' if present, otherwise fall back to SubnetId
        local subnet_name
        subnet_name=$(echo "$out" | jq -r '.Subnets[0] | (.Tags[]? | select(.Key=="Name") | .Value) // .SubnetId // ""' 2> /dev/null || true)
        if [[ -n "$subnet_name" ]]; then
            echo "$subnet_name"
            return 0
        fi
        # Fallback to returning the ID
        echo "$subnet_id"
        return 0
    else
        # If AWS call failed, return the ID to avoid breaking consumers
        echo "$subnet_id"
        return 1
    fi
}

#######################################
# Function to resolve a WAFv2 WebACL ARN to its human friendly Name
# Arguments:
#   $1 - WAF WebACL ARN (e.g. arn:aws:wafv2:us-east-1:123456789012:regional/webacl/MyWebACL/uuid)
#   $2 - optional region (falls back to get_aws_region)
# Outputs:
#   WebACL Name or ARN (if unresolved)
# Returns:
#   0 on success (found name), 1 on failure (returns ARN or N/A)
#######################################
function get_waf_name {
    local waf_arn="$1"
    local region="${2:-$(get_aws_region)}"

    if [[ -z "$waf_arn" || "$waf_arn" == "N/A" ]]; then
        echo "N/A"
        return 1
    fi

    # Try to extract the name directly from ARN: /webacl/<name>/...
    local waf_name
    waf_name=$(echo "$waf_arn" | sed -n 's#.*/webacl/\([^/]*\)/.*#\1#p' || true)
    if [[ -n "$waf_name" ]]; then
        echo "$waf_name"
        return 0
    fi

    # Fallback: query list-web-acls for REGIONAL and CLOUDFRONT scopes and match ARN
    for scope in REGIONAL CLOUDFRONT; do
        local out
        # WAF CLOUDFRONT scope is global and should be queried in us-east-1 for reliability
        local list_region
        if [[ "$scope" == "CLOUDFRONT" ]]; then
            list_region="us-east-1"
        else
            list_region="$region"
        fi

        out=$(aws wafv2 list-web-acls --scope "$scope" --region "$list_region" --output json 2> /dev/null || echo '{}')
        waf_name=$(echo "$out" | jq -r --arg arn "$waf_arn" '.WebACLs[]? | select(.ARN==$arn) | .Name' 2> /dev/null || true)
        if [[ -n "$waf_name" && "$waf_name" != "null" ]]; then
            echo "$waf_name"
            return 0
        fi
    done

    # As a last resort, return the ARN so callers have something to show
    echo "$waf_arn"
    return 1
}

#######################################
# Function to get WAF association for a resource
# Arguments:
#   $1 - Resource ARN
#   $2 - AWS region
# Outputs:
#   WAF Web ACL ARN or "N/A"
# Returns:
#   0 on success
#######################################
function get_waf_association {
    local resource_arn="$1"
    local region="$2"

    # Try WAFv2 first, then WAF Classic
    local waf_result
    waf_result=$(aws wafv2 get-web-acl-for-resource --resource-arn "$resource_arn" --region "$region" 2> /dev/null || echo '{}')

    extract_jq_value "$waf_result" '.WebACL.ARN'
}

#######################################
# Function to check if AWS service is available in region
# Arguments:
#   $1 - AWS service name
#   $2 - AWS region (optional, uses current region)
# Returns:
#   0 if service is available, 1 if not
#######################################
function is_service_available_in_region {
    local service="$1"
    local region="${2:-$(get_aws_region)}"

    # Some services are global
    case "$service" in
        iam | cloudfront | route53 | waf)
            return 0
            ;;
    esac

    # Check if service is available in region by attempting to list resources
    case "$service" in
        ec2)
            # Check that the region exists in the aws ec2 describe-regions output
            if aws ec2 describe-regions --output json 2> /dev/null | jq -r '.Regions[].RegionName' 2> /dev/null | grep -qx "$region"; then
                return 0
            else
                return 1
            fi
            ;;
        s3)
            # S3 is global but we can attempt a head-bucket call is not appropriate here; assume S3 generally available
            return 0
            ;;
        lambda)
            # Try a lightweight call to list-functions in the target region
            if aws lambda list-functions --region "$region" --max-items 1 > /dev/null 2>&1; then
                return 0
            else
                return 1
            fi
            ;;
        *)
            # Generic check - try to call a basic list operation
            log "DEBUG" "Service availability check not implemented for: $service"
            return 0
            ;;
    esac
}

#######################################
# Function to parse ARN components
# Arguments:
#   $1 - AWS ARN
# Outputs:
#   JSON object with ARN components
# Returns:
#   0 on success, 1 on invalid ARN
#######################################
function parse_arn {
    local arn="$1"

    # Validate ARN format
    if [[ ! "$arn" =~ ^arn:aws[^:]*:[^:]*:[^:]*:[^:]*:.+ ]]; then
        log "ERROR" "Invalid ARN format: $arn"
        return 1
    fi

    # Split ARN into components
    IFS=':' read -ra ARN_PARTS <<< "$arn"

    # Create JSON object
    jq -n \
        --arg partition "${ARN_PARTS[1]}" \
        --arg service "${ARN_PARTS[2]}" \
        --arg region "${ARN_PARTS[3]}" \
        --arg account "${ARN_PARTS[4]}" \
        --arg resource "${ARN_PARTS[5]}" \
        '{
            partition: $partition,
            service: $service,
            region: $region,
            account: $account,
            resource: $resource
        }'
}

#######################################
# Function to validate AWS CLI configuration
# Returns:
#   0 if AWS CLI is properly configured, exits on error
#######################################
function validate_aws_config {
    validate_dependencies "aws" "jq"

    # Test AWS CLI access
    if ! aws sts get-caller-identity > /dev/null 2>&1; then
        error_exit "AWS CLI is not properly configured. Run 'aws configure' or set AWS credentials."
    fi

    log "INFO" "AWS CLI configuration validated"
}
