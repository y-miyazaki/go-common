# Global / Base (G)

**G-01 (SHOULD): No API keys/passwords/tokens in source**

Check: Are API keys, passwords, and tokens not embedded in source code?
Why: Embedded secrets cause security breaches, credential leakage, audit violations
Fix: Use environment variables or AWS Secrets Manager, remove constants

**G-02 (SHOULD): Keep init() free of I/O, panics, and heavy side effects**

Check: Does init() avoid panics, external I/O, and non-trivial side effects? Is it minimal and deterministic?
Why: Complex init() hides initialization failures, causes unpredictable startup order across packages, and makes unit testing difficult
Fix: Limit init() to simple variable assignments; move complex initialization to explicit constructors or main()

**G-03 (SHOULD): Prefer types whose zero value is usable**

Check: Are types designed so their zero value is a valid and useful state where possible? (Types requiring mandatory initialization such as DB clients or config structs are exempt when documented.)
Why: Types with unusable zero values require mandatory initialization guards and cause subtle nil-dereference bugs when forgotten
Fix: Design structs so the zero value represents a valid empty state (e.g., sync.Mutex zero value is an unlocked mutex); document when zero value is not valid

**G-04 (SHOULD): Copy slices/maps at API boundaries**

Check: Are slices and maps copied when accepting from or returning to external callers?
Why: Shared references to slices/maps allow external callers to mutate internal state, violating encapsulation and causing hard-to-reproduce data corruption
Fix: Copy incoming slices/maps before storing in structs; return copies rather than direct internal references to callers
