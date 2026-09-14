# Go Validation — Testing (TEST)

Authoring conventions live in companion Go Test rules (stem `go-test`). This file interprets **validation** `TEST-*` failures from `scripts/validate.sh`, not how to write tests.

## go test (TEST)

- TEST-01: All tests pass (exit code 0). Reproduce with the failing package path the script printed; do not rewrite suites to skip failures.
- TEST-02: Race detector reports no races (`-race` when the environment supports it). Treat races as blocking; fix shared-state ownership (companion Go rules CON-*).
- TEST-03: Coverage meets the repository threshold. Raise coverage by tests for new behavior; do not lower the gate.
- TEST-04: No test binary build errors. Fix compile errors in `*_test.go` using companion Go Test rules (stem `go-test`) for naming, tables, and fixtures.

When `go test` fails, load companion Go Test rules for table-driven layout, `errors.Is`, and stubs. Do not copy alternate assertion stacks into a package that already uses stdlib plus go-cmp.
