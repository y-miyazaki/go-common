## improvements

Templates are guidance and review rubrics, not rigid prose generators.
Adapt structure and depth to the repository and context.
Replace all `<placeholder>` values with actual content.
Remove sections that do not apply rather than leaving them empty.

```markdown
# Improvements

Tracks planned, ongoing, and completed improvements for <project>.

Focus on:
- operational impact
- architectural improvements
- technical debt reduction
- risk reduction

Avoid:
- vague wishlist items
- low-context task lists

## Improvement Priorities

Document:
- current operational pain points
- architectural limitations
- maintenance burden
- scaling or security concerns

## Planned Improvements

| Title | Priority          | Impact   | Status                |
| ----- | ----------------- | -------- | --------------------- |
| <A>   | <high/medium/low> | <impact> | <planned/in-progress> |

Document:
- expected benefits
- migration risks
- dependencies
- rollout considerations

## Completed Improvements

| Date       | Improvement | Operational Outcome |
| ---------- | ----------- | ------------------- |
| YYYY-MM-DD | <change>    | <result>            |

## Deferred or Rejected Improvements (Optional)

| Proposal | Reason Deferred/Rejected |
| -------- | ------------------------ |
| <idea>   | <reason>                 |

## Decision Prompts

Consider:
- Which improvements reduce operational risk most?
- Which technical debt blocks future changes?
- Which improvements increase long-term maintainability?
```
