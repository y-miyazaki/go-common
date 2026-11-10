# Output Format Specification

Use the following report structure for validation output.

````markdown
# Validation Result: <target>

## Summary

- Status: ✅ PASS | ❌ FAIL
- Tools run: <count>
- Tools passed: <count>
- Tools failed: <count>

## Tool Results

| Tool        | Status                        | Message                                    |
| ----------- | ----------------------------- | ------------------------------------------ |
| <tool-name> | ✅ Pass / ❌ Fail / ⊘ Skipped | OK / <brief error summary> / <skip reason> |

## Error Details

### <Tool Name>

```
<Full command output>
```
````

## Rules

- Run tools in the execution order of `scripts/validate.sh`: go mod tidy → go build → golangci-lint → go test → go test -race → go test -cover → govulncheck.
- The script continues after a failing step (no fail-fast). Report every step that ran, then set the overall status from all of them.
- Report a step that did not run as `⊘ Skipped` with the reason: `--skip-tests` (test, race, cover), no `go.mod` in a scoped directory (go mod tidy), or `govulncheck` not installed. When govulncheck is not installed the script's final Result block still prints ✅ for it; use the earlier `govulncheck not installed` warning instead.
- Report coverage below the threshold as `⚠️ Warning`; it does not make Status ❌ FAIL.
- If `go` or `golangci-lint` is missing, the script stops before any check: set Status ❌ FAIL and name the missing binary in `## Error Details`.
- In `## Tool Results`, list every tool in execution order, including skipped ones.
- In `## Error Details`, include only failed tools with their full output.
- If all tools pass: write "All validations passed." in Summary; omit `## Error Details`.
- Status is ❌ FAIL if any tool exits with a non-zero exit code or the script stops early; ✅ PASS otherwise.

## Status Symbols

| Symbol | Meaning | When to Use                     |
| ------ | ------- | ------------------------------- |
| ✅     | Pass    | Tool exited 0 (no issues found) |
| ❌     | Fail    | Tool exited non-zero            |
| ⊘      | Skipped | Tool did not run; state reason  |
| ⚠️     | Warning | Non-blocking finding (coverage) |
