---
paths:
  - "**/*.go"
---

# Go Test Instructions

## Scope

- Scope covers authoring and updating `*_test.go` files and test doubles (stubs/fakes) in the same package.
- Production Go rules remain in companion Go rules (stem `go`); this file is the test companion.
- When adding or changing behavior in a `.go` source file, add or update the paired `*_test.go` in the same change (MUST) — same obligation as companion Go rules (stem `go`).
- Content here is **Go test practice**, not domain/business rules. Domain fixtures (AWS ARNs, API payloads) belong in each test; patterns belong here.

### Rule application (`globs` / `paths`)

| Pattern                 | Use when                                                           | Result                                                                            |
| ----------------------- | ------------------------------------------------------------------ | --------------------------------------------------------------------------------- |
| `**/*.go` (recommended) | Test rules must apply while editing **production** `.go` files too | Suite pairing and test templates inject on source edits (G-05 companion coverage) |
| `**/*_test.go` only     | Test rules should inject **only** when a test file is open         | Pairing obligation is easy to miss when changing `foo.go`                         |

**Why Bats lists two globs but Go uses one:** Bats pairs `*.sh` with `*.bats` — different extensions, so both must be listed. Go tests live in `*_test.go`, which already matches `**/*.go`. Listing `**/*_test.go` in addition is optional documentation only; it does not change Cursor behavior.

## Standards

### Naming Conventions

| Component      | Rule                                                         | Example                                                |
| -------------- | ------------------------------------------------------------ | ------------------------------------------------------ |
| Test file      | `<source-stem>_test.go`; splits keep `<source-stem>_` prefix | `parser.go` → `parser_test.go`, `parser_error_test.go` |
| Test func      | `Test<TypeOrFunc>_<Behavior>`                                | `TestParser_NormalizesInput`                           |
| Subtest name   | Lowercase descriptive phrase                                 | `empty input returns error`                            |
| Sentinel error | `errTest<Name>` at package scope                             | `errTestNotFound`                                      |
| Stub type      | `stub<Role>`                                                 | `stubStore`, `stubRunCollector`                        |

### Default Test Stack

Use this stack unless the **package** has already adopted optional tools (see Optional Tools):

| Layer                | Tool                              | Role                                        |
| -------------------- | --------------------------------- | ------------------------------------------- |
| Test structure       | table-driven `[]struct` + `t.Run` | one `Test*` function loops rows (TBL-01)    |
| Assertions (simple)  | stdlib `testing`                  | `t.Fatalf` / `t.Errorf` with `got` / `want` |
| Assertions (complex) | `github.com/google/go-cmp/cmp`    | struct, slice, map deep equality            |
| Test doubles         | hand-written stub/fake            | consumer-interface substitution             |

### File Header

Lint `*_test.go` with the same toolchain the repository uses for tests. Do not disable revive or golangci as a whole. When revive `comments-density` is enabled on tests, a new table-driven file may start with:

```go
//revive:disable:comments-density reason: table-driven tests are self-explanatory via subtest names.
package <same-as-source>
```

Omit this directive when that rule is not enabled.

### Table-Driven Test Template

```go
var errTestInvalid = errors.New("invalid") // FIX-02: package-level sentinel (err113)

func TestExample(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{name: "empty input", input: "", want: ""},
		{name: "invalid input", input: "!", wantErr: errTestInvalid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Example(tt.input)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("Example(%q) error = nil, want %v", tt.input, tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Example(%q) error = %v, want %v", tt.input, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Example(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("Example(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
```

On Go before 1.22, large rows, or gocritic `rangeValCopy`, use `for i := range tests { tt := tests[i] }`. Omit `t.Parallel()` only when the test shares mutable state (TBL-06). Complex compares: `cmp.Diff(tt.want, got)`.

### Optional Tools (OPT)

Apply only when the **package** has adopted the tool. Do not mix testify with the default stdlib plus go-cmp stack (CONS-01).

| Tool              | When to use                                      | Rule                                         |
| ----------------- | ------------------------------------------------ | -------------------------------------------- |
| testify `require` | package already uses testify                     | `require` only, never `assert`               |
| mockery + gomock  | call count, order, or arguments are the behavior | keep consumer interfaces small (1–3 methods) |
| hand-written stub | default doubles                                  | `func` fields on `stub<Role>`                |

## Guidelines

### Naming (NAME)

- NAME-01 (MUST): Every `*_test.go` filename must start with the stem of the production file it tests (`<source-stem>_test.go`). When splitting a suite, keep that prefix (`<source-stem>_<aspect>_test.go`) — do not use unrelated names (`error_test.go`, `helpers_test.go`) that hide the target file
- NAME-02 (SHOULD): Reserve `example_test.go` and `export_test.go` only for godoc `Example` functions and external-test export wiring — not general unit tests

### Table-driven tests (TBL)

- TBL-01 (MUST): Default to `[]struct` + `for _, tt := range tests` + `t.Run` for unit tests — including a single case (table of one row is fine). Requires Go 1.22+ when combining `for _, tt := range tests` with `t.Parallel()` inside subtests
- TBL-02 (MUST): Put `name` as the first struct field; use lowercase descriptive subtest names
- TBL-03 (SHOULD): When success and error paths share setup, combine them in one table with `wantErr` (or equivalent). `TestFoo_OK` and `TestFoo_Error` may be separate when fixtures or assertions differ; do not add a top-level `Test*` per table row
- TBL-04 (SHOULD): Add edge cases as new table rows, not new top-level `Test*` functions
- TBL-05 (SHOULD): Skip the table only when setup differs materially per scenario (heavy fixtures, integration wiring) — use a dedicated `Test*` with a comment explaining why
- TBL-06 (SHOULD): Call `t.Parallel()` on the parent test and each subtest. Omit it only when the test shares mutable state with other tests or rows (shared fixture, package vars, reused stub, `t.Setenv`, working directory)

### Package choice (PKG)

- PKG-01 (MUST): Use `package foo` (same package) when testing unexported functions or types
- PKG-02 (SHOULD): Use `package foo_test` only when the suite exercises the public API exclusively

### Fixtures (FIX)

- FIX-01 (MUST): Prefer inline literals or small `const`/`var` fixtures in the test file for unit tests
- FIX-02 (MUST): Use package-level `var errTest… = errors.New("…")` for stub errors — never `errors.New` inside a test body (`err113`)
- FIX-03 (MUST): Initialize maps with `make(map[K]V)`
- FIX-04 (SHOULD): Put large payloads in `testdata/` only when a sibling suite in the same package already does; otherwise keep fixtures inline

### Assertions (ASSERT)

- ASSERT-01 (MUST): Inside `t.Run` subtests, use `t.Fatalf` for a failed case. Without `t.Run`, use `t.Errorf` so later table rows still run
- ASSERT-02 (MUST): Format failure messages as `Func(args…) = got, want want` (or `error = %v, want %v` for errors)
- ASSERT-03 (MUST): Use `errors.Is` / `errors.As` for error identity — not string comparison on `err.Error()`
- ASSERT-04 (SHOULD): Use `go-cmp` when comparing structs, slices, or maps; use direct `!=` for scalars

### Stubs and fakes (STUB)

- STUB-01 (MUST): Implement the consumer interface in the test file; match production method signatures exactly
- STUB-02 (MUST): Name unused parameters `_` or assign to `_` to avoid compile errors
- STUB-03 (SHOULD): Drive behavior with struct fields or `func` fields before adopting generated mocks
- STUB-04 (SHOULD): Add `//nolint:gocritic` only when the stub must match a `hugeParam` interface signature

### Optional mocks (MOCK)

- MOCK-01 (SHOULD): Adopt mockery/gomock only when call contract (count, order, arguments) is what the test must prove
- MOCK-02 (SHOULD): Keep generated mocks in the package or `mocks/` subtree the repository already uses — do not introduce a new layout without cause
- MOCK-03 (SHOULD): Do not mock types you do not own wholesale — define a small consumer interface instead

### Consistency (CONS)

- CONS-01 (MUST): Use one assertion stack per package — default is stdlib plus go-cmp; do not mix that default with testify
- CONS-02 (MUST): Match the nearest sibling `*_test.go` layout before introducing a new style
- CONS-03 (SHOULD): Adopt optional tools at package scope, not file-by-file

### Boundary without live I/O (BOUND)

- BOUND-01 (MUST): Unit-test pure functions, parsers, classifiers, and cache keys without network/SDK clients
- BOUND-02 (MUST): Do not call constructors that perform I/O against zero-value or live clients in default unit tests
- BOUND-03 (SHOULD): Gate real I/O behind `//go:build integration` only when the repository already uses that tag

### Testing (TEST)

- TEST-01 (SHOULD): Prefer table-driven tests with subtests and edges (see TBL-01)
- TEST-02 (SHOULD): Design testable APIs; inject time and rand through interfaces or function fields
- TEST-03 (SHOULD): Stub external deps through consumer interfaces (see STUB/MOCK)
- TEST-04 (SHOULD): Share helpers and fixtures in `testing_test.go` or package test helpers — not production code
- TEST-05 (SHOULD): Isolate integration tests with `//go:build integration` (see BOUND-03)
- TEST-06 (SHOULD): Call `t.Helper()` as the first statement in every test helper function
- TEST-07 (SHOULD): Keep one assertion stack per package; match sibling `*_test.go` layout (see CONS-01)

### Test Anti-Patterns (TAP)

- TAP-01 (SHOULD): One `Test*` function per table row when a table-driven loop would suffice
- TAP-02 (SHOULD): Mixing testify `require`/`assert` with the default stdlib plus go-cmp stack in the same package
- TAP-03 (SHOULD): Using `testify/assert` (non-fatal) in table-driven subtests
- TAP-04 (MUST): Inline `errors.New` in test bodies (`err113`)
- TAP-05 (SHOULD): Copying large table rows with `for _, tt := range tests` when gocritic `rangeValCopy` is enabled
- TAP-06 (MUST): Split `*_test.go` files named without the source stem prefix
- TAP-07 (SHOULD): Disabling revive or golangci as a whole on `*_test.go`
- TAP-08 (SHOULD): Adding `//revive:disable:comments-density` when that rule is not enabled on tests
- TAP-09 (SHOULD): Refactoring production code solely to inject mocks unless a test already justifies the seam
- TAP-10 (MUST): Real credentials, tokens, or production identifiers in fixtures

### Code Modification Guidelines

- Open the nearest sibling `*_test.go` in the same package and mirror its layout before writing new tests.

## Testing and Validation

Operational notes:

- Separate integration tests with `//go:build integration` when the repository uses that tag.
- Target meaningful coverage on new behavior.

On-demand validation: see go-validation skill SKILL.md.

## Security Guidelines

- Do not embed real credentials, tokens, or production identifiers in fixtures. Use obvious placeholders.
- Validate external inputs in test helpers and table rows; do not pass untrusted strings into production parsers without boundary cases.
- Do not log or print sensitive values in failure messages; redact tokens and identifiers in `t.Fatalf` output.
