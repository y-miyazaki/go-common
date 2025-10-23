---
applyTo: "**/.github/workflows/*.yaml,**/.github/workflows/*.yml"
description: "AI Assistant Instructions for GitHub Actions Workflows"
---

# AI Assistant Instructions for GitHub Actions Workflows

**言語ポリシー**: ドキュメントは日本語、コード・コメントは英語。

このリポジトリは GitHub Actions を使用した CI/CD ワークフローを含みます。

| workflow files                 | Purpose / Description 明                     |
| ------------------------------ | -------------------------------------------- |
| ci-push-\*.yaml                | CI（言語テスト、lint、カバレッジ）           |
| ci-push-markdown.yaml          | Markdown ファイル専用 CI（markdownlint）     |
| reusable-cd-go-aws-ecr.yaml    | Go 用 ECR デプロイ用再利用可能ワークフロー   |
| reusable-cd-terraform-aws.yaml | Terraform 用デプロイ用再利用可能ワークフロー |
| reusable-ci-go-common.yaml     | Go 用再利用可能ワークフロー                  |
| reusable-ci-markdown.yaml      | Markdown 用再利用可能ワークフロー            |
| reusable-ci-terraform-aws.yaml | Terraform 用再利用可能ワークフロー           |

## Standards

### Workflow Design Patterns

### Reusable Workflow Architecture

このプロジェクトでは**再利用可能ワークフロー**パターンを採用しています：

```yaml
# Caller workflow pattern (ci-push-dev.yaml)
jobs:
  go-ci:
    uses: ./.github/workflows/reusable-ci-go-common.yaml
    with:
      environment: "dev"
      component: "go-common"
      go_path: "."
    permissions:
      id-token: write
      contents: read
      pull-requests: write
```

### Trigger Path Management

ファイル種別ごとに独立した CI トリガーを設定：

- **Go コード**: `pkg/**`, `example/**`, `*.go`, `go.mod`, `go.sum`
- **Markdown**: `**/*.md`, `.markdownlint.yaml`
- **Terraform**: `terraform/**`, `**/*.tf`, `**/*.tfvars`
- **Workflows**: `.github/workflows/**/*.yaml`

### Workflow Standards

### Input/Output Conventions

```yaml
# Reusable workflow inputs (standardized)
on:
  workflow_call:
    inputs:
      environment:
        description: "Target environment (dev, qa, stg, prd, etc.)"
        required: true
        type: string
      component:
        description: "Logical component name (e.g. go-common)"
        required: true
        type: string
      go_path: # Go-specific
        description: "Path to Go project root (containing go.mod)"
        required: true
        type: string
```

### Permission Management

#### OIDC Token Requirements

Codecov upload 等で OIDC を使用する場合：

```yaml
jobs:
  caller-job:
    permissions:
      id-token: write # OIDC token issuance for reusable workflows
      contents: read
      pull-requests: write
```

#### Minimal Permission Principle

```yaml
permissions:
  contents: read # Repository checkout
  pull-requests: write # PR comments (reviewdog)
  id-token: write # OIDC authentication (if needed)
  # Other permissions only when explicitly required
```

### Environment Variable Management

```yaml
env:
  ENV: ${{ inputs.environment }}
  COMPONENT: ${{ inputs.component }}
  GO_PROJECT_PATH: ${{ inputs.go_path }}
  GO_VERSION: "1.25.3" # Centralized version management
```

## Guidelines

### Tool Integration Patterns

### Reviewdog Integration (PR Comments)

```yaml
- name: Run linter with reviewdog (PR comments)
  if: github.event_name == 'pull_request' && github.actor != 'dependabot[bot]'
  uses: reviewdog/action-golangci-lint@v2
  with:
    github_token: ${{ secrets.GITHUB_TOKEN }}
    reporter: github-pr-review
    fail_level: error
    filter_mode: added
```

### Codecov Integration (OIDC)

```yaml
- name: Upload to Codecov
  if: ${{ (steps.repo_visibility.outputs.result != 'true') || (inputs.codecov_token != '') }}
  uses: codecov/codecov-action@v5
  with:
    use_oidc: true
    files: ${{ env.GO_PROJECT_PATH }}/coverage/coverage.out
    token: ${{ inputs.codecov_token }} # Fallback for private repos
    fail_ci_if_error: false
```

### Artifact Management

```yaml
- name: Upload Coverage Artifact
  uses: actions/upload-artifact@v4
  with:
    name: coverage-${{ env.COMPONENT }}-${{ env.ENV }}
    path: ${{ env.GO_PROJECT_PATH }}/coverage
```

### Maintenance Guidelines

#### Version Management

- アクションのバージョンは定期的に更新する
- `@v1`のようなメジャーバージョン指定は避け、`@v1.2.3`形式を使用
- 破壊的変更のあるアクション更新時は段階的にテストする

#### Performance Optimization

- キャッシュを活用してビルド時間を短縮する
- 不要な step は条件分岐で除外する
- 並列実行可能なジョブは`needs`で適切に制御する

## Testing and Validation

### Workflow Testing Guidelines

- 新しいワークフローを作成時は最低 1 回のテスト実行を行う
- reusable ワークフローの変更時は全ての caller workflow でテスト実行する
- 権限変更時は特に OIDC 関連の動作を確認する

### Validation Requirements

- すべてのワークフローが構文的に正しい YAML であること
- 必要な`permissions`が設定されていること
- `if`条件が適切に設定されていること（dependabot 除外等）

## Security Guidelines

### Secret Management

- **公開リポジトリ**: OIDC 認証を優先使用
- **プライベートリポジトリ**: 必要に応じて token を`secrets`で管理
- **Fork PR**: secrets は利用不可のため条件分岐で対処

### Dependabot Handling

```yaml
if: github.actor != 'dependabot[bot]' # Skip for dependabot PRs
```

### Repository Visibility Detection

```yaml
- name: Determine repository visibility
  id: repo_visibility
  uses: actions/github-script@v8
  with:
    script: |
      const { data } = await github.rest.repos.get({
        owner: context.repo.owner,
        repo: context.repo.repo,
      });
      return data.private ? 'true' : 'false';
```

### Error Handling and Reporting

### Job Failure Handling

```yaml
- name: Summary
  if: always()
  run: |
    echo "CI completed for component: ${{ env.COMPONENT }} in env: ${{ env.ENV }}" | tee -a $GITHUB_STEP_SUMMARY
    if [ "${{ job.status }}" != "success" ]; then
      echo "❌ One or more steps failed. See logs above for details." >> $GITHUB_STEP_SUMMARY
    fi
```

### PR Failure Comments

```yaml
- name: Post PR failure comment
  if: failure() && github.event_name == 'pull_request' && github.actor != 'dependabot[bot]'
  uses: actions/github-script@v8
  with:
    script: |
      const runUrl = `${context.serverUrl}/${context.repo.owner}/${context.repo.repo}/actions/runs/${context.runId}`;
      const body = `❌ CI failed for component: ${process.env.COMPONENT}\n\nRun page: ${runUrl}`;
      github.rest.issues.createComment({
        issue_number: context.issue.number,
        owner: context.repo.owner,
        repo: context.repo.repo,
        body
      });
```

## MCP Tools

**詳細な MCP Tools の設定・使用方法は `.github/copilot-instructions.md` を参照。**

### GitHub Actions 特有の活用パターン

**ワークフロー内での AWS CLI 実行例:**

```yaml
- name: AWS Resource Check
  run: |
    aws sts get-caller-identity
    aws s3 ls
  env:
    AWS_REGION: ${{ secrets.AWS_REGION }}
```

**ワークフロー設計時のベストプラクティス参照:**

```
# GitHub Actions 情報取得
resolve: "github-actions" → get-docs: topic="security best practices"

# CI/CD パフォーマンス最適化
resolve: "github-actions" → get-docs: topic="workflow optimization"
```
