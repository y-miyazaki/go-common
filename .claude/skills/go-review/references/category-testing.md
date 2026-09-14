# Testing (TEST)

**TEST-00 (MUST): Add or update \*_test.go in the same change as behavior**

Check: When adding or materially changing exported behavior, are corresponding `*_test.go` files added or updated in the same change?
Why: Untested behavior changes are hard to review and regress silently (see [Google eng-practices: Keep related test code in the same CL](https://google.github.io/eng-practices/review/developer/small-cls.html#test_code) and [go.dev: Add a test](https://go.dev/doc/tutorial/add-a-test))
Fix: Add TestXxx functions in `*_test.go` next to the package under test; cover new paths with table-driven subtests

**TEST-01 (SHOULD): Prefer table-driven tests with subtests and edges**

Check: Are []struct format table-driven tests, subtests, and edge cases covered?
Why: Duplicate test cases and Go idiom violations cause test omissions, increased maintenance cost
Fix: []struct format table-driven, use subtests, cover edge cases

**TEST-02 (SHOULD): Design testable APIs; inject time and rand**

Check: Is the API designed for testability, and are time/rand injected rather than called directly?
Why: Direct time/randomness usage and hard-to-substitute dependencies cause flaky tests and block isolated unit verification
Fix: Inject clock/rand through interfaces or function fields; design APIs so tests can substitute dependencies without live I/O

**TEST-03 (SHOULD): Stub external deps through consumer interfaces**

Check: Are external dependencies stubbed or mocked through small consumer-side interfaces with dependency injection?
Why: Real network/SDK calls make tests slow, flaky, and risky in CI
Fix: Define small consumer-side interfaces; substitute hand-written stubs/fakes by default, or generated gomock mocks when call-contract testing is required

**TEST-04 (SHOULD): Share helpers/fixtures outside production packages**

Check: Are testing_test.go separated, common helper functions, and fixture management present?
Why: Duplicate test code and scattered setup/teardown make maintenance difficult, increase test addition cost
Fix: Separate testing_test.go, common helper functions, fixture management

**TEST-05 (SHOULD): Isolate integration tests with build tags**

Check: Are build tags separated, // +build integration, and parallel execution configured?
Why: Mixed unit/integration tests and long execution time delay CI/CD, feedback
Fix: Separate build tags, // +build integration, configure parallel execution

**TEST-06 (SHOULD): Call t.Helper() first in test helpers**

Check: Do test helper functions call t.Helper() as their first statement?
Why: Without t.Helper(), test failure line numbers point to the helper function body rather than the test call site, making failures harder to trace and diagnose
Fix: Add t.Helper() as the first line of every test helper function that calls t.Fatal/t.Error/t.Log

**TEST-07 (SHOULD): Keep one assertion stack per package; match sibling tests**

Check: Does the package use one assertion stack (default stdlib plus go-cmp, or testify at package scope) and match sibling *_test.go layout?
Why: Mixed assertion or mock styles in the same package increase review cost and produce inconsistent failure output
Fix: Default to stdlib `t.Fatalf`/`t.Errorf` plus go-cmp for complex compares; adopt testify `require` or mockery/gomock only at package scope and match existing suite style

**TEST-08 (MUST): Prefix every \*_test.go filename with the source stem under test**

Check: Does every `*_test.go` file name start with the stem of the production file it tests (for example `parser.go` → `parser_test.go` or `parser_error_test.go`)?
Why: Unrelated split names (`error_test.go`, `helpers_test.go`) hide which source file a suite covers and make navigation harder during review
Fix: Rename to `<source-stem>_test.go` or `<source-stem>_<aspect>_test.go`; reserve `example_test.go` and `export_test.go` for godoc examples and external-test export wiring only
