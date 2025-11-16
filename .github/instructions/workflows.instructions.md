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

- inputs, secrets, permissions は A-Z 順に整理
- job 名、step 名は一貫性を持たせる
- 再利用可能ワークフローは汎用的に設計し、特定プロジェクト依存を避ける
- 再利用可能ワークフローでは環境変数(例: `OUTDIR`, `ARTIFACT_NAME`, `ARTIFACT_BASE`)の定義は一箇所に集約してください。通常は `Setup Parameters` のような専用 step で `echo "NAME=value" >> $GITHUB_ENV` を行い、以降の step はその環境変数を直接参照するだけにします。
- 他の step で同じ変数を再定義したり、冗長にフォールバック式（例: `OUTDIR="${OUTDIR:-${{ inputs.output_dir }}}"`）を多用するのは避けてください。ドキュメントや repository の `workflows.instructions.md` で明示的に『一箇所で定義し、再定義は行わない』方針を示すべきです。
- 例外: 呼び出し元が明示的に上書きできる設計にする場合は、`Setup Parameters` で入力（`inputs`）を優先するか、`setup-parameters` step に `id` を付けて `outputs` で上書き値を渡す方法を採用してください。
  理由: `GITHUB_ENV` に書き出すことで別 step が同じ値を確実に読み取れるようになります。逆に同じ変数を複数箇所で定義すると可読性が低下し、意図しない値の衝突やデバッグ困難さを招きます。

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
  COMPONENT: ${{ inputs.component }}
  ENVIRONMENT: ${{ inputs.environment }}
  GO_PROJECT_PATH: ${{ inputs.go_path }}
  GO_VERSION: "1.25.3" # Centralized version management
```

## Guidelines

### Tool Integration Patterns

### Reviewdog Integration (PR Comments)

```yaml
- name: Run linter with reviewdog (PR comments)
  if: github.event_name == 'pull_request' && github.actor != 'dependabot[bot]'
  uses: reviewdog/action-golangci-lint@f9bba13753278f6a73b27a56a3ffb1bfda90ed71 # v2.8.0
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
  uses: codecov/codecov-action@5a1091511ad55cbe89839c7260b706298ca349f7 # v5.5.1
  with:
    use_oidc: true
    files: ${{ env.GO_PROJECT_PATH }}/coverage/coverage.out
    token: ${{ inputs.codecov_token }} # Fallback for private repos
    fail_ci_if_error: false
```

### Artifact Management

```yaml
- name: Upload Coverage Artifact
  uses: actions/upload-artifact@330a01c490aca151604b8cf639adc76d48f6c5d4 # v5.0.0
  with:
  name: coverage-${{ env.COMPONENT }}-${{ env.ENVIRONMENT }}
    path: ${{ env.GO_PROJECT_PATH }}/coverage
```

### Maintenance Guidelines

#### Version Management

- GitHub Actions のバージョンは定期的に更新する
- uses で利用するバージョン指定は、@v1.2.3 のようなバージョン指定は避け、Commit SHA を使用する
  - Commit SHA だけだと可読性が低いため、コメントでバージョン情報を併記する
- 破壊的変更のある GitHub Actions 更新時は段階的にテストする

#### Performance Optimization

- キャッシュを活用してビルド時間を短縮する
- 不要な step は条件分岐で除外する
- 並列実行可能なジョブは`needs`で適切に制御する

## Testing and Validation

### Code Modification Guidelines

コード修正時は以下コマンドで一括検証する：

```bash
actionlint
```

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
  uses: actions/github-script@ed597411d8f924073f98dfc5c65a23a2325f34cd # v8
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
  echo "CI completed for component: ${{ env.COMPONENT }} in env: ${{ env.ENVIRONMENT }}" | tee -a $GITHUB_STEP_SUMMARY
    if [ "${{ job.status }}" != "success" ]; then
      echo "❌ One or more steps failed. See logs above for details." >> $GITHUB_STEP_SUMMARY
    fi
```

### PR Failure Comments

```yaml
- name: Post PR failure comment
  if: failure() && github.event_name == 'pull_request' && github.actor != 'dependabot[bot]'
  uses: actions/github-script@ed597411d8f924073f98dfc5c65a23a2325f34cd # v8.0.0
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
