# Function Design (FUNC)

**FUNC-01 (SHOULD): Split mixed-responsibility or multi-layer functions**

Check: Are there no multiple responsibilities or mixed business/infrastructure layers in single functions?
Why: Mixed responsibilities make testing difficult, prevent reuse, increase maintenance cost
Fix: Apply single responsibility principle, separate layers, extract helper functions

**FUNC-02 (SHOULD): Unify pointer vs value receivers; avoid large values**

Check: Are there no mixed pointer/value receivers or large value receivers?
Why: Mixed receivers cause copy costs, changes not reflected, reduced readability
Fix: Pointer receiver principle, unify receiver name to 1-2 characters

**FUNC-03 (SHOULD): Keep generic constraints minimal and locally scoped**

Check: Are type abstraction boundaries explicit and are generic constraints minimal but sufficient?
Why: Ambiguous type boundaries and over-generalized constraints reduce readability, obscure contracts, and increase maintenance cost
Fix: Express domain contracts with clear constraints, keep abstraction local to required scope, and avoid speculative generic layers
