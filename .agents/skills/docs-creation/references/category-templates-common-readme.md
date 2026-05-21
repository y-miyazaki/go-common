## readme

Templates are guidance and review rubrics, not rigid prose generators.
Adapt structure and depth to the repository and context.
Replace all `<placeholder>` values with actual content.
Remove sections that do not apply rather than leaving them empty.

````markdown
# <Repository Name>

<One to three sentences describing the repository purpose and intended operational scope.>

README is the repository entry point for humans and AI assistants.

Focus on:
- repository identification
- quick operational entry points
- high-level navigation
- contributor guidance
- links to deeper documentation

Deprioritize:
- deep architectural explanations
- exhaustive configuration references
- large troubleshooting sections
- implementation-specific details
- duplicated documentation

Avoid:
- duplicating `docs/`
- embedding large generated tables
- storing operational runbooks
- documenting every module/package/resource
- excessive badge-only sections without operational value

## Status

| Item              | Value                            |
| ----------------- | -------------------------------- |
| Lifecycle         | <experimental/stable/deprecated> |
| Primary Stack     | <Go/Terraform/etc>               |
| Deployment Target | <AWS/Kubernetes/etc>             |

Document:
- operational maturity
- intended usage scope
- important limitations if applicable

## Repository Purpose

Describe:
- what this repository manages
- who operates or consumes it
- intended operational boundaries
- primary use cases

Avoid:
- implementation-level explanations
- detailed architecture descriptions

## Quick Start

### Prerequisites

- <required tool>
- <required tool>

### Common Commands

```sh
<setup/build/test/deploy commands>
```

Focus on:
- common contributor/operator workflows
- minimal onboarding friction

## Repository Structure

| Path     | Purpose   |
| -------- | --------- |
| `<path>` | <purpose> |
| `<path>` | <purpose> |

Document:
- major ownership boundaries
- important entrypoint directories
- where to find deeper documentation

Avoid:
- exhaustive directory trees
- documenting every internal directory

## Documentation Index

| Document                   | Purpose                     |
| -------------------------- | --------------------------- |
| `docs/index.md`            | Documentation catalog       |
| `docs/architecture.md`     | System architecture         |
| `docs/design-decisions.md` | Architectural rationale     |
| `docs/troubleshooting.md`  | Operational troubleshooting |

Focus on:
- discoverability
- navigation
- operational relevance

## Operational Notes (Optional)

Document only:
- critical operational assumptions
- deployment constraints
- environment limitations
- important warnings

Move detailed operational procedures into:
- runbooks
- troubleshooting documents
- architecture/design documents

## Contribution Guidance

Reference:
- `CONTRIBUTING.md`
- `AGENTS.md`
- `docs/index.md`

Document:
- contribution entry points
- repository conventions
- important review expectations

## Related Repositories or Services (Optional)

| Repository/Service | Purpose        |
| ------------------ | -------------- |
| `<name>`           | <relationship> |

## Decision Prompts

Consider:
- Can a new contributor understand repository purpose quickly?
- Are common workflows discoverable within minutes?
- Are deep details delegated to specialized docs?
- Is documentation navigation obvious?
- Are operationally critical constraints visible?
- Does README avoid becoming a duplicate of `docs/`?
````
