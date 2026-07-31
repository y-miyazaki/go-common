# Context Handling (CTX)

**CTX-01 (MUST): Exported I/O APIs take context.Context as first param**

Check: Do public functions and methods accept context.Context as first argument?
Why: Missing context prevents timeout control, cancellation propagation, difficult testing
Fix: Add context.Context as first argument to all public APIs, unify ctx variable name

**CTX-02 (SHOULD): Do not store ambiguous request-scoped contexts**

Check: Is context origin and propagation path explicit across layer boundaries?
Why: Ambiguous context ownership causes cancellation gaps, timeout inconsistencies, and difficult incident analysis
Fix: Define context ownership per layer, propagate caller context, and document intentional boundary resets

**CTX-03 (SHOULD): Pass context into goroutines that do I/O or wait**

Check: Is context passed when launching goroutines?
Why: Missing context causes goroutine leaks, no cancellation propagation, resource exhaustion
Fix: Always pass context when launching goroutines, monitor context.Done()

**CTX-04 (SHOULD): Call cancel promptly; do not leak derived contexts**

Check: Is cancel from WithCancel/WithTimeout called with defer?
Why: Uncalled cancel causes resource leaks, goroutine leaks, memory growth
Fix: Mandatory defer cancel(), recommend defer even with WithTimeout
