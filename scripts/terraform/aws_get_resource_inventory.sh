#!/bin/bash
# shellcheck disable=SC2004  # Disable complexity warnings only
#######################################
# Description: Collect AWS resource information and output as CSV with resource categories.
# Usage: ./aws_get_resource_inventory.sh [-h] [-v] [-d] [-r REGION] [-c CATEGORIES] [-p]
#   options:
#     -h, --help       Display this help message
#     -v, --verbose    Enable verbose output
#     -d, --dry-run    Run in dry-run mode (no changes made)
#     -r, --region     AWS region to use (default: $AWS_DEFAULT_REGION or ap-northeast-1)
#     -c, --categories Comma-separated list of categories to collect (optional)
#     -p, --preserve-newlines  Preserve newlines in CSV output (better for Excel/Numbers)
#
# Output:
# - Generates CSV file containing AWS resource inventory across multiple resource types
# - Groups output by resource categories (alb, apigateway, cloudfront, etc.)
# - CSV output is Excel/Numbers compatible with proper quoting for fields containing commas/newlines
# - Only includes category headers when resources exist for that category
# - Uses standardized column structure with clear column names
#
# CSV Column Order Standards:
# - Core columns (always first): Category,Subcategory,Subsubcategory,Region,Name
# - Primary identifier: ARN/ID (immediately after Name)
# - Related ARNs: Role,TaskDefinition (when applicable, always as full ARNs)
# - Resource-specific attributes: Type,Performance,etc. (functional properties)
# - Status/operational data: Status,State,etc. (moved to later columns)
# - Date/time data: Created_Date,Last_Modified_Date,etc. (last columns)
# - Schedule/timing data: CronSchedule,etc. (last)
# - Removed columns: Tasks_Running (eliminated from all outputs)
# - Multiple resource types: Use consistent columns with empty values for non-applicable fields
# - Value policy: Empty for non-applicable fields, "N/A" only when data retrieval fails
#
# CSV Format Notes:
# - By default, newlines are sanitized to spaces for maximum compatibility
# - Use --preserve-newlines option to keep newlines for better Excel/Numbers display
# - All values with commas, quotes, or newlines are properly quoted per RFC 4180
# - Works with Excel, Numbers, Google Sheets, and other CSV-compatible applications
#
# Design Rules:
# - Modular design with separate collection functions per resource type
# - Extensible to add new resource types easily
# - Process substitution used for buffer population to maintain performance
# - Category-based output with consistent column alignment
# - Resource categories sorted alphabetically (A-Z order)
# - No empty category headers - headers only appear when resources exist
# - All functions follow consistent naming pattern: collect_<type>_inventory
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
OUTPUT_FILE="aws_resource_inventory.csv"
AWS_REGION="${AWS_DEFAULT_REGION:-ap-northeast-1}"
CATEGORIES=""
REGIONS_TO_CHECK=()
SORT_OUTPUT=true
PRESERVE_NEWLINES=false

# AWS resource categories list (A-Z order)
AWS_RESOURCE_CATEGORIES=(
    "acm"
    "alb"
    "apigateway"
    "batch"
    "bedrock"
    "cloudfront"
    "cognito"
    "cloudwatch_alarms"
    "cloudwatch_logs"
    "ec2"
    "ecr"
    "ecs"
    "efs"
    "eventbridge"
    "glue"
    "iam"
    "kinesis"
    "kms"
    "lambda"
    "quicksight"
    "rds"
    "redshift"
    "route53"
    "s3"
    "secretsmanager"
    "sns"
    "sqs"
    "transferfamily"
    "vpc"
    "waf"
)

# Categories that need to maintain grouping structure (no sorting)
NO_SORT_CATEGORIES=(
    "acm"     # Sort by DomainName within function
    "alb"     # Maintain LoadBalancer grouping structure
    "ecs"     # Maintain cluster grouping structure
    "lambda"  # Maintain environment variables integrity (quoted CSV fields)
    "rds"     # Maintain cluster grouping structure
    "route53" # Maintain HostedZone grouping structure
    "vpc"     # Maintain VPC grouping structure
)

#######################################
# Display usage information
#######################################
function show_usage {
    echo "Usage: $(basename "$0") [options]"
    echo ""
    echo "Description: Generate AWS resource inventory in CSV format"
    echo ""
    echo "Options:"
    echo "  -h, --help       Display this help message"
    echo "  -v, --verbose    Enable verbose output"
    echo "  -o, --output     Specify output file path (default: aws_resource_inventory.csv)"
    echo "  -r, --region     AWS region to query (default: \$AWS_DEFAULT_REGION or ap-northeast-1)"
    echo "  -d, --dry-run    Run in dry-run mode (no changes made)"
    echo "  -t, --test      Enable test mode (header names as values for automated testing)"
    echo "  -c, --categories Comma-separated list of categories to collect (optional)"
    echo "  -n, --no-sort    Disable sorting for all outputs (preserve original order)"
    echo "  -p, --preserve-newlines  Preserve newlines in CSV output (better for Excel/Numbers)"
    echo ""
    echo "Available categories:"
    echo "  acm, alb, apigateway, batch, bedrock, cloudfront, cognito, ec2, ecr, ecs,"
    echo "  efs, glue, iam, kms, lambda, quicksight, rds, redshift, route53, s3,"
    echo "  secretsmanager, sns, sqs, transferfamily, vpc, waf"
    echo ""
    echo "Examples:"
    echo "  $(basename "$0") -v -o my_inventory.csv"
    echo "  $(basename "$0") -r us-east-1 -o us_inventory.csv"
    echo "  $(basename "$0") -p -c s3,route53  # Preserve newlines for Route53 TXT records"
    echo "  $(basename "$0") -c vpc,s3 -o vpc_s3_only.csv"
    echo "  $(basename "$0") -n -o unsorted_inventory.csv"
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
            -t | --test)
                # Test mode for header validation (reserved for future use)
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
            -n | --no-sort)
                SORT_OUTPUT=false
                shift
                ;;
            -p | --preserve-newlines)
                PRESERVE_NEWLINES=true
                shift
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
# Function to initialize regions to check
#######################################
function initialize_regions {
    REGIONS_TO_CHECK=("$AWS_REGION")
    if [[ "$AWS_REGION" != "us-east-1" ]]; then
        REGIONS_TO_CHECK+=("us-east-1")
    fi
    log "INFO" "Regions to check: ${REGIONS_TO_CHECK[*]}"
}

#######################################
# Function to collect ACM inventory (with categories)
#######################################
function collect_acm_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ARN,Status,Type,Created_Date"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    while IFS= read -r cert_data; do
        [[ -z "$cert_data" ]] && continue
        local cert_arn cert_name cert_status cert_type cert_domain cert_created
        cert_arn=$(extract_jq_value "$cert_data" '.CertificateArn')
        cert_domain=$(extract_jq_value "$cert_data" '.DomainName')
        cert_status=$(extract_jq_value "$cert_data" '.Status')
        cert_type=$(extract_jq_value "$cert_data" '.Type')
        cert_name="${cert_domain}"

        # Get certificate details for creation date
        local cert_details cert_arn
        cert_arn=$(extract_jq_value "$cert_data" '.CertificateArn')
        cert_details=$(aws acm describe-certificate --certificate-arn "$cert_arn" --region "$region" 2>/dev/null || echo '{}')
        cert_created=$(extract_jq_value "$cert_details" '.Certificate.CreatedAt')

        buffer+="acm,Certificate,,${region},${cert_name},${cert_arn},${cert_status},${cert_type},${cert_created}\n"
    done < <(aws acm list-certificates --region "$region" | jq -c '.CertificateSummaryList | sort_by(.DomainName)[] ')

    echo "$buffer"
}

#######################################
# Function to collect ALB inventory (with categories)
#######################################
function collect_alb_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ARN,DNS_Name,Type,WAF,Protocol,Port,HealthCheck,SSLPolicy,State,Created_Date"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    while IFS= read -r alb_data; do
        [[ -z "$alb_data" ]] && continue
        local alb_arn alb_name alb_dns alb_type alb_state alb_waf alb_created
        alb_arn=$(extract_jq_value "$alb_data" '.LoadBalancerArn')
        alb_name=$(extract_jq_value "$alb_data" '.LoadBalancerName')
        alb_dns=$(extract_jq_value "$alb_data" '.DNSName')
        alb_type=$(extract_jq_value "$alb_data" '.Type')
        alb_state=$(extract_jq_value "$alb_data" '.State.Code')
        alb_created=$(extract_jq_value "$alb_data" '.CreatedTime')
        alb_waf=$(extract_jq_value "$(aws wafv2 get-web-acl-for-resource --resource-arn "$alb_arn" --region "$region" 2>/dev/null || echo '{}')" '.WebACL.ARN')
        buffer+="alb,LoadBalancer,,${region},$alb_name,$alb_arn,$alb_dns,$alb_type,$alb_waf,,,,,$alb_state,$alb_created\n"

        # Target Groups
        while IFS= read -r tg_data; do
            [[ -z "$tg_data" ]] && continue
            local tg_arn tg_name tg_type tg_protocol tg_port tg_health_check
            tg_arn=$(extract_jq_value "$tg_data" '.TargetGroupArn')
            tg_name=$(extract_jq_value "$tg_data" '.TargetGroupName')
            tg_type=$(extract_jq_value "$tg_data" '.TargetType')
            tg_protocol=$(extract_jq_value "$tg_data" '.Protocol')
            tg_port=$(extract_jq_value "$tg_data" '.Port')
            tg_health_check=$(extract_jq_value "$tg_data" '.HealthCheckPath')
            buffer+="alb,,TargetGroup,${region},$tg_name,$tg_arn,,$tg_type,,$tg_protocol,$tg_port,$tg_health_check,,,\n"
        done < <(aws elbv2 describe-target-groups --load-balancer-arn "$alb_arn" --region "$region" 2>/dev/null | jq -c '.TargetGroups[]?' || echo "")

        # Listeners
        while IFS= read -r listener_data; do
            [[ -z "$listener_data" ]] && continue
            local listener_arn listener_protocol listener_port listener_ssl_policy
            listener_arn=$(extract_jq_value "$listener_data" '.ListenerArn')
            listener_protocol=$(extract_jq_value "$listener_data" '.Protocol')
            listener_port=$(extract_jq_value "$listener_data" '.Port')
            listener_ssl_policy=$(extract_jq_value "$listener_data" '.SslPolicy')
            buffer+="alb,,Listener,${region},${listener_protocol}:${listener_port},$listener_arn,,,,$listener_protocol,$listener_port,,$listener_ssl_policy,,\n"
        done < <(aws elbv2 describe-listeners --load-balancer-arn "$alb_arn" --region "$region" 2>/dev/null | jq -c '.Listeners[]?' || echo "")
    done < <(aws elbv2 describe-load-balancers --region "$region" | jq -c '.LoadBalancers[]')

    echo "$buffer"
}

#######################################
# Function to collect API Gateway inventory (with categories)
#######################################
function collect_apigateway_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ID,Protocol_Type,WAF"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    # REST APIs
    while IFS= read -r api_data; do
        [[ -z "$api_data" ]] && continue
        local api_id api_name api_type api_waf
        api_id=$(extract_jq_value "$api_data" '.id')
        api_name=$(extract_jq_value "$api_data" '.name')
        api_type="REST"

        # Get WAF association for REST API - directly from stage data
        api_waf=""
        local stages_data
        stages_data=$(aws apigateway get-stages --rest-api-id "$api_id" --region "$region" 2>/dev/null || echo '{"item":[]}')
        if [[ "$stages_data" != '{"item":[]}' ]]; then
            api_waf=$(extract_jq_value "$stages_data" '.item[]?.webAclArn')
        fi

        buffer+="apigateway,RestAPI,,${region},$api_name,$api_id,$api_type,$api_waf\n"
    done < <(aws apigateway get-rest-apis --region "$region" | jq -c '.items[]')

    # HTTP APIs
    while IFS= read -r api_data; do
        [[ -z "$api_data" ]] && continue
        local api_id api_name api_type api_waf
        api_id=$(extract_jq_value "$api_data" '.ApiId')
        api_name=$(extract_jq_value "$api_data" '.Name')
        api_type=$(extract_jq_value "$api_data" '.ProtocolType')

        # Get WAF association for HTTP API - directly from stage data
        api_waf=""
        local http_stages_data
        http_stages_data=$(aws apigatewayv2 get-stages --api-id "$api_id" --region "$region" 2>/dev/null || echo '{"Items":[]}')
        if [[ "$http_stages_data" != '{"Items":[]}' ]]; then
            api_waf=$(extract_jq_value "$http_stages_data" '.Items[]?.WebAclArn')
        fi

        buffer+="apigateway,HttpAPI,,${region},$api_name,$api_id,$api_type,$api_waf\n"
    done < <(aws apigatewayv2 get-apis --region "$region" | jq -c '.Items[]')

    echo "$buffer"
}

#######################################
# Function to collect Batch inventory (with categories)
#######################################
function collect_batch_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ARN,Type_Priority,Status"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    while IFS= read -r queue_data; do
        [[ -z "$queue_data" ]] && continue
        local queue_name queue_arn queue_state queue_priority
        queue_name=$(extract_jq_value "$queue_data" '.jobQueueName')
        queue_arn=$(extract_jq_value "$queue_data" '.jobQueueArn')
        queue_state=$(extract_jq_value "$queue_data" '.state')
        queue_priority=$(extract_jq_value "$queue_data" '.priority')
        buffer+="batch,JobQueue,,${region},$queue_name,$queue_arn,$queue_priority,$queue_state\n"
    done < <(aws batch describe-job-queues --region "$region" | jq -c '.jobQueues[]')

    while IFS= read -r compute_data; do
        [[ -z "$compute_data" ]] && continue
        local compute_name compute_arn compute_state compute_type
        compute_name=$(extract_jq_value "$compute_data" '.computeEnvironmentName')
        compute_arn=$(extract_jq_value "$compute_data" '.computeEnvironmentArn')
        compute_state=$(extract_jq_value "$compute_data" '.state')
        compute_type=$(extract_jq_value "$compute_data" '.type')
        buffer+="batch,ComputeEnvironment,,${region},$compute_name,$compute_arn,$compute_type,$compute_state\n"
    done < <(aws batch describe-compute-environments --region "$region" | jq -c '.computeEnvironments[]')

    echo "$buffer"
}

#######################################
# Function to collect Bedrock inventory (with categories)
#######################################
function collect_bedrock_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,Identifier,Provider,Input_Modalities,Output_Modalities"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    # List foundation models
    while IFS= read -r model_data; do
        [[ -z "$model_data" ]] && continue
        local model_id model_name model_provider model_input_modalities model_output_modalities
        model_id=$(extract_jq_value "$model_data" '.modelId')
        model_name=$(extract_jq_value "$model_data" '.modelName')
        model_provider=$(extract_jq_value "$model_data" '.providerName')
        model_input_modalities=$(extract_jq_array "$model_data" '.inputModalities')
        model_output_modalities=$(extract_jq_array "$model_data" '.outputModalities')

        buffer+="bedrock,FoundationModel,,${region},$model_name,$model_id,$model_provider,$model_input_modalities,$model_output_modalities\n"
        # Removed count increment (not used)
    done < <(aws bedrock list-foundation-models --region "$region" 2>/dev/null | jq -c '.modelSummaries[]?' || true)

    # List custom models
    while IFS= read -r custom_model_data; do
        [[ -z "$custom_model_data" ]] && continue
        local custom_model_arn custom_model_name custom_model_status
        custom_model_arn=$(extract_jq_value "$custom_model_data" '.modelArn')
        custom_model_name=$(extract_jq_value "$custom_model_data" '.modelName')
        custom_model_status=$(extract_jq_value "$custom_model_data" '.status')

        buffer+="bedrock,CustomModel,,${region},$custom_model_name,$custom_model_arn,$custom_model_status,\n"
        # Removed count increment (not used)
    done < <(aws bedrock list-custom-models --region "$region" 2>/dev/null | jq -c '.modelSummaries[]?' || true)

    # List model customization jobs
    while IFS= read -r job_data; do
        [[ -z "$job_data" ]] && continue
        local job_arn job_name job_status
        job_arn=$(extract_jq_value "$job_data" '.jobArn')
        job_name=$(extract_jq_value "$job_data" '.jobName')
        job_status=$(extract_jq_value "$job_data" '.status')

        buffer+="bedrock,CustomizationJob,,${region},$job_name,$job_arn,$job_status,\n"
        # Removed count increment (not used)
    done < <(aws bedrock list-model-customization-jobs --region "$region" 2>/dev/null | jq -c '.modelCustomizationJobSummaries[]?' || true)

    echo "$buffer"
}

#######################################
# Function to collect CloudFront inventory (with categories)
#######################################
function collect_cloudfront_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ID,Alternate_Domain,Origin,Price_Class,WAF,Status"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    # CloudFront is a global service, only process from us-east-1
    if [[ "$region" != "us-east-1" ]]; then
        echo ""
        return 0
    fi

    local buffer=""

    while IFS= read -r dist_data; do
        [[ -z "$dist_data" ]] && continue
        local dist_id dist_domain dist_status dist_price_class dist_aliases dist_origin dist_waf_arn
        dist_id=$(extract_jq_value "$dist_data" '.Id')
        dist_domain=$(extract_jq_value "$dist_data" '.DomainName')
        dist_status=$(extract_jq_value "$dist_data" '.Status')
        dist_price_class=$(extract_jq_value "$dist_data" '.PriceClass')
        dist_aliases=$(extract_jq_value "$dist_data" '.Aliases.Items | join(";")')
        dist_origin=$(extract_jq_value "$dist_data" '.Origins.Items[0].DomainName')

        # Get WAF WebACLId from distribution config
        local dist_config_json dist_web_acl_id
        dist_config_json=$(aws cloudfront get-distribution-config --id "$dist_id" 2>/dev/null)
        dist_web_acl_id=$(extract_jq_value "$dist_config_json" '.DistributionConfig.WebACLId')
        dist_waf_arn="$dist_web_acl_id"

        buffer+="cloudfront,Distribution,,Global,$dist_domain,$dist_id,$dist_aliases,$dist_origin,$dist_price_class,$dist_waf_arn,$dist_status\n"
    done < <(aws cloudfront list-distributions | jq -c '.DistributionList.Items[]')

    echo "$buffer"
}

#######################################
# Function to collect CloudWatch Alarms inventory
#######################################
function collect_cloudwatch_alarms_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ARN,MetricName,Namespace,Statistic,Threshold,ComparisonOperator,EvaluationPeriods,Period,TreatMissingData,State"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    while IFS= read -r alarm_data; do
        [[ -z "$alarm_data" ]] && continue
        local alarm_name alarm_arn alarm_state alarm_metric_name alarm_namespace alarm_statistic alarm_threshold alarm_comparison_operator alarm_evaluation_periods alarm_period alarm_treat_missing_data

        alarm_name=$(extract_jq_value "$alarm_data" '.AlarmName')
        alarm_arn=$(extract_jq_value "$alarm_data" '.AlarmArn')
        alarm_state=$(extract_jq_value "$alarm_data" '.StateValue')
        alarm_metric_name=$(extract_jq_value "$alarm_data" '.MetricName')
        alarm_namespace=$(extract_jq_value "$alarm_data" '.Namespace')
        alarm_statistic=$(extract_jq_value "$alarm_data" '.Statistic')
        alarm_threshold=$(extract_jq_value "$alarm_data" '.Threshold')
        alarm_comparison_operator=$(extract_jq_value "$alarm_data" '.ComparisonOperator')
        alarm_evaluation_periods=$(extract_jq_value "$alarm_data" '.EvaluationPeriods')
        alarm_period=$(extract_jq_value "$alarm_data" '.Period')
        alarm_treat_missing_data=$(extract_jq_value "$alarm_data" '.TreatMissingData')

        buffer+="cloudwatch,Alarm,,${region},$alarm_name,$alarm_arn,$alarm_metric_name,$alarm_namespace,$alarm_statistic,$alarm_threshold,$alarm_comparison_operator,$alarm_evaluation_periods,$alarm_period,$alarm_treat_missing_data,$alarm_state\n"

    done < <(aws cloudwatch describe-alarms --region "$region" 2>/dev/null | jq -c '.MetricAlarms[], .CompositeAlarms[]' || true)

    echo "$buffer"
}

#######################################
# Function to collect CloudWatch Logs inventory
#######################################
function collect_cloudwatch_logs_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ARN,RetentionInDays,StoredBytes,MetricFilterCount,SubscriptionFilterCount,KmsKeyId,Created_Date"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    while IFS= read -r log_group_data; do
        [[ -z "$log_group_data" ]] && continue
        local log_group_name log_group_arn log_group_retention log_group_stored_bytes log_group_metric_filter_count log_group_subscription_filter_count log_group_kms_key_id log_group_creation_time

        log_group_name=$(extract_jq_value "$log_group_data" '.logGroupName')
        log_group_arn=$(extract_jq_value "$log_group_data" '.arn')
        log_group_retention=$(extract_jq_value "$log_group_data" '.retentionInDays' 'Never Expire')
        log_group_stored_bytes=$(extract_jq_value "$log_group_data" '.storedBytes' '0')
        log_group_metric_filter_count=$(extract_jq_value "$log_group_data" '.metricFilterCount' '0')
        log_group_subscription_filter_count=$(extract_jq_value "$log_group_data" '.subscriptionFilterCount' '0')
        log_group_kms_key_id=$(extract_jq_value "$log_group_data" '.kmsKeyId')
        log_group_creation_time=$(extract_jq_value "$log_group_data" '.creationTime')

        # Convert Unix timestamp to readable format if needed
        if [[ "$log_group_creation_time" =~ ^[0-9]+$ ]]; then
            log_group_creation_time=$(date -d "@$((log_group_creation_time / 1000))" '+%Y-%m-%d %H:%M:%S' 2>/dev/null || echo "$log_group_creation_time")
        fi

        buffer+="cloudwatch,LogGroup,,${region},$log_group_name,$log_group_arn,$log_group_retention,$log_group_stored_bytes,$log_group_metric_filter_count,$log_group_subscription_filter_count,$log_group_kms_key_id,$log_group_creation_time\n"

    done < <(aws logs describe-log-groups --region "$region" 2>/dev/null | jq -c '.logGroups[]' || true)

    echo "$buffer"
}

#######################################
# Function to collect Cognito inventory (with categories)
#######################################
function collect_cognito_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ID,Status_AllowUnauthenticated,Created_Date"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    # UserPool - has Status and CreationDate
    while IFS= read -r pool_data; do
        [[ -z "$pool_data" ]] && continue
        local pool_id pool_name pool_status pool_creation_date
        pool_id=$(extract_jq_value "$pool_data" '.Id')
        pool_name=$(extract_jq_value "$pool_data" '.Name')
        pool_status=$(extract_jq_value "$pool_data" '.Status')
        pool_creation_date=$(extract_jq_value "$pool_data" '.CreationDate')
        buffer+="cognito,UserPool,,${region},$pool_name,$pool_id,$pool_status,$pool_creation_date\n"
    done < <(aws cognito-idp list-user-pools --max-results 60 --region "$region" | jq -c '.UserPools[]')

    # IdentityPool - has AllowUnauthenticatedIdentities
    while IFS= read -r pool_data; do
        [[ -z "$pool_data" ]] && continue
        local pool_id pool_name pool_allow_unauthenticated
        pool_id=$(extract_jq_value "$pool_data" '.IdentityPoolId')
        pool_name=$(extract_jq_value "$pool_data" '.IdentityPoolName')
        pool_allow_unauthenticated=$(extract_jq_value "$pool_data" '.AllowUnauthenticatedIdentities')
        buffer+="cognito,IdentityPool,,${region},$pool_name,$pool_id,,$pool_allow_unauthenticated\n"
    done < <(aws cognito-identity list-identity-pools --max-results 60 --region "$region" | jq -c '.IdentityPools[]')

    echo "$buffer"
}

#######################################
# Function to collect EC2 inventory (with categories)
#######################################
function collect_ec2_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,Instance_ID,Instance_Type,State"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    while IFS= read -r instance_data; do
        [[ -z "$instance_data" ]] && continue
        local instance_id instance_name instance_type instance_state
        instance_id=$(extract_jq_value "$instance_data" '.InstanceId')
        instance_type=$(extract_jq_value "$instance_data" '.InstanceType')
        instance_state=$(extract_jq_value "$instance_data" '.State.Name')
        instance_name=$(normalize_csv_value "$(extract_jq_value "$instance_data" '.Tags[]? | select(.Key=="Name") | .Value')")
        buffer+="ec2,Instance,,${region},${instance_name},${instance_id},${instance_type},${instance_state}\n"
    done < <(aws ec2 describe-instances --region "$region" | jq -c '.Reservations[].Instances[]')

    echo "$buffer"
}

#######################################
# Function to collect ECR inventory (with categories)
#######################################
function collect_ecr_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,URI,Image_Count,Created_Date"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    while IFS= read -r repo_data; do
        [[ -z "$repo_data" ]] && continue
        local repo_name repo_uri repo_created repo_image_count
        repo_name=$(extract_jq_value "$repo_data" '.repositoryName')
        repo_uri=$(extract_jq_value "$repo_data" '.repositoryUri')
        repo_created=$(extract_jq_value "$repo_data" '.createdAt')
        repo_image_count=$(extract_jq_value "$(aws ecr describe-images --repository-name "$repo_name" --region "$region" 2>/dev/null || echo '{}')" '.imageDetails | length' '0')
        buffer+="ecr,Repository,,${region},$repo_name,$repo_uri,$repo_image_count,$repo_created\n"
    done < <(aws ecr describe-repositories --region "$region" | jq -c '.repositories[]')

    echo "$buffer"
}

#######################################
# Function to collect ECS inventory (grouped by Cluster)
#######################################
function collect_ecs_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ARN,Role,TaskDefinition,LaunchType,Status,CronSchedule,Spec,RuntimePlatform,PortMappings,Environment"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    # Get AWS Account ID once for efficiency
    local aws_account_id
    if ! aws_account_id=$(get_aws_account_id); then
        aws_account_id="unknown"
        log "WARN" "Could not retrieve AWS account ID"
    fi

    # Get all clusters first
    local cluster_arns_raw
    cluster_arns_raw=$(aws ecs list-clusters --region "$region" --query 'clusterArns[]' --output text 2>/dev/null || true)
    if [[ -z "$cluster_arns_raw" || "$cluster_arns_raw" == "None" ]]; then
        echo ""
        return 0
    fi
    local cluster_arns
    cluster_arns=$(echo "$cluster_arns_raw" | tr '\t' '\n')

    # 事前にEventBridgeルールをregion単位で収集
    local eventbridge_rules
    eventbridge_rules=$(aws events list-rules --region "$region" --query "Rules[?State==\`ENABLED\`]" --output json 2>/dev/null || echo "[]")

    # 事前に全ScheduledTaskのターゲット情報を取得
    declare -A scheduled_tasks_by_cluster
    if [[ "$eventbridge_rules" != "[]" ]]; then
        while IFS= read -r rule_data; do
            [[ -z "$rule_data" ]] && continue
            local rule_name rule_schedule rule_state rule_arn
            rule_name=$(extract_jq_value "$rule_data" '.Name')
            rule_schedule=$(normalize_csv_value "$(extract_jq_value "$rule_data" '.ScheduleExpression')")
            rule_state=$(extract_jq_value "$rule_data" '.State')
            rule_arn="arn:aws:events:${region}:${aws_account_id}:rule/${rule_name}"
            local rule_targets
            rule_targets=$(aws events list-targets-by-rule --rule "$rule_name" --region "$region" 2>/dev/null || echo '{"Targets":[]}')
            # ECSターゲットを配列化してforループで処理
            mapfile -t ecs_targets < <(echo "$rule_targets" | jq -c '.Targets[]? | select(.EcsParameters?)')
            for ecs_target in "${ecs_targets[@]}"; do
                local cluster_arn task_def_arn task_launch_type task_def_name task_role_arn
                cluster_arn=$(extract_jq_value "$ecs_target" '.EcsParameters.ClusterArn // .Arn')
                task_def_arn=$(extract_jq_value "$ecs_target" '.EcsParameters.TaskDefinitionArn')
                task_launch_type=$(extract_jq_value "$ecs_target" '.EcsParameters.LaunchType')
                if [[ -n "$task_def_arn" && "$task_def_arn" != "N/A" ]]; then
                    task_def_name=$(basename "$task_def_arn" | cut -d':' -f1)
                    local task_def_details
                    task_def_details=$(aws ecs describe-task-definition --task-definition "$task_def_name" --region "$region" --query 'taskDefinition' 2>/dev/null || echo '{}')
                    task_role_arn=$(extract_jq_value "$task_def_details" '.taskRoleArn // .executionRoleArn')
                else
                    task_def_arn="N/A"
                    task_def_name="N/A"
                    task_role_arn="N/A"
                fi
                local scheduled_row="ecs,,ScheduledTask,${region},$rule_name,$rule_arn,$task_role_arn,$task_def_arn,$task_launch_type,$rule_state,$rule_schedule,,,,\n"
                if [[ -n "$cluster_arn" ]]; then
                    scheduled_tasks_by_cluster[$cluster_arn]+="$scheduled_row"
                fi
            done
        done < <(echo "$eventbridge_rules" | jq -c '.[]')
    fi

    # 各クラスタごとに出力
    while IFS= read -r cluster_arn; do
        [[ -z "$cluster_arn" || "$cluster_arn" == "None" ]] && continue
        local cluster_data
        cluster_data=$(aws ecs describe-clusters --clusters "$cluster_arn" --region "$region" --query 'clusters[0]' 2>/dev/null) || continue
        local cluster_name cluster_status
        cluster_name=$(extract_jq_value "$cluster_data" '.clusterName')
        cluster_status=$(extract_jq_value "$cluster_data" '.status')
        local current_cluster_output=""
        current_cluster_output+="ecs,Cluster,,${region},$cluster_name,$cluster_arn,,,,$cluster_status,,,,,\n"

        # Get services for this cluster
        local service_arns_raw
        service_arns_raw=$(aws ecs list-services --cluster "$cluster_arn" --region "$region" --query 'serviceArns[]' --output text 2>/dev/null || true)

        if [[ -n "$service_arns_raw" && "$service_arns_raw" != "None" ]]; then
            # Convert tab-separated output to newline-separated
            local service_arns
            service_arns=$(echo "$service_arns_raw" | tr '\t' '\n')

            while IFS= read -r service_arn; do
                [[ -z "$service_arn" || "$service_arn" == "None" ]] && continue

                # Get service details
                local service_data
                service_data=$(aws ecs describe-services --cluster "$cluster_arn" --services "$service_arn" --region "$region" --query 'services[0]' 2>/dev/null) || continue

                local service_name service_status service_desired service_running service_task_def service_role service_launch_type
                service_name=$(extract_jq_value "$service_data" '.serviceName')
                service_status=$(extract_jq_value "$service_data" '.status')
                service_desired=$(extract_jq_value "$service_data" '.desiredCount')
                service_running=$(extract_jq_value "$service_data" '.runningCount')
                service_launch_type=$(extract_jq_value "$service_data" '.launchType')
                local service_status_detail="${service_status} (${service_running}/${service_desired})"

                # Get task definition details
                local task_def_arn service_task_def service_role service_task_def_arn
                task_def_arn=$(extract_jq_value "$service_data" '.taskDefinition')

                if [[ -n "$task_def_arn" && "$task_def_arn" != "null" ]]; then
                    service_task_def_arn="$task_def_arn"
                    service_task_def=$(echo "$task_def_arn" | sed 's/.*task-definition\/\([^:]*\).*/\1/')

                    # Get task role from task definition
                    local task_def_details task_role_arn
                    task_def_details=$(aws ecs describe-task-definition --task-definition "$service_task_def" --region "$region" --query 'taskDefinition' 2>/dev/null || echo '{}')
                    task_role_arn=$(extract_jq_value "$task_def_details" '.taskRoleArn // .executionRoleArn')
                    service_role="$task_role_arn"
                else
                    service_task_def_arn="N/A"
                    service_task_def="N/A"
                    service_role="N/A"
                fi

                # Service row: Category,Subcategory,Subsubcategory,Region,Name,ARN,Role,TaskDefinition,LaunchType,Status,CronSchedule,Spec,RuntimePlatform
                # For services, CronSchedule and RuntimePlatform are empty
                current_cluster_output+="ecs,,Service,${region},$service_name,$service_arn,$service_role,$service_task_def_arn,$service_launch_type,$service_status_detail,,,,,\n"
            done <<<"$service_arns"
        fi

        # クラスタ直下にScheduledTask出力
        if [[ -n "${scheduled_tasks_by_cluster[$cluster_arn]:-}" ]]; then
            current_cluster_output+="${scheduled_tasks_by_cluster[$cluster_arn]}"
        fi
        buffer+="$current_cluster_output"
    done <<<"$cluster_arns"

    # Collect all Task Definition families (latest revision only, no duplicates)
    # Get unique task definition families first, then get latest revision for each
    declare -A task_def_families
    while IFS= read -r task_def_arn; do
        [[ -z "$task_def_arn" || "$task_def_arn" == "None" ]] && continue
        local family_name
        family_name=$(basename "$task_def_arn" | cut -d':' -f1)
        # Store only the latest revision for each family
        if [[ -z "${task_def_families[$family_name]:-}" ]]; then
            task_def_families[$family_name]="$task_def_arn"
        fi
    done < <(aws ecs list-task-definitions --region "$region" --query 'taskDefinitionArns[]' --output text 2>/dev/null | tr '\t' '\n')

    # Process each unique task definition family (latest revision only)
    for family_name in "${!task_def_families[@]}"; do
        local task_def_arn="${task_def_families[$family_name]}"

        # Get task definition details
        local task_def_details task_def_status task_role_arn task_def_cpu task_def_memory task_def_network_mode task_def_requires_attributes task_def_revision

        # Get detailed task definition information for the latest revision
        task_def_details=$(aws ecs describe-task-definition --task-definition "$family_name" --region "$region" --query 'taskDefinition' 2>/dev/null || echo '{}')

        if [[ "$task_def_details" != "{}" ]]; then
            task_def_status=$(extract_jq_value "$task_def_details" '.status')
            task_role_arn=$(extract_jq_value "$task_def_details" '.taskRoleArn // .executionRoleArn')
            task_def_cpu=$(extract_jq_value "$task_def_details" '.cpu')
            task_def_memory=$(extract_jq_value "$task_def_details" '.memory')
            task_def_network_mode=$(extract_jq_value "$task_def_details" '.networkMode')
            task_def_requires_attributes=$(extract_jq_array "$task_def_details" '.requiresAttributes[].name')
            task_def_revision=$(extract_jq_value "$task_def_details" '.revision')
            task_def_arn=$(extract_jq_value "$task_def_details" '.taskDefinitionArn')

            # Get runtime platform information
            local task_def_runtime_os task_def_runtime_arch task_def_runtime_platform
            task_def_runtime_os=$(extract_jq_value "$task_def_details" '.runtimePlatform.operatingSystemFamily')
            task_def_runtime_arch=$(extract_jq_value "$task_def_details" '.runtimePlatform.cpuArchitecture')

            # Format runtime platform as OS/ARCH or fallback to default
            if [[ -n "$task_def_runtime_os" && "$task_def_runtime_os" != "N/A" ]] && [[ -n "$task_def_runtime_arch" && "$task_def_runtime_arch" != "N/A" ]]; then
                task_def_runtime_platform="${task_def_runtime_os}/${task_def_runtime_arch}"
            else
                task_def_runtime_platform="LINUX/X86_64" # Default for ECS
            fi

            # Get port mappings information
            local task_def_port_mappings
            task_def_port_mappings=$(echo "$task_def_details" | jq -r '.containerDefinitions[]?.portMappings[]? | "\(.containerPort):\(.hostPort // "dynamic"):\(.protocol // "tcp")"' 2>/dev/null | paste -sd ',' - || echo "")

            # Apply normalize_csv_value to port mappings
            task_def_port_mappings=$(normalize_csv_value "$task_def_port_mappings")

            # Get environment variables information
            local task_def_environment
            task_def_environment=$(echo "$task_def_details" | jq -r '.containerDefinitions[]?.environment[]? | "\(.name)=\(.value)"' 2>/dev/null | paste -sd $'\n' - || echo "")

            # Apply normalize_csv_value to environment variables
            task_def_environment=$(normalize_csv_value "$task_def_environment")

            # Format requires attributes for display
            if [[ -n "$task_def_requires_attributes" && "$task_def_requires_attributes" != "N/A" ]]; then
                task_def_requires_attributes=$(echo "$task_def_requires_attributes" | tr ',' ';')
            fi

            # TaskDefinition row: Category,Subcategory,Subsubcategory,Region,Name,ARN,Role,TaskDefinition,LaunchType,Status,CronSchedule,Spec,RuntimePlatform,PortMappings,Environment
            # For TaskDefinitions: LaunchType is empty, Status shows actual status, CronSchedule is empty, Spec shows CPU/Memory/NetworkMode, RuntimePlatform shows OS/ARCH, PortMappings shows container:host:protocol, Environment shows name=value pairs
            local task_def_details_summary="${task_def_cpu}CPU/${task_def_memory}MB/${task_def_network_mode}"
            buffer+="ecs,TaskDefinition,,${region},$family_name:${task_def_revision},$task_def_arn,$task_role_arn,$task_def_arn,,$task_def_status,,$task_def_details_summary,$task_def_runtime_platform,$task_def_port_mappings,$task_def_environment\n"
        fi
    done

    echo "$buffer"
}

#######################################
# Function to collect EFS inventory (with categories)
#######################################
function collect_efs_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ID,Type,Performance,Throughput,Encrypted,Size,SubnetID,IPAddress,SecurityGroups,Path,UID,GID,State,Created_Date"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    # Get File Systems
    while IFS= read -r fs_data; do
        [[ -z "$fs_data" ]] && continue
        local fs_id fs_name fs_state fs_creation_time fs_size fs_performance_mode fs_throughput_mode fs_encrypted
        fs_id=$(extract_jq_value "$fs_data" '.FileSystemId')
        fs_name=$(extract_jq_value "$fs_data" '.Name')
        fs_state=$(extract_jq_value "$fs_data" '.LifeCycleState')
        fs_creation_time=$(extract_jq_value "$fs_data" '.CreationTime')
        fs_size=$(extract_jq_value "$fs_data" '.SizeInBytes.Value')
        fs_performance_mode=$(extract_jq_value "$fs_data" '.PerformanceMode')
        fs_throughput_mode=$(extract_jq_value "$fs_data" '.ThroughputMode')
        fs_encrypted=$(extract_jq_value "$fs_data" '.Encrypted')

        # FileSystem: Core info with empty fields for mount-target-specific data
        buffer+="efs,FileSystem,,${region},$fs_name,$fs_id,FileSystem,$fs_performance_mode,$fs_throughput_mode,$fs_encrypted,$fs_size,,,,,,$fs_state,$fs_creation_time\n"

        # Get Mount Targets for this File System
        while IFS= read -r mt_data; do
            [[ -z "$mt_data" ]] && continue
            local mt_id mt_subnet_id mt_ip mt_state mt_security_groups
            mt_id=$(extract_jq_value "$mt_data" '.MountTargetId')
            mt_subnet_id=$(extract_jq_value "$mt_data" '.SubnetId')
            mt_ip=$(extract_jq_value "$mt_data" '.IpAddress')
            mt_state=$(extract_jq_value "$mt_data" '.LifeCycleState')

            # Get security groups for this mount target
            mt_security_groups=$(extract_jq_value "$(aws efs describe-mount-target-security-groups --mount-target-id "$mt_id" --region "$region" 2>/dev/null || echo '{}')" '.SecurityGroups | join(",")')

            # MountTarget: Core info with empty fields for filesystem-specific and accesspoint-specific data
            buffer+="efs,MountTarget,,${region},$mt_id,$mt_id,MountTarget,,,,,$mt_subnet_id,$mt_ip,$mt_security_groups,,,,$mt_state,\n"
        done < <(aws efs describe-mount-targets --file-system-id "$fs_id" --region "$region" 2>/dev/null | jq -c '.MountTargets[]' || true)

        # Get Access Points for this File System
        while IFS= read -r ap_data; do
            [[ -z "$ap_data" ]] && continue
            local ap_id ap_name ap_state ap_path ap_uid ap_gid
            ap_id=$(extract_jq_value "$ap_data" '.AccessPointId')
            ap_name=$(extract_jq_value "$ap_data" '.Name')
            ap_state=$(extract_jq_value "$ap_data" '.LifeCycleState')
            ap_path=$(extract_jq_value "$ap_data" '.RootDirectory.Path')
            ap_uid=$(extract_jq_value "$ap_data" '.PosixUser.Uid')
            ap_gid=$(extract_jq_value "$ap_data" '.PosixUser.Gid')

            # Use default "/" for path if empty
            if [[ -z "$ap_path" || "$ap_path" == "N/A" ]]; then
                ap_path="/"
            fi

            # AccessPoint: Core info with empty fields for filesystem-specific and mount-target-specific data
            buffer+="efs,AccessPoint,,${region},$ap_name,$ap_id,AccessPoint,,,,,,,,$ap_path,$ap_uid,$ap_gid,$ap_state,\n"
        done < <(aws efs describe-access-points --file-system-id "$fs_id" --region "$region" 2>/dev/null | jq -c '.AccessPoints[]' || true)

    done < <(aws efs describe-file-systems --region "$region" 2>/dev/null | jq -c '.FileSystems[]' || true)

    echo "$buffer"
}

#######################################
# Function to collect EventBridge inventory (with categories)
#######################################
function collect_eventbridge_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ARN,State,Description,ScheduleExpression,RoleArn,Target"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    # Collect EventBridge Rules
    while IFS= read -r rule_data; do
        [[ -z "$rule_data" ]] && continue
        local rule_name rule_arn rule_state rule_description rule_schedule rule_role
        rule_name=$(extract_jq_value "$rule_data" '.Name')
        rule_arn=$(extract_jq_value "$rule_data" '.Arn')
        rule_state=$(extract_jq_value "$rule_data" '.State')
        rule_description=$(normalize_csv_value "$(extract_jq_value "$rule_data" '.Description')")
        rule_schedule=$(normalize_csv_value "$(extract_jq_value "$rule_data" '.ScheduleExpression')")

        # Get targets for this rule to determine role
        local targets_data rule_target_info
        targets_data=$(aws events list-targets-by-rule --rule "$rule_name" --region "$region" 2>/dev/null || echo '{"Targets":[]}')

        # Get first target's role and ARN (automatically get N/A if empty)
        rule_role=$(extract_jq_value "$targets_data" '.Targets[0].RoleArn')
        rule_target_info=$(extract_jq_value "$targets_data" '.Targets[0].Arn')

        # Apply normalize_csv_value to fields that may contain commas
        rule_description=$(normalize_csv_value "$rule_description")
        rule_target_info=$(normalize_csv_value "$rule_target_info")

        # Rule row: EventBridge Rule
        buffer+="eventbridge,Rule,,${region},$rule_name,$rule_arn,$rule_state,$rule_description,$rule_schedule,$rule_role,$rule_target_info\n"

    done < <(aws events list-rules --region "$region" 2>/dev/null | jq -c '.Rules[]' || true)

    # Collect EventBridge Schedules (from EventBridge Scheduler)
    while IFS= read -r schedule_data; do
        [[ -z "$schedule_data" ]] && continue
        local schedule_name schedule_arn schedule_state schedule_description schedule_expression schedule_role schedule_target
        schedule_name=$(extract_jq_value "$schedule_data" '.Name')
        schedule_arn=$(extract_jq_value "$schedule_data" '.Arn')
        schedule_state=$(extract_jq_value "$schedule_data" '.State')

        # Get detailed schedule information
        local schedule_details
        schedule_details=$(aws scheduler get-schedule --name "$schedule_name" --region "$region" 2>/dev/null || echo '{}')

        schedule_description=$(normalize_csv_value "$(extract_jq_value "$schedule_details" '.Description')")
        schedule_expression=$(extract_jq_value "$schedule_details" '.ScheduleExpression')

        # Get target role and target ARN
        schedule_role=$(extract_jq_value "$schedule_details" '.Target.RoleArn')
        schedule_target=$(extract_jq_value "$schedule_details" '.Target.Arn')

        # Apply normalize_csv_value to fields that may contain commas
        schedule_description=$(normalize_csv_value "$schedule_description")
        schedule_expression=$(normalize_csv_value "$schedule_expression")
        schedule_target=$(normalize_csv_value "$schedule_target")

        # Schedule row: EventBridge Scheduler
        buffer+="eventbridge,Scheduler,,${region},$schedule_name,$schedule_arn,$schedule_state,$schedule_description,$schedule_expression,$schedule_role,$schedule_target\n"

    done < <(aws scheduler list-schedules --region "$region" 2>/dev/null | jq -c '.Schedules[]' || true)

    echo "$buffer"
}

#######################################
# Function to collect Glue inventory (with categories)
#######################################
function collect_glue_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ID,Description,Role,Timeout,WorkerType,NumberOfWorkers,MaxRetries,GlueVersion,Language,ScriptLocation"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    # Collect Databases
    while IFS= read -r db_data; do
        [[ -z "$db_data" ]] && continue
        local db_name db_description
        db_name=$(extract_jq_value "$db_data" '.Name')
        db_description=$(normalize_csv_value "$(extract_jq_value "$db_data" '.Description')")
        # Database: Description,Location filled, job-specific fields empty
        buffer+="glue,Database,,${region},$db_name,$db_name,$db_description,,,,,,,\n"
    done < <(aws glue get-databases --region "$region" | jq -c '.DatabaseList[]')

    # Collect Jobs
    while IFS= read -r job_data; do
        [[ -z "$job_data" ]] && continue
        local job_name job_role job_timeout job_worker_type job_max_retries job_script_location job_num_workers job_glue_version job_language job_command_type
        job_name=$(extract_jq_value "$job_data" '.Name')
        job_role=$(extract_jq_value "$job_data" '.Role')
        job_timeout=$(extract_jq_value "$job_data" '.Timeout')
        job_worker_type=$(extract_jq_value "$job_data" '.WorkerType')
        job_max_retries=$(extract_jq_value "$job_data" '.MaxRetries')
        job_num_workers=$(extract_jq_value "$job_data" '.NumberOfWorkers')
        job_glue_version=$(extract_jq_value "$job_data" '.GlueVersion')
        job_command_type=$(extract_jq_value "$job_data" '.Command.Name')

        # Get Python version or determine language from job type
        local python_version
        python_version=$(extract_jq_value "$job_data" '.Command.PythonVersion')
        if [[ -n "$python_version" && "$python_version" != "N/A" ]]; then
            job_language="Python${python_version}"
        elif [[ "$job_command_type" == "glueetl" || "$job_command_type" == "pythonshell" ]]; then
            job_language="Python3" # Default for ETL and Python shell jobs
        else
            job_language="$job_command_type" # Use command type as fallback
        fi

        job_script_location=$(extract_jq_value "$job_data" '.Command.ScriptLocation')
        # Job: database-specific fields empty, job-specific fields filled
        buffer+="glue,Job,,${region},$job_name,$job_name,,$job_role,$job_timeout,$job_worker_type,$job_num_workers,$job_max_retries,$job_glue_version,$job_language,$job_script_location\n"
    done < <(aws glue get-jobs --region "$region" | jq -c '.Jobs[]')

    echo "$buffer"
}

#######################################
# Function to collect IAM inventory (with categories)
#######################################
function collect_iam_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ARN,Path,Created_Date"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    # IAM is a global service, only process from us-east-1
    if [[ "$region" != "us-east-1" ]]; then
        echo ""
        return 0
    fi

    local buffer=""

    while IFS= read -r role_data; do
        [[ -z "$role_data" ]] && continue
        local role_name role_arn role_path role_created
        role_name=$(extract_jq_value "$role_data" '.RoleName')
        role_arn=$(extract_jq_value "$role_data" '.Arn')
        role_path=$(extract_jq_value "$role_data" '.Path')
        role_created=$(extract_jq_value "$role_data" '.CreateDate')
        buffer+="iam,Role,,Global,$role_name,$role_arn,$role_path,$role_created\n"
    done < <(aws iam list-roles | jq -c '.Roles[]')

    while IFS= read -r user_data; do
        [[ -z "$user_data" ]] && continue
        local user_name user_arn user_path user_created
        user_name=$(extract_jq_value "$user_data" '.UserName')
        user_arn=$(extract_jq_value "$user_data" '.Arn')
        user_path=$(extract_jq_value "$user_data" '.Path')
        user_created=$(extract_jq_value "$user_data" '.CreateDate')
        buffer+="iam,User,,Global,$user_name,$user_arn,$user_path,$user_created\n"
    done < <(aws iam list-users | jq -c '.Users[]')

    while IFS= read -r policy_data; do
        [[ -z "$policy_data" ]] && continue
        local policy_name policy_arn policy_path policy_created
        policy_name=$(extract_jq_value "$policy_data" '.PolicyName')
        policy_arn=$(extract_jq_value "$policy_data" '.Arn')
        policy_path=$(extract_jq_value "$policy_data" '.Path')
        policy_created=$(extract_jq_value "$policy_data" '.CreateDate')
        buffer+="iam,Policy,,Global,$policy_name,$policy_arn,$policy_path,$policy_created\n"
    done < <(aws iam list-policies --scope Local | jq -c '.Policies[]')

    echo "$buffer"
}

#######################################
# Function to collect Kinesis inventory (with categories)
#######################################
function collect_kinesis_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ARN,Status,ShardsOrDest,Retention,Created_Date,Last_Modified_Date"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    # Kinesis Streams
    while IFS= read -r stream_name; do
        [[ -z "$stream_name" ]] && continue
        local stream_desc stream_arn stream_status stream_shards stream_retention stream_created
        stream_desc=$(aws kinesis describe-stream --stream-name "$stream_name" --region "$region" 2>/dev/null || echo "{}")
        stream_arn=$(extract_jq_value "$stream_desc" '.StreamDescription.StreamARN')
        stream_status=$(extract_jq_value "$stream_desc" '.StreamDescription.StreamStatus')
        stream_shards=$(extract_jq_value "$stream_desc" '.StreamDescription.Shards | length' '0')
        stream_retention=$(extract_jq_value "$stream_desc" '.StreamDescription.RetentionPeriodHours')
        stream_created=$(extract_jq_value "$stream_desc" '.StreamDescription.StreamCreationTimestamp')
        buffer+="kinesis,Stream,,${region},$stream_name,$stream_arn,$stream_status,$stream_shards,$stream_retention,$stream_created,N/A\n"
    done < <(aws kinesis list-streams --region "$region" | jq -r '.StreamNames[]?')

    # Kinesis Data Firehose
    while IFS= read -r ds_name; do
        [[ -z "$ds_name" ]] && continue
        local firehose_data
        firehose_data=$(aws firehose describe-delivery-stream --delivery-stream-name "$ds_name" --region "$region" 2>/dev/null)
        [[ -z "$firehose_data" ]] && continue
        local firehose_name firehose_arn firehose_status firehose_dest firehose_created firehose_last_update
        firehose_name=$(extract_jq_value "$firehose_data" '.DeliveryStreamDescription.DeliveryStreamName')
        firehose_arn=$(extract_jq_value "$firehose_data" '.DeliveryStreamDescription.DeliveryStreamARN')
        firehose_status=$(extract_jq_value "$firehose_data" '.DeliveryStreamDescription.DeliveryStreamStatus')
        firehose_dest=$(extract_jq_value "$firehose_data" '.DeliveryStreamDescription.Destinations[0].DestinationId')
        firehose_created=$(extract_jq_value "$firehose_data" '.DeliveryStreamDescription.CreateTimestamp')
        firehose_last_update=$(extract_jq_value "$firehose_data" '.DeliveryStreamDescription.LastUpdateTimestamp')
        buffer+="kinesis,Firehose,,${region},$firehose_name,$firehose_arn,$firehose_status,$firehose_dest,,$firehose_created,$firehose_last_update\n"
    done < <(aws firehose list-delivery-streams --region "$region" | jq -r '.DeliveryStreamNames[]?')

    echo "$buffer"
}

#######################################
# Function to collect KMS inventory (with categories)
#######################################
function collect_kms_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ARN,Description,Usage,State,Manager"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    while IFS= read -r key_data; do
        [[ -z "$key_data" ]] && continue
        local key_id key_arn key_description key_usage key_state key_name
        key_id=$(extract_jq_value "$key_data" '.KeyId')

        # Get detailed key information
        local key_details
        key_details=$(aws kms describe-key --key-id "$key_id" --region "$region" 2>/dev/null || echo '{"KeyMetadata":{}}')

        key_arn=$(extract_jq_value "$key_details" '.KeyMetadata.Arn')
        key_description=$(normalize_csv_value "$(extract_jq_value "$key_details" '.KeyMetadata.Description')")
        key_usage=$(extract_jq_value "$key_details" '.KeyMetadata.KeyUsage')
        key_state=$(extract_jq_value "$key_details" '.KeyMetadata.KeyState')

        # Skip AWS managed keys for cleaner output (optional)
        local key_manager
        key_manager=$(extract_jq_value "$key_details" '.KeyMetadata.KeyManager')

        # Check for aliases and use alias name if available
        key_aliases=$(extract_jq_value "$(aws kms list-aliases --key-id "$key_id" --region "$region" 2>/dev/null || echo '{}')" '.Aliases[0].AliasName')

        if [[ -n "$key_aliases" && "$key_aliases" != "N/A" ]]; then
            key_name="$key_aliases"
        else
            key_name="$key_id"
        fi

        buffer+="kms,Key,,${region},$key_name,$key_arn,$key_description,$key_usage,$key_state,$key_manager\n"
        # Removed count increment (not used)
    done < <(aws kms list-keys --region "$region" | jq -c '.Keys[]')

    echo "$buffer"
}

#######################################
# Function to collect Lambda inventory (with categories)
#######################################
function collect_lambda_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ARN,Role,Type,Runtime,Architecture,Memory,Timeout,EnvVars,Last_Modified_Date"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    while IFS= read -r func_data; do
        [[ -z "$func_data" ]] && continue
        local func_name func_arn func_runtime func_arch func_role func_memory func_timeout func_last_modified func_env_vars
        func_name=$(extract_jq_value "$func_data" '.FunctionName')
        func_arn=$(extract_jq_value "$func_data" '.FunctionArn')
        func_runtime=$(extract_jq_value "$func_data" '.Runtime')
        func_arch=$(extract_jq_value "$func_data" '.Architectures[0]')
        func_role=$(extract_jq_value "$func_data" '.Role')
        func_memory=$(extract_jq_value "$func_data" '.MemorySize')
        func_timeout=$(extract_jq_value "$func_data" '.Timeout')
        func_last_modified=$(extract_jq_value "$func_data" '.LastModified')

        # 環境変数を1セルにまとめて出力（改行区切り、PRIVATE_KEYはマスク）
        local func_env_vars_raw=""
        if echo "$func_data" | jq -e '.Environment.Variables' >/dev/null 2>&1; then
            while IFS= read -r env_pair; do
                [[ -z "$env_pair" ]] && continue
                local env_key env_value
                env_key=$(echo "$env_pair" | cut -d'=' -f1)
                env_value=$(echo "$env_pair" | cut -d'=' -f2-)
                if [[ "$env_key" =~ PRIVATE_KEY ]]; then
                    env_value="*****"
                fi
                if [[ -n "$func_env_vars_raw" ]]; then
                    func_env_vars_raw+=$'\n'
                fi
                func_env_vars_raw+="${env_key}=${env_value}"
            done < <(echo "$func_data" | jq -r '.Environment.Variables | to_entries[] | "\(.key)=\(.value)"' 2>/dev/null)
        fi
        func_env_vars=$(normalize_csv_value "$func_env_vars_raw")

        buffer+="lambda,Function,,${region},$func_name,$func_arn,$func_role,Function,$func_runtime,$func_arch,$func_memory,$func_timeout,$func_env_vars,$func_last_modified\n"
    done < <(aws lambda list-functions --region "$region" | jq -c '.Functions | sort_by(.FunctionName) | .[]')

    echo "$buffer"
}

#######################################
# Function to collect QuickSight inventory (with categories)
#######################################
function collect_quicksight_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ID,Type,Status,Created_Date"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    while IFS= read -r ds_data; do
        [[ -z "$ds_data" ]] && continue
        local ds_id ds_name ds_type ds_status
        ds_id=$(extract_jq_value "$ds_data" '.DataSourceId')
        ds_name=$(extract_jq_value "$ds_data" '.Name')
        ds_type=$(extract_jq_value "$ds_data" '.Type')
        ds_status=$(extract_jq_value "$ds_data" '.Status')
        buffer+="quicksight,DataSource,,${region},$ds_name,$ds_id,$ds_type,$ds_status,N/A\n"
    done < <(aws quicksight list-data-sources --aws-account-id "$(get_aws_account_id 2>/dev/null || echo 'unknown')" --region "$region" 2>/dev/null | jq -c '.DataSources[]' || true)

    while IFS= read -r analysis_data; do
        [[ -z "$analysis_data" ]] && continue
        local analysis_id analysis_name analysis_status analysis_created
        analysis_id=$(extract_jq_value "$analysis_data" '.AnalysisId')
        analysis_name=$(extract_jq_value "$analysis_data" '.Name')
        analysis_name=$(normalize_csv_value "$analysis_name")
        analysis_status=$(extract_jq_value "$analysis_data" '.Status')
        analysis_created=$(extract_jq_value "$analysis_data" '.CreatedTime')
        buffer+="quicksight,Analysis,,${region},$analysis_name,$analysis_id,$analysis_status,$analysis_created\n"
    done < <(aws quicksight list-analyses --aws-account-id "$(get_aws_account_id 2>/dev/null || echo 'unknown')" --region "$region" 2>/dev/null | jq -c '.AnalysisSummaryList[]' || true)

    echo "$buffer"
}

#######################################
# Function to collect RDS inventory (grouped by Cluster where applicable)
# - DBClusters are shown first with their member DBInstances as children
# - DBInstances show Writer/Reader role when part of a cluster
# - Standalone DBInstances (not in clusters) are listed separately
# - Added columns: EngineLifecycleSupport, IAM_DB_Auth, Kerberos_Auth, KMS_Key, AZ, BackupRetentionPeriod
#######################################
#######################################
# Function to collect RDS inventory (grouped by Cluster)
#######################################
function collect_rds_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ID,Type,Engine,Version,InstanceClass,Storage,MultiAZ,Members,EngineLifecycleSupport,IAM_DB_Auth,Kerberos_Auth,KMS_Key,AZ,BackupRetentionPeriod"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    # First, collect all clusters and their instances
    declare -A cluster_instances_map

    # Get all clusters first
    while IFS= read -r cluster_data; do
        [[ -z "$cluster_data" ]] && continue
        local cluster_id cluster_engine cluster_members cluster_version cluster_multi_az
        local cluster_extended_support cluster_iam_auth cluster_kerberos_auth cluster_kms_key cluster_az

        cluster_id=$(extract_jq_value "$cluster_data" '.DBClusterIdentifier')
        cluster_engine=$(extract_jq_value "$cluster_data" '.Engine')
        cluster_members=$(extract_jq_value "$cluster_data" '.DBClusterMembers | length')
        cluster_version=$(extract_jq_value "$cluster_data" '.EngineVersion')

        # MultiAZ for cluster
        cluster_multi_az=$(extract_jq_value "$cluster_data" '.MultiAZ')

        # Extended Support (available for certain engines) - use EngineLifecycleSupport
        cluster_extended_support=$(extract_jq_value "$cluster_data" '.EngineLifecycleSupport')

        # IAM Database Authentication
        cluster_iam_auth=$(extract_jq_value "$cluster_data" '.IAMDatabaseAuthenticationEnabled')

        # Kerberos Authentication
        kerberos_config=$(extract_jq_value "$cluster_data" '.DomainMemberships[]? | select(.Status == "joined") | .Domain')
        cluster_kerberos_auth=$(if [[ -n "$kerberos_config" && "$kerberos_config" != "N/A" ]]; then echo "true"; else echo "false"; fi)

        # KMS Key
        cluster_kms_key=$(extract_jq_value "$cluster_data" '.KmsKeyId')

        # Availability Zones
        cluster_az=$(extract_jq_value "$cluster_data" '.AvailabilityZones | join(";")')

        # Backup Retention Period
        cluster_backup_retention=$(extract_jq_value "$cluster_data" '.BackupRetentionPeriod')

        local current_cluster_output=""
        # DBCluster: Core info with empty fields for instance-specific data (InstanceClass, Storage)
        current_cluster_output+="rds,DBCluster,,${region},$cluster_id,$cluster_id,DBCluster,$cluster_engine,$cluster_version,,,$cluster_multi_az,$cluster_members,$cluster_extended_support,$cluster_iam_auth,$cluster_kerberos_auth,$cluster_kms_key,$cluster_az,$cluster_backup_retention\n"

        # Get cluster members and their roles
        while IFS= read -r member_data; do
            [[ -z "$member_data" ]] && continue
            local member_id member_role
            member_id=$(extract_jq_value "$member_data" '.DBInstanceIdentifier')
            member_role=$(extract_jq_value "$member_data" '.IsClusterWriter')
            if [[ "$member_role" == "true" ]]; then
                member_role="Writer"
            else
                member_role="Reader"
            fi

            # Store cluster member information
            cluster_instances_map["$member_id"]="$cluster_id|$member_role"
        done < <(echo "$cluster_data" | jq -c '.DBClusterMembers[]?' || echo "")

        # Get detailed instance information for cluster members
        while IFS= read -r instance_data; do
            [[ -z "$instance_data" ]] && continue
            local instance_id instance_cluster_id instance_engine instance_version instance_class instance_storage instance_multi_az
            local instance_extended_support instance_iam_auth instance_kerberos_auth instance_kms_key instance_az

            instance_id=$(extract_jq_value "$instance_data" '.DBInstanceIdentifier')
            instance_cluster_id=$(extract_jq_value "$instance_data" '.DBClusterIdentifier')

            # Only process instances that belong to this cluster
            if [[ "$instance_cluster_id" == "$cluster_id" ]]; then
                instance_engine=$(extract_jq_value "$instance_data" '.Engine')
                instance_version=$(extract_jq_value "$instance_data" '.EngineVersion')
                instance_class=$(extract_jq_value "$instance_data" '.DBInstanceClass')
                instance_storage=$(extract_jq_value "$instance_data" '.AllocatedStorage')
                instance_multi_az=$(extract_jq_value "$instance_data" '.MultiAZ')

                # Extended Support (inherit from cluster or check instance) - use EngineLifecycleSupport
                instance_extended_support=$(extract_jq_value "$instance_data" '.EngineLifecycleSupport')
                if [[ -z "$instance_extended_support" || "$instance_extended_support" == "null" ]]; then
                    instance_extended_support="$cluster_extended_support"
                fi

                # IAM Database Authentication (inherit from cluster)
                instance_iam_auth="$cluster_iam_auth"

                # Kerberos Authentication (inherit from cluster)
                instance_kerberos_auth="$cluster_kerberos_auth"

                # KMS Key (inherit from cluster)
                instance_kms_key="$cluster_kms_key"

                # Availability Zone
                instance_az=$(extract_jq_value "$instance_data" '.AvailabilityZone')

                # Backup Retention Period (inherit from cluster)
                instance_backup_retention="$cluster_backup_retention"

                # Get role from cluster members map
                local instance_role=""
                if [[ -n "${cluster_instances_map[$instance_id]:-}" ]]; then
                    instance_role="${cluster_instances_map[$instance_id]#*|}"
                fi

                # DBInstance (cluster member): Core info with empty Members field (cluster-level attribute)
                current_cluster_output+="rds,,DBInstance,${region},$instance_id,$instance_id,DBInstance ($instance_role),$instance_engine,$instance_version,$instance_class,$instance_storage,$instance_multi_az,,$instance_extended_support,$instance_iam_auth,$instance_kerberos_auth,$instance_kms_key,$instance_az,$instance_backup_retention\n"
            fi
        done < <(aws rds describe-db-instances --region "$region" | jq -c '.DBInstances[]')

        buffer+="$current_cluster_output"
    done < <(aws rds describe-db-clusters --region "$region" | jq -c '.DBClusters[]')

    # Now get standalone DB instances (not part of any cluster)
    while IFS= read -r db_data; do
        [[ -z "$db_data" ]] && continue
        local db_id db_cluster_id db_engine db_version db_class db_storage db_multi_az
        local db_extended_support db_iam_auth db_kerberos_auth db_kms_key db_az db_backup_retention

        db_id=$(extract_jq_value "$db_data" '.DBInstanceIdentifier')
        db_cluster_id=$(extract_jq_value "$db_data" '.DBClusterIdentifier')

        # Only process standalone instances (not part of a cluster)
        if [[ -z "$db_cluster_id" || "$db_cluster_id" == "null" ]]; then
            db_engine=$(extract_jq_value "$db_data" '.Engine')
            db_version=$(extract_jq_value "$db_data" '.EngineVersion')
            db_class=$(extract_jq_value "$db_data" '.DBInstanceClass')
            db_storage=$(extract_jq_value "$db_data" '.AllocatedStorage')
            db_multi_az=$(extract_jq_value "$db_data" '.MultiAZ')

            # Extended Support - use EngineLifecycleSupport
            db_extended_support=$(extract_jq_value "$db_data" '.EngineLifecycleSupport')

            # IAM Database Authentication
            db_iam_auth=$(extract_jq_value "$db_data" '.IAMDatabaseAuthenticationEnabled')

            # Kerberos Authentication
            kerberos_config=$(extract_jq_value "$db_data" '.DomainMemberships[]? | select(.Status == "joined") | .Domain')
            db_kerberos_auth=$(if [[ -n "$kerberos_config" && "$kerberos_config" != "N/A" ]]; then echo "true"; else echo "false"; fi)

            # KMS Key
            db_kms_key=$(extract_jq_value "$db_data" '.KmsKeyId')

            # Availability Zone
            db_az=$(extract_jq_value "$db_data" '.AvailabilityZone')

            # Backup Retention Period
            db_backup_retention=$(extract_jq_value "$db_data" '.BackupRetentionPeriod')

            # Standalone DBInstance: Core info with empty Members field (not applicable to standalone instances)
            buffer+="rds,DBInstance,,${region},$db_id,$db_id,DBInstance,$db_engine,$db_version,$db_class,$db_storage,$db_multi_az,,$db_extended_support,$db_iam_auth,$db_kerberos_auth,$db_kms_key,$db_az,$db_backup_retention\n"
        fi
    done < <(aws rds describe-db-instances --region "$region" | jq -c '.DBInstances[]')

    # Clear the associative array for the current region
    unset cluster_instances_map

    echo "$buffer"
}

#######################################
# Function to collect Redshift inventory (with categories)
#######################################
function collect_redshift_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ID,Status,Node_Type"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    while IFS= read -r cluster_data; do
        [[ -z "$cluster_data" ]] && continue
        local cluster_id cluster_status cluster_node_type
        cluster_id=$(extract_jq_value "$cluster_data" '.ClusterIdentifier')
        cluster_status=$(extract_jq_value "$cluster_data" '.ClusterStatus')
        cluster_node_type=$(extract_jq_value "$cluster_data" '.NodeType')
        buffer+="redshift,Cluster,,${region},$cluster_id,$cluster_id,$cluster_status,$cluster_node_type\n"
    done < <(aws redshift describe-clusters --region "$region" | jq -c '.Clusters[]')

    echo "$buffer"
}

#######################################
# Function to collect Route53 inventory (with categories)
# Note: TXT records may contain multiple lines (e.g., SPF, DKIM records)
# Use --preserve-newlines option for better display in Excel/Numbers
#######################################
#######################################
# Function to collect Route53 inventory (grouped by HostedZone)
#######################################
function collect_route53_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ID,Type,Comment,TTL,RecordType,Value,Record_Count"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    # Route53 is a global service, only process from us-east-1
    if [[ "$region" != "us-east-1" ]]; then
        echo ""
        return 0
    fi

    local buffer=""

    while IFS= read -r zone_data; do
        [[ -z "$zone_data" ]] && continue
        local zone_id zone_name zone_type zone_records zone_comment
        zone_id=$(extract_jq_value "$zone_data" '.Id')
        # Extract zone ID from the full path (e.g., /hostedzone/Z123456789 -> Z123456789)
        zone_id=$(basename "$zone_id")
        zone_name=$(extract_jq_value "$zone_data" '.Name')
        zone_type=$(extract_jq_value "$zone_data" '.Config.PrivateZone')
        if [[ "$zone_type" == "true" ]]; then
            zone_type="Private"
        else
            zone_type="Public"
        fi
        zone_records=$(extract_jq_value "$zone_data" '.ResourceRecordSetCount')
        zone_comment=$(extract_jq_value "$zone_data" '.Config.Comment')
        zone_comment=$(normalize_csv_value "$zone_comment")

        # Add current HostedZone data to output
        local current_zone_output=""
        # HostedZone: Type, Comment, Record_Count (empty fields for record-specific data)
        current_zone_output+="route53,HostedZone,,Global,$zone_name,$zone_id,$zone_type,$zone_comment,,,,$zone_records\n"

        # Get Resource Record Sets for this Hosted Zone
        while IFS= read -r record_data; do
            [[ -z "$record_data" ]] && continue
            local record_name record_type record_ttl record_values
            record_name=$(extract_jq_value "$record_data" '.Name')
            record_type=$(extract_jq_value "$record_data" '.Type')
            record_ttl=$(extract_jq_value "$record_data" '.TTL')

            # Handle alias records vs regular records
            local alias_target
            alias_target=$(extract_jq_value "$record_data" '.AliasTarget.DNSName')
            if [[ -n "$alias_target" && "$alias_target" != "null" ]]; then
                record_values=$(normalize_csv_value "$alias_target")
                record_ttl="Alias"
            else
                # Get regular record values
                local raw_values
                raw_values=$(extract_jq_value "$record_data" '.ResourceRecords[]?.Value')
                record_values=$(normalize_csv_value "$raw_values")
            fi

            # Skip default NS and SOA records for the zone itself to reduce noise
            # These are automatically created records and typically not of interest for inventory
            # Note: Custom A, CNAME, MX, etc. records including those with zone name will still be shown
            if [[ "$record_name" == "$zone_name" && ("$record_type" == "NS" || "$record_type" == "SOA") ]]; then
                continue
            fi

            # Record: TTL, RecordType, Value (empty fields for hostedzone-specific data)
            current_zone_output+="route53,,Record,Global,$record_name,,,,$record_ttl,$record_type,$record_values,\n"
        done < <(aws route53 list-resource-record-sets --hosted-zone-id "$zone_id" 2>/dev/null | jq -c '.ResourceRecordSets[]?' || echo "")

        # Add this zone's data to overall output
        buffer+="$current_zone_output"

    done < <(aws route53 list-hosted-zones | jq -c '.HostedZones[]')

    echo "$buffer"
}

#######################################
# Function to collect S3 inventory (with categories)
# Note: Lifecycle rule names may contain multiple entries separated by semicolons
#######################################
function collect_s3_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ARN,Encryption,Versioning,PAB_BlockPublicACLs,PAB_IgnorePublicACLs,PAB_BlockPublicPolicy,PAB_RestrictPublicBuckets,AccessLogARN,LifecycleRules,Created_Date"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    # S3 is a global service, only process from us-east-1
    if [[ "$region" != "us-east-1" ]]; then
        echo ""
        return 0
    fi

    local buffer=""

    while IFS= read -r bucket_data; do
        [[ -z "$bucket_data" ]] && continue
        local bucket_name bucket_created bucket_region bucket_encryption bucket_arn
        local bucket_versioning bucket_pab_block_public_acls bucket_pab_ignore_public_acls bucket_pab_block_public_policy bucket_pab_restrict_public_buckets bucket_lifecycle bucket_access_log_arn

        bucket_name=$(extract_jq_value "$bucket_data" '.Name')
        bucket_created=$(extract_jq_value "$bucket_data" '.CreationDate')
        bucket_region=$(extract_jq_value "$(aws s3api get-bucket-location --bucket "$bucket_name" 2>/dev/null || echo '{}')" '.LocationConstraint' 'us-east-1')

        # Construct bucket ARN
        bucket_arn="arn:aws:s3:::${bucket_name}"

        # Get encryption
        bucket_encryption=$(extract_jq_value "$(aws s3api get-bucket-encryption --bucket "$bucket_name" 2>/dev/null || echo '{}')" '.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm' 'None')

        # Get versioning
        bucket_versioning=$(extract_jq_value "$(aws s3api get-bucket-versioning --bucket "$bucket_name" 2>/dev/null || echo '{}')" '.Status' 'Disabled')

        # Get Public Access Block settings
        local pab_data
        pab_data=$(aws s3api get-public-access-block --bucket "$bucket_name" 2>/dev/null | jq '.PublicAccessBlockConfiguration' 2>/dev/null || echo "null")

        # Parse each Public Access Block setting with defaults
        bucket_pab_block_public_acls=$(extract_jq_value "$pab_data" '.BlockPublicAcls' 'false')
        bucket_pab_ignore_public_acls=$(extract_jq_value "$pab_data" '.IgnorePublicAcls' 'false')
        bucket_pab_block_public_policy=$(extract_jq_value "$pab_data" '.BlockPublicPolicy' 'false')
        bucket_pab_restrict_public_buckets=$(extract_jq_value "$pab_data" '.RestrictPublicBuckets' 'false')

        # Get access logging configuration
        local access_log_data bucket_access_log_bucket
        access_log_data=$(aws s3api get-bucket-logging --bucket "$bucket_name" 2>/dev/null | jq '.LoggingEnabled' 2>/dev/null || echo "null")
        if [[ "$access_log_data" != "null" && -n "$access_log_data" ]]; then
            bucket_access_log_bucket=$(extract_jq_value "$access_log_data" '.TargetBucket')
            if [[ -n "$bucket_access_log_bucket" && "$bucket_access_log_bucket" != "N/A" ]]; then
                bucket_access_log_arn="arn:aws:s3:::${bucket_access_log_bucket}"
            else
                bucket_access_log_arn=""
            fi
        else
            bucket_access_log_arn=""
        fi
        bucket_access_log_arn=$(normalize_csv_value "$bucket_access_log_arn")

        # Get lifecycle configuration with actual rule names
        local lifecycle_config lifecycle_rule_names
        lifecycle_config=$(aws s3api get-bucket-lifecycle-configuration --bucket "$bucket_name" 2>/dev/null | jq '.Rules' 2>/dev/null || echo "null")
        if [[ "$lifecycle_config" != "null" && "$lifecycle_config" != "[]" ]]; then
            # Extract rule names or IDs
            lifecycle_rule_names=$(extract_jq_value "$lifecycle_config" '.[] | .ID' 'Unnamed')
            bucket_lifecycle="$lifecycle_rule_names"
        else
            bucket_lifecycle="N/A"
        fi
        bucket_lifecycle=$(normalize_csv_value "$bucket_lifecycle")

        buffer+="s3,Bucket,,${bucket_region},$bucket_name,$bucket_arn,$bucket_encryption,$bucket_versioning,$bucket_pab_block_public_acls,$bucket_pab_ignore_public_acls,$bucket_pab_block_public_policy,$bucket_pab_restrict_public_buckets,$bucket_access_log_arn,$bucket_lifecycle,$bucket_created\n"
    done < <(aws s3api list-buckets | jq -c '.Buckets[]')

    echo "$buffer"
}

#######################################
# Function to collect Secrets Manager inventory (with categories)
#######################################
function collect_secretsmanager_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ARN,Description,Last_Modified_Date"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    while IFS= read -r secret_data; do
        [[ -z "$secret_data" ]] && continue
        local secret_name secret_arn secret_description secret_last_changed
        secret_name=$(extract_jq_value "$secret_data" '.Name')
        secret_arn=$(extract_jq_value "$secret_data" '.ARN')
        secret_description=$(normalize_csv_value "$(extract_jq_value "$secret_data" '.Description')")
        secret_last_changed=$(extract_jq_value "$secret_data" '.LastChangedDate')
        buffer+="secretsmanager,Secret,,${region},$secret_name,$secret_arn,$secret_description,$secret_last_changed\n"
        # Removed count increment (not used)
    done < <(aws secretsmanager list-secrets --region "$region" | jq -c '.SecretList[]')

    echo "$buffer"
}

#######################################
# Function to collect SNS inventory (with categories)
#######################################
function collect_sns_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ARN,Display_Name,Subscription_Count"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    while IFS= read -r topic_data; do
        [[ -z "$topic_data" ]] && continue
        local topic_arn topic_name topic_display_name topic_subscriptions
        topic_arn=$(extract_jq_value "$topic_data" '.TopicArn')
        # Extract topic name from ARN (last part after the last colon)
        topic_name="${topic_arn##*:}"
        topic_display_name=$(extract_jq_value "$(aws sns get-topic-attributes --topic-arn "$topic_arn" --region "$region" 2>/dev/null || echo '{}')" '.Attributes.DisplayName')
        topic_subscriptions=$(extract_jq_value "$(aws sns list-subscriptions-by-topic --topic-arn "$topic_arn" --region "$region" 2>/dev/null || echo '{}')" '.Subscriptions | length' '0')
        buffer+="sns,Topic,,${region},${topic_name},${topic_arn},${topic_display_name},${topic_subscriptions}\n"
        # Removed count increment (not used)
    done < <(aws sns list-topics --region "$region" | jq -c '.Topics[]')

    echo "$buffer"
}

#######################################
# Function to collect SQS inventory (with categories)
#######################################
function collect_sqs_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,URL,Type,VisibilityTimeout,DelaySeconds,MessageRetentionPeriod,MaxReceiveCount,DLQ_TargetARN,Created_Date,Last_Modified_Date"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    while IFS= read -r queue_url; do
        [[ -z "$queue_url" ]] && continue
        local queue_name queue_type visibility_timeout delay_seconds max_receive_count dlq_target_arn message_retention_period created_timestamp last_modified_timestamp
        queue_name=$(basename "$queue_url")

        # Get queue attributes
        local queue_attrs
        queue_attrs=$(aws sqs get-queue-attributes --queue-url "$queue_url" --attribute-names All --region "$region" 2>/dev/null || echo '{"Attributes":{}}')

        queue_type=$(extract_jq_value "$queue_attrs" '.Attributes.FifoQueue' 'false')
        if [[ "$queue_type" == "true" ]]; then
            queue_type="FIFO"
        else
            queue_type="Standard"
        fi

        visibility_timeout=$(extract_jq_value "$queue_attrs" '.Attributes.VisibilityTimeout')
        delay_seconds=$(extract_jq_value "$queue_attrs" '.Attributes.DelaySeconds')
        message_retention_period=$(extract_jq_value "$queue_attrs" '.Attributes.MessageRetentionPeriod')
        created_timestamp=$(extract_jq_value "$queue_attrs" '.Attributes.CreatedTimestamp')
        last_modified_timestamp=$(extract_jq_value "$queue_attrs" '.Attributes.LastModifiedTimestamp')

        # Dead Letter Queue related attributes
        local redrive_policy
        redrive_policy=$(extract_jq_value "$queue_attrs" '.Attributes.RedrivePolicy')
        if [[ -n "$redrive_policy" && "$redrive_policy" != "N/A" ]]; then
            max_receive_count=$(extract_jq_value "$redrive_policy" '.maxReceiveCount')
            dlq_target_arn=$(extract_jq_value "$redrive_policy" '.deadLetterTargetArn')
        else
            max_receive_count=""
            dlq_target_arn=""
        fi
        max_receive_count=$(normalize_csv_value "$max_receive_count")
        dlq_target_arn=$(normalize_csv_value "$dlq_target_arn")

        # Convert timestamps to readable format
        if [[ -n "$created_timestamp" && "$created_timestamp" =~ ^[0-9]+$ ]]; then
            created_timestamp=$(date -d "@$created_timestamp" "+%Y-%m-%d %H:%M:%S" 2>/dev/null || echo "$created_timestamp")
        fi
        if [[ -n "$last_modified_timestamp" && "$last_modified_timestamp" =~ ^[0-9]+$ ]]; then
            last_modified_timestamp=$(date -d "@$last_modified_timestamp" "+%Y-%m-%d %H:%M:%S" 2>/dev/null || echo "$last_modified_timestamp")
        fi
        created_timestamp=$(normalize_csv_value "$created_timestamp")
        last_modified_timestamp=$(normalize_csv_value "$last_modified_timestamp")

        buffer+="sqs,Queue,,${region},$queue_name,$queue_url,$queue_type,$visibility_timeout,$delay_seconds,$message_retention_period,$max_receive_count,$dlq_target_arn,$created_timestamp,$last_modified_timestamp\n"
    done < <(aws sqs list-queues --region "$region" | jq -r '.QueueUrls[]?')

    echo "$buffer"
}

#######################################
# Function to collect Transfer Family inventory (with categories)
#######################################
function collect_transferfamily_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,Server_ID,Protocol,State"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    while IFS= read -r server_data; do
        [[ -z "$server_data" ]] && continue
        local server_id server_state server_protocol
        server_id=$(extract_jq_value "$server_data" '.ServerId')
        server_state=$(extract_jq_value "$server_data" '.State')
        server_protocol=$(extract_jq_value "$server_data" '.Protocols[0]')
        buffer+="transferfamily,Server,,${region},$server_id,$server_id,$server_protocol,$server_state\n"
    done < <(aws transfer list-servers --region "$region" | jq -c '.Servers[]')

    echo "$buffer"
}

#######################################
# Function to collect VPC inventory (grouped by VPC)
#######################################
#######################################
# Function to collect VPC inventory (grouped by VPC)
#######################################
function collect_vpc_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ID,Description,CIDR,State,Settings"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    # Collect all VPC data in a single call with tags
    local vpc_json
    vpc_json=$(aws ec2 describe-vpcs --region "$region" 2>/dev/null) || return 0

    if [[ $(extract_jq_value "$vpc_json" '.Vpcs | length' '0') -eq 0 ]]; then
        echo ""
        return 0
    fi

    # Process each VPC and collect all its subresources
    while IFS= read -r vpc_id; do
        [[ -z "$vpc_id" || "$vpc_id" == "null" ]] && continue

        # Get VPC details
        local vpc_data
        vpc_data=$(echo "$vpc_json" | jq -r --arg vpc_id "$vpc_id" '.Vpcs[] | select(.VpcId == $vpc_id)')
        local cidr vpc_name vpc_state
        cidr=$(extract_jq_value "$vpc_data" '.CidrBlock')
        vpc_state=$(extract_jq_value "$vpc_data" '.State')

        # Get VPC name from tags
        vpc_name=$(normalize_csv_value "$(extract_jq_value "$vpc_data" '.Tags[]? | select(.Key == "Name") | .Value')")

        # Start with main VPC entry
        local current_vpc_output=""
        current_vpc_output+="vpc,VPC,,${region},${vpc_name},$vpc_id,,$cidr,$vpc_state\n"

        # Get route tables for this VPC to help classify subnets
        local rt_json
        rt_json=$(aws ec2 describe-route-tables --filters "Name=vpc-id,Values=$vpc_id" --region "$region" 2>/dev/null) || true

        # Collect and classify subnets for this VPC
        local subnet_json
        subnet_json=$(aws ec2 describe-subnets --filters "Name=vpc-id,Values=$vpc_id" --region "$region" 2>/dev/null) || true

        local public_subnets="" private_subnets=""
        while IFS= read -r subnet_data; do
            [[ -z "$subnet_data" ]] && continue
            local subnet_id subnet_cidr subnet_name
            subnet_id=$(extract_jq_value "$subnet_data" '.SubnetId')
            subnet_cidr=$(extract_jq_value "$subnet_data" '.CidrBlock')
            subnet_name=$(normalize_csv_value "$(extract_jq_value "$subnet_data" '.Tags[]? | select(.Key == "Name") | .Value')")

            # Check if subnet is public by looking at associated route table
            local is_public=false
            local subnet_rt_associations
            subnet_rt_associations=$(echo "$rt_json" | jq -r --arg subnet_id "$subnet_id" '.RouteTables[] | select(.Associations[]?.SubnetId == $subnet_id) | .RouteTableId')

            # If no explicit association, check main route table
            if [[ -z "$subnet_rt_associations" ]]; then
                subnet_rt_associations=$(echo "$rt_json" | jq -r '.RouteTables[] | select(.Associations[]?.Main == true) | .RouteTableId')
            fi

            # Check if any associated route table has route to internet gateway
            while IFS= read -r rt_id; do
                [[ -z "$rt_id" ]] && continue
                local has_igw_route
                has_igw_route=$(echo "$rt_json" | jq -r --arg rt_id "$rt_id" '.RouteTables[] | select(.RouteTableId == $rt_id) | .Routes[] | select(.GatewayId? // empty | startswith("igw-")) | .GatewayId')
                if [[ -n "$has_igw_route" ]]; then
                    is_public=true
                    break
                fi
            done <<<"$subnet_rt_associations"

            if [[ "$is_public" == "true" ]]; then
                public_subnets+="vpc,,PublicSubnet,${region},${subnet_name},$subnet_id,,$subnet_cidr,\n"
            else
                private_subnets+="vpc,,PrivateSubnet,${region},${subnet_name},$subnet_id,,$subnet_cidr,\n"
            fi
        done < <(echo "$subnet_json" | jq -c '.Subnets[]?')

        # Collect route tables for this VPC
        local route_tables=""
        while IFS= read -r rt_data; do
            [[ -z "$rt_data" ]] && continue
            local rt_id rt_name
            rt_id=$(extract_jq_value "$rt_data" '.RouteTableId')
            rt_name=$(normalize_csv_value "$(extract_jq_value "$rt_data" '.Tags[]? | select(.Key == "Name") | .Value')")
            route_tables+="vpc,,RouteTable,${region},${rt_name},$rt_id,,,\n"
        done < <(echo "$rt_json" | jq -c '.RouteTables[]?')

        # Collect internet gateways for this VPC
        local internet_gateways=""
        local igw_json
        igw_json=$(aws ec2 describe-internet-gateways --filters "Name=attachment.vpc-id,Values=$vpc_id" --region "$region" 2>/dev/null) || true

        while IFS= read -r igw_data; do
            [[ -z "$igw_data" ]] && continue
            local igw_id igw_name
            igw_id=$(extract_jq_value "$igw_data" '.InternetGatewayId')
            igw_name=$(normalize_csv_value "$(extract_jq_value "$igw_data" '.Tags[]? | select(.Key == "Name") | .Value')")
            internet_gateways+="vpc,,InternetGateway,${region},${igw_name},$igw_id,,,attached\n"
        done < <(echo "$igw_json" | jq -c '.InternetGateways[]?')

        # Collect NAT gateways for this VPC
        local nat_gateways=""
        local nat_json
        nat_json=$(aws ec2 describe-nat-gateways --filter "Name=vpc-id,Values=$vpc_id" --region "$region" 2>/dev/null) || true

        while IFS= read -r nat_data; do
            [[ -z "$nat_data" ]] && continue
            local nat_id nat_name nat_state
            nat_id=$(extract_jq_value "$nat_data" '.NatGatewayId')
            nat_state=$(extract_jq_value "$nat_data" '.State')
            nat_name=$(normalize_csv_value "$(extract_jq_value "$nat_data" '.Tags[]? | select(.Key == "Name") | .Value')")
            nat_gateways+="vpc,,NATGateway,${region},${nat_name},$nat_id,,,$nat_state\n"
        done < <(echo "$nat_json" | jq -c '.NatGateways[]?')

        # Collect Network ACLs for this VPC
        local network_acls=""
        local nacl_json
        nacl_json=$(aws ec2 describe-network-acls --filters "Name=vpc-id,Values=$vpc_id" --region "$region" 2>/dev/null) || true
        while IFS= read -r nacl_data; do
            [[ -z "$nacl_data" ]] && continue
            local nacl_id nacl_name nacl_desc nacl_entries nacl_associations
            nacl_id=$(extract_jq_value "$nacl_data" '.NetworkAclId')
            nacl_name=$(normalize_csv_value "$(extract_jq_value "$nacl_data" '.Tags[]? | select(.Key == "Name") | .Value')")
            nacl_desc=""
            # Summarize entries (rules) with line breaks
            nacl_entries=$(echo "$nacl_data" | jq -r '.Entries[]? | "Rule#: " + (.RuleNumber|tostring) + " | Protocol: " + (.Protocol // "-") + " | RuleAction: " + (.RuleAction // "-") + " | Egress: " + (.Egress|tostring) + " | CIDR: " + (.CidrBlock // "-") + " | PortRange: " + (if .PortRange? then (.PortRange.From|tostring) + "-" + (.PortRange.To|tostring) else "-" end)' | paste -sd "\n" -)
            nacl_entries=$(normalize_csv_value "$nacl_entries")
            # Summarize associations (subnets)
            nacl_associations=$(echo "$nacl_data" | jq -r '.Associations[]? | .SubnetId' | paste -sd ";" -)
            nacl_associations=$(normalize_csv_value "$nacl_associations")
            network_acls+="vpc,,NetworkACL,${region},${nacl_name},${nacl_id},${nacl_desc},,,${nacl_entries}\n"
        done < <(echo "$nacl_json" | jq -c '.NetworkAcls[]?')

        # Collect security groups for this VPC
        local security_groups=""
        local sg_json
        sg_json=$(aws ec2 describe-security-groups --filters "Name=vpc-id,Values=$vpc_id" --region "$region" 2>/dev/null) || true
        while IFS= read -r sg_data; do
            [[ -z "$sg_data" ]] && continue
            local sg_id sg_name sg_desc sg_inbound sg_outbound
            sg_id=$(extract_jq_value "$sg_data" '.GroupId')
            sg_name=$(normalize_csv_value "$(extract_jq_value "$sg_data" '.GroupName')")
            sg_desc=$(normalize_csv_value "$(extract_jq_value "$sg_data" '.Description')")
            # Inbound rules summary (改行区切り、説明付き)
            sg_inbound=$(echo "$sg_data" | jq -r '.IpPermissions[]? | "Protocol: " + (.IpProtocol // "-") + " | FromPort: " + ((.FromPort // "-") | tostring) + " | ToPort: " + ((.ToPort // "-") | tostring) + " | CIDR: " + (.IpRanges[]?.CidrIp // "-")' | paste -sd "\n" -)
            # Outbound rules summary (改行区切り、説明付き)
            sg_outbound=$(echo "$sg_data" | jq -r '.IpPermissionsEgress[]? | "Protocol: " + (.IpProtocol // "-") + " | FromPort: " + ((.FromPort // "-") | tostring) + " | ToPort: " + ((.ToPort // "-") | tostring) + " | CIDR: " + (.IpRanges[]?.CidrIp // "-")' | paste -sd "\n" -)
            sg_settings=$(normalize_csv_value "Inbound:\n$sg_inbound\nOutbound:\n$sg_outbound")
            # Output as Subsubcategory for VPC
            security_groups+="vpc,,SecurityGroup,${region},${sg_name},${sg_id},${sg_desc},,,${sg_settings}\n"
        done < <(echo "$sg_json" | jq -c '.SecurityGroups[]?')

        # Collect VPC endpoints for this VPC
        local endpoints=""
        local endpoint_json
        endpoint_json=$(aws ec2 describe-vpc-endpoints --filters "Name=vpc-id,Values=$vpc_id" --region "$region" 2>/dev/null) || true
        while IFS= read -r ep_data; do
            [[ -z "$ep_data" ]] && continue
            local ep_id ep_name ep_type ep_state ep_service ep_subnets ep_route_tables ep_security_groups ep_desc
            ep_id=$(extract_jq_value "$ep_data" '.VpcEndpointId')
            ep_type=$(extract_jq_value "$ep_data" '.VpcEndpointType')
            ep_state=$(extract_jq_value "$ep_data" '.State')
            ep_service=$(extract_jq_value "$ep_data" '.ServiceName')
            ep_subnets=$(extract_jq_value "$ep_data" '.SubnetIds | join(";")')
            ep_route_tables=$(extract_jq_value "$ep_data" '.RouteTableIds | join(";")')
            ep_security_groups=$(extract_jq_value "$ep_data" '.Groups[]?.GroupId' 'N/A')
            ep_security_groups=$(normalize_csv_value "$ep_security_groups")
            ep_name=$(normalize_csv_value "$(extract_jq_value "$ep_data" '.Tags[]? | select(.Key == "Name") | .Value')")
            ep_desc="Type: $ep_type\nService: $ep_service\nSubnets: $ep_subnets\nRouteTables: $ep_route_tables\nSecurityGroups: $ep_security_groups"
            ep_desc=$(normalize_csv_value "$ep_desc")
            endpoints+="vpc,,Endpoint,${region},${ep_name},${ep_id},,,${ep_state},${ep_desc}\n"
        done < <(echo "$endpoint_json" | jq -c '.VpcEndpoints[]?')

        # Group VPC subresources before adding to buffer
        current_vpc_output+="$public_subnets"
        current_vpc_output+="$private_subnets"
        current_vpc_output+="$route_tables"
        current_vpc_output+="$internet_gateways"
        current_vpc_output+="$nat_gateways"
        current_vpc_output+="$network_acls"
        current_vpc_output+="$security_groups"
        current_vpc_output+="$endpoints"

        # Add this VPC's data to overall output
        buffer+="$current_vpc_output"

    done < <(echo "$vpc_json" | jq -r '.Vpcs[]?.VpcId')

    echo "$buffer"
}

#######################################
# Function to collect WAF inventory (with categories)
#######################################
function collect_waf_inventory {
    local region=$1
    local header="Category,Subcategory,Subsubcategory,Region,Name,ID,Description,Scope"

    # Return header if requested
    if [[ "$region" == "header" ]]; then
        echo "$header"
        return 0
    fi

    local buffer=""

    # Regional WAF
    while IFS= read -r waf_data; do
        [[ -z "$waf_data" ]] && continue
        local waf_name waf_id waf_description waf_scope
        waf_name=$(extract_jq_value "$waf_data" '.Name')
        waf_id=$(extract_jq_value "$waf_data" '.Id')
        waf_description=$(normalize_csv_value "$(extract_jq_value "$waf_data" '.Description')")
        waf_scope="REGIONAL"
        buffer+="waf,WebACL,,${region},$waf_name,$waf_id,$waf_description,$waf_scope\n"
    done < <(aws wafv2 list-web-acls --scope REGIONAL --region "$region" | jq -c '.WebACLs[]')

    # CloudFront WAF (only from us-east-1)
    if [[ "$region" == "us-east-1" ]]; then
        while IFS= read -r waf_data; do
            [[ -z "$waf_data" ]] && continue
            local waf_name waf_id waf_description waf_scope
            waf_name=$(extract_jq_value "$waf_data" '.Name')
            waf_id=$(extract_jq_value "$waf_data" '.Id')
            waf_description=$(normalize_csv_value "$(extract_jq_value "$waf_data" '.Description')")
            waf_scope="CLOUDFRONT"
            buffer+="waf,WebACL,,Global,$waf_name,$waf_id,$waf_description,$waf_scope\n"
        done < <(aws wafv2 list-web-acls --scope CLOUDFRONT --region us-east-1 | jq -c '.WebACLs[]')
    fi

    echo "$buffer"
}

#######################################
# Common utility functions for AWS resource collection
#######################################

# Generic function to collect AWS resources across regions
function collect_aws_resources {
    local category=$1

    log "INFO" "Collecting $category information from AWS..."

    # Get header from the first call
    local collect_function="collect_${category}_inventory"
    if ! declare -f "$collect_function" >/dev/null; then
        log "WARN" "Collection function $collect_function not found for category $category"
        return 1
    fi

    # Get header
    local csv_header
    csv_header=$($collect_function "header")

    local buffer=""
    for region in "${REGIONS_TO_CHECK[@]}"; do
        log "INFO" "Checking $category resources in region: $region"
        buffer+=$($collect_function "$region")
    done
    # Check if this category should maintain grouping structure (no sorting)
    # Check if this category should maintain grouping structure (no sorting)
    local sort_output="true"
    for no_sort_category in "${NO_SORT_CATEGORIES[@]}"; do
        if [[ "$category" == "$no_sort_category" ]]; then
            sort_output="false"
        fi
    done

    output_csv_data "$category" "$csv_header" "$buffer" "$sort_output"
}

# Function to output CSV data with standard formatting
function output_csv_data {
    local category=$1
    local header=$2
    local buffer=$3
    local sort_output=${4:-"$SORT_OUTPUT"} # Use explicit parameter or global setting

    if [[ -n "$buffer" ]]; then
        {
            echo "$header"
            if [[ "$sort_output" == "true" ]]; then
                csv_sort "$buffer"
            else
                printf "%b" "$buffer"
            fi
            echo ""
        } >>"$OUTPUT_FILE"
    fi
    log "INFO" "$category inventory written to $OUTPUT_FILE"
}

# Function to iterate through regions and collect resources
function collect_regional_resources {
    local category=$1
    local collect_function=$2
    local buffer=""

    for region in "${REGIONS_TO_CHECK[@]}"; do
        log "INFO" "Checking $category resources in region: $region"
        buffer+=$($collect_function "$region")
    done

    echo "$buffer"
}

#######################################
# Main execution function
#######################################
function main {
    # Parse arguments
    parse_arguments "$@"

    # Record start time
    start_time=$(date +%s)

    # Validate dependencies
    validate_dependencies "aws" "jq"

    # Check AWS credentials before any AWS CLI usage
    if ! check_aws_credentials; then
        error_exit "AWS credentials are not set or invalid."
    fi

    # Initialize regions to check
    initialize_regions

    # Log script start
    echo_section "Starting AWS resource inventory collection"
    log "INFO" "Output file: $OUTPUT_FILE"
    log "INFO" "AWS region: $AWS_REGION"
    if [[ "$PRESERVE_NEWLINES" == "true" ]]; then
        log "INFO" "CSV newlines: Preserved (better for Excel/Numbers)"
    else
        log "INFO" "CSV newlines: Sanitized (maximum compatibility)"
    fi

    if [[ -n "$CATEGORIES" ]]; then
        log "INFO" "Collecting categories: $CATEGORIES"
    else
        log "INFO" "Collecting all categories"
    fi

    if [[ "$DRY_RUN" == "true" ]]; then
        log "INFO" "Running in dry-run mode, no changes will be made"
        exit 0
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
            for valid in "${AWS_RESOURCE_CATEGORIES[@]}"; do
                if [[ "$category" == "$valid" ]]; then
                    valid_category=true
                    break
                fi
            done
            if [[ "$valid_category" == "false" ]]; then
                error_exit "Invalid category: $category. Valid categories are: ${AWS_RESOURCE_CATEGORIES[*]}"
            fi
        done
    else
        # Use all categories if none specified
        categories_to_process=("${AWS_RESOURCE_CATEGORIES[@]}")
        log "INFO" "Processing all categories: ${categories_to_process[*]}"
    fi

    # Collect inventory for specified categories
    for resource_category in "${categories_to_process[@]}"; do
        collect_aws_resources "$resource_category"
    done

    # Record end time and calculate elapsed time
    end_time=$(date +%s)
    elapsed=$((end_time - start_time))
    log "INFO" "AWS resource inventory completed in ${elapsed} seconds"

    echo_section "AWS resource inventory collection completed successfully"
    echo "Results written to: $OUTPUT_FILE"
}

# Only call main function if script is executed directly, not sourced
if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi
