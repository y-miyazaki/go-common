#!/bin/bash
#######################################
# Description: Comprehensive script validation tool for all scripts in the workspace
# Usage: ./validate_all_scripts.sh [options]
#   options:
#     -h, --help     Display this help message
#     -v, --verbose  Enable verbose output
#     -f, --fix      Auto-fix issues where possible
#     -q, --quiet    Suppress non-error output
#######################################

#######################################
# Global variables and default values
#######################################

VERBOSE=false
AUTO_FIX=false
QUIET=false
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORKSPACE_ROOT="$(dirname "$SCRIPT_DIR")"
# Global variable for script search paths
SEARCH_PATHS=("$WORKSPACE_ROOT/scripts" "$WORKSPACE_ROOT/env")

# Load all-in-one library
# shellcheck source=../lib/all.sh
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/lib/all.sh"

# Counters for statistics
TOTAL_SCRIPTS=0
PASSED_SCRIPTS=0
FAILED_SCRIPTS=0
WARNINGS_COUNT=0

# Arrays for tracking results
declare -a PASSED_SCRIPTS_LIST=()
declare -a FAILED_SCRIPTS_LIST=()
declare -a WARNING_SCRIPTS_LIST=()

#######################################
# Override log function to handle warning counter and quiet mode
#######################################
function custom_log {
    local level=$1
    local message=$2

    # Handle quiet mode
    if [[ "$QUIET" == "true" ]] && [[ "$level" == "INFO" || "$level" == "DEBUG" ]]; then
        return 0
    fi

    # Call the library log function
    log "$level" "$message"

    # Increment warning counter for WARN level
    if [[ "$level" == "WARN" ]]; then
        ((WARNINGS_COUNT++))
    fi
}

#######################################
# Override echo_section for quiet mode support
#######################################
function custom_echo_section {
    if [[ "$QUIET" != "true" ]]; then
        echo_section "$1"
    fi
}
#######################################
# Display usage information
#######################################
function show_usage {
    show_help_header "$(basename "$0")" "Comprehensive validation tool for all shell scripts in the workspace" "[options]"
    echo "Options:"
    echo "  -h, --help     Display this help message"
    echo "  -v, --verbose  Enable verbose output and detailed analysis"
    echo "  -f, --fix      Auto-fix issues where possible (formatting, permissions)"
    echo "  -q, --quiet    Suppress non-error output (only show summary)"
    echo ""
    echo "Features:"
    echo "  - Recursive script discovery"
    echo "  - Syntax validation (bash -n)"
    echo "  - Shellcheck static analysis"
    echo "  - Permission validation"
    echo "  - Shebang validation"
    echo "  - Function dependency analysis"
    echo "  - Performance analysis"
    show_help_footer
    echo "Examples:"
    echo "  $(basename "$0")           # Run all validations"
    echo "  $(basename "$0") -v        # Verbose output with detailed analysis"
    echo "  $(basename "$0") -f        # Auto-fix fixable issues"
    echo "  $(basename "$0") -q        # Quiet mode, show only summary"
    exit 0
}

#######################################
# Function to find all shell scripts recursively
#######################################
function find_shell_scripts {
    local scripts=()

    for search_path in "${SEARCH_PATHS[@]}"; do
        if [[ -d "$search_path" ]]; then
            # Send debug log to stderr to avoid interfering with function output
            custom_log "DEBUG" "Searching for scripts in: $search_path" >&2

            # Find files with .sh extension
            while IFS= read -r -d '' script; do
                # Double-check it's a file (not directory)
                if [[ -f "$script" ]]; then
                    scripts+=("$script")
                fi
            done < <(find "$search_path" -type f -name "*.sh" -print0 2>/dev/null)

            # Find files with shell shebang but no .sh extension
            while IFS= read -r -d '' script; do
                if [[ -f "$script" ]] && head -1 "$script" 2>/dev/null | grep -q "^#!/.*sh"; then
                    # Only add if not already in scripts array
                    local already_added=false
                    for existing in "${scripts[@]}"; do
                        if [[ "$existing" == "$script" ]]; then
                            already_added=true
                            break
                        fi
                    done
                    if [[ "$already_added" == "false" ]]; then
                        scripts+=("$script")
                    fi
                fi
            done < <(find "$search_path" -type f ! -name "*.sh" -executable -print0 2>/dev/null)
        else
            # Send debug log to stderr
            custom_log "DEBUG" "Search path does not exist or is not a directory: $search_path" >&2
        fi
    done

    # Output only valid files, filtered again to ensure no directories slip through
    for script in "${scripts[@]}"; do
        if [[ -f "$script" && ! -d "$script" ]]; then
            echo "$script"
        fi
    done | sort -u
}

#######################################
# Function to validate script syntax
#######################################
function validate_syntax {
    local script="$1"
    local script_name
    script_name="$(basename "$script")"

    custom_log "DEBUG" "Validating syntax for: $script_name"

    if bash -n "$script" 2>/dev/null; then
        custom_log "DEBUG" "✅ Syntax validation passed: $script_name"
        return 0
    else
        local error_output
        error_output=$(bash -n "$script" 2>&1)
        custom_log "ERROR" "❌ Syntax validation failed: $script_name"
        if [[ "$VERBOSE" == "true" ]]; then
            echo "Syntax errors:"
            # shellcheck disable=SC2001
            echo "$error_output" | sed 's/^/  /'
        fi
        return 1
    fi
}

#######################################
# Function to auto-fix script formatting with shfmt
#######################################
function auto_fix_formatting {
    local script="$1"
    local script_name
    script_name="$(basename "$script")"

    if [[ "$AUTO_FIX" != "true" ]]; then
        return 0
    fi

    custom_log "DEBUG" "Auto-fixing formatting for: $script_name"

    if command -v shfmt &>/dev/null; then
        if shfmt -w -i 4 -ci "$script" 2>/dev/null; then
            custom_log "INFO" "✅ Formatted with shfmt: $script_name"
            return 0
        else
            custom_log "WARN" "⚠️  shfmt formatting failed: $script_name"
            return 1
        fi
    else
        custom_log "DEBUG" "shfmt not available, skipping formatting: $script_name"
        return 0
    fi
}

#######################################
# Function to auto-fix shellcheck issues
#######################################
function auto_fix_shellcheck {
    local script="$1"
    local script_name
    script_name="$(basename "$script")"

    if [[ "$AUTO_FIX" != "true" ]]; then
        return 0
    fi

    custom_log "DEBUG" "Auto-fixing shellcheck issues for: $script_name"

    # Get shellcheck suggestions in diff format for auto-fixable issues
    local shellcheck_fixes
    if shellcheck_fixes=$(shellcheck -e SC1091 -f diff "$script" 2>/dev/null); then
        if [[ -n "$shellcheck_fixes" ]]; then
            # Apply the diff patches
            if echo "$shellcheck_fixes" | patch -p1 --silent 2>/dev/null; then
                custom_log "INFO" "✅ Applied shellcheck fixes: $script_name"
                return 0
            else
                custom_log "WARN" "⚠️  Failed to apply shellcheck fixes: $script_name"
                return 1
            fi
        fi
    fi

    return 0
}

#######################################
# Function to run shellcheck
#######################################
function run_shellcheck {
    local script="$1"
    local script_name
    script_name="$(basename "$script")"

    custom_log "DEBUG" "Running shellcheck for: $script_name"

    local shellcheck_output
    if shellcheck_output=$(shellcheck -e SC1091 -f gcc "$script" 2>&1); then
        custom_log "DEBUG" "✅ Shellcheck passed: $script_name"
        return 0
    else
        custom_log "WARN" "⚠️  Shellcheck issues found: $script_name"
        if [[ "$VERBOSE" == "true" ]]; then
            echo "Shellcheck output:"
            # shellcheck disable=SC2001
            echo "$shellcheck_output" | sed 's/^/  /'
        fi
        return 1
    fi
}

#######################################
# Function to validate shebang
#######################################
function validate_shebang {
    local script="$1"
    local script_name
    script_name="$(basename "$script")"

    custom_log "DEBUG" "Validating shebang for: $script_name"

    local first_line
    first_line=$(head -1 "$script" 2>/dev/null)

    if [[ "$first_line" =~ ^#!/.*sh ]]; then
        custom_log "DEBUG" "✅ Valid shebang: $script_name ($first_line)"
        return 0
    elif [[ "$first_line" =~ ^#! ]]; then
        custom_log "WARN" "⚠️  Non-shell shebang: $script_name ($first_line)"
        return 1
    else
        custom_log "WARN" "⚠️  Missing or invalid shebang: $script_name"
        if [[ "$AUTO_FIX" == "true" ]]; then
            custom_log "INFO" "Auto-fixing shebang for: $script_name"
            # Create backup and add shebang
            cp "$script" "$script.bak"
            echo "#!/bin/bash" >"$script.tmp"
            cat "$script.bak" >>"$script.tmp"
            mv "$script.tmp" "$script"
            rm "$script.bak"
            custom_log "INFO" "✅ Added shebang to: $script_name"
        fi
        return 1
    fi
}

#######################################
# Function to validate file permissions
#######################################
function validate_permissions {
    local script="$1"
    local script_name
    script_name="$(basename "$script")"

    custom_log "DEBUG" "Validating permissions for: $script_name"

    if [[ -x "$script" ]]; then
        custom_log "DEBUG" "✅ Script is executable: $script_name"
        return 0
    else
        custom_log "WARN" "⚠️  Script is not executable: $script_name"
        if [[ "$AUTO_FIX" == "true" ]]; then
            custom_log "INFO" "Auto-fixing permissions for: $script_name"
            chmod +x "$script"
            custom_log "INFO" "✅ Made executable: $script_name"
        fi
        return 1
    fi
}

#######################################
# Function to analyze script functions
#######################################
function analyze_functions {
    local script="$1"
    local script_name
    script_name="$(basename "$script")"

    if [[ "$VERBOSE" != "true" ]]; then
        return 0
    fi

    custom_log "DEBUG" "Analyzing functions in: $script_name"

    local functions
    functions=$(grep -n "^function\|^[a-zA-Z_][a-zA-Z0-9_]*\s*()" "$script" 2>/dev/null | head -10)

    if [[ -n "$functions" ]]; then
        custom_log "INFO" "Functions found in $script_name:"
        echo "$functions" | while IFS= read -r func; do
            echo "  $func"
        done
    else
        custom_log "DEBUG" "No functions found in: $script_name"
    fi
}

#######################################
# Function to check script complexity
#######################################
function check_complexity {
    local script="$1"
    local script_name
    script_name="$(basename "$script")"

    if [[ "$VERBOSE" != "true" ]]; then
        return 0
    fi

    custom_log "DEBUG" "Checking complexity for: $script_name"

    local line_count
    line_count=$(wc -l <"$script" 2>/dev/null)

    local function_count
    function_count=$(grep -c "^function\|^[a-zA-Z_][a-zA-Z0-9_]*\s*()" "$script" 2>/dev/null)

    local complexity_score=0

    # Line count scoring
    if [[ $line_count -gt 500 ]]; then
        complexity_score=$((complexity_score + 3))
    elif [[ $line_count -gt 200 ]]; then
        complexity_score=$((complexity_score + 2))
    elif [[ $line_count -gt 100 ]]; then
        complexity_score=$((complexity_score + 1))
    fi

    # Function count scoring
    if [[ $function_count -gt 20 ]]; then
        complexity_score=$((complexity_score + 2))
    elif [[ $function_count -gt 10 ]]; then
        complexity_score=$((complexity_score + 1))
    fi

    custom_log "INFO" "Complexity analysis for $script_name:"
    custom_log "INFO" "  Lines: $line_count, Functions: $function_count, Score: $complexity_score"

    if [[ $complexity_score -ge 4 ]]; then
        custom_log "WARN" "  High complexity script detected"
    fi
}

#######################################
# Function to validate a single script
#######################################
function validate_script {
    local script="$1"
    local script_name
    script_name="$(basename "$script")"
    local relative_path
    # shellcheck disable=SC2295
    relative_path="${script#$WORKSPACE_ROOT/}"
    local validation_passed=true

    ((TOTAL_SCRIPTS++))

    custom_log "INFO" "Validating script: $relative_path"

    # Skip if it's a directory (should not happen after filtering, but safety check)
    if [[ -d "$script" ]]; then
        custom_log "DEBUG" "Skipping directory: $script_name"
        ((TOTAL_SCRIPTS--)) # Don't count directories
        return 0
    fi

    # Check if file exists and is readable
    if [[ ! -f "$script" ]] || [[ ! -r "$script" ]]; then
        custom_log "ERROR" "❌ Script not accessible: $script_name"
        FAILED_SCRIPTS_LIST+=("$relative_path (not accessible)")
        ((FAILED_SCRIPTS++))
        return 1
    fi

    # Apply auto-formatting with shfmt if in fix mode
    auto_fix_formatting "$script"

    # Validate shebang
    if ! validate_shebang "$script"; then
        validation_passed=false
    fi

    # Validate permissions
    if ! validate_permissions "$script"; then
        validation_passed=false
    fi

    # Validate syntax
    if ! validate_syntax "$script"; then
        validation_passed=false
    fi

    # Run shellcheck
    local shellcheck_passed=true
    if ! run_shellcheck "$script"; then
        shellcheck_passed=false
        validation_passed=false
    fi

    # Analyze functions (verbose mode only)
    analyze_functions "$script"

    # Check complexity (verbose mode only)
    check_complexity "$script"

    # Update counters and arrays
    if [[ "$validation_passed" == "true" ]]; then
        custom_log "INFO" "✅ All validations passed: $script_name"
        PASSED_SCRIPTS_LIST+=("$relative_path")
        ((PASSED_SCRIPTS++))
    else
        custom_log "ERROR" "❌ Validation failed: $script_name"
        FAILED_SCRIPTS_LIST+=("$relative_path")
        ((FAILED_SCRIPTS++))
    fi

    # Track warnings separately
    if [[ "$shellcheck_passed" == "false" ]] && [[ "$validation_passed" == "true" ]]; then
        WARNING_SCRIPTS_LIST+=("$relative_path (shellcheck warnings)")
    fi
}

#######################################
# Function to generate summary report
#######################################
function generate_summary_report {
    custom_echo_section "Validation Summary Report"

    echo "📊 Statistics:"
    echo "  Total scripts validated: $TOTAL_SCRIPTS"
    echo "  Scripts passed: $PASSED_SCRIPTS"
    echo "  Scripts failed: $FAILED_SCRIPTS"
    echo "  Warnings issued: $WARNINGS_COUNT"
    echo ""

    if [[ ${#PASSED_SCRIPTS_LIST[@]} -gt 0 ]]; then
        echo "✅ Scripts that passed all validations:"
        for script in "${PASSED_SCRIPTS_LIST[@]}"; do
            echo "  - $script"
        done
        echo ""
    fi

    if [[ ${#FAILED_SCRIPTS_LIST[@]} -gt 0 ]]; then
        echo "❌ Scripts that failed validation:"
        for script in "${FAILED_SCRIPTS_LIST[@]}"; do
            echo "  - $script"
        done
        echo ""
    fi

    if [[ ${#WARNING_SCRIPTS_LIST[@]} -gt 0 ]]; then
        echo "⚠️  Scripts with warnings:"
        for script in "${WARNING_SCRIPTS_LIST[@]}"; do
            echo "  - $script"
        done
        echo ""
    fi

    # Overall result
    if [[ $FAILED_SCRIPTS -eq 0 ]]; then
        if [[ $WARNINGS_COUNT -eq 0 ]]; then
            echo "🎉 All scripts passed validation with no warnings!"
        else
            echo "✅ All scripts passed validation (with $WARNINGS_COUNT warnings)"
        fi
    else
        echo "❌ $FAILED_SCRIPTS script(s) failed validation"
    fi
}

#######################################
# Function to generate recommendations
#######################################
function generate_recommendations {
    if [[ "$VERBOSE" != "true" ]]; then
        return 0
    fi

    custom_echo_section "Recommendations"

    echo "🔧 Improvement suggestions:"
    echo ""

    if [[ $FAILED_SCRIPTS -gt 0 ]]; then
        echo "🔨 To auto-fix common issues, run:"
        echo "  $(basename "$0") --fix"
        echo ""
    fi

    echo "📚 Best practices:"
    echo "  - Use 'set -euo pipefail' for strict error handling"
    echo "  - Quote variables to prevent word splitting"
    echo "  - Use local variables in functions"
    echo "  - Add comprehensive documentation headers"
    echo "  - Implement proper error handling and logging"
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
                shift
                ;;
            -f | --fix)
                AUTO_FIX=true
                shift
                ;;
            -q | --quiet)
                QUIET=true
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
# Main process
#######################################
function main {
    # Parse arguments
    parse_arguments "$@"

    # Validate dependencies
    validate_dependencies "bash" "find" "grep" "sed" "shellcheck" "shfmt"

    # Log script start
    custom_echo_section "Comprehensive Script Validation Tool"
    custom_log "INFO" "Workspace root: $WORKSPACE_ROOT"
    custom_log "INFO" "Auto-fix mode: $AUTO_FIX"
    custom_log "INFO" "Verbose mode: $VERBOSE"
    custom_log "INFO" "Quiet mode: $QUIET"

    # Find all shell scripts
    custom_echo_section "Script Discovery"
    local scripts
    readarray -t scripts < <(find_shell_scripts)

    if [[ ${#scripts[@]} -eq 0 ]]; then
        custom_log "WARN" "No shell scripts found in the workspace"
        exit 0
    fi

    custom_log "INFO" "Found ${#scripts[@]} shell scripts to validate"

    # Validate each script
    custom_echo_section "Script Validation"
    for script in "${scripts[@]}"; do
        validate_script "$script"
    done

    # Generate summary report
    generate_summary_report

    # Generate recommendations
    generate_recommendations

    custom_echo_section "Validation Complete"

    # Exit with appropriate code
    if [[ $FAILED_SCRIPTS -gt 0 ]]; then
        exit 1
    fi
}

# Only call main function if script is executed directly, not sourced
if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi
