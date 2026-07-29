# Go Validation Checklist

## Execution Order

Run tools in this order (fail-fast: stop on first failure):

1. `go mod tidy` — dependency consistency check
2. `gofumpt` — formatting compliance
3. `go vet` — static analysis
4. `golangci-lint` — linting and style checks
5. `go test` — unit tests and race condition detection
6. `govulncheck` — known vulnerability scan

## go mod tidy (MOD)

- MOD-01 (SHOULD): go.mod and go.sum are consistent with source imports
- MOD-02 (SHOULD): No extraneous or missing dependencies

## gofumpt (FMT)

- FMT-01 (SHOULD): All .go files are gofumpt-formatted
- FMT-02 (SHOULD): Import blocks organized per goimports conventions

## go vet (VET)

- VET-01 (SHOULD): No static analysis errors reported
- VET-02 (SHOULD): No type-safety violations detected
- VET-03 (SHOULD): No suspicious constructs (unreachable code, printf mismatches)

## golangci-lint (LINT)

- LINT-01 (SHOULD): All enabled linters pass with zero findings
- LINT-02 (SHOULD): No style violations above configured severity
- LINT-03 (SHOULD): Cyclomatic complexity within threshold
- LINT-04 (SHOULD): No deprecated constructs used

## go test (TEST)

- TEST-01 (SHOULD): All tests pass (exit code 0)
- TEST-02 (SHOULD): Race detector reports no races (`-race` flag)
- TEST-03 (SHOULD): Coverage meets project threshold
- TEST-04 (SHOULD): No test binary build errors

## govulncheck (SEC)

- SEC-01 (SHOULD): No known vulnerabilities in direct or transitive dependencies
- SEC-02 (SHOULD): All flagged CVEs reviewed and acknowledged if suppressed

## Pass Criteria

- All tools exit with code 0
- No errors or warnings above configured thresholds
- See [common-output-format.md](common-output-format.md) for output structure
