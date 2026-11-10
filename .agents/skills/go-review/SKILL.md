---
name: go-review
description: >-
  Review Go code for security, correctness, and maintainability.
  Use when reviewing Go PRs requiring judgment beyond automated checks.
license: Apache-2.0
metadata:
  author: y-miyazaki
  version: "1.2.4"
---

## Input

- Go files and PR context or file diffs (required); for ad-hoc file review without a PR, the Go files alone
- Validation output from the companion `go-validation` report or CI (optional): used to skip tool-detectable checks

## Output Specification

Return structured Markdown in accordance with [references/common-output-format.md](references/common-output-format.md). That file is the source of truth for the output contract.

Each issue must include file path with line reference (`<path>#L<line>`), risk summary, and remediation guidance.

## Execution Scope

- Systematically apply review checklist from [references/common-checklist.md](references/common-checklist.md)
- Focus on checks requiring human/AI judgment (design, concurrency, security patterns)
- Do not modify code files or approve/merge PRs

### USE FOR:

- review Go PRs where `.go` files are in the changeset
- assess design, security, and concurrency risks not covered by static checks
- perform risk-focused review on multi-package changes
- ad-hoc review of Go source files outside a PR context

### DO NOT USE FOR:

- run formatting/lint/test/vulnerability command pipelines (`gofumpt`, `go vet`, `golangci-lint`, `go test`, `govulncheck`)
- implement code fixes directly
- changesets containing only non-Go files (e.g., docs-only, CI config-only)

## Reference Files Guide

- [common-checklist.md](references/common-checklist.md) (always read)
- [common-output-format.md](references/common-output-format.md) (always read)
- [common-troubleshooting.md](references/common-troubleshooting.md) (read on failure)
- [category-global.md](references/category-global.md) (always read)
- [category-security.md](references/category-security.md) (always read)
- [category-concurrency.md](references/category-concurrency.md) (always read)
- [category-error-handling.md](references/category-error-handling.md) (always read)
- [category-anti-patterns.md](references/category-anti-patterns.md) (always read)
- [category-architecture.md](references/category-architecture.md) (always read)
- [category-code-standards.md](references/category-code-standards.md) (always read)
- [category-context.md](references/category-context.md) (always read)
- [category-dependencies.md](references/category-dependencies.md) (always read)
- [category-documentation.md](references/category-documentation.md) (always read)
- [category-function-design.md](references/category-function-design.md) (always read)
- [category-testing.md](references/category-testing.md) (always read)
- [category-test-anti-patterns.md](references/category-test-anti-patterns.md) (always read)

## Workflow

1. Read PR context and change intent. If PR context is unavailable, review file diffs only and defer PR-context-dependent checks.
2. Scope findings to changed files and hunks. Report violations in unchanged code outside the diff only under `## Expert Assessment`, not `## Issues`. For ad-hoc file review without a diff, treat the whole file as in scope.
3. Exclude generated files and vendored third-party code (for example `vendor/`) from primary findings unless they introduce security-critical risk.
4. For very large PRs (>50 changed Go files), prioritize security/correctness checks first and defer low-risk style checks if evidence is insufficient.
5. When validation output (companion `go-validation` report or equivalent CI results) is supplied, do not re-report failures those tools detect: take the overlapping ItemID's status from the tool result, cite the tool, and spend the review on judgment checks. Category items with a validation `Note:` name the known overlaps.
6. Apply the full review checklist and collect failed/deferred ItemIDs.
7. Output required report sections per [references/common-output-format.md](references/common-output-format.md). Prioritize `SEC-*` findings first.

### Severity and Status Rules

| Status   | When to use                                                                                       |
| -------- | ------------------------------------------------------------------------------------------------- |
| Failed   | Finding is confirmed from source code with concrete evidence (file + line)                        |
| Deferred | Check cannot be evaluated — file too large to fully analyze, or ambiguous without runtime context |
| Passed   | Check evaluated and no issue found (counted in summary only)                                      |

Severity priority for Issues section ordering: `SEC-*` > `CON-*` > `ERR-*` > all others.

### Error Handling

| Condition                               | Severity    | Action                                                                                                |
| --------------------------------------- | ----------- | ----------------------------------------------------------------------------------------------------- |
| `common-checklist.md` unavailable       | Fatal       | Stop, report missing dependency                                                                       |
| `common-output-format.md` unavailable   | Recoverable | Note missing file; emit `## Checks Summary`, `## Checks (Failed/Deferred Only)`, and `## Issues` only |
| PR contains only generated/vendor files | Recoverable | Report "no reviewable Go source" and stop                                                             |

### Examples

- Prompt: `Review Go code changes for design and correctness`
- Result: Structured report per [references/common-output-format.md](references/common-output-format.md); prioritize `SEC-*` findings.
