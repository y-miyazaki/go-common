## design-decisions

Templates are guidance and review rubrics, not rigid prose generators.
Adapt structure and depth to the repository and context.
Replace all `<placeholder>` values with actual content.
Remove sections that do not apply rather than leaving them empty.

```markdown
# Design Decisions

Records major technical and operational decisions to preserve rationale
and prevent repeated re-investigation.

Focus on:
- why a decision exists
- tradeoffs
- rejected alternatives
- long-term operational implications

Avoid:
- obvious decisions requiring no rationale
- implementation-only details without architectural impact

## <Decision Title>

### Context

<What problem or constraint existed?>

### Decision

<What was chosen?>

### Rationale

Explain:
- why this approach was selected
- operational or architectural benefits
- constraints influencing the decision

### Alternatives Rejected

| Alternative | Reason Rejected |
| ----------- | --------------- |
| <A>         | <reason>        |
| <B>         | <reason>        |

### Consequences

Document:
- tradeoffs
- future limitations
- migration implications
- operational impact

## Decision Prompts

Consider:
- What future flexibility was sacrificed?
- Which assumptions could become invalid later?
- What operational burden does this introduce?
- Which alternatives may become viable later?
```
