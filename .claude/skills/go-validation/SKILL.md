---
name: go-validation
description: >-
  Validate Go formatting, linting, tests, and vulnerabilities for maintainable and secure code delivery.
  Use when committing Go changes, running CI validation, or debugging failing checks in repositories.
license: Apache-2.0
metadata:
  author: y-miyazaki
  version: "1.0.2"
---

## Input

- Go path or directory (required)
- Validation script (required): `scripts/validate.sh` (run from the go-validation skill directory)
- `--fix` (optional): applies `go mod tidy` and golangci-lint auto-fixes in place
- `--verbose` (optional): tool-level diagnostics and benchmark tests
- `--skip-tests` (optional): skip go test, race, and coverage checks
- `--dry-run` (optional): print planned commands without running them

## Output Specification

Return structured Markdown in accordance with [references/common-output-format.md](references/common-output-format.md). That file is the source of truth for the output contract.

Structured validation results in fixed tool order.

## Execution Scope

- **Always use `scripts/validate.sh`** for comprehensive validation. Do not run individual commands.
- Individual commands are for debugging only (see [references/common-individual-commands.md](references/common-individual-commands.md)).
- The script runs every step even after a failure (no fail-fast); see [references/common-checklist.md](references/common-checklist.md) for order and skip conditions.
- **Do not modify source files** except through `--fix`.
- **Do not create or delete files**.

### USE FOR:

- run deterministic Go validation before commit or merge
- reproduce CI failures for Go format/lint/test/vulnerability checks
- check blocking validation status and the (non-blocking) coverage threshold warning

### DO NOT USE FOR:

- perform architecture or design reviews
- generate new source code as a primary task
- validate non-Go projects

## Reference Files Guide

- [common-checklist.md](references/common-checklist.md) (always read)
- [common-output-format.md](references/common-output-format.md) (always read)
- [common-troubleshooting.md](references/common-troubleshooting.md) (read on failure)
- [common-individual-commands.md](references/common-individual-commands.md) (read on debugging)
- [category-security.md](references/category-security.md) (read on failure)
- [category-testing.md](references/category-testing.md) (read on failure)

## Workflow

1. Run `bash scripts/validate.sh`.
2. For fast iteration, run `bash scripts/validate.sh <path>` where `<path>` is a Go package or directory.
3. Use `--verbose` to collect tool-level diagnostics.
4. Use `--fix` only for formatting, `go mod tidy`, and linter auto-fixes, then review diffs.
5. Retry at most 2 times after fixes; if checks still fail, return blocking failures and stop.

### Error Handling

| Condition                             | Severity    | Action                                                                                                                                                                                    |
| ------------------------------------- | ----------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `scripts/validate.sh` missing         | Fatal       | Stop; report missing script                                                                                                                                                               |
| No Go files under target path         | Info        | Script exits non-zero before any check (`No Go files found`); report no reviewable Go code and stop                                                                                       |
| `go` or `golangci-lint` missing       | Fatal       | Stop; report the missing binary; set `Status: ❌ FAIL`                                                                                                                                    |
| `govulncheck` missing                 | Recoverable | Script logs `govulncheck not installed` and skips it, but its final Result block still prints ✅ for govulncheck; trust the warning and report govulncheck as `⊘ Skipped` (not installed) |
| Single tool fails, others succeed     | Recoverable | Report every tool; mark the failed tool ❌ Fail with exit status                                                                                                                          |
| All tools fail                        | Fatal       | Set `Status: ❌ FAIL` with per-tool stderr summaries in `## Error Details`                                                                                                                |
| `common-checklist.md` unavailable     | Fatal       | Stop; report missing dependency                                                                                                                                                           |
| `common-output-format.md` unavailable | Recoverable | Note missing file; emit `## Summary`, `## Tool Results`, and `## Error Details` only                                                                                                      |

### Examples

- Prompt: `Validate Go formatting, lint, tests, and vulnerabilities`
- Command: `bash scripts/validate.sh --verbose`
- Result: Structured report per [references/common-output-format.md](references/common-output-format.md) (per-tool pass/fail).
