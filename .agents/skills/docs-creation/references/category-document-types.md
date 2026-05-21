# `document_type` Mapping

Use this mapping to resolve `document_type` to the default target file under `docs/`.
If a user provides `target_file`, prioritize that explicit path after schema validation.

## Core Types

| document_type      | Default file               | Purpose                                           | Required (minimum sections)                                                  |
| ------------------ | -------------------------- | ------------------------------------------------- | ---------------------------------------------------------------------------- |
| `specification`    | `docs/specification.md`    | Record behavioral requirements and expected flows | H1, feature/behavior sections, configuration defaults (if present)           |
| `architecture`     | `docs/architecture.md`     | Explain system structure and boundaries           | H1, System/Component Structure section, Project Layout, Key Design Decisions |
| `design`           | `docs/design.md`           | Describe module-level implementation design       | H1, Overview, Module Architecture, Parameters                                |
| `design-decisions` | `docs/design-decisions.md` | Track major decisions and rejected alternatives   | H1, Decision, Rationale, Alternatives Rejected                               |
| `troubleshooting`  | `docs/troubleshooting.md`  | Provide issue diagnostics and recovery steps      | H1, Symptoms, Root Cause, Resolution, Prevention (if applicable)             |
| `general`          | no fixed file (ask user)   | Capture documentation outside predefined types    | H1 and purpose-aligned sections based on selected template                   |

## Extension Types

| document_type       | Default file                | Purpose                                             | Required (minimum sections)                                 |
| ------------------- | --------------------------- | --------------------------------------------------- | ----------------------------------------------------------- |
| `module-catalog`    | `docs/module-catalog.md`    | Catalog modules with key inputs/outputs             | H1, at least one module entry with Purpose, Inputs, Outputs |
| `monitoring`        | `docs/monitoring.md`        | Define alerts, dashboards, and operational runbook  | H1, Alerts, Dashboards, Runbooks                            |
| `performance`       | `docs/performance.md`       | Record bottlenecks, benchmarks, and tuning actions  | H1, Benchmarks, Known Bottlenecks, Tuning                   |
| `security-coverage` | `docs/security-coverage.md` | Summarize security control and service coverage     | H1, Coverage Matrix, Legend                                 |
| `maintenance-notes` | `docs/maintenance-notes.md` | Capture periodic operations and maintenance history | H1, Periodic Tasks, Known Quirks, Change Log                |
| `improvements`      | `docs/improvements.md`      | Track planned and completed improvement work        | H1, Backlog/Planned Items, Completed Items, Next Actions    |
