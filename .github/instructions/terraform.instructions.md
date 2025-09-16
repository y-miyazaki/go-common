---
applyTo: "**/*.tf,**/*.tfvars,**/*.tfstate,**/*.tfbackend"
---

# GitHub Copilot Instructions for Terraform

**Language Note**: This document is written in Japanese, but all generated code and comments must be in English.

## Overview

このリポジトリは AWS のセキュリティ・監視インフラを管理する Terraform コードを含みます。

| ディレクトリ/ファイル | 役割・説明                                 |
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

| コンポーネント | 規則                | 例                                        |
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

コード修正時は以下を順次実行する：

```bash
# 環境変数設定（例：dev環境）
ENV=dev

# 初期化・検証・プラン
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

**詳細な security guidelines は `.github/instructions/general.instructions.md` を参照。**

### Terraform Specific Security

- IAM ポリシーは最小権限で設計する
- S3 等はデフォルトで暗号化を有効化する
- API Gateway 等は WAF・レート制限を設定する
- モジュール・依存関係は最新バージョンを使用する
- リソースベースポリシー・VPC セキュリティグループを適切に設定する

## MCP Tools

**詳細な MCP Tools の設定は `.github/instructions/general.instructions.md` を参照。**

Terraform 作業での主な活用：

### awslabs.aws-api-mcp-server (AWS リソース確認・操作)

**Terraform 設計時のリソース確認:**

```
# 既存リソースの確認
mcp_awslabs_aws-a_suggest_aws_commands with query="List all VPCs with their CIDR blocks in us-east-1"
mcp_awslabs_aws-a_call_aws with cli_command="aws ec2 describe-vpcs --region us-east-1"

# Terraformインポート用の情報取得
mcp_awslabs_aws-a_suggest_aws_commands with query="Get security group details for Terraform import"
```

**デプロイ前後の状態確認:**

```
# リソース作成確認
mcp_awslabs_aws-a_call_aws with cli_command="aws cloudformation describe-stacks --stack-name terraform-managed-stack --region us-east-1"

# セキュリティ設定の確認
mcp_awslabs_aws-a_call_aws with cli_command="aws s3api get-bucket-versioning --bucket terraform-state-bucket"
```

### aws-knowledge-mcp-server (公式リファレンス)

**Terraform ベストプラクティス確認:**

```
# AWSサービス固有の制限・要件確認
mcp_aws-knowledge_aws___search_documentation with search_phrase="S3 bucket policy terraform cross-account access"

# セキュリティ設定の公式ガイダンス
mcp_aws-knowledge_aws___read_documentation with url="https://docs.aws.amazon.com/s3/latest/userguide/bucket-policies.html"

# 関連する設定オプションの発見
mcp_aws-knowledge_aws___recommend with url="https://docs.aws.amazon.com/s3/latest/userguide/security.html"
```

### context7 (Terraform プロバイダー・モジュール情報)

**Terraform プロバイダー使用方法:**

```
# AWSプロバイダーの最新機能確認
mcp_context7_resolve-library-id with libraryName="terraform-aws-provider"
mcp_context7_get-library-docs with context7CompatibleLibraryID="/hashicorp/terraform-provider-aws" and topic="s3 bucket configuration"

# Terraformモジュール使用例
mcp_context7_resolve-library-id with libraryName="terraform"
mcp_context7_get-library-docs with context7CompatibleLibraryID="/hashicorp/terraform" and topic="module composition"
```

**Infrastructure as Code での品質向上:**

- **設計検証**: AWS 公式ドキュメントでサービス制限・要件を事前確認
- **セキュリティ**: 最新のセキュリティベストプラクティスを Terraform コードに反映
- **メンテナンス**: プロバイダーアップデート時の破壊的変更を事前に把握
