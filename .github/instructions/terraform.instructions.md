---
applyTo: "**/*.tf,**/*.tfvars,**/*.tfstate,**/*.tfbackend"
---

<!-- omit in toc -->

# GitHub Copilot Instructions for Terraform

**Language Note**: This document is written in Japanese, but all generated code and comments must be in English.

## Project Overview

このリポジトリは AWS のセキュリティ・監視インフラを管理する Terraform コードを含みます。

| ディレクトリ/ファイル | 役割・説明                                 |
| --------------------- | ------------------------------------------ |
| terraform/application | 基本セキュリティ設定・AWS アカウント初期化 |
| terraform/base        | 基本セキュリティ設定・AWS アカウント初期化 |
| terraform/monitor     | CloudWatch 監視・アラート・可観測性        |
| terraform/management  | 組織管理・コンプライアンス・ガバナンス     |
| lambda/               | サーバレス監視・自動化                     |
| modules/              | AWS リソース共通モジュール                 |

## Coding Standards

## Naming Conventions

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

## Documentation and Comments

- すべての関数・リソースは詳細な説明を含める
- 目的・機能は冒頭で明記する
- 複雑なロジックは詳細なコメントで説明する
- コメント・ドキュメントは英語で記載する
- 複雑なモジュールは使用例も記載する

## Error Handling

## Testing and Validation

### Code Modification Guidelines

- コード修正時は必ず以下を実施する
  - ENV=dev; terraform init -reconfigure -backend-config=terraform.${ENV}.tfbackend
  - ENV=dev; terraform fmt --recursive && terraform validate && tflint -f compact --var-file=terraform.${ENV}.tfvars && trivy fs . --format table --config /workspace/trivy.yaml --secret-config /workspace/trivy-secret.yaml
  - ENV=dev; terraform plan -lock=false -var-file=terraform.${ENV}.tfvars
- テスト結果は明示し、失敗時は修正案を提示する

### Validation Requirements

- すべての検証に合格してからコード提出する
- 実行可能な Terraform コードは plan 実行で検証する
- セキュリティスキャン結果を確認し、問題があれば修正する

## Security Guidelines

### Terraform Security Best Practices

- 秘密情報はハードコーディングせず、AWS Secrets Manager 等を利用する
- IAM ポリシーは最小権限で設計する
- S3 等はデフォルトで暗号化を有効化する
- API Gateway 等は WAF・レート制限を設定する
- モジュール・依存関係は最新バージョンを使用する
- リソースベースポリシー・VPC セキュリティグループを適切に設定する

### AWS Environment Security

- CloudTrail 等の監査ログを有効化する
- trivy/AWS Inspector 等で定期的に脆弱性スキャンを行う
- ネットワークアクセスは必要最小限に制限する
- データ保存・転送時は暗号化を有効化する

## MCP Tools

- 本ドキュメントの MCP に関する詳細は `.github/instructions/general.instructions.md` の「MCP Tools」を正本として参照すること。
- Terraform 作業での補助的な利用方針:
  - awslabs.aws-api-mcp-server: 必要に応じて AWS CLI の提案・実行に使用（明示的なリージョン指定・最小スコープ運用）。
  - aws-knowledge-mcp-server: 公式ドキュメントの検索・参照に使用。
  - context7: コンテキスト情報の管理・操作を支援。
  - terraform: Terraform コードの検索・構造把握・シンボル/正規表現ベースの安全な編集に使用（.tf/.tfvars 等）。
