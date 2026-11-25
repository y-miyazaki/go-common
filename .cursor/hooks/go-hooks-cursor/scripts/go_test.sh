#!/bin/bash
#######################################
# Description: Hook for go test.
#              Tests packages containing changed Go files and reports failures
#              in the appropriate format for the active AI agent.
#
# Usage: Called by apm hook runner (not invoked directly).
#        Receives hook event JSON via stdin.
#
# Design Rules:
#   - Exit 0 if tool not found or no changed files (silent skip)
#   - Call report_failure on test failure (agent-aware error signal)
#   - Skip (exit 0) once the same working-tree state was reported HOOK_BLOCK_LIMIT times
#     (default 2) in this session; any file change re-runs the tool
#   - Supports Kiro CLI, Claude Code, Copilot CLI, Cursor, Antigravity, VS Code
#######################################

# Error handling: exit on error, unset variable, or failed pipeline
set -euo pipefail

# Secure defaults
umask 027
export LC_ALL=C.UTF-8

# Capture stdin (hook event JSON) for agent detection.
# Pipe is consumed once; must be read before any other stdin operation.
HOOK_STDIN_DATA=""
if [[ ! -t 0 ]]; then
    HOOK_STDIN_DATA=$(cat)
fi

#######################################
# get_changed_dirs: Collect changed Go locations as "kind<TAB>directory" lines
#
# Description:
#   Gathers modified/added/removed/untracked Go files, go.mod and go.sum from git.
#   Kind P (package) is emitted for a changed or removed .go file whose directory still
#   exists. Kind M (whole module) is emitted for go.mod/go.sum changes and for removed files
#   whose directory is gone, because importers of a removed package can break. testdata and
#   vendor are skipped. Each git command is guarded with || true to prevent pipefail from
#   terminating the script.
#
# Globals:
#   None
#
# Arguments:
#   None
#
# Outputs:
#   Sorted unique "kind<TAB>directory" lines to stdout
#
# Returns:
#   0 on success
#
# Usage:
#   mapfile -t changes < <(get_changed_dirs)
#
#######################################
function get_changed_dirs {
    local path dir kind

    while IFS= read -r path; do
        [[ -z ${path} ]] && continue
        dir=$(dirname "${path}")
        case "/${dir}/" in
            */testdata/* | */vendor/*) continue ;;
        esac
        case "${path##*/}" in
            go.mod | go.sum) kind=M ;;
            *)
                if [[ -d ${dir} ]]; then
                    kind=P
                else
                    kind=M
                fi
                ;;
        esac
        printf '%s\t%s\n' "${kind}" "${dir}"
    done < <(
        {
            git diff --name-only --no-renames -- '*.go' ':(glob)**/go.mod' ':(glob)**/go.sum' 2> /dev/null || true
            git diff --cached --name-only --no-renames -- '*.go' ':(glob)**/go.mod' ':(glob)**/go.sum' 2> /dev/null || true
            git ls-files --others --exclude-standard -- '*.go' ':(glob)**/go.mod' ':(glob)**/go.sum' 2> /dev/null || true
        } | awk 'NF' | sort -u
    ) | sort -u
}

#######################################
# find_go_module_root: Find nearest ancestor directory containing go.mod
#
# Description:
#   A starting directory that no longer exists is replaced by its nearest existing ancestor.
#
# Globals:
#   None
#
# Arguments:
#   $1 - Starting directory path (relative or absolute)
#
# Outputs:
#   Module root absolute path to stdout
#
# Returns:
#   0 if found, 1 otherwise
#
# Usage:
#   find_go_module_root "test/go/server"
#
#######################################
function find_go_module_root {
    local start_dir="$1"
    local dir

    while [[ ! -d $start_dir && $start_dir != "." && $start_dir != "/" ]]; do
        start_dir=$(dirname "$start_dir")
    done

    if ! dir=$(cd "$start_dir" 2> /dev/null && pwd); then
        return 1
    fi

    while [[ $dir != "/" ]]; do
        if [[ -f $dir/go.mod ]]; then
            printf '%s\n' "$dir"
            return 0
        fi
        dir=$(dirname "$dir")
    done

    return 1
}

#######################################
# collect_go_module_roots: Map changed locations to unique Go module roots
#
# Globals:
#   None
#
# Arguments:
#   $@ - "kind<TAB>directory" lines from get_changed_dirs
#
# Outputs:
#   Newline-separated unique module root paths to stdout
#
# Returns:
#   0 on success
#
# Usage:
#   mapfile -t module_roots < <(collect_go_module_roots "${changes[@]}")
#
#######################################
function collect_go_module_roots {
    local line
    local mod_root

    for line in "$@"; do
        [[ -n $line ]] || continue
        if mod_root=$(find_go_module_root "${line#*$'\t'}"); then
            printf '%s\n' "$mod_root"
        fi
    done | sort -u
}

#######################################
# collect_module_packages: Map changed locations to package patterns inside one Go module
#
# Description:
#   Emits ./... when a location of kind M belongs to the module; otherwise one ./pkg pattern
#   per changed package directory.
#
# Globals:
#   None
#
# Arguments:
#   $1 - Module root absolute path
#   $@ - "kind<TAB>directory" lines from get_changed_dirs
#
# Outputs:
#   Newline-separated unique patterns relative to the module root (./pkg or ./...) to stdout
#
# Returns:
#   0 on success
#
# Usage:
#   mapfile -t pkgs < <(collect_module_packages "$mod_root" "${changes[@]}")
#
#######################################
function collect_module_packages {
    local mod_root="$1"
    shift
    local line kind dir abs

    for line in "$@"; do
        [[ -n $line ]] || continue
        kind="${line%%$'\t'*}"
        dir="${line#*$'\t'}"
        [[ $(find_go_module_root "$dir") == "$mod_root" ]] || continue
        if [[ $kind == "M" ]]; then
            printf './...\n'
            continue
        fi
        abs=$(cd "$dir" 2> /dev/null && pwd) || continue
        printf '.%s\n' "${abs#"$mod_root"}"
    done | sort -u
}

#######################################
# filter_buildable_packages: Keep only packages that have Go files for this platform
#
# Description:
#   Packages whose files are all excluded by build constraints make go tools fail, so they
#   are dropped. ./... is kept as is. When go list itself fails, the patterns are returned
#   unfiltered so the real tool reports the error.
#
# Globals:
#   None
#
# Arguments:
#   $1 - Module root absolute path
#   $@ - Package patterns relative to the module root
#
# Outputs:
#   Newline-separated patterns to stdout (none when no package applies)
#
# Returns:
#   0 on success
#
# Usage:
#   mapfile -t pkgs < <(filter_buildable_packages "$mod_root" "${pkgs[@]}")
#
#######################################
function filter_buildable_packages {
    local mod_root="$1"
    shift

    if [[ $1 == "./..." ]] || ! command -v go > /dev/null 2>&1; then
        printf '%s\n' "$@"
        return 0
    fi

    local dirs
    if ! dirs=$(cd "$mod_root" && go list -e -f '{{if or .GoFiles .CgoFiles .InvalidGoFiles .TestGoFiles .XTestGoFiles}}{{.Dir}}{{end}}' "$@" 2> /dev/null); then
        printf '%s\n' "$@"
        return 0
    fi

    local dir
    while IFS= read -r dir; do
        [[ -n $dir ]] || continue
        if [[ $dir == "$mod_root" ]]; then
            printf '.\n'
        else
            printf '.%s\n' "${dir#"$mod_root"}"
        fi
    done <<< "$dirs"
}

#######################################
# truncate_reason_text: Cap reason size for agent responses
#
# Globals:
#   REASON_TRUNCATE_CHARS - optional max characters (default 32768)
#
# Arguments:
#   $1 - text to truncate (printed to stdout)
#
# Outputs:
#   Truncated text to stdout
#
# Returns:
#   None
#######################################
function truncate_reason_text {
    local text="$1"
    local max_chars="${REASON_TRUNCATE_CHARS:-32768}"
    if ((${#text} > max_chars)); then
        printf '%s\n...[truncated]' "${text:0:max_chars}"
    else
        printf '%s' "$text"
    fi
}
#######################################
# collapse_reason_for_cursor_display: Flatten multiline hook reasons for Cursor stop UI
#
# Globals:
#   None
#
# Arguments:
#   $1 - reason text
#
# Outputs:
#   Single-line reason on stdout
#
# Returns:
#   None
#
#######################################
function collapse_reason_for_cursor_display {
    local text="$1"
    printf '%s' "$text" | tr '\n' ' ' | sed 's/[[:space:]]\+/ /g;s/^[[:space:]]*//;s/[[:space:]]*$//'
}

#######################################
# emit_json_with_reason: Build hook JSON via stdin (avoids ARG_MAX)
#
# Globals:
#   None
#
# Arguments:
#   $1 - reason text
#   $2 - jq filter using . for the reason value
#
# Outputs:
#   JSON object to stdout
#
# Returns:
#   None
#######################################
function emit_json_with_reason {
    local text="$1"
    local jq_filter="$2"
    truncate_reason_text "$text" | jq -Rs "$jq_filter"
}

#######################################
# report_failure: Emit error in the format the current agent expects, then exit.
#
# Description:
#   Identifies the AI agent from HOOK_STDIN_DATA structure, then returns
#   the agent-specific response format:
#     - Kiro CLI: stop → {"decision":"block","reason":"..."}
#     - Claude Code: Stop → {"decision":"block"}, PostToolUse → hookSpecificOutput
#     - GitHub Copilot: agentStop → {"decision":"block"}, postToolUse → additionalContext
#     - Antigravity: Stop → {"decision":"continue","reason":"..."}
#     - Cursor: stop → followup_message, other events → exit 2 + stderr
#     - unknown: exit 2 + stderr
#
# Globals:
#   None
#
# Arguments:
#   $1 - reason: Human-readable description of what failed and how to fix it
#
# Outputs:
#   Writes JSON to stdout or errors to stderr
#
# Returns:
#   Exits with 0 or 2
#
# Usage:
#   report_failure "go test failed: ..."
#
#######################################
function report_failure {
    local reason="$1"
    record_blocked_state
    local agent=""
    local hook_event=""

    # Step 1: Detect agent (agent-first strategy)
    if [[ -n $HOOK_STDIN_DATA ]]; then
        # 1. Antigravity (highest priority - unique fields)
        if echo "$HOOK_STDIN_DATA" | jq -e ".terminationReason" > /dev/null 2>&1; then
            agent="antigravity"
        elif echo "$HOOK_STDIN_DATA" | jq -e ".toolCall" > /dev/null 2>&1; then
            agent="antigravity"

        # 2. Check for hook_event_name first (most reliable discriminator)
        elif echo "$HOOK_STDIN_DATA" | jq -e 'has("cursor_version") or has("generation_id") or has("workspace_roots")' > /dev/null 2>&1; then
            agent="cursor"
            hook_event=$(echo "$HOOK_STDIN_DATA" | jq -r '.hook_event_name // "stop"' 2> /dev/null)
        elif echo "$HOOK_STDIN_DATA" | jq -e ".hook_event_name" > /dev/null 2>&1; then
            hook_event=$(echo "$HOOK_STDIN_DATA" | jq -r '.hook_event_name' 2> /dev/null)

            # Check event name pattern to determine agent type
            if echo "$hook_event" | grep -qE '^(afterFileEdit|beforeShellExecution|beforeMCPExecution|beforeReadFile)$'; then
                # camelCase with Cursor-only event names
                agent="cursor"
            elif echo "$hook_event" | grep -qE '^(Stop|PostToolUse|PreToolUse)$'; then
                # PascalCase = Claude Code
                agent="claude_code"
            elif echo "$hook_event" | grep -qE '^(stop|postToolUse|preToolUse|agentSpawn|userPromptSubmit)$'; then
                # camelCase with Kiro values = Kiro
                agent="kiro"
            else
                # Default to Claude Code for unknown PascalCase
                agent="claude_code"
            fi

        # 3. Copilot CLI (env var or Copilot-unique fields, no hook_event_name)
        elif [[ -n ${GITHUB_COPILOT_API_TOKEN:-} ]] \
            || echo "$HOOK_STDIN_DATA" | jq -e '.transcriptPath // .stopReason // .stop_reason // .toolResult // .tool_result' > /dev/null 2>&1; then
            agent="copilot"
            if echo "$HOOK_STDIN_DATA" | jq -e ".stopReason" > /dev/null 2>&1; then
                hook_event="agentStop"
            elif echo "$HOOK_STDIN_DATA" | jq -e ".toolResult" > /dev/null 2>&1; then
                hook_event="postToolUse"
            elif echo "$HOOK_STDIN_DATA" | jq -e ".toolName" > /dev/null 2>&1; then
                hook_event="preToolUse"
            fi

        # 4. VS Code extension (fallback based on vscode-specific fields)
        elif echo "$HOOK_STDIN_DATA" | jq -e 'has("stop_hook_active") or has("tool_use_id")' > /dev/null 2>&1; then
            agent="vscode"
            if echo "$HOOK_STDIN_DATA" | jq -e 'has("stop_hook_active")' > /dev/null 2>&1; then
                hook_event="Stop"
            elif echo "$HOOK_STDIN_DATA" | jq -e 'has("tool_use_id")' > /dev/null 2>&1; then
                hook_event="PostToolUse"
            fi
        fi
    fi

    # Final fallback: env var check
    if [[ -z $agent && -n ${GITHUB_COPILOT_API_TOKEN:-} ]]; then
        agent="copilot"
    fi

    # Step 2: Build response per agent spec (A-Z order)
    if [[ $agent == "cursor" && $hook_event == "stop" ]]; then
        reason=$(collapse_reason_for_cursor_display "$reason")
    fi

    case "$agent" in
        antigravity)
            emit_json_with_reason "$reason" '{decision: "continue", reason: .}'
            exit 0
            ;;
        claude_code)
            if [[ $hook_event == "Stop" ]]; then
                emit_json_with_reason "$reason" '{decision: "block", reason: .}'
                exit 0
            elif [[ $hook_event == "PostToolUse" ]]; then
                emit_json_with_reason "$reason" '{hookSpecificOutput: {hookEventName: "PostToolUse", additionalContext: .}}'
                exit 0
            else
                echo "$reason" >&2
                exit 2
            fi
            ;;
        copilot)
            case "$hook_event" in
                Stop | agentStop)
                    emit_json_with_reason "$reason" '{decision: "block", reason: .}'
                    exit 0
                    ;;
                PostToolUse | postToolUse)
                    emit_json_with_reason "$reason" '{additionalContext: .}'
                    exit 0
                    ;;
                *)
                    echo "$reason" >&2
                    exit 2
                    ;;
            esac
            ;;
        cursor)
            if [[ $hook_event == "stop" ]]; then
                emit_json_with_reason "$reason" '{followup_message: .}'
                exit 0
            else
                echo "$reason" >&2
                exit 2
            fi
            ;;
        kiro)
            if [[ $hook_event == "stop" ]]; then
                emit_json_with_reason "$reason" '{decision: "block", reason: .}'
                exit 0
            else
                echo "$reason" >&2
                exit 2
            fi
            ;;
        vscode)
            if [[ $hook_event == "Stop" ]]; then
                emit_json_with_reason "$reason" '{hookSpecificOutput: {hookEventName: "Stop", decision: "block", reason: .}}'
                exit 0
            elif [[ $hook_event == "PostToolUse" ]]; then
                emit_json_with_reason "$reason" '{decision: "block", reason: ., hookSpecificOutput: {hookEventName: "PostToolUse", additionalContext: .}}'
                exit 0
            else
                echo "$reason" >&2
                exit 2
            fi
            ;;
        *)
            echo "$reason" >&2
            exit 2
            ;;
    esac
}

#######################################
# blocked_state_file: Resolve the per-session, per-hook file that stores the blocked state
#
# Globals:
#   HOOK_STDIN_DATA - hook event JSON (session_id / conversation_id / sessionId / conversationId)
#   TMPDIR          - optional base directory (default /tmp)
#
# Arguments:
#   None
#
# Outputs:
#   State file path to stdout
#
# Returns:
#   0 always
#
#######################################
function blocked_state_file {
    local session name
    session=$(printf '%s' "${HOOK_STDIN_DATA:-}" | jq -r '.session_id // .conversation_id // .sessionId // .conversationId // empty' 2> /dev/null || true)
    session=$(printf '%s' "${session:-nosession}" | tr -c 'A-Za-z0-9_-' '_')
    name=$(basename "${BASH_SOURCE[0]}" .sh)
    printf '%s/hook-state-%s/%s__%s\n' "${TMPDIR:-/tmp}" "$(id -u)" "${session:0:64}" "$name"
}

#######################################
# clear_blocked_state: Forget the blocked state (hash and session counts) after a passing run
#
# Globals:
#   None
#
# Arguments:
#   None
#
# Outputs:
#   None
#
# Returns:
#   0 always
#
#######################################
function clear_blocked_state {
    rm -f "$(blocked_state_file)" 2> /dev/null || true
}

#######################################
# compute_state_hash: Hash the working-tree change state (tracked diff + untracked content)
#
# Globals:
#   None
#
# Arguments:
#   None
#
# Outputs:
#   Hex digest to stdout (empty when no hash command is available)
#
# Returns:
#   0 always
#
#######################################
function compute_state_hash {
    local -a hasher
    if command -v sha256sum > /dev/null 2>&1; then
        hasher=(sha256sum)
    elif command -v shasum > /dev/null 2>&1; then
        hasher=(shasum -a 256)
    else
        return 0
    fi

    {
        git diff HEAD 2> /dev/null || true
        git ls-files --others --exclude-standard -z 2> /dev/null | xargs -0 -r "${hasher[@]}" 2> /dev/null || true
    } | "${hasher[@]}" | cut -d' ' -f1
}

#######################################
# is_blocked_state: Decide from the saved state whether the hook should let the agent stop
#
# Description:
#   One shared decision, controlled by variables the hook sets in main (all optional):
#   the same state reported HOOK_BLOCK_LIMIT times, or the session reported HOOK_SESSION_LIMIT
#   times in total. stop_hook_active is not consulted: a changed state is asked about again.
#   After the session limit the hook no longer runs its check unless it clears the state earlier
#   (a hook with a cheap clean-tree test, such as review_gate.sh, does).
#
# Globals:
#   HOOK_BLOCK_LIMIT - max reports per unchanged state (default 2)
#   HOOK_SESSION_LIMIT - max reports per session across changing states (default 0 = unlimited)
#   HOOK_STATE_HASH - state hash captured at hook start
#
# Arguments:
#   None
#
# Outputs:
#   None
#
# Returns:
#   0 when the hook should pass (limit reached), 1 when it should report
#
#######################################
function is_blocked_state {
    local current="${HOOK_STATE_HASH:-}"
    [[ -n $current ]] || return 1

    local limit="${HOOK_BLOCK_LIMIT:-2}" session_limit="${HOOK_SESSION_LIMIT:-0}"
    [[ $limit =~ ^[0-9]+$ ]] || limit=2
    [[ $session_limit =~ ^[0-9]+$ ]] || session_limit=0

    local file saved="" count=0 total=0
    file=$(blocked_state_file)
    [[ -f $file ]] || return 1
    read -r saved count total < "$file" || true
    [[ ${count:-0} =~ ^[0-9]+$ ]] || count=0
    [[ ${total:-0} =~ ^[0-9]+$ ]] || total=0
    if ((session_limit > 0 && total >= session_limit)); then
        return 0
    fi
    [[ $saved == "$current" ]] || return 1
    ((count >= limit))
}

#######################################
# record_blocked_state: Count one more report for the post-run state
#
# Description:
#   Re-hashes after the tool ran so auto-fixes made by the tool itself are included.
#   The per-state count restarts at 1 when the state differs from the stored one;
#   the session total keeps counting until clear_blocked_state.
#
# Globals:
#   None
#
# Arguments:
#   None
#
# Outputs:
#   None
#
# Returns:
#   0 always
#
#######################################
function record_blocked_state {
    local hash file saved="" count=0 total=0
    hash=$(compute_state_hash)
    [[ -n $hash ]] || return 0
    file=$(blocked_state_file)
    if [[ -f $file ]]; then
        read -r saved count total < "$file" || true
    fi
    [[ ${count:-0} =~ ^[0-9]+$ ]] || count=0
    [[ ${total:-0} =~ ^[0-9]+$ ]] || total=0
    if [[ $saved == "$hash" ]]; then
        count=$((count + 1))
    else
        count=1
    fi
    total=$((total + 1))
    {
        mkdir -p "$(dirname "$file")" \
            && printf '%s %s %s\n' "$hash" "$count" "$total" > "$file"
    } 2> /dev/null || true
}

#######################################
# main: Entry point
#
# Description:
#   Runs go test (-short, bounded by GO_TEST_TIMEOUT, default 90s) for the changed
#   packages of each Go module root that contains changed files.
#   Collects failures and calls report_failure with a summary.
#
# Globals:
#   None
#
# Arguments:
#   None
#
# Outputs:
#   None
#
# Returns:
#   0 on success or skip
#
# Usage:
#   main
#
#######################################
function main {
    command -v jq > /dev/null 2>&1 || exit 0
    command -v go > /dev/null 2>&1 || exit 0

    local root
    root=$(git rev-parse --show-toplevel 2> /dev/null) || exit 0
    cd "$root" || exit 0

    HOOK_STATE_HASH=$(compute_state_hash)
    if is_blocked_state; then
        exit 0
    fi

    local changes=()
    mapfile -t changes < <(get_changed_dirs)

    if ((${#changes[@]} == 0)); then
        exit 0
    fi

    local module_roots=()
    mapfile -t module_roots < <(collect_go_module_roots "${changes[@]}")

    if ((${#module_roots[@]} == 0)); then
        exit 0
    fi

    local fails=0
    local output=""
    for mod_root in "${module_roots[@]}"; do
        [[ -n $mod_root && -d $mod_root ]] || continue
        local rel_mod="${mod_root#"$root"/}"
        local pkgs=()
        mapfile -t pkgs < <(collect_module_packages "$mod_root" "${changes[@]}")
        ((${#pkgs[@]} > 0)) || continue
        local pkg
        for pkg in "${pkgs[@]}"; do
            if [[ $pkg == "./..." ]]; then
                pkgs=("./...")
                break
            fi
        done
        mapfile -t pkgs < <(filter_buildable_packages "$mod_root" "${pkgs[@]}")
        ((${#pkgs[@]} > 0)) || continue
        local result
        result=$(cd "$mod_root" && go test -short -timeout "${GO_TEST_TIMEOUT:-90s}" "${pkgs[@]}" 2>&1) || {
            fails=$((fails + 1))
            output+="=== ${rel_mod:-.} ==="$'\n'
            output+="${result}"$'\n'
        }
    done

    if [[ $fails -gt 0 ]]; then
        report_failure "go test failed in changed Go packages:
${output}"
    fi

    clear_blocked_state
}

if [[ ${BASH_SOURCE[0]} == "${0}" ]]; then
    main "$@"
fi
