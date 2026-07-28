## Code Standards (CODE)

**CODE-01 (MUST): Keep interfaces small (1-3 methods) on the consumer side**

Check: Are interfaces kept small (1-3 methods) and defined on the consumer side?
Why: Too many methods and provider-side definitions make mocking difficult, increase test burden, reduce flexibility
Fix: Split interfaces with 5+ methods into focused roles, define interfaces where they are consumed

**CODE-02 (SHOULD): Minimize exported API surface; hide internals with internal/**

Check: Are there no excessive exports, unclear package name responsibilities, or unused internal/?
Why: Excessive exports increase API surface area, make maintenance difficult, increase breaking change risk
Fix: Minimize public APIs, express responsibility in package names, hide internal implementation with internal/

**CODE-03 (SHOULD): Unexport invariant fields and mutexes; split oversized structs**

Check: Are fields with invariants unexported and protected by methods? Are mutexes unexported? Are structs with 20+ fields split?
Why: Exposing fields that maintain invariants breaks encapsulation and causes race conditions; exported mutexes leak synchronization details
Fix: Unexport fields that enforce invariants and provide accessor methods. Keep exported fields for DTOs, config structs, and serialization targets. Split large structs by responsibility
