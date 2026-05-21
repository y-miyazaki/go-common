## design

Templates are guidance and review rubrics, not rigid prose generators.
Adapt structure and depth to the repository and context.
Replace all `<placeholder>` values with actual content.
Remove sections that do not apply rather than leaving them empty.

```markdown
# Design Document — <Project or Module Name>

This document describes implementation structure, ownership boundaries,
runtime behavior, and design constraints.

Focus on:
- component interactions
- ownership boundaries
- state management
- operational behavior
- lifecycle expectations

Avoid:
- documenting every internal detail equally
- duplicating source code structure without rationale

## Overview

<Describe purpose and design scope.>

## Component Design

### <Component Name>

#### Responsibilities

- <responsibility>
- <responsibility>

#### Dependencies

- <internal/external dependency>

#### Lifecycle

Document:
- initialization
- runtime behavior
- shutdown/cleanup
- failure handling

## State and Data Flow

Describe:
- state ownership
- mutation boundaries
- synchronization assumptions
- data lifecycle

## Interfaces and Contracts

Document:
- externally consumed interfaces
- compatibility guarantees
- validation rules
- failure semantics

## Design Policies

- <policy>
- <policy>

## Constraints and Tradeoffs

Document:
- known limitations
- intentional compromises
- operational tradeoffs

## Naming Conventions (Optional)

Explain naming patterns if operationally or architecturally important.

## Parameters (Optional)

| Parameter | Type     | Description   |
| --------- | -------- | ------------- |
| `<name>`  | `<type>` | <description> |

## Decision Prompts

Consider:
- Which components own mutable state?
- Which interfaces are externally consumed?
- Which failures are hardest to recover from?
- Which implementation details are intentionally hidden?
```
