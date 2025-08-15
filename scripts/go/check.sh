#!/bin/bash
#######################################
# Description: Comprehensive Go code quality and testing script for Lambda functions
# Usage: ./check.sh [options] <directory>
#   options:
#     -h, --help     Display this help message
#     -v, --verbose  Enable verbose output
#     -d, --dry-run  Run in dry-run mode (no changes made)
#     -f, --fix      Automatically fix issues where possible
#######################################

# Error handling: exit on error, unset variable, or failed pipeline
set -euo pipefail

# Get script directory for library loading
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export SCRIPT_DIR

# Load common libraries - ALWAYS use this pattern
# shellcheck source=../lib/all.sh
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/../lib/all.sh"

#######################################
# Global variables and default values
#######################################
VERBOSE=false
DRY_RUN=false
FIX_MODE=false
TARGET_DIR=""
TARGET_PATTERN="./..."
COVERAGE_THRESHOLD=80
EXIT_CODE=0
# 新規: スコープ実行フラグ (特定ディレクトリのみ)
IS_SCOPED=false

# Flags for individual checks
GO_FMT_FAILED=0
GO_VET_FAILED=0
LINT_FAILED=0
TEST_FAILED=0
RACE_FAILED=0
COVERAGE_FAILED=0
SECURITY_FAILED=0

# Counters for issues
LINT_ISSUES_COUNT=0
TEST_FAIL_COUNT=0
COVERAGE_PERCENT=""

# Functions are now provided by common.sh library

#######################################
# Display usage information
#######################################
function show_usage {
    echo "Usage: $(basename "$0") [options]"
    echo ""
    echo "Description: Comprehensive Go code quality and testing script for Lambda functions"
    echo "Automatically detects Go directories if no specific target is provided"
    echo ""
    echo "Options:"
    echo "  -h, --help     Display this help message"
    echo "  -v, --verbose  Enable verbose output and benchmark tests"
    echo "  -d, --dry-run  Run in dry-run mode (no changes made)"
    echo "  -f, --fix      Automatically fix issues where possible"
    echo ""
    echo "Arguments:"
    echo "  directory      Target directory to check (optional)"
    echo "                 If not provided, auto-detects all Go directories"
    echo "                 Use './...' to check all packages recursively"
    echo ""
    echo "Examples:"
    echo "  $(basename "$0")                                    # Auto-detect all Go directories"
    echo "  $(basename "$0") -v                                 # Auto-detect with verbose output"
    echo "  $(basename "$0") -f ./cmd/cloudwatch                # Check specific directory with auto-fix"
    echo "  $(basename "$0") -v -f ./...                        # Check all packages recursively"
    exit 0
}

# Dependencies validation is now handled by common.sh library

#######################################
# Function to check if directory contains Go files
#######################################
function has_go_files {
    local dir=$1
    find "$dir" -name "*.go" | head -1 | wc -l
}

#######################################
# Function to determine target pattern for Go tools
#######################################
function determine_target_pattern {
    local base_dir=${1:-.}

    # If ./... is specified, use it as-is
    if [[ "$base_dir" == "./..." ]]; then
        echo "./..."
        return 0
    fi

    # If a specific directory is provided and it exists
    if [[ -d "$base_dir" && "$base_dir" != "." ]]; then
        local go_files_count
        go_files_count=$(has_go_files "$base_dir")
        if [[ "$go_files_count" -gt 0 ]]; then
            echo "$base_dir/..."
            return 0
        fi
    fi

    # For current directory, check if there are Go files
    local go_files_count
    go_files_count=$(find . -name "*.go" -not -path "./vendor/*" -not -path "./.*" | wc -l)

    if [[ "$go_files_count" -eq 0 ]]; then
        error_exit "No Go files found in $base_dir or its subdirectories"
    fi

    # Use ./... for recursive processing from current directory
    echo "./..."
}

#######################################
# Function to run go mod tidy
#######################################
function run_go_mod_tidy {
    echo_section "Running go mod tidy"

    if [[ "$DRY_RUN" == "true" ]]; then
        log "INFO" "DRY-RUN: Would run 'go mod tidy'"
        return 0
    fi

    # 特定ディレクトリ指定時はそのディレクトリ内に go.mod がある場合のみ実行
    local mod_dir="."
    if [[ "$IS_SCOPED" == "true" ]]; then
        if [[ -f "$TARGET_DIR/go.mod" ]]; then
            mod_dir="$TARGET_DIR"
        else
            log "INFO" "Skipping go mod tidy (no go.mod in scoped directory: $TARGET_DIR)"
            return 0
        fi
    fi

    pushd "$mod_dir" >/dev/null || true
    if go mod tidy; then
        log "INFO" "go mod tidy completed successfully (dir=$mod_dir)"
    else
        log "ERROR" "go mod tidy failed (dir=$mod_dir)"
        EXIT_CODE=1
    fi
    popd >/dev/null || true
}

#######################################
# Function to run go fmt
#######################################
function run_go_fmt {
    echo_section "Running go fmt"

    if [[ "$DRY_RUN" == "true" ]]; then
        log "INFO" "DRY-RUN: Would run 'go fmt $TARGET_PATTERN'"
        return 0
    fi

    # Use go fmt instead of gofmt for better Go module support
    if [[ "$FIX_MODE" == "true" ]]; then
        log "INFO" "Automatically formatting files..."
        if go fmt "$TARGET_PATTERN"; then
            log "INFO" "Files formatted successfully"
        else
            log "ERROR" "go fmt failed"
            EXIT_CODE=1
            GO_FMT_FAILED=1
        fi
    else
        # Check if formatting is needed by using gofmt -l
        local fmt_output
        fmt_output=$(go fmt -n "$TARGET_PATTERN" 2>/dev/null || true)

        if [[ -n "$fmt_output" ]]; then
            echo "Files that would be formatted:"
            echo "$fmt_output"
            log "WARN" "Some files need formatting. Use -f flag to auto-fix"
            EXIT_CODE=1
            GO_FMT_FAILED=1
        else
            log "INFO" "All files are properly formatted"
        fi
    fi
}

#######################################
# Function to run go vet
#######################################
function run_go_vet {
    echo_section "Running go vet"

    if [[ "$DRY_RUN" == "true" ]]; then
        log "INFO" "DRY-RUN: Would run 'go vet $TARGET_PATTERN'"
        return 0
    fi

    if go vet "$TARGET_PATTERN"; then
        log "INFO" "go vet passed"
    else
        log "ERROR" "go vet found issues"
        EXIT_CODE=1
        GO_VET_FAILED=1
    fi
}

#######################################
# Function to run golangci-lint
#######################################
function run_golangci_lint {
    echo_section "Running golangci-lint"

    if [[ "$DRY_RUN" == "true" ]]; then
        log "INFO" "DRY-RUN: Would run 'golangci-lint run $TARGET_PATTERN'"
        return 0
    fi

    local lint_args=("run")

    if [[ "$FIX_MODE" == "true" ]]; then
        lint_args+=("--fix")
    fi

    if [[ "$VERBOSE" == "true" ]]; then
        lint_args+=("-v")
    fi

    lint_args+=("$TARGET_PATTERN")

    if golangci-lint "${lint_args[@]}" 2>&1 | tee /tmp/golint_output.txt; then
        log "INFO" "golangci-lint passed"
        LINT_ISSUES_COUNT=0
    else
        log "ERROR" "golangci-lint found issues"
        EXIT_CODE=1
        LINT_FAILED=1
        # Count issues
        LINT_ISSUES_COUNT=$(grep -cE '^\S+\.go:' /tmp/golint_output.txt || echo 0)
    fi
    rm -f /tmp/golint_output.txt
}

#######################################
# Function to run tests
#######################################
function run_tests {
    echo_section "Running go test"

    if [[ "$DRY_RUN" == "true" ]]; then
        log "INFO" "DRY-RUN: Would run 'go test -v $TARGET_PATTERN'"
        return 0
    fi

    local test_args=("-v")

    if [[ "$VERBOSE" == "true" ]]; then
        test_args+=("-x")
    fi

    test_args+=("$TARGET_PATTERN")

    if go test "${test_args[@]}" 2>&1 | tee /tmp/gotest_output.txt; then
        log "INFO" "All tests passed"
        TEST_FAIL_COUNT=0
    else
        log "ERROR" "Some tests failed"
        EXIT_CODE=1
        TEST_FAILED=1
        # Count failed tests
        TEST_FAIL_COUNT=$(grep -c '^--- FAIL:' /tmp/gotest_output.txt || echo 0)
    fi
    rm -f /tmp/gotest_output.txt
}

#######################################
# Function to run race condition tests
#######################################
function run_race_tests {
    echo_section "Running race condition tests"

    if [[ "$DRY_RUN" == "true" ]]; then
        log "INFO" "DRY-RUN: Would run 'CGO_ENABLED=1 go test -race $TARGET_PATTERN'"
        return 0
    fi

    if CGO_ENABLED=1 go test -race "$TARGET_PATTERN"; then
        log "INFO" "Race condition tests passed"
    else
        log "ERROR" "Race condition tests failed"
        EXIT_CODE=1
        RACE_FAILED=1
    fi
}

#######################################
# Function to run coverage tests
#######################################
function run_coverage_tests {
    echo_section "Running coverage tests"

    if [[ "$DRY_RUN" == "true" ]]; then
        log "INFO" "DRY-RUN: Would run coverage tests for $TARGET_PATTERN"
        return 0
    fi

    local coverage_file
    coverage_file="/tmp/coverage_$(date +%s).out"
    if go test -coverprofile="$coverage_file" "$TARGET_PATTERN"; then
        local coverage_percent
        local coverage_output
        coverage_output=$(go tool cover -func="$coverage_file" | grep total | awk '{print $3}')
        coverage_percent=${coverage_output%\%}
        COVERAGE_PERCENT="$coverage_percent"

        if [[ -n "$coverage_percent" ]]; then
            echo "Coverage: ${coverage_percent}%"

            # Compare coverage using awk for better compatibility
            if awk "BEGIN {exit !($coverage_percent >= $COVERAGE_THRESHOLD)}"; then
                log "INFO" "Coverage ($coverage_percent%) meets threshold ($COVERAGE_THRESHOLD%)"
            else
                log "WARN" "Coverage ($coverage_percent%) below threshold ($COVERAGE_THRESHOLD%)"
                EXIT_CODE=1
                COVERAGE_FAILED=1
            fi

            if [[ "$VERBOSE" == "true" ]]; then
                echo "Detailed coverage report:"
                go tool cover -func="$coverage_file"
            fi
        else
            log "WARN" "Could not determine coverage percentage"
        fi

        # Cleanup
        rm -f "$coverage_file"
    else
        log "ERROR" "Coverage tests failed"
        EXIT_CODE=1
        COVERAGE_FAILED=1
    fi
}

#######################################
# Function to run security checks
#######################################
function run_security_checks {
    echo_section "Running security checks"

    if [[ "$DRY_RUN" == "true" ]]; then
        log "INFO" "DRY-RUN: Would run security checks"
        return 0
    fi

    local search_root="."
    if [[ "$IS_SCOPED" == "true" ]]; then
        search_root="$TARGET_DIR"
    fi

    # Check for govulncheck
    if command -v govulncheck &>/dev/null; then
        log "INFO" "Running govulncheck... (root=$search_root)"
        if govulncheck "$TARGET_PATTERN"; then
            log "INFO" "No known vulnerabilities found"
        else
            log "WARN" "Potential vulnerabilities found"
            EXIT_CODE=1
            SECURITY_FAILED=1
        fi
    else
        log "WARN" "govulncheck not installed, skipping vulnerability check"
        log "INFO" "Install with: go install golang.org/x/vuln/cmd/govulncheck@latest"
    fi

    # Check for hardcoded secrets (basic patterns)
    log "INFO" "Checking for potential hardcoded secrets... (root=$search_root)"
    local secret_patterns=(
        "password.*=.*['\"][^'\"]*['\"]"
        "secret.*=.*['\"][^'\"]*['\"]"
        "token.*=.*['\"][^'\"]*['\"]"
        "key.*=.*['\"][^'\"]*['\"]"
        "aws_access_key_id.*=.*['\"][^'\"]*['\"]"
        "aws_secret_access_key.*=.*['\"][^'\"]*['\"]"
    )

    local secrets_found=false
    for pattern in "${secret_patterns[@]}"; do
        if find "$search_root" -name "*.go" -not -path "*/vendor/*" -not -path "*/.*" -exec grep -l -i -E "$pattern" {} + 2>/dev/null | head -1; then
            secrets_found=true
        fi
    done

    if [[ "$secrets_found" == "true" ]]; then   # pragma: allowlist secret   # pragma: allowlist secret
        log "WARN" "Potential hardcoded secrets found. Please review and use environment variables instead."
        EXIT_CODE=1
        SECURITY_FAILED=1
    else
        log "INFO" "No obvious hardcoded secrets found"
    fi
}

#######################################
# Function to run benchmark tests (optional)
#######################################
function run_benchmark_tests {
    echo_section "Running benchmark tests"

    if [[ "$DRY_RUN" == "true" ]]; then
        log "INFO" "DRY-RUN: Would run benchmark tests"
        return 0
    fi

    # Check if there are any benchmark tests
    local has_benchmarks
    has_benchmarks=$(find . -name "*_test.go" -not -path "./vendor/*" -not -path "./.*" -exec grep -l "func Benchmark" {} + 2>/dev/null | head -1)

    if [[ -n "$has_benchmarks" ]]; then
        log "INFO" "Running benchmark tests..."
        if go test -bench=. -benchmem "$TARGET_PATTERN" >/dev/null 2>&1; then
            log "INFO" "Benchmark tests completed"
            if [[ "$VERBOSE" == "true" ]]; then
                go test -bench=. -benchmem "$TARGET_PATTERN"
            fi
        else
            log "WARN" "Some benchmark tests failed"
        fi
    else
        log "INFO" "No benchmark tests found"
    fi
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
            -d | --dry-run)
                DRY_RUN=true
                shift
                ;;
            -f | --fix)
                FIX_MODE=true
                shift
                ;;
            -*)
                error_exit "Unknown option: $1"
                ;;
            *)
                if [[ -z "${TARGET_DIR:-}" ]]; then
                    TARGET_DIR="$1"
                else
                    error_exit "Unexpected argument: $1"
                fi
                shift
                ;;
        esac
    done
    if [[ -z "${TARGET_DIR:-}" ]]; then
        TARGET_DIR="."
        log "INFO" "No target specified, will auto-detect Go directories from current location"
    fi
    if [[ "$TARGET_DIR" != "./..." && ! -d "$TARGET_DIR" && "$TARGET_DIR" != "." ]]; then
        error_exit "Target directory does not exist: $TARGET_DIR"
    fi
    if [[ "$TARGET_DIR" != "." && "$TARGET_DIR" != "./..." ]]; then
        IS_SCOPED=true
    fi
}

#######################################
# Main process
#######################################
function main {
    parse_arguments "$@"
    start_time=$(date +%s)
    validate_dependencies "go" "golangci-lint"
    TARGET_PATTERN=$(determine_target_pattern "$TARGET_DIR")
    echo_section "Starting Go code quality checks"
    log "INFO" "Target input: $TARGET_DIR"
    log "INFO" "Target pattern: $TARGET_PATTERN"
    log "INFO" "Scoped mode: $IS_SCOPED"
    log "INFO" "Verbose mode: $VERBOSE"
    log "INFO" "Dry-run mode: $DRY_RUN"
    log "INFO" "Fix mode: $FIX_MODE"
    run_go_mod_tidy
    run_go_fmt
    run_go_vet
    run_golangci_lint
    run_tests
    run_race_tests
    run_coverage_tests
    run_security_checks
    if [[ "$VERBOSE" == "true" ]]; then
        run_benchmark_tests
    fi
    # Record end time and calculate elapsed time
    end_time=$(date +%s)
    elapsed=$((end_time - start_time))

    if [[ "$EXIT_CODE" -eq 0 ]]; then
        echo_section "All checks completed successfully in ${elapsed} seconds"
        log "INFO" "✅ All validations passed"
    else
        echo_section "Result (completed in ${elapsed} seconds)"
        echo "Result:" >&2
        [[ "$GO_FMT_FAILED" == "1" ]] && echo "❌ go fmt" >&2 || echo "✅ go fmt" >&2
        if [[ "$GO_VET_FAILED" == "1" ]]; then
            echo "❌ go vet" >&2
        else
            echo "✅ go vet" >&2
        fi
        if [[ "$LINT_FAILED" == "1" ]]; then
            echo -n "❌ golangci-lint" >&2
            if [[ "$LINT_ISSUES_COUNT" -gt 0 ]]; then
                echo " ($LINT_ISSUES_COUNT issues)" >&2
            else
                echo "" >&2
            fi
        else
            echo "✅ golangci-lint" >&2
        fi
        if [[ "$TEST_FAILED" == "1" ]]; then
            echo -n "❌ go test" >&2
            if [[ "$TEST_FAIL_COUNT" -gt 0 ]]; then
                echo " ($TEST_FAIL_COUNT failed)" >&2
            else
                echo "" >&2
            fi
        else
            echo "✅ go test" >&2
        fi
        [[ "$RACE_FAILED" == "1" ]] && echo "❌ go test -race" >&2 || echo "✅ go test -race" >&2
        if [[ "$COVERAGE_FAILED" == "1" ]]; then
            if [[ -n "$COVERAGE_PERCENT" ]]; then
                echo "❌ go test -cover ($COVERAGE_PERCENT%)" >&2
            else
                echo "❌ go test -cover" >&2
            fi
        else
            if [[ -n "$COVERAGE_PERCENT" ]]; then
                echo "✅ go test -cover ($COVERAGE_PERCENT%)" >&2
            else
                echo "✅ go test -cover" >&2
            fi
        fi
        [[ "$SECURITY_FAILED" == "1" ]] && echo "❌ security checks (govulncheck or secrets)" >&2 || echo "✅ security checks (govulncheck or secrets)" >&2
        log "ERROR" "❌ Some validations failed"
    fi

    exit $EXIT_CODE
}

# Only call main function if script is executed directly, not sourced
if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    main "$@"
fi
