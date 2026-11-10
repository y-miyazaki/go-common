# Go Validation Checklist

## Execution Order

`scripts/validate.sh` runs these steps in this order and continues after a failing step (no fail-fast); the exit code is non-zero when any blocking step failed:

1. `go mod tidy` — dependency consistency check (`-diff` check; rewrites only with `--fix`; skipped when a scoped directory has no `go.mod`)
2. `go build` — compile check
3. `golangci-lint run` — linting, plus formatting and vet checks when the repository's golangci-lint config enables them (`--fix` applies linter auto-fixes)
4. `go test` — unit tests (skipped with `--skip-tests`)
5. `go test -race` — race detection (skipped with `--skip-tests`)
6. `go test -cover` — coverage against the threshold; below threshold is a non-blocking warning (skipped with `--skip-tests`)
7. `govulncheck` — known vulnerability scan (skipped with a warning when `govulncheck` is not installed)

`go` and `golangci-lint` are required; the script stops before any check when either is missing.

## go mod tidy (MOD)

- MOD-01 (SHOULD): go.mod and go.sum are consistent with source imports
- MOD-02 (SHOULD): No extraneous or missing dependencies

## go build (BUILD)

- BUILD-01 (SHOULD): Target packages compile without errors

## golangci-lint (LINT)

- LINT-01 (SHOULD): All enabled linters pass with zero findings
- LINT-02 (SHOULD): No style violations above configured severity
- LINT-03 (SHOULD): Cyclomatic complexity within threshold
- LINT-04 (SHOULD): No deprecated constructs used
- LINT-05 (SHOULD): Formatting (for example gofumpt/goimports) and vet findings are clean when enabled in the golangci-lint config

## go test (TEST)

- TEST-01 (SHOULD): All tests pass (exit code 0)
- TEST-02 (SHOULD): Race detector reports no races (`-race` flag)
- TEST-03 (SHOULD): Coverage meets project threshold (non-blocking warning when below)
- TEST-04 (SHOULD): No test binary build errors

## govulncheck (SEC)

- SEC-01 (SHOULD): No known vulnerabilities in direct or transitive dependencies
- SEC-02 (SHOULD): All flagged CVEs reviewed and acknowledged if suppressed

## Pass Criteria

- All blocking steps exit with code 0 (a coverage warning does not fail the run)
- Skipped steps are reported as `⊘ Skipped` with the reason
- See [common-output-format.md](common-output-format.md) for output structure
