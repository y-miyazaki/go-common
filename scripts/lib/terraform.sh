#!/bin/bash
#######################################
# Description: Terraform-specific utility functions for shell scripts
# Usage: source /path/to/scripts/lib/terraform.sh
#
# This library provides Terraform-related functions:
# - Terraform environment validation
# - Standard Terraform workflow operations
# - Backend configuration management
# - Plugin cache management
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
# Function to validate Terraform environment variables
# Arguments:
#   $1 - Environment name (optional, uses ENV if not provided)
# Globals:
#   ENV - Environment name
#   TF_PLUGIN_CACHE_DIR - Terraform plugin cache directory
# Returns:
#   0 on success, exits on error
#######################################
function validate_terraform_env {
    local env_name="${1:-${ENV:-}}"

    # Validate required environment variables
    validate_env_vars "ENV" "TF_PLUGIN_CACHE_DIR"

    # Validate dependencies
    validate_dependencies "terraform" "tfenv"

    # Create plugin cache directory
    if [[ -n "${TF_PLUGIN_CACHE_DIR}" ]]; then
        execute_command "mkdir -p '${TF_PLUGIN_CACHE_DIR}'"
        log "INFO" "Terraform plugin cache directory: ${TF_PLUGIN_CACHE_DIR}"
    fi

    log "INFO" "Terraform environment validated for: $env_name"
}

#######################################
# Function to install Terraform version using tfenv
# Arguments:
#   None (uses .terraform-version file if present)
# Returns:
#   0 on success, exits on error
#######################################
function terraform_install {
    echo_section "Installing Terraform version"

    if ! execute_command "tfenv install"; then
        error_exit "Failed to install Terraform with tfenv"
    fi

    # Show installed version
    local tf_version
    tf_version=$(terraform version -json 2> /dev/null | jq -r '.terraform_version' 2> /dev/null || terraform version | head -n1)
    log "INFO" "Terraform version: $tf_version"
}

#######################################
# Function to initialize Terraform with backend configuration
# Arguments:
#   $1 - Environment name (optional, uses ENV if not provided)
#   $2 - Additional init options (optional)
# Globals:
#   ENV - Environment name
# Returns:
#   0 on success, exits on error
#######################################
function terraform_init {
    local env_name="${1:-${ENV}}"
    local additional_options="${2:-}"
    local backend_config="terraform.${env_name}.tfbackend"

    echo_section "Terraform initialization"

    # Check if backend config file exists
    if [[ ! -f "$backend_config" ]]; then
        error_exit "Backend configuration file not found: $backend_config"
    fi

    local init_cmd="terraform init -reconfigure -backend-config=$backend_config"
    if [[ -n "$additional_options" ]]; then
        init_cmd="$init_cmd $additional_options"
    fi

    if ! execute_command "$init_cmd"; then
        error_exit "Failed to initialize Terraform"
    fi

    log "INFO" "Terraform initialized with backend config: $backend_config"
}

#######################################
# Function to validate Terraform configuration
# Arguments:
#   None
# Returns:
#   0 on success, exits on error
#######################################
function terraform_validate {
    echo_section "Terraform validation"

    if ! execute_command "terraform validate"; then
        error_exit "Terraform configuration validation failed"
    fi

    log "INFO" "Terraform configuration is valid"
}

#######################################
# Function to format Terraform configuration files
# Arguments:
#   $1 - Check only mode (optional, "check" to only check formatting)
# Returns:
#   0 on success, 1 if formatting issues found (in check mode)
#######################################
function terraform_format {
    local check_mode="${1:-}"

    if [[ "$check_mode" == "check" ]]; then
        echo_section "Terraform format check"
        if execute_command "terraform fmt -check -diff"; then
            log "INFO" "Terraform files are properly formatted"
            return 0
        else
            log "ERROR" "Terraform files need formatting"
            return 1
        fi
    else
        echo_section "Terraform format"
        execute_command "terraform fmt -recursive"
        log "INFO" "Terraform files formatted"
    fi
}

#######################################
# Function to create Terraform plan
# Arguments:
#   $1 - Environment name (optional, uses ENV if not provided)
#   $2 - Plan file path (optional, defaults to terraform.tfplan)
#   $3 - Additional plan options (optional)
# Globals:
#   ENV - Environment name
# Returns:
#   0 on success, exits on error
#######################################
function terraform_plan {
    local env_name="${1:-${ENV}}"
    local plan_file="${2:-terraform.tfplan}"
    local additional_options="${3:-}"
    local var_file="terraform.${env_name}.tfvars"

    echo_section "Terraform plan"

    # Check if variables file exists
    if [[ ! -f "$var_file" ]]; then
        error_exit "Variables file not found: $var_file"
    fi

    local plan_cmd="terraform plan -out=$plan_file -var-file=$var_file"
    if [[ -n "$additional_options" ]]; then
        plan_cmd="$plan_cmd $additional_options"
    fi

    if ! execute_command "$plan_cmd"; then
        error_exit "Failed to create Terraform plan"
    fi

    log "INFO" "Terraform plan created: $plan_file"
}

#######################################
# Function to apply Terraform configuration
# Arguments:
#   $1 - Environment name (optional, uses ENV if not provided)
#   $2 - Plan file path (optional, if provided, applies from plan file)
#   $3 - Auto-approve flag (optional, "auto-approve" to skip confirmation)
# Globals:
#   ENV - Environment name
# Returns:
#   0 on success, exits on error
#######################################
function terraform_apply {
    local env_name="${1:-${ENV}}"
    local plan_file="${2:-}"
    local auto_approve="${3:-}"

    echo_section "Terraform apply"

    local apply_cmd
    if [[ -n "$plan_file" && -f "$plan_file" ]]; then
        # Apply from plan file
        apply_cmd="terraform apply $plan_file"
    else
        # Direct apply with variables
        local var_file="terraform.${env_name}.tfvars"
        if [[ ! -f "$var_file" ]]; then
            error_exit "Variables file not found: $var_file"
        fi
        apply_cmd="terraform apply -var-file=$var_file"
    fi

    # Add auto-approve if specified
    if [[ "$auto_approve" == "auto-approve" ]]; then
        apply_cmd="$apply_cmd --auto-approve"
    fi

    if ! execute_command "$apply_cmd"; then
        error_exit "Failed to apply Terraform configuration"
    fi

    log "INFO" "Terraform configuration applied successfully"
}

#######################################
# Function to destroy Terraform resources
# Arguments:
#   $1 - Environment name (optional, uses ENV if not provided)
#   $2 - Auto-approve flag (optional, "auto-approve" to skip confirmation)
# Globals:
#   ENV - Environment name
# Returns:
#   0 on success, exits on error
#######################################
function terraform_destroy {
    local env_name="${1:-${ENV}}"
    local auto_approve="${2:-}"
    local var_file="terraform.${env_name}.tfvars"

    echo_section "Terraform destroy"

    # Check if variables file exists
    if [[ ! -f "$var_file" ]]; then
        error_exit "Variables file not found: $var_file"
    fi

    local destroy_cmd="terraform destroy -var-file=$var_file"

    # Add auto-approve if specified
    if [[ "$auto_approve" == "auto-approve" ]]; then
        destroy_cmd="$destroy_cmd --auto-approve"
    fi

    if ! execute_command "$destroy_cmd"; then
        error_exit "Failed to destroy Terraform resources"
    fi

    log "INFO" "Terraform resources destroyed successfully"
}

#######################################
# Function to run complete Terraform workflow
# Arguments:
#   $1 - Environment name (optional, uses ENV if not provided)
#   $2 - Workflow type (optional: "plan", "apply", "destroy", defaults to "apply")
#   $3 - Auto-approve flag (optional, "auto-approve" to skip confirmation)
# Globals:
#   ENV - Environment name
# Returns:
#   0 on success, exits on error
#######################################
function terraform_workflow {
    local env_name="${1:-${ENV}}"
    local workflow_type="${2:-apply}"
    local auto_approve="${3:-}"

    # Validate environment
    validate_terraform_env "$env_name"

    # Install Terraform
    terraform_install

    # Initialize Terraform
    terraform_init "$env_name"

    # Validate configuration
    terraform_validate

    # Execute workflow based on type
    case "$workflow_type" in
        "plan")
            terraform_plan "$env_name"
            ;;
        "apply")
            terraform_apply "$env_name" "" "$auto_approve"
            ;;
        "destroy")
            terraform_destroy "$env_name" "$auto_approve"
            ;;
        *)
            error_exit "Unknown workflow type: $workflow_type. Use 'plan', 'apply', or 'destroy'."
            ;;
    esac

    log "INFO" "Terraform workflow completed: $workflow_type"
}

#######################################
# Function to check Terraform workspace
# Outputs:
#   Current workspace name
# Returns:
#   0 on success
#######################################
function terraform_get_workspace {
    terraform workspace show 2> /dev/null || echo "default"
}

#######################################
# Function to switch Terraform workspace
# Arguments:
#   $1 - Workspace name
# Returns:
#   0 on success, exits on error
#######################################
function terraform_select_workspace {
    local workspace="$1"

    if [[ -z "$workspace" ]]; then
        error_exit "Workspace name is required"
    fi

    # Check if workspace exists, create if it doesn't
    if ! terraform workspace select "$workspace" 2> /dev/null; then
        log "INFO" "Creating new workspace: $workspace"
        execute_command "terraform workspace new '$workspace'"
    fi

    log "INFO" "Switched to workspace: $workspace"
}
