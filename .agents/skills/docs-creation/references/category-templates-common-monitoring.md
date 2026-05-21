## monitoring

Templates are guidance and review rubrics, not rigid prose generators.
Adapt structure and depth to the repository and context.
Replace all `<placeholder>` values with actual content.
Remove sections that do not apply rather than leaving them empty.

```markdown
# Monitoring

Monitoring, alerting, diagnostics, and operational visibility guidance
for <project>.

Focus on:
- operational visibility
- actionable alerts
- ownership and escalation
- incident investigation support

Avoid:
- dashboards without operational purpose
- noisy non-actionable alerts

## Monitoring Strategy

Describe:
- critical signals
- SLO/SLA assumptions
- incident detection goals

## Alerts

| Alert | Trigger | Severity | Owner  | Runbook |
| ----- | ------- | -------- | ------ | ------- |
| <A>   | <cond>  | <level>  | <team> | <link>  |

Document:
- why alert matters
- expected operator action
- escalation expectations

## Dashboards

Describe:
- dashboard purpose
- intended audience
- operational usage

## Runbooks

### <Alert or Failure Scenario>

1. <investigation step>
2. <validation step>
3. <mitigation step>

## Decision Prompts

Consider:
- Which failures are hardest to detect?
- Which alerts require immediate action?
- Which metrics predict incidents early?
```
