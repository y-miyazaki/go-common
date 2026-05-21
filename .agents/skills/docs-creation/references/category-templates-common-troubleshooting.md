## troubleshooting

Templates are guidance and review rubrics, not rigid prose generators.
Adapt structure and depth to the repository and context.
Replace all `<placeholder>` values with actual content.
Remove sections that do not apply rather than leaving them empty.

```markdown
# Troubleshooting

Common operational issues, investigation workflows,
and recovery guidance for <project>.

Focus on:
- symptom-driven diagnosis
- operational recovery
- root-cause guidance

Avoid:
- shallow "restart and retry" guidance
- undocumented assumptions

## <Symptom or Error Message>

### Symptoms

- <observable behavior>

### Likely Causes

- <cause>
- <cause>

### Investigation

```sh
<commands or investigation steps>
```

### Resolution

```sh
<resolution steps>
```

### Prevention (Optional)

<How to reduce recurrence>

### Escalation Notes (Optional)

<When to escalate or involve another team>

## Decision Prompts

Consider:
- Which failures are hardest to diagnose?
- Which symptoms indicate systemic issues?
- Which recovery actions are destructive?
```
