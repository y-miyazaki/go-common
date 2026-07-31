# Concurrency (CON)

**CON-01 (SHOULD): Ensure goroutines exit (watch context.Done / completion)**

Check: Do goroutines terminate properly and monitor context.Done()?
Why: Unterminated goroutines cause memory leaks, resource exhaustion, performance degradation
Fix: Clarify termination conditions, monitor context.Done(), use explicit completion signaling, verify with pprof

**CON-02 (SHOULD): Only the sender closes a channel (once)**

Check: Is channel close responsibility on the sender side?
Why: Receiver-side close, multiple closes, or forgotten close cause panics, goroutine leaks, deadlocks
Fix: Sender has close responsibility, prohibit receiver close, defer close, only once

**CON-03 (SHOULD): Define lock/completion ownership for shared state**

Check: Are synchronization boundaries and ownership rules clear and consistently applied?
Why: Unclear locking or completion ownership causes race conditions, deadlocks, and hidden coupling
Fix: Define lock ownership per shared state, avoid mixed synchronization models, and keep synchronization intent explicit
