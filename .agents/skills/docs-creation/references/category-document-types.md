# Document Types

Use this mapping to resolve `document_type` to the default target file, template, and determine whether the document is mandatory.
If a user provides `target_file`, prioritize that explicit path after schema validation.

## Core Types

| document_type      | Default file               | Template file                                        | Purpose                                           | Mandatory |
| ------------------ | -------------------------- | ---------------------------------------------------- | ------------------------------------------------- | --------- |
| `readme`           | `README.md`                | `category-templates-common-readme.md`                | Repository entry point for humans and AI agents   | yes       |
| `specification`    | `docs/specification.md`    | `category-templates-common-specification.md`         | Record behavioral requirements and expected flows | no        |
| `architecture`     | `docs/architecture.md`     | `category-templates-common-architecture.md`          | Explain system structure and boundaries           | no        |
| `design`           | `docs/design.md`           | `category-templates-common-design.md`                | Describe module-level implementation design       | no        |
| `design-decisions` | `docs/design-decisions.md` | `category-templates-common-design-decisions.md`      | Track major decisions and rejected alternatives   | no        |
| `troubleshooting`  | `docs/troubleshooting.md`  | `category-templates-common-troubleshooting.md`       | Provide issue diagnostics and recovery steps      | no        |
| `general`          | (ask user)                 | `category-templates-common-general.md`               | Capture documentation outside predefined types    | no        |

## Extension Types

| document_type       | Default file                | Template file                                        | Purpose                                             | Mandatory |
| ------------------- | --------------------------- | ---------------------------------------------------- | --------------------------------------------------- | --------- |
| `module-catalog`    | `docs/module-catalog.md`    | `category-templates-common-module-catalog.md`        | Catalog modules with key inputs/outputs             | no        |
| `monitoring`        | `docs/monitoring.md`        | `category-templates-common-monitoring.md`            | Define alerts, dashboards, and operational runbook  | no        |
| `performance`       | `docs/performance.md`       | `category-templates-common-performance.md`           | Record bottlenecks, benchmarks, and tuning actions  | no        |
| `security-coverage` | `docs/security-coverage.md` | `category-templates-common-security-coverage.md`     | Summarize security control and service coverage     | no        |
| `maintenance-notes` | `docs/maintenance-notes.md` | `category-templates-common-maintenance-notes.md`     | Capture periodic operations and maintenance history | no        |
| `improvements`      | `docs/improvements.md`      | `category-templates-common-improvements.md`          | Track planned and completed improvement work        | no        |
