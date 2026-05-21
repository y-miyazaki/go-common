## performance

Templates are guidance and review rubrics, not rigid prose generators.
Adapt structure and depth to the repository and context.
Replace all `<placeholder>` values with actual content.
Remove sections that do not apply rather than leaving them empty.

```markdown
# Performance

Performance expectations, bottlenecks, scaling behavior,
and optimization guidance for <project>.

Focus on:
- operational bottlenecks
- scaling boundaries
- latency/throughput constraints
- resource-sensitive behavior

Avoid:
- micro-optimizations without operational relevance
- benchmark numbers without context

## Performance Goals

Document:
- latency expectations
- throughput expectations
- scaling assumptions
- resource constraints

## Known Bottlenecks

### <Bottleneck>

#### Impact

<Operational or user impact>

#### Cause

<Root cause or limitation>

#### Mitigation

<Current or planned mitigation>

## Resource Characteristics

Document:
- memory-sensitive operations
- CPU-intensive paths
- I/O bottlenecks
- concurrency limitations

## Tuning Guidance

| Parameter | Impact   | Tradeoff   |
| --------- | -------- | ---------- |
| <param>   | <effect> | <tradeoff> |

## Decision Prompts

Consider:
- Which workloads scale poorly?
- Which operations are latency-sensitive?
- Which optimizations increase operational complexity?
```
