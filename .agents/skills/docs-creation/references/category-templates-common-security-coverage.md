## security-coverage

Templates are guidance and review rubrics, not rigid prose generators.
Adapt structure and depth to the repository and context.
Replace all `<placeholder>` values with actual content.
Remove sections that do not apply rather than leaving them empty.

```markdown
# Security Coverage

Documents implemented, partial, and intentionally excluded
security controls within <project>.

Focus on:
- security boundaries
- ownership
- enforcement mechanisms
- intentional exclusions

Avoid:
- vague "secured" statements
- undocumented trust assumptions

## Coverage Matrix

| Control/Service | Status                         | Owner   | Notes   |
| --------------- | ------------------------------ | ------- | ------- |
| <control>       | <implemented/partial/excluded> | <owner> | <notes> |

## Security Boundaries

Document:
- trust boundaries
- privileged components
- externally exposed interfaces
- secret handling assumptions

## Known Gaps and Exclusions

Document:
- intentionally unsupported controls
- operational constraints
- accepted risks

## Validation and Enforcement

Document:
- policy enforcement tools
- scanning/linting
- audit expectations
- review requirements

## Legend

- implemented — control active and managed by this repository
- partial — control partially implemented; gaps noted
- excluded — intentionally out of scope; reason noted

## Decision Prompts

Consider:
- Which trust assumptions are implicit?
- Which controls rely on human process?
- Which failures could expose sensitive systems?
```
