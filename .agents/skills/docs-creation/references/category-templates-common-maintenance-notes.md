## maintenance-notes

Templates are guidance and review rubrics, not rigid prose generators.
Adapt structure and depth to the repository and context.
Replace all `<placeholder>` values with actual content.
Remove sections that do not apply rather than leaving them empty.

```markdown
# Maintenance Notes

Operational maintenance procedures, recurring tasks,
and long-term maintenance considerations for <project>.

Focus on:
- operational continuity
- recurring maintenance
- lifecycle management
- operational quirks

Avoid:
- temporary issue tracking
- undocumented operational assumptions

## Recurring Tasks

| Task | Frequency | Owner  | Notes   |
| ---- | --------- | ------ | ------- |
| <A>  | <freq>    | <team> | <notes> |

Document:
- prerequisites
- operational risks
- rollback expectations

## Known Operational Quirks

### <Issue>

Describe:
- trigger conditions
- operational impact
- workaround
- escalation guidance

## Lifecycle Considerations

Document:
- upgrade expectations
- deprecation concerns
- dependency maintenance
- operational freeze constraints

## Change History (Optional)

| Date       | Change   | Notes   |
| ---------- | -------- | ------- |
| YYYY-MM-DD | <change> | <notes> |

## Decision Prompts

Consider:
- Which tasks are operationally risky?
- Which dependencies require proactive maintenance?
- Which operational knowledge is tribal today?
```
