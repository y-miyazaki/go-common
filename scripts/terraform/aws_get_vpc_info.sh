#!/bin/bash
#######################################
# Description: Get VPC-related resource information for a specified VPC
# Usage: ./aws_get_vpc_info.sh <VPC_ID>
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
VPC_ID=""
REGION="${AWS_DEFAULT_REGION:-ap-northeast-1}"

#######################################
# Display usage information
#######################################
function show_usage {
    echo "Usage: $(basename "$0") [options]"
    echo ""
    echo "Description: Get VPC-related resource information for a specified VPC"
    echo ""
    echo "Options:"
    echo "  -h, --help    Display this help message"
    echo ""
    echo "Arguments:"
    echo "  VPC_ID        VPC ID to query resources for"
    echo ""
    echo "Examples:"
    echo "  $(basename "$0") vpc-1234567890abcdef0"
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
                # Process VPC ID
                if [[ -z "${VPC_ID:-}" ]]; then
                    VPC_ID="$1"
                else
                    error_exit "Unexpected argument: $1"
                fi
                shift
                ;;
        esac
    done

    # Validate required arguments
    if [[ -z "${VPC_ID:-}" ]]; then
        echo "Error: VPC ID is required" >&2
        show_usage
    fi
}

#######################################
# Check VPC peering connections
#######################################
function check_vpc_peering {
    echo_section "VPC Peering Connections"
    if aws ec2 describe-vpc-peering-connections --region "$REGION" --filters 'Name=requester-vpc-info.vpc-id,Values='"$VPC_ID" | grep VpcPeeringConnectionId; then
        log "INFO" "VPC peering connections found"
    else
        log "INFO" "No VPC peering connections found"
    fi
}

#######################################
# Check NAT gateways
#######################################
function check_nat_gateways {
    echo_section "NAT Gateways"
    if aws ec2 describe-nat-gateways --region "$REGION" --filter 'Name=vpc-id,Values='"$VPC_ID" | grep NatGatewayId; then
        log "INFO" "NAT gateways found"
    else
        log "INFO" "No NAT gateways found"
    fi
}

#######################################
# Check EC2 instances
#######################################
function check_ec2_instances {
    echo_section "EC2 Instances"
    if aws ec2 describe-instances --region "$REGION" --filters 'Name=vpc-id,Values='"$VPC_ID" | grep InstanceId; then
        log "INFO" "EC2 instances found"
    else
        log "INFO" "No EC2 instances found"
    fi
}

#######################################
# Check VPN gateways
#######################################
function check_vpn_gateways {
    echo_section "VPN Gateways"
    if aws ec2 describe-vpn-gateways --region "$REGION" --filters 'Name=attachment.vpc-id,Values='"$VPC_ID" | grep VpnGatewayId; then
        log "INFO" "VPN gateways found"
    else
        log "INFO" "No VPN gateways found"
    fi
}

#######################################
# Check network interfaces
#######################################
function check_network_interfaces {
    echo_section "Network Interfaces"
    if aws ec2 describe-network-interfaces --region "$REGION" --filters 'Name=vpc-id,Values='"$VPC_ID" | grep NetworkInterfaceId; then
        log "INFO" "Network interfaces found"
    else
        log "INFO" "No network interfaces found"
    fi
}

#######################################
# Main execution function
#######################################
function main {
    # Parse arguments
    parse_arguments "$@"

    # Validate required dependencies
    validate_dependencies "aws"

    # Check AWS credentials before any AWS CLI usage
    if ! check_aws_credentials; then
        error_exit "AWS credentials are not set or invalid."
    fi

    # Log script start
    echo_section "VPC Resource Information"
    log "INFO" "Target VPC: $VPC_ID"
    log "INFO" "Target region: $REGION"

    # Check all VPC-related resources
    check_vpc_peering
    check_nat_gateways
    check_ec2_instances
    check_vpn_gateways
    check_network_interfaces

    echo_section "Process completed successfully"
    log "INFO" "VPC resource information retrieval completed"
}

# Only call main function if script is executed directly, not sourced
if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi
