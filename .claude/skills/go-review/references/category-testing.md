# Testing (TEST)

*_TEST-00 (MUST): Add/update *\_test.go in the same change as behavior*_

Check: When adding or materially changing exported behavior, are corresponding `*_test.go` files added or updated in the same change?
Why: Untested behavior changes are hard to review and regress silently (see [Google eng-practices: Keep related test code in the same CL](https://google.github.io/eng-practices/review/developer/small-cls.html#test_code) and [go.dev: Add a test](https://go.dev/doc/tutorial/add-a-test))
Fix: Add TestXxx functions in `*_test.go` next to the package under test; cover new paths with table-driven subtests

**TEST-01 (SHOULD): Prefer table-driven tests with subtests and edges**

Check: Are []struct format table-driven tests, subtests, and edge cases covered?
Why: Duplicate test cases and Go idiom violations cause test omissions, increased maintenance cost
Fix: []struct format table-driven, use subtests, cover edge cases

**TEST-02 (SHOULD): Use assert vs require correctly; inject time/rand**

Check: Are assert for non-fatal and require for fatal checks used, testable API designed, and time/rand injected?
Why: Excessive testify dependency, untestable APIs, and direct time/randomness usage increase external dependencies, unstable tests
Fix: Decide testify dependency project policy, consider testability, inject time.Now/rand interfaces

**TEST-03 (SHOULD): Mock external deps through consumer interfaces**

Check: Are gomock/testify mock used, interfaces segregated, and dependency injection present?
Why: Real calls to external dependencies cause unstable tests, long execution time, production impact
Fix: Use gomock/testify mock, segregate interfaces, dependency injection

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
