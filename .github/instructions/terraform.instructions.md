---
applyTo: "**/*.tf,**/*.tfvars,**/*.hcl"
description: "AI Assistant Instructions for Terraform"
---

# AI Assistant Instructions for Terraform

**言語ポリシー**: ドキュメントは日本語、コード・コメントは英語。

このリポジトリは Terraform で AWS インフラを管理するためのプロジェクトです。

| Directory/File        | Purpose / Description                      |
| --------------------- | ------------------------------------------ |
| terraform/application | 基本セキュリティ設定・AWS アカウント初期化 |
| terraform/base        | 基本セキュリティ設定・AWS アカウント初期化 |
| terraform/monitor     | CloudWatch 監視・アラート・可観測性        |
| terraform/management  | 組織管理・コンプライアンス・ガバナンス     |
| lambda/               | サーバレス監視・自動化                     |
| modules/              | AWS リソース共通モジュール                 |

## Standards

### Coding Standards

### Naming Conventions

命名規則は以下を参照：
https://www.terraform-best-practices.com/naming

| Component      | Rule                | Example                                   |
| -------------- | ------------------- | ----------------------------------------- |
| モジュール変数 | snake_case          | resource_name, is_enabled, vpc_cidr_block |
| ファイル名     | snake_case          | main_security.tf, iam_roles.tf            |
| AWS リソース名 | kebab-case + prefix | ${var.name_prefix}example-role            |
| ローカル変数   | snake_case          | local_config, security_group_ids          |
| 出力変数       | snake_case          | vpc_id, subnet_ids, security_group_id     |
| タグ名         | PascalCase          | Environment, Project, Owner               |

### Terraform Standards

以下の内容は tflint,trivy で指摘される項目以外の内容を記載する。

- module の variables.tf
  - 変数コメントは先頭に(Optional) or (Required)を明記
  - validation チェック
    - 仕様変更でトラブルになる可能性考慮し、特定のパラメータ値を必須とするようなチェックは行わない
    - Required な変数で `length > 0` のようなチェックはしない

## Guidelines

### Documentation and Comments

- すべての関数・リソースは詳細な説明を含める
- 目的・機能は冒頭で明記する
- 複雑なロジックは詳細なコメントで説明する
- コメント・ドキュメントは英語で記載する
- 複雑なモジュールは使用例も記載する

### Error Handling

- Terraform エラーは詳細なログと共に適切に処理する
- plan/apply 実行時のエラーは必ず原因を特定してから修正する
- State lock エラー時は適切な解除手順を実行する
- リソース依存関係エラーは depends_on で明示的に解決する

## Testing and Validation

### Code Modification Guidelines

コード修正時は以下コマンドで一括検証する：

```bash
# Environment variable setup (example: dev environment)
ENV=dev

# Initialize, validate, and plan
terraform init -reconfigure -backend-config=terraform.${ENV}.tfbackend
terraform fmt --recursive && terraform validate
tflint -f compact --var-file=terraform.${ENV}.tfvars
trivy fs . --format table --config /workspace/trivy.yaml --secret-config /workspace/trivy-secret.yaml
terraform plan -lock=false -var-file=terraform.${ENV}.tfvars
```

### Validation Requirements

- すべての検証（fmt, validate, tflint, trivy, plan）に合格してからコード提出する
- セキュリティスキャン結果を確認し、問題があれば修正する

## Security Guidelines

### Terraform Specific Security

- IAM ポリシーは最小権限で設計する
- S3 等はデフォルトで暗号化を有効化する
- API Gateway 等は WAF・レート制限を設定する
- モジュール・依存関係は最新バージョンを使用する
- リソースベースポリシー・VPC セキュリティグループを適切に設定する

## MCP Tools

**詳細な MCP Tools の設定・使用方法は `.github/copilot-instructions.md` を参照。**

### Terraform 作業特有の活用パターン

**AWS リソース設計前の事前調査:**

```
# 既存リソースとの整合性確認
aws ec2 describe-vpcs --region us-east-1
aws rds describe-db-instances --region us-east-1

# Terraform Import 用の情報収集
aws s3api get-bucket-versioning --bucket terraform-state-bucket
```

**公式ドキュメントによるベストプラクティス確認:**

```
# AWS サービス固有の制約・要件確認
search: "S3 bucket policy terraform cross-account access"

# セキュリティ設定の公式ガイダンス参照
url: "https://docs.aws.amazon.com/s3/latest/userguide/bucket-policies.html"
```

**プロバイダー・モジュール情報の活用:**

```
# AWS Provider 最新機能確認
resolve: "terraform-aws-provider" → get-docs: topic="s3 bucket configuration"

# Terraform モジュール構成例
resolve: "terraform" → get-docs: topic="module composition"
```
