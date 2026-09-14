# Test Anti-Patterns (TAP)

**TAP-01 (SHOULD): One `Test*` function per table row when a table-driven loop would suffice**

Check: Are separate top-level `Test*` functions used for cases that share setup and differ only by inputs (TBL-01)?
Why: Duplicate setup hides edge cases and makes suites harder to extend during review
Fix: Consolidate into one `Test*` with `[]struct` and `t.Run`; see companion Go Test rules (stem `go-test`)

**TAP-02 (SHOULD): Mixing testify `require`/`assert` with the default stdlib plus go-cmp stack in the same package**

Check: Does the package mix testify `require`/`assert` with stdlib `t.Fatalf`/`t.Errorf` and go-cmp (CONS-01, TEST-07)?
Why: Mixed assertion stacks produce inconsistent failure output and increase review cost
Fix: Adopt one stack at package scope; default to stdlib plus go-cmp unless the package already standardized on testify

**TAP-03 (SHOULD): Using `testify/assert` (non-fatal) in table-driven subtests**

Check: Are non-fatal `testify/assert` calls used inside `t.Run` subtests?
Why: Non-fatal assertions allow subtests to continue after failure, producing mixed failure modes
Fix: Use `t.Fatalf` inside `t.Run`, or testify `require` only — not `assert`

**TAP-04 (MUST): Inline `errors.New` in test bodies (`err113`)**

Check: Are sentinel errors created with inline `errors.New` inside test functions instead of package-level `errTest*` variables (FIX-02)?
Why: Inline sentinels trigger err113 and break stable `errors.Is` identity across table rows
Fix: Declare `var errTest<Name> = errors.New("...")` at package scope

**TAP-05 (SHOULD): Copying large table rows with `for _, tt := range tests` when gocritic `rangeValCopy` is enabled**

Check: Does the table loop copy large structs when `rangeValCopy` is enabled (TBL-02)?
Why: Silent copies inflate test memory and can mask performance regressions
Fix: Use `for i := range tests { tt := tests[i] }` or store pointers in table rows

**TAP-06 (MUST): Split `*_test.go` files named without the source stem prefix**

Check: Are test files named without the production source stem prefix (TEST-08, NAME-01)?
Why: Names like `error_test.go` hide which source file the suite covers
Fix: Rename to `<source-stem>_test.go` or `<source-stem>_<aspect>_test.go`

**TAP-07 (SHOULD): Disabling revive or golangci as a whole on `*_test.go`**

Check: Are file-wide `//nolint` or linter-disable directives applied to entire `*_test.go` files?
Why: Blanket disables hide real test bugs (unchecked errors, races, secrets in fixtures)
Fix: Target specific linters or rules with a documented reason; keep checks that catch real test bugs

**TAP-08 (SHOULD): Adding `//revive:disable:comments-density` when that rule is not enabled on tests**

Check: Is a `comments-density` revive disable present when that rule is not enabled for tests?
Why: Unnecessary directives add noise and suggest copy-paste without local lint context
Fix: Remove the directive when the rule is not enabled; add it only when comments-density is active on tests

**TAP-09 (SHOULD): Refactoring production code solely to inject mocks unless a test already justifies the seam**

Check: Was production code restructured only to inject mocks when hand-written stubs would suffice (MOCK-01)?
Why: Unjustified seams widen the API surface and increase maintenance cost
Fix: Prefer consumer-side stubs in the test file; adopt mockery/gomock only when call contract is what the test must prove

**TAP-10 (MUST): Real credentials, tokens, or production identifiers in fixtures**

Check: Do fixtures or table rows contain real API keys, passwords, tokens, or production identifiers (G-01)?
Why: Secrets in the tree leak via review, CI logs, and clones
Fix: Use obvious placeholders; load real secrets from the environment only in gated integration tests
