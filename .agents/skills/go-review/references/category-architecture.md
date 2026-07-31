# Architecture (ARCH)

**ARCH-01 (SHOULD): Separate handler/usecase/repository from infrastructure**

Check: Are handler/usecase/repository separated and business/infrastructure layers separated?
Why: Mixed business logic and infrastructure layers make testing difficult, technology stack changes difficult
Fix: Apply Clean Architecture, separate handler/usecase/repository

**ARCH-02 (SHOULD): Inject deps via constructor interfaces (not package globals)**

Check: Are dependencies passed via constructor arguments as interfaces rather than accessed as global variables? (Package-level variables are acceptable in single-binary Lambda/CLI tools when interfaces are defined and tests use save/restore helpers.)
Why: Global variable dependencies and hardcoded dependencies prevent mocking, parallel testing
Fix: Accept dependencies as interface arguments in constructors; use wire/dig only when constructor graphs become complex

**ARCH-03 (SHOULD): Keep domain logic free of DB/HTTP/API concerns**

Check: Is business logic free from infrastructure concerns (DB, HTTP, external APIs)?
Why: Scattered business logic mixed with infrastructure makes testing difficult and technology changes expensive
Fix: Keep domain logic in pure Go types and functions; access infrastructure through interfaces defined in the domain layer

**ARCH-04 (SHOULD): Avoid circular deps; use standard layout and internal/**

Check: Are there no circular dependencies, standard layout compliance, and internal/ utilization?
Why: Circular dependencies and package bloat make builds difficult, understanding difficult
Fix: Control dependency direction, comply with standard layout, utilize internal/

**ARCH-05 (SHOULD): Abstract external integrations behind consumer interfaces**

Check: Are adapter patterns, interface definitions, and abstraction layers implemented?
Why: Direct external API calls and no abstraction layer cause vendor lock-in, difficult testing
Fix: Adapter pattern, define interfaces, implement abstraction layer
