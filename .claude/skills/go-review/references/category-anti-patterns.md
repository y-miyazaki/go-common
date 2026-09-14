# Anti-Patterns (AP)

**AP-01 (SHOULD): Discard errors with `_` without a justifying comment**

Check: Are returned errors discarded with `_` and no comment (ERR-06)?
Why: Silent discards hide failed I/O and security-sensitive failures
Fix: Handle the error or assign to `_` with a comment that states why ignore is safe

**AP-02 (SHOULD): Store request-scoped context on structs**

Check: Is `context.Context` stored on a struct as request-scoped state (CTX-02)?
Why: Stale cancelation and values leak across requests
Fix: Pass `ctx` as the first argument; do not persist ambiguous request contexts

**AP-03 (SHOULD): Concatenate untrusted input into SQL**

Check: Is SQL built with string concatenation of untrusted input (SEC-01)?
Why: Concatenation enables injection
Fix: Use parameterized queries or a query builder with placeholders

**AP-04 (SHOULD): Embed secrets in source or test data**

Check: Are API keys, passwords, or tokens present in source, logs, or fixtures (G-01)?
Why: Secrets in the tree leak via review, CI logs, and clones
Fix: Use placeholders in examples; load real secrets from the environment or a secret manager
