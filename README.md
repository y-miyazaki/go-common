![Go](https://custom-icon-badges.herokuapp.com/badge/Go-00ADD8.svg?logo=Go&logoColor=white)
![Apache-2.0](https://custom-icon-badges.herokuapp.com/badge/license-Apache%202.0-8BB80A.svg?logo=law&logoColor=white)
[![Go Report Card](https://goreportcard.com/badge/github.com/y-miyazaki/go-common)](https://goreportcard.com/report/github.com/y-miyazaki/go-common)
[![GitHub release](https://img.shields.io/github/release/y-miyazaki/go-common.svg)](https://github.com/y-miyazaki/go-common/releases/latest)
[![ci-push-dev](https://github.com/y-miyazaki/go-common/actions/workflows/ci-push-dev.yaml/badge.svg)](https://github.com/y-miyazaki/go-common/actions/workflows/ci-push-dev.yaml)
[![Codecov](https://codecov.io/gh/y-miyazaki/go-common/branch/develop/graph/badge.svg)](https://codecov.io/gh/y-miyazaki/go-common)

<!-- omit in toc -->
# go-common

This repository provides common libraries and example applications for Go language, utilizing libraries such as AWS SDK v2, Gin, and GORM with practical samples.

<!-- omit in toc -->
## Table of Contents

- [Project Overview](#project-overview)
  - [Directory Structure](#directory-structure)
- [Installation](#installation)
- [Quick code example](#quick-code-example)
- [Development](#development)
  - [Local development environment](#local-development-environment)
    - [Required](#required)
    - [Setup](#setup)
  - [Local development environment (devcontainer)](#local-development-environment-devcontainer)
    - [Required](#required-1)
    - [Setup](#setup-1)
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
| LICENSE               | License file                  |
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

## Quick code example

Here is a tiny example showing how to use `pkg/aws` and `pkg/db` helpers (adjust imports and usage to your needs):

```go
package main

import (
  "context"
  "fmt"
  "github.com/y-miyazaki/go-common/pkg/aws"
  "github.com/y-miyazaki/go-common/pkg/db"
)

func main() {
  ctx := context.Background()
  // Example: initialize AWS client helper
  awsCfg, err := aws.NewConfig(ctx)
  if err != nil {
    panic(err)
  }
  fmt.Println("AWS region:", awsCfg.Region)

  // Example: open database connection
  dsn := "user:pass@tcp(localhost:3306)/example"
  conn, err := db.Open(ctx, dsn)
  if err != nil {
    panic(err)
  }
  defer conn.Close()
  fmt.Println("DB connected")
}
```

## Development

This section describes how to set up a local development environment and run example applications.

### Local development environment

This section explains how to set up a local, non-containerized development workspace on your machine. It includes quick, platform-agnostic steps to clone the repository, install dependencies, and run example applications for development and testing.

#### Required

- Go 1.24 or higher  
- Git  
- GitHub CLI (gh)  

#### Setup

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

### Local development environment (devcontainer)

This section describes how to use a VS Code devcontainer to create a reproducible development environment. The devcontainer provides consistent tool versions and preconfigured mounts so contributors can work with the same development setup.

#### Required

- Visual Studio Code (VSCode)  
- Docker (when using devcontainer)  

#### Setup

1. Clone the repository (if not already done):  

    ```bash
    git clone https://github.com/y-miyazaki/go-common.git
    cd go-common
    ```

2. Create devcontainer  

    ```bash
    mkdir -p .devcontainer
    mkdir -p env/common/tmp/gh
    touch env/common/tmp/.gitconfig
    cp -p env/example/.devcontainer/devcontainer.json .devcontainer/devcontainer.json
    ```

3. Adjust devcontainer.json  
    The following excerpt is a locally mounted configuration; update mount paths to match your environment. Replace `${env:HOME}/workspace/go-common` with your actual local path to this repository (e.g., `/Users/yourname/projects/go-common` on macOS or `C:\Users\yourname\projects\go-common` on Windows).

    ```bash
    cat .devcontainer/devcontainer.json
    ```

    ```json
    {
        "runArgs": [
            "-v",
            "${env:HOME}/workspace/go-common:/workspace",
            "-v",
            "${env:HOME}/workspace/go-common/env/common/.bashrc:/home/vscode/.bashrc",
            "-v",
            "${env:HOME}/workspace/go-common/env/common/tmp/.gitconfig:/home/vscode/.gitconfig",
            "-v",
            "${env:HOME}/workspace/go-common/env/common/tmp/.aws:/home/vscode/.aws",
            "-v",
            "${env:HOME}/workspace/go-common/env/common/gh:/home/vscode/.config/gh",
            "-v",
            "/var/run/docker.sock:/var/run/docker.sock"
        ]
    }
    ```

4. Configure .gitconfig  
    The following excerpt is a locally mounted file; update values for your name and email.
    ```bash
    cat env/common/tmp/.gitconfig
    ```

    ```gitconfig:env/common/tmp/.gitconfig
    [user]
        name = Your Name
        email = your.email@example.com
    [init]
        defaultBranch = main
    [credential]
        helper = !gh auth git-credential
    [safe]
        directory = /workspace
    ```

5. Launch devcontainer from Visual Studio Code  
    Open the command palette with `F1` or `Ctrl+Shift+P`, then run `Dev Containers: Open folder in Container` (or `Dev containers: Reopen in Container` / `Dev Containers: Rebuild Container`).  

6. Download Go modules:  
    ```bash
    go mod download
    ```

7. GitHub CLI login  
    After launching the devcontainer, run the following to log in to GitHub CLI. Skip if already logged in.

    ```bash
    gh auth login
    ```

8. Run the Gin example (example/gin1):  
    ```bash
    cd example/gin1
    go run ./...
    # or from repo root
    # bash ./scripts/go/check.sh -f ./example/gin1/
    ```

9. Test the server (default port 8080):  
    ```bash
    curl http://localhost:8080/health
    ```

### Commands

#### Build and Test
```bash
# Batch verification (recommended)
bash ./scripts/go/check.sh

# Specific directory only
bash ./scripts/go/check.sh -f ./example/gin1/

# Individual execution (when necessary)
go build ./...
go test ./...
go test -cover ./...
```

#### Code Quality Check
```bash
# Batch verification (recommended)
bash ./scripts/go/check.sh

# Individual execution (when necessary)
go mod tidy
go fmt ./...
go vet ./...
golangci-lint run
govulncheck ./...
```

#### Development Support
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

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Note

- This project is under development, so APIs or structures may change.
- Security note: Use environment variables or secret management tools for sensitive information, avoid hardcoding.
- Utilize MCP Tools for support in AWS documentation search and context management.
