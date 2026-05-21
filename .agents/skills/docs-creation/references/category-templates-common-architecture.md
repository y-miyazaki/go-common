## architecture

Templates are guidance and review rubrics, not rigid prose generators.
Adapt structure and depth to the repository and context.
Replace all `<placeholder>` values with actual content.
Remove sections that do not apply rather than leaving them empty.

```markdown
# Architecture Overview

This document provides structural and operational context for <project>.

Focus on:
- system boundaries
- ownership relationships
- runtime topology
- dependency flow
- operational responsibilities

Avoid:
- low-level implementation details
- directory listings without explanation

## System Overview

Describe:
- major components
- interaction patterns
- trust boundaries
- deployment topology
- external dependencies

## Structural Layout

### <Component or Domain>

| Item | Responsibility |
| ---- | -------------- |
| <A>  | <role>         |
| <B>  | <role>         |

Document:
- ownership boundaries
- communication paths
- lifecycle relationships

## Runtime and Operational Topology

Describe:
- runtime dependencies
- scaling boundaries
- network or service relationships
- operational segregation

## Architectural Constraints

Document:
- intentional limitations
- coupling constraints
- compatibility requirements
- operational assumptions

## Key Design Decisions

Reference:
- ADRs
- design-decisions.md
- security or operational rationale

## Decision Prompts

Consider:
- Which components are tightly coupled?
- Which boundaries are operationally sensitive?
- What assumptions exist around scaling or failure isolation?
- Which dependencies are hardest to replace?
```
