[![ci-push-dev](https://github.com/y-miyazaki/go-common/actions/workflows/ci-push-dev.yaml/badge.svg)](https://github.com/y-miyazaki/go-common/actions/workflows/ci-push-dev.yaml)

<!-- omit in toc -->
# go-common

This repository provides common libraries and example applications for Go language, utilizing libraries such as AWS SDK v2, Gin, and GORM with practical samples.

<!-- omit in toc -->
## Table of Contents

- [Project Overview](#project-overview)
  - [Directory Structure](#directory-structure)
- [Installation](#installation)
- [Local Development Environment](#local-development-environment)
  - [Required](#required)
  - [Setting](#setting)
  - [Create Local Development Environment](#create-local-development-environment)
- [Commands](#commands)
  - [Build and Test](#build-and-test)
  - [Code Quality Check](#code-quality-check)
  - [Development Support](#development-support)
- [Troubleshooting](#troubleshooting)
  - [Common Issues](#common-issues)
  - [Getting Help](#getting-help)
- [License](#license)
- [Note](#note)

## Project Overview

This project provides common libraries and practical example applications to learn Go language best practices. The main technology stack includes:

- Go 1.24
- AWS SDK v2
- Gin (Web Framework)
- GORM (ORM)
- golangci-lint (Lint Tool)
- Delve (Debugger)

### Directory Structure

| Path                  | Description                   |
| --------------------- | ----------------------------- |
| .github/              | GitHub related files          |
| .github/instructions/ | Copilot instruction files     |
| .github/workflows/    | GitHub Actions CI/CD          |
| .vscode/              | VS Code settings              |
| env/                  | Environment related files     |
| env/common/           | Common environment settings   |
| env/common/scripts/   | Common scripts (init.sh etc.) |
| example/              | Example applications          |
| example/gin1/         | Sample using Gin              |
| example/...           | Other examples                |
| pkg/                  | Common libraries              |
| pkg/aws/              | AWS related utilities         |
| pkg/db/               | Database related              |
| pkg/...               | Other common modules          |
| go.mod                | Go module definition          |
| go.sum                | Go dependencies               |
| README.md             | This file                     |

## Installation

To use this library in your Go project, add it as a dependency:

```bash
go get github.com/y-miyazaki/go-common
```

Then, import the required packages in your code:

```go
import (
    "github.com/y-miyazaki/go-common/pkg/aws"  // AWS utilities
    "github.com/y-miyazaki/go-common/pkg/db"   // Database utilities
    // Add other packages as needed
)
```

After adding the dependency, run:

```bash
go mod tidy
```

This will download and install all necessary dependencies.

## Local Development Environment

### Required

- Go 1.24 or higher
- Git
- GitHub CLI (gh)
- Docker (when using devcontainer)
- Node.js (optional, when using ESLint)

### Setting

1. Clone the repository:
   ```bash
   git clone https://github.com/y-miyazaki/go-common.git
   cd go-common
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up authentication with GitHub CLI:
   ```bash
   gh auth login
   ```

### Create Local Development Environment

This project supports devcontainer. Open in VS Code and select "Reopen in Container" to build an environment with pre-installed tools.

The devcontainer includes the following tools:
- Go and common utilities
- Git (built from source)
- Node.js, npm, ESLint
- Docker CLI
- GitHub CLI

## Commands

### Build and Test
```bash
# Batch verification (recommended)
bash /workspace/scripts/go/check.sh

# Specific directory only
bash /workspace/scripts/go/check.sh -f ./example/gin1/

# Individual execution (when necessary)
go build ./...
go test ./...
go test -cover ./...
```

### Code Quality Check
```bash
# Batch verification (recommended)
bash /workspace/scripts/go/check.sh

# Individual execution (when necessary)
go mod tidy
go fmt ./...
go vet ./...
golangci-lint run
govulncheck ./...
```

### Development Support
```bash
# Module organization
go mod tidy

# Update dependencies
go get -u ./...

# Debug execution (Delve)
dlv debug ./example/gin1
```

## Troubleshooting

### Common Issues

- **Git push returns 403 error**  
  Check if GitHub CLI authentication is set correctly. Use `gh auth status` to verify status, and re-run `gh auth login` if necessary.
- **devcontainer startup failure**  
  Check execution permissions of `env/common/scripts/init.sh` and verify `/bin/sh` compatibility issues. Run with `bash` if needed.
- **Dependency errors**  
  Run `go mod tidy` and ensure Go version is 1.24 or higher.
- **Lint errors**  
  Check output of `golangci-lint run` and fix pointed locations. Refer to configuration file `.golangci.yml`.

### Getting Help

- **Documentation**  
  Refer to this README and documentation in each directory.
- **Issues**  
  Report issues via GitHub Issues.
- **Contribution**  
  Refer to CONTRIBUTING.md and create pull requests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Note

- This project is under development, so APIs or structures may change.
- Security note: Use environment variables or secret management tools for sensitive information, avoid hardcoding.
- Utilize MCP Tools for support in AWS documentation search and context management.
