# Go Review Checklist

## Architecture (ARCH)

- ARCH-01 (SHOULD): Separate handler/usecase/repository from infrastructure
- ARCH-02 (SHOULD): Inject deps via constructor interfaces (not package globals)
- ARCH-03 (SHOULD): Keep domain logic free of DB/HTTP/API concerns
- ARCH-04 (SHOULD): Avoid circular deps; use standard layout and internal/
- ARCH-05 (SHOULD): Abstract external integrations behind consumer interfaces

## Code Standards (CODE)

- CODE-01 (MUST): Keep interfaces small (1-3 methods) on the consumer side
- CODE-02 (SHOULD): Minimize exported API surface; hide internals with internal/
- CODE-03 (SHOULD): Unexport invariant fields and mutexes; split oversized structs

## Concurrency (CON)

- CON-01 (SHOULD): Ensure goroutines exit (watch context.Done / completion)
- CON-02 (SHOULD): Only the sender closes a channel (once)
- CON-03 (SHOULD): Define lock/completion ownership for shared state

## Context Handling (CTX)

- CTX-01 (MUST): Exported I/O APIs take context.Context as first param
- CTX-02 (SHOULD): Do not store ambiguous request-scoped contexts
- CTX-03 (SHOULD): Pass context into goroutines that do I/O or wait
- CTX-04 (SHOULD): Call cancel promptly; do not leak derived contexts

## Dependencies (DEP)

- DEP-01 (SHOULD): List direct deps in go.mod with pinned versions

## Documentation (DOC)

- DOC-01 (SHOULD): Package has a doc comment stating purpose
- DOC-02 (MUST): Public APIs have godoc covering args/returns/errors
- DOC-03 (SHOULD): Comments are consistently English

## Error Handling (ERR)

- ERR-01 (MUST): Wrap errors with fmt.Errorf %w and context
- ERR-02 (SHOULD): Use distinct sentinel/custom errors per failure mode
- ERR-03 (SHOULD): Panic only for fatal bugs; recover at boundaries
- ERR-04 (SHOULD): Timeouts/retries and classify external errors
- ERR-05 (SHOULD): Do not leak internals in user-facing error messages
- ERR-06 (MUST): Never discard errors with _ unless commented

## Function Design (FUNC)

- FUNC-01 (SHOULD): Split mixed-responsibility or multi-layer functions
- FUNC-02 (SHOULD): Unify pointer vs value receivers; avoid large values
- FUNC-03 (SHOULD): Keep generic constraints minimal and locally scoped

## Global / Base (G)

- G-01 (SHOULD): No API keys/passwords/tokens in source
- G-02 (SHOULD): Keep init() free of I/O, panics, and heavy side effects
- G-03 (SHOULD): Prefer types whose zero value is usable
- G-04 (SHOULD): Copy slices/maps at API boundaries

## Security (SEC)

- SEC-01 (SHOULD): Validate inputs; ban string-concat SQL
- SEC-02 (SHOULD): Escape/sanitize outputs for HTML/JSON/CRLF sinks
- SEC-03 (SHOULD): Authenticate endpoints; verify JWT; enforce RBAC
- SEC-04 (SHOULD): Mask passwords/tokens in logs
- SEC-05 (SHOULD): Least privilege; no production debug; explicit CORS

## Testing (TEST)

- TEST-00 (MUST): Add or update \*_test.go in the same change as behavior
- TEST-01 (SHOULD): Prefer table-driven tests with subtests and edges
- TEST-02 (SHOULD): Design testable APIs; inject time and rand
- TEST-03 (SHOULD): Stub external deps through consumer interfaces
- TEST-04 (SHOULD): Share helpers/fixtures outside production packages
- TEST-05 (SHOULD): Isolate integration tests with build tags
- TEST-06 (SHOULD): Call t.Helper() first in test helpers
- TEST-07 (SHOULD): Keep one assertion stack per package; match sibling tests
- TEST-08 (MUST): Prefix every \*_test.go filename with the source stem under test

## Anti-Patterns (AP)

- AP-01 (SHOULD): Discard errors with `_` without a justifying comment
- AP-02 (SHOULD): Store request-scoped context on structs
- AP-03 (SHOULD): Concatenate untrusted input into SQL
- AP-04 (SHOULD): Embed secrets in source or test data

## Test Anti-Patterns (TAP)

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
