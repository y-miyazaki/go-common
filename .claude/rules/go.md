---
paths:
  - "**/*.go"
---

# Go Development Instructions

## Scope

- Scope covers implementing, testing, and validating Go source code.
- When adding or changing behavior, add or update `*_test.go` in the same change (TEST-00).
- Handle every returned `error` explicitly; blank-identifier discard needs a justifying comment (ERR-06).
- Wrap errors with `fmt.Errorf("…: %w", err)` to preserve the causal chain (ERR-01).
- Exported functions that perform I/O or may block take `context.Context` as the first parameter (CTX-01).

## Standards

### Naming Conventions

| Component                 | Rule       | Example        |
| ------------------------- | ---------- | -------------- |
| Interface (single-method) | -er suffix | Reader, Closer |
| Interface (multi-method)  | role name  | UserRepository |

### Unexported Helper Placement

| Condition                             | Preferred placement              |
| ------------------------------------- | -------------------------------- |
| Single-struct-specific responsibility | Unexported method on that struct |
| Pure helper shared across types/files | Package-level free function      |

## Guidelines

### Architecture (ARCH)
- ARCH-01 (SHOULD): Separate handler/usecase/repository from infrastructure
- ARCH-02 (SHOULD): Inject deps via constructor interfaces (not package globals)
- ARCH-03 (SHOULD): Keep domain logic free of DB/HTTP/API concerns
- ARCH-04 (SHOULD): Avoid circular deps; use standard layout and internal/
- ARCH-05 (SHOULD): Abstract external integrations behind consumer interfaces

### Code Standards (CODE)
- CODE-01 (MUST): Keep interfaces small (1-3 methods) on the consumer side
- CODE-02 (SHOULD): Minimize exported API surface; hide internals with internal/
- CODE-03 (SHOULD): Unexport invariant fields and mutexes; split oversized structs

### Concurrency (CON)
- CON-01 (SHOULD): Ensure goroutines exit (watch context.Done / completion)
- CON-02 (SHOULD): Only the sender closes a channel (once)
- CON-03 (SHOULD): Define lock/completion ownership for shared state

### Context Handling (CTX)
- CTX-01 (MUST): Exported I/O APIs take context.Context as first param
- CTX-02 (SHOULD): Do not store ambiguous request-scoped contexts
- CTX-03 (SHOULD): Pass context into goroutines that do I/O or wait
- CTX-04 (SHOULD): Call cancel promptly; do not leak derived contexts

### Dependencies (DEP)
- DEP-01 (SHOULD): List direct deps in go.mod with pinned versions

### Documentation (DOC)
- DOC-01 (SHOULD): Package has a doc comment stating purpose
- DOC-02 (MUST): Public APIs have godoc covering args/returns/errors
- DOC-03 (SHOULD): Comments are consistently English

### Error Handling (ERR)
- ERR-01 (MUST): Wrap errors with fmt.Errorf %w and context
- ERR-02 (SHOULD): Use distinct sentinel/custom errors per failure mode
- ERR-03 (SHOULD): Panic only for fatal bugs; recover at boundaries
- ERR-04 (SHOULD): Timeouts/retries and classify external errors
- ERR-05 (SHOULD): Do not leak internals in user-facing error messages
- ERR-06 (MUST): Never discard errors with _ unless commented

### Function Design (FUNC)
- FUNC-01 (SHOULD): Split mixed-responsibility or multi-layer functions
- FUNC-02 (SHOULD): Unify pointer vs value receivers; avoid large values
- FUNC-03 (SHOULD): Keep generic constraints minimal and locally scoped

### Global / Base (G)
- G-01 (SHOULD): No API keys/passwords/tokens in source
- G-02 (SHOULD): Keep init() free of I/O, panics, and heavy side effects
- G-03 (SHOULD): Prefer types whose zero value is usable
- G-04 (SHOULD): Copy slices/maps at API boundaries

### Security (SEC)
- SEC-01 (SHOULD): Validate inputs; ban string-concat SQL
- SEC-02 (SHOULD): Escape/sanitize outputs for HTML/JSON/CRLF sinks
- SEC-03 (SHOULD): Authenticate endpoints; verify JWT; enforce RBAC
- SEC-04 (SHOULD): Mask passwords/tokens in logs
- SEC-05 (SHOULD): Least privilege; no production debug; explicit CORS

### Testing (TEST)
- TEST-01 (SHOULD): Prefer table-driven tests with subtests and edges
- TEST-02 (SHOULD): Use assert vs require correctly; inject time/rand
- TEST-03 (SHOULD): Mock external deps through consumer interfaces
- TEST-04 (SHOULD): Share helpers/fixtures outside production packages
- TEST-05 (SHOULD): Isolate integration tests with build tags
- TEST-06 (SHOULD): Call t.Helper() first in test helpers

### Code Modification Guidelines

- When adding or changing behavior, add or update *_test.go files in the same change.


## Testing and Validation

Operational notes:

- Target at least 80% test coverage when verifying behavior changes.
- Separate integration tests with `//go:build integration`.

On-demand validation: see go-validation skill SKILL.md.

## Security Guidelines

- Do not hardcode secrets or credentials in source code, logs, or test data.
- Validate external inputs and use explicit error handling for privilege-boundary operations.
- Wrap error outputs to avoid sensitive data leakage and log only the minimum required information.
