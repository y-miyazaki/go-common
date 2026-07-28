## Error Handling (ERR)

**ERR-01 (MUST): Wrap errors with fmt.Errorf %w and context**

Check: Are errors wrapped with fmt.Errorf("%w", err) and context information included?
Why: Returning only error strings makes debugging difficult, lacks stack traces, root cause unclear
Fix: Wrap with fmt.Errorf("%w", err), add context information

**ERR-02 (SHOULD): Use distinct sentinel/custom errors per failure mode**

Check: Are sentinel errors defined with distinct semantics, not reused across unrelated failure modes? Are custom errors compatible with errors.Is/As?
Why: Reusing a sentinel error for unrelated conditions (e.g., using ErrNotFound for both "user not found" and "invalid input") defeats programmatic error handling and obscures root cause
Fix: Define distinct sentinel errors per failure category; ensure custom errors implement Unwrap for errors.Is/As compatibility

**ERR-03 (SHOULD): Panic only for fatal bugs; recover at boundaries**

Check: Are panics only for fatal errors and defer+recover implemented?
Why: Panic overuse and missing recover cause sudden application termination, data inconsistency
Fix: Panic only for fatal errors, implement defer+recover, return error for normal errors

**ERR-04 (SHOULD): Timeouts/retries and classify external errors**

Check: Are timeouts set, retries implemented, and errors classified?
Why: Missing timeouts and retries cause infinite waits, failure propagation
Fix: Set context timeout, exponential backoff, classify transient/permanent errors

**ERR-05 (SHOULD): Do not leak internals in user-facing error messages**

Check: Are there no internal implementation exposure, external stack trace disclosure, or SQL statement exposure?
Why: Internal information exposure causes information leakage, provides attack clues, security risks
Fix: Separate user-facing messages and internal logs, don't disclose detailed information

**ERR-06 (MUST): Never discard errors with _ unless commented**

Check: Is every returned `error` handled explicitly — no blank identifier assignment (`_ =` / `_, _ =`) unless a comment on the same line or immediately above justifies the discard?
Why: Silently discarding errors hides failures and makes production incidents hard to diagnose
Fix: Check and wrap errors; when discard is intentional, document why next to the assignment
