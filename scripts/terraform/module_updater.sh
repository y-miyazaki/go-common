#!/bin/bash
#######################################
# Description: Complete Module-by-module Terraform updater with individual validation
# Usage: ./module_updater.sh [options] <terraform_directory>
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
# Global variables
#######################################
VERBOSE=false
DRY_RUN=false
CHECK_ONLY=false
RECURSIVE_SEARCH=false
TERRAFORM_DIR=""
DEFAULT_ENV="${ENV:-dev}"
BACKUP_DIR=""
NO_PLAN=false

# Counters
TOTAL_MODULES=0
UPDATED_MODULES=0
FAILED_MODULES=0

# Arrays for tracking
declare -a UPDATED_MODULES_LIST=()
declare -a FAILED_MODULES_LIST=()

# Global variables using common library defaults
DEFAULT_ENV="${ENV:-dev}"

# Arrays for tracking module updates
declare -A MODULE_FILES_MAP=()    # module_source -> list of files
declare -A MODULE_VERSIONS_MAP=() # module_source -> current_version|latest_version
declare -A PROJECT_DIRS_MAP=()    # project_dir -> 1 (for tracking unique project directories)

#######################################
# Display usage information
#######################################
function show_usage {
    echo "Usage: $(basename "$0") [options]"
    echo ""
    echo "Description: Update Terraform module versions and validate configurations"
    echo ""
    echo "Options:"
    echo "  -h, --help                 Display this help message"
    echo "  -v, --verbose              Enable verbose output (shows detailed plan content differences)"
    echo "  -d, --dry-run              Run in dry-run mode"
    echo "  -c, --check-only           Check only mode"
    echo "  -r, --recursive            Search recursively"
    echo "  --no-plan                  Skip plan/validation (no plan mode)"
    echo ""
    echo "Arguments:"
    echo "  terraform_directory        Target Terraform directory to process"
    echo ""
    echo "Environment Variables:"
    echo "  ENV              Environment for tfvars file selection (default: dev)"
    echo ""
    echo "Verbose Mode Features:"
    echo "  • Shows detailed terraform plan content differences (baseline vs current)"
    echo "  • Displays error details for failed operations"
    echo "  • Saves detailed diff logs to terraform_show_diff_YYYYMMDD_HHMMSS.log"
    echo "  • Provides summary of resource changes (added/removed/modified)"
    echo ""
    echo "Examples:"
    echo "  $(basename "$0") -v terraform/application"
    echo "  $(basename "$0") --dry-run --check-only terraform/base"
    echo "  $(basename "$0") --no-plan terraform/base"
    echo "  $(basename "$0") -v -r terraform/                  # Recursive with verbose output"
    exit 0
}

#######################################
# Parse command line arguments
#######################################
function parse_arguments {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h | --help) show_usage ;;
            -v | --verbose)
                VERBOSE=true
                shift
                ;;
            -d | --dry-run)
                DRY_RUN=true
                shift
                ;;
            -c | --check-only)
                CHECK_ONLY=true
                shift
                ;;
            -r | --recursive)
                RECURSIVE_SEARCH=true
                shift
                ;;
            --no-plan)
                NO_PLAN=true
                shift
                ;;
            -*) error_exit "Unknown option: $1" ;;
            *)
                if [[ -z "$TERRAFORM_DIR" ]]; then
                    TERRAFORM_DIR="$1"
                else
                    error_exit "Unexpected argument: $1"
                fi
                shift
                ;;
        esac
    done

    if [[ -z "$TERRAFORM_DIR" ]]; then
        echo "Error: Terraform directory is required" >&2
        show_usage
    fi

    TERRAFORM_DIR=$(realpath "$TERRAFORM_DIR")
    if [[ ! -d "$TERRAFORM_DIR" ]]; then
        error_exit "Terraform directory does not exist: $TERRAFORM_DIR"
    fi
}

function get_latest_version {
    local module_source=$1
    local base_module="$module_source"

    # Handle submodules
    if [[ "$module_source" =~ ^(.+)//modules/ ]]; then
        base_module="${BASH_REMATCH[1]}"
    fi

    local api_url="https://registry.terraform.io/v1/modules/${base_module}"

    if response=$(curl -s --connect-timeout 5 --max-time 10 "$api_url" 2>/dev/null); then
        if echo "$response" | jq -e '.version' >/dev/null 2>&1; then
            echo "$response" | jq -r '.version' 2>/dev/null || echo "unknown"
        else
            echo "unknown"
        fi
    else
        echo "unknown"
    fi
}

function setup_backup_directory {
    BACKUP_DIR="${TERRAFORM_DIR}/.terraform_module_backups/$(date +'%Y%m%d_%H%M%S')"
    mkdir -p "$BACKUP_DIR"
    log "INFO" "Backup directory created: $BACKUP_DIR"
}

function backup_file {
    local file=$1
    local backup_file
    backup_file="${BACKUP_DIR}/$(basename "$file").$(date +'%H%M%S_%N').bak"
    cp "$file" "$backup_file"
    echo "$backup_file"
}

function find_terraform_modules {
    local dir="$1"
    local search_recursive="${2:-false}"

    if [[ "$search_recursive" == "true" ]]; then
        find "$dir" -name "*.tf" -type f -not -path "*/.terraform/*" -exec grep -l "module.*{" {} \;
    else
        find "$dir" -maxdepth 1 -name "*.tf" -type f -exec grep -l "module.*{" {} \;
    fi
}

function extract_modules_from_file {
    local file="$1"
    # Extract module blocks and find source/version pairs robustly (POSIX awk, nested block aware)
    awk '
    BEGIN { in_module=0; block=""; nest=0; }
    /^[[:space:]]*module[[:space:]]+"[^"]+"[[:space:]]*{/ {
      in_module=1; block = $0 "\n"; nest=1; next;
    }
    in_module {
      block = block $0 "\n";
      if ($0 ~ /{/) nest++;
      if ($0 ~ /}/) nest--;
      if (nest == 0) {
        source=""; version="";
        if (match(block, /source[[:space:]]*=[[:space:]]*\"[^\"]+\"/)) {
          s = substr(block, RSTART, RLENGTH);
          gsub(/.*=[[:space:]]*\"/, "", s);
          gsub(/\".*/, "", s);
          source = s;
        }
        if (match(block, /version[[:space:]]*=[[:space:]]*\"[^\"]+\"/)) {
          s = substr(block, RSTART, RLENGTH);
          gsub(/.*=[[:space:]]*\"/, "", s);
          gsub(/\".*/, "", s);
          version = s;
        }
        if (length(source) > 0 && length(version) > 0) { print source "||" version "||" FILENAME; }
        in_module=0; block="";
      }
      next;
    }
  ' "$file"
}

function update_module_version {
    local file="$1"
    local module_source="$2"
    local current_version="$3"
    local new_version="$4"

    log "INFO" "Updating $module_source: $current_version -> $new_version in $(basename "$file")"

    # Create backup
    local backup_file
    backup_file=$(backup_file "$file")

    # Update version using sed
    local escaped_current
    # shellcheck disable=SC2016
    escaped_current=$(printf '%s\n' "$current_version" | sed 's/[[\.*^$()+?{|]/\\&/g')
    local escaped_new
    # shellcheck disable=SC2016
    escaped_new=$(printf '%s\n' "$new_version" | sed 's/[[\.*^$()+?{|]/\\&/g')

    if sed -i.tmp "s/version[[:space:]]*=[[:space:]]*\"${escaped_current}\"/version = \"${escaped_new}\"/" "$file"; then
        rm -f "$file.tmp"
        log "INFO" "✅ Successfully updated $module_source in $(basename "$file")"
        return 0
    else
        # Restore from backup on failure
        cp "$backup_file" "$file"
        rm -f "$file.tmp"
        log "ERROR" "❌ Failed to update $module_source in $(basename "$file")"
        return 1
    fi
}

function process_single_module_update {
    local file="$1"
    local module_source="$2"
    local current_version="$3"
    local latest_version="$4"

    if [[ "$DRY_RUN" == "true" ]]; then
        #    echo "📦 Would update: $module_source: $current_version -> $latest_version in $(basename "$file")"
        ((UPDATED_MODULES++))
        return 0
    fi

    # For batch processing, just collect module information
    collect_module_for_batch_update "$file" "$module_source" "$current_version" "$latest_version"
}

#######################################
# Collect module information for batch update
#######################################
function collect_module_for_batch_update {
    local file="$1"
    local module_source="$2"
    local current_version="$3"
    local latest_version="$4"

    # Store file information for this module
    if [[ -z "${MODULE_FILES_MAP[$module_source]}" ]]; then
        MODULE_FILES_MAP[$module_source]="$file"
    else
        MODULE_FILES_MAP[$module_source]="${MODULE_FILES_MAP[$module_source]}|$file"
    fi

    # Store version information
    MODULE_VERSIONS_MAP[$module_source]="$current_version|$latest_version"

    # Track project directories
    local terraform_project_dir
    terraform_project_dir=$(find_terraform_project_root "$file")
    PROJECT_DIRS_MAP[$terraform_project_dir]=1

    log "INFO" "Collected: $module_source ($current_version -> $latest_version) in $(basename "$file")"
}

#######################################
# Process all collected modules in batch
#######################################
function process_batch_module_updates {
    local total_modules=${#MODULE_FILES_MAP[@]}
    local current_module=0

    if [[ $total_modules -eq 0 ]]; then
        log "INFO" "No modules to update"
        return 0
    fi

    echo_section "Batch Module Updates ($total_modules modules)"

    if [[ "$NO_PLAN" != "true" ]]; then
        # Create baseline plans for all affected projects
        create_baseline_plans_for_projects
    else
        log "INFO" "--no-plan specified: Skipping baseline plan creation and validation."
    fi

    # Process each module
    for module_source in "${!MODULE_FILES_MAP[@]}"; do
        ((current_module++))
        echo "🔄 [$current_module/$total_modules] Processing module: $module_source"

        local version_info="${MODULE_VERSIONS_MAP[$module_source]}"
        local current_version="${version_info%|*}"
        local latest_version="${version_info#*|}"
        local files_list="${MODULE_FILES_MAP[$module_source]}"

        # Split files and update all at once
        IFS='|' read -ra files_array <<<"$files_list"

        local update_success=true
        local updated_files=()

        # Update all files for this module
        for file in "${files_array[@]}"; do
            if update_module_version "$file" "$module_source" "$current_version" "$latest_version"; then
                updated_files+=("$file")
                log "INFO" "✅ Updated $module_source in $(basename "$file")"
            else
                log "ERROR" "❌ Failed to update $module_source in $(basename "$file")"
                update_success=false
                break
            fi
        done

        if [[ "$update_success" == "true" ]]; then
            if [[ "$NO_PLAN" != "true" ]]; then
                # Validate all affected projects
                if validate_all_affected_projects "${files_array[@]}"; then
                    ((UPDATED_MODULES++))
                    UPDATED_MODULES_LIST+=("$module_source ($current_version -> $latest_version) [${#files_array[@]} files]")
                    log "INFO" "✅ Successfully updated and validated $module_source in ${#files_array[@]} files"
                else
                    log "WARN" "❌ Validation failed for $module_source, rolling back all files..."
                    rollback_module_files "$module_source" "${updated_files[@]}"
                    ((FAILED_MODULES++))
                    FAILED_MODULES_LIST+=("$module_source")
                fi
            else
                # No plan/validation mode: treat as success
                ((UPDATED_MODULES++))
                UPDATED_MODULES_LIST+=("$module_source ($current_version -> $latest_version) [${#files_array[@]} files] (no plan mode)")
                log "INFO" "✅ Updated $module_source in ${#files_array[@]} files (no plan mode)"
            fi
        else
            # Rollback any successful updates for this module
            if [[ ${#updated_files[@]} -gt 0 ]]; then
                log "WARN" "❌ Partial update failure for $module_source, rolling back successful updates..."
                rollback_module_files "$module_source" "${updated_files[@]}"
            fi
            ((FAILED_MODULES++))
            FAILED_MODULES_LIST+=("$module_source")
        fi
    done
}

#######################################
# Create baseline plans for all affected projects
#######################################
function create_baseline_plans_for_projects {
    echo_section "Creating baseline plans for affected projects"

    for project_dir in "${!PROJECT_DIRS_MAP[@]}"; do
        log "INFO" "Creating baseline plan for: $project_dir"
        if ! create_baseline_plan "$project_dir" "$DEFAULT_ENV"; then
            log "ERROR" "Failed to create baseline plan for: $project_dir"
            error_exit "Cannot proceed without baseline plans"
        fi
    done
}

#######################################
# Validate all affected projects
#######################################
function validate_all_affected_projects {
    local files=("$@")
    local project_dirs_to_validate=()

    # Get unique project directories from the files
    for file in "${files[@]}"; do
        local project_dir
        project_dir=$(find_terraform_project_root "$file")

        # Check if project_dir already exists in array
        local dir_exists=false
        for existing_dir in "${project_dirs_to_validate[@]}"; do
            if [[ "$existing_dir" == "$project_dir" ]]; then
                dir_exists=true
                break
            fi
        done

        if [[ "$dir_exists" == "false" ]]; then
            project_dirs_to_validate+=("$project_dir")
        fi
    done

    # Validate each project
    for project_dir in "${project_dirs_to_validate[@]}"; do
        log "INFO" "Validating project: $project_dir"
        if ! validate_terraform_with_plan_comparison "$project_dir" "$DEFAULT_ENV"; then
            log "ERROR" "Validation failed for project: $project_dir"
            return 1
        fi
    done

    log "INFO" "✅ All affected projects validated successfully"
    return 0
}

# shellcheck disable=SC2317
function rollback_module_files {
    local module_source="$1"
    shift
    local files=("$@")

    for file in "${files[@]}"; do
        local backup_files
        mapfile -t backup_files < <(find "${BACKUP_DIR}" -name "$(basename "$file").*.bak" -type f | sort)
        if [[ ${#backup_files[@]} -gt 0 ]] && [[ -f "${backup_files[-1]}" ]]; then
            cp "${backup_files[-1]}" "$file"
            log "INFO" "Rolled back $module_source in $(basename "$file")"
        else
            log "ERROR" "❌ No backup file found for rollback: $(basename "$file")"
        fi
    done
}

# shellcheck disable=SC2317
function process_terraform_directory {
    local terraform_dir="$1"

    echo_section "Processing: $(basename "$terraform_dir")"
    cd "$terraform_dir" || {
        error_exit "Failed to change to directory: $terraform_dir"
        return 1
    }

    # Setup backup directory for actual updates
    if [[ "$DRY_RUN" == "false" ]] && [[ "$CHECK_ONLY" == "false" ]]; then
        setup_backup_directory
    fi

    # Find terraform files with modules
    local terraform_files
    readarray -t terraform_files < <(find_terraform_modules "$terraform_dir" "$RECURSIVE_SEARCH")

    if [[ ${#terraform_files[@]} -eq 0 ]]; then
        log "INFO" "No Terraform files with modules found"
        return 0
    fi

    log "INFO" "Found ${#terraform_files[@]} files with modules"

    # Process each file
    for file in "${terraform_files[@]}"; do
        local file_basename
        file_basename=$(basename "$file")
        log "INFO" "Scanning file: $file_basename"

        # Extract modules from file
        local modules_info
        modules_info=$(extract_modules_from_file "$file")
        if [[ "$VERBOSE" == "true" ]]; then
            echo "[DEBUG] modules_info for $file:" >&2
            echo "$modules_info" >&2
        fi
        while IFS= read -r module_line; do
            [[ -z "$module_line" ]] && continue
            ((TOTAL_MODULES++))
            if [[ "$VERBOSE" == "true" ]]; then
                echo "[DEBUG] TOTAL_MODULES incremented: $TOTAL_MODULES ($module_line)" >&2
            fi

            local module_source="${module_line%%||*}"
            local remaining="${module_line#*||}"
            local current_version="${remaining%%||*}"

            log "INFO" "Found module: $module_source (current: $current_version)"

            if [[ "$CHECK_ONLY" == "true" ]]; then
                echo "📦 Found: $module_source (current: $current_version) in $file_basename"
                continue
            fi

            # Get latest version
            local latest_version
            latest_version=$(get_latest_version "$module_source")

            if [[ "$latest_version" == "unknown" ]]; then
                log "WARN" "Could not determine latest version for $module_source"
                continue
            fi

            if [[ "$current_version" == "$latest_version" ]]; then
                log "INFO" "Module $module_source is already up to date ($current_version)"
                continue
            fi

            # Process individual module update
            echo "📦 Update available: $module_source: $current_version -> $latest_version in $file_basename"
            process_single_module_update "$file" "$module_source" "$current_version" "$latest_version"

        done <<<"$modules_info"
    done

    # After scanning all files, process batch updates
    if [[ "$DRY_RUN" == "false" ]] && [[ "$CHECK_ONLY" == "false" ]]; then
        process_batch_module_updates
    fi

    # Note: Backup files will be cleaned up by cleanup_all_artifacts function at the end
    # This preserves them during execution for rollback purposes
}

function generate_summary_report {
    echo_section "Update Summary Report"

    echo "📊 Statistics:"
    echo "  Total modules scanned: $TOTAL_MODULES"
    echo "  Modules updated: $UPDATED_MODULES"
    echo "  Modules failed: $FAILED_MODULES"
    echo ""

    if [[ ${#UPDATED_MODULES_LIST[@]} -gt 0 ]]; then
        echo "✅ Successfully updated modules:"
        for module in "${UPDATED_MODULES_LIST[@]}"; do
            echo "  - $module"
        done
        echo ""
    fi

    if [[ ${#FAILED_MODULES_LIST[@]} -gt 0 ]]; then
        echo "❌ Failed to update modules:"
        for module in "${FAILED_MODULES_LIST[@]}"; do
            echo "  - $module"
        done
        echo ""
    fi
}

#######################################
# Create baseline terraform plan before module updates
#######################################
function create_baseline_plan {
    local terraform_dir="$1"
    local env="${2:-$DEFAULT_ENV}"

    cd "$terraform_dir" || return 1

    log "INFO" "Creating baseline plan before module update..."

    # Debug output for troubleshooting
    echo "[DEBUG] Current directory: $(pwd)"
    echo "[DEBUG] Directory listing:" && ls -l
    echo "[DEBUG] ENV: $env"
    set -x

    # Check for backend configuration file
    local backend_config="terraform.${env}.tfbackend"
    local init_options=()

    if [[ -f "$backend_config" ]]; then
        init_options+=("-reconfigure" "-backend-config=$backend_config")
    fi

    # Initialize terraform
    if ! terraform init "${init_options[@]}" >/dev/null 2>&1; then
        set +x
        log "ERROR" "Initial terraform init failed"
        if [[ "$VERBOSE" == "true" ]]; then
            log "ERROR" "Init error - running in verbose mode"
            # shellcheck disable=SC2086
            terraform init "${init_options[@]}"
        fi
        return 1
    fi
    set +x

    # Create baseline plan
    local tfvars_file="terraform.${env}.tfvars"
    local plan_options=""
    if [[ -f "$tfvars_file" ]]; then
        plan_options="-var-file=$tfvars_file"
    fi

    local baseline_plan_command="terraform plan -lock=false -out=.terraform_baseline.plan"
    if [[ -n "$plan_options" ]]; then
        # shellcheck disable=SC2086
        baseline_plan_command="$baseline_plan_command $plan_options"
    fi

    if ! $baseline_plan_command >.terraform_baseline.log 2>&1; then
        log "ERROR" "Failed to create baseline plan"
        if [[ "$VERBOSE" == "true" ]]; then
            echo_section "BASELINE PLAN ERROR OUTPUT"
            log "ERROR" "Command failed: $baseline_plan_command"
            log "ERROR" "Error details:"
            cat .terraform_baseline.log
            echo ""
        fi
        return 1
    fi

    log "INFO" "✅ Baseline plan created successfully"
    # Clean up init and validate logs on success
    rm -f .terraform_init.log .terraform_validate.log
    return 0
}

#######################################
# Validate terraform with plan content comparison
#######################################
function validate_terraform_with_plan_comparison {
    local terraform_dir="$1"
    local env="${2:-$DEFAULT_ENV}"

    cd "$terraform_dir" || return 1

    log "INFO" "Validating Terraform configuration with plan content comparison..."

    # Check for backend configuration file
    local backend_config="terraform.${env}.tfbackend"
    local init_options=("-input=false" "-upgrade")

    if [[ -f "$backend_config" ]]; then
        init_options+=("-reconfigure" "-backend-config=$backend_config")
    fi

    # Re-initialize to ensure new module version is downloaded
    log "INFO" "Re-initializing Terraform to download updated modules..."
    if ! terraform init "${init_options[@]}" >.terraform_init.log 2>&1; then
        log "ERROR" "Terraform init failed after module update"
        if [[ "$VERBOSE" == "true" ]]; then
            echo_section "TERRAFORM INIT ERROR OUTPUT"
            log "ERROR" "Init command failed with options: ${init_options[*]}"
            log "ERROR" "Error details:"
            cat .terraform_init.log
            echo ""
        fi
        return 1
    fi

    # Validate syntax after init
    if ! terraform validate >.terraform_validate.log 2>&1; then
        log "ERROR" "Terraform syntax validation failed"
        if [[ "$VERBOSE" == "true" ]]; then
            echo_section "TERRAFORM VALIDATE ERROR OUTPUT"
            log "ERROR" "Validation failed, error details:"
            cat .terraform_validate.log
            echo ""
        fi
        return 1
    fi

    # Create current plan and compare with baseline
    local tfvars_file="terraform.${env}.tfvars"
    local plan_options=""
    if [[ -f "$tfvars_file" ]]; then
        plan_options="-var-file=$tfvars_file"
    fi

    local current_plan_command="terraform plan -lock=false -out=.terraform_current.plan"
    if [[ -n "$plan_options" ]]; then
        current_plan_command="$current_plan_command $plan_options"
    fi

    if ! $current_plan_command >.terraform_current.log 2>&1; then
        log "ERROR" "Failed to create current plan"
        if [[ "$VERBOSE" == "true" ]]; then
            echo_section "CURRENT PLAN ERROR OUTPUT"
            log "ERROR" "Command failed: $current_plan_command"
            log "ERROR" "Error details:"
            cat .terraform_current.log
            echo ""
        fi
        return 1
    fi

    # Compare plan file contents using 'terraform show'
    terraform show -no-color .terraform_baseline.plan >.terraform_baseline.txt 2>/dev/null
    terraform show -no-color .terraform_current.plan >.terraform_current.txt 2>/dev/null

    if diff -q .terraform_baseline.txt .terraform_current.txt >/dev/null 2>&1; then
        log "INFO" "✅ No infrastructure changes detected"
        # Clean up plan files
        rm -f .terraform_baseline.plan .terraform_baseline.log .terraform_baseline.txt
        rm -f .terraform_current.plan .terraform_current.log .terraform_current.txt
        rm -f .terraform_init.log .terraform_validate.log
        return 0
    else
        log "WARN" "⚠️  Infrastructure changes detected between baseline and current plan"
        if [[ "$VERBOSE" == "true" ]]; then
            echo_section "TERRAFORM SHOW DIFFERENCES DETECTED"
            log "INFO" "Detailed plan content differences (baseline vs current):"
            echo ""

            # Show unified diff with context
            if command -v colordiff >/dev/null 2>&1; then
                diff -u .terraform_baseline.txt .terraform_current.txt | colordiff
            else
                diff -u .terraform_baseline.txt .terraform_current.txt
            fi

            echo ""
            echo_section "SUMMARY OF CHANGES"

            # Count and summarize changes
            local added_resources
            local removed_resources
            local modified_resources

            added_resources=$(diff .terraform_baseline.txt .terraform_current.txt | grep -c "^+.*resource\|^+.*data\." || true)
            removed_resources=$(diff .terraform_baseline.txt .terraform_current.txt | grep -c "^-.*resource\|^-.*data\." || true)
            modified_resources=$(diff .terraform_baseline.txt .terraform_current.txt | grep -c "^[+-].*~\|^[+-].*-/+\|^[+-].*force replacement" || true)

            log "INFO" "Resources to be added: $added_resources"
            log "INFO" "Resources to be removed: $removed_resources"
            log "INFO" "Resources to be modified: $modified_resources"

            # Save detailed diff to log file for later review
            local diff_log_file
            diff_log_file="terraform_show_diff_$(date +%Y%m%d_%H%M%S).log"

            {
                echo "# Terraform Show Differences - $(date)"
                echo "# Directory: $(pwd)"
                echo "# Baseline vs Current Plan Content Comparison"
                echo ""
                diff -u .terraform_baseline.txt .terraform_current.txt
            } >"$diff_log_file"

            log "INFO" "Detailed diff saved to: $diff_log_file"
            echo ""
        else
            log "INFO" "Use -v option to see detailed plan content differences"
        fi
        # Clean up plan files
        rm -f .terraform_baseline.plan .terraform_baseline.log .terraform_baseline.txt
        rm -f .terraform_current.plan .terraform_current.log .terraform_current.txt
        rm -f .terraform_init.log .terraform_validate.log
        return 1
    fi
}

# shellcheck disable=SC2317
function cleanup_plan_files {
    local terraform_dir="$1"

    cd "$terraform_dir" || return 1

    # Remove plan comparison files
    rm -f .terraform_baseline.plan .terraform_baseline.log .terraform_baseline.txt
    rm -f .terraform_current.plan .terraform_current.log .terraform_current.txt
    rm -f .terraform_init.log .terraform_validate.log

    # Remove terraform show diff log files (both old and new naming)
    find . -maxdepth 1 -name "terraform_show_diff_*.log" -type f -delete 2>/dev/null || true
    find . -maxdepth 1 -name "terraform_plan_diff_*.log" -type f -delete 2>/dev/null || true

    log "INFO" "Plan comparison files cleaned up"
}

# shellcheck disable=SC2317
function cleanup_all_artifacts {
    log "INFO" "Cleaning up all artifacts..."

    # Cleanup backup directories
    if [[ -n "$BACKUP_DIR" ]] && [[ -d "$BACKUP_DIR" ]]; then
        log "INFO" "Removing backup directory: $BACKUP_DIR"
        rm -rf "$BACKUP_DIR"
    fi

    # Find and cleanup all backup directories in terraform directory
    local backup_dirs
    readarray -t backup_dirs < <(find "$TERRAFORM_DIR" -type d -name ".terraform_module_backups" 2>/dev/null)

    for backup_dir in "${backup_dirs[@]}"; do
        if [[ -d "$backup_dir" ]]; then
            log "INFO" "Removing backup directory: $backup_dir"
            rm -rf "$backup_dir"
        fi
    done

    if [[ "$NO_PLAN" == "true" ]]; then
        log "INFO" "--no-plan specified: Skipping plan file cleanup."
        return
    fi

    # Cleanup plan comparison files in all terraform directories
    local terraform_dirs
    if [[ "$RECURSIVE_SEARCH" == "true" ]]; then
        readarray -t terraform_dirs < <(find "$TERRAFORM_DIR" -name "*.tf" -type f -exec dirname {} \; 2>/dev/null | sort -u)
    else
        terraform_dirs=("$TERRAFORM_DIR")
    fi

    for dir in "${terraform_dirs[@]}"; do
        if [[ -d "$dir" ]]; then
            cd "$dir" || continue
            local plan_files=(
                .terraform_baseline.plan .terraform_baseline.log .terraform_baseline.txt
                .terraform_current.plan .terraform_current.log .terraform_current.txt
                .terraform_init.log .terraform_validate.log
            )

            for file in "${plan_files[@]}"; do
                if [[ -f "$file" ]]; then
                    rm -f "$file"
                fi
            done

            # Remove terraform show diff log files (both old and new naming)
            find . -maxdepth 1 -name "terraform_show_diff_*.log" -type f -delete 2>/dev/null || true
            find . -maxdepth 1 -name "terraform_plan_diff_*.log" -type f -delete 2>/dev/null || true
        fi
    done

    log "INFO" "✅ All artifacts cleaned up successfully"
}

# shellcheck disable=SC2317
function find_terraform_project_root {
    local file_path="$1"
    local current_dir
    current_dir=$(dirname "$file_path")

    # Look for terraform configuration files that indicate a project root
    while [[ "$current_dir" != "/" ]]; do
        # Check for backend config files or variables files
        if [[ -f "$current_dir/terraform.dev.tfbackend" ]] || [[ -f "$current_dir/terraform.dev.tfvars" ]] || [[ -f "$current_dir/main.tf" ]]; then
            echo "$current_dir"
            return 0
        fi

        # Also check for terraform directory structure patterns
        if [[ -f "$current_dir/versions.tf" ]] || [[ -f "$current_dir/providers.tf" ]]; then
            echo "$current_dir"
            return 0
        fi

        current_dir=$(dirname "$current_dir")
    done

    # If no project root found, return the file's directory
    # shellcheck disable=SC2086
    # shellcheck disable=SC2005
    echo "$(dirname \"$file_path\")"
}

#######################################
# Main execution function
#######################################
function main {
    parse_arguments "$@"

    # Validate dependencies
    validate_dependencies "terraform" "curl" "jq"

    echo_section "Complete Module-by-Module Updater"
    log "INFO" "Target directory: $TERRAFORM_DIR"
    log "INFO" "Recursive search: $RECURSIVE_SEARCH"
    log "INFO" "Mode: $(
        if [[ "$CHECK_ONLY" == "true" ]]; then
            echo "CHECK-ONLY"
        elif [[ "$DRY_RUN" == "true" ]]; then
            echo "DRY-RUN"
        else
            echo "UPDATE"
        fi
    )"

    local start_time
    start_time=$(date +%s)

    process_terraform_directory "$TERRAFORM_DIR"

    generate_summary_report

    # Cleanup backup directories and plan files
    cleanup_all_artifacts

    local end_time
    end_time=$(date +%s)
    local elapsed=$((end_time - start_time))
    log "INFO" "Process completed in ${elapsed} seconds"

    echo_section "Process completed"

    # Exit with appropriate code
    if [[ $FAILED_MODULES -gt 0 ]]; then
        exit 1
    fi
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi
