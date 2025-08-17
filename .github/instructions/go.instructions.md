---
applyTo: "**/*.go"
---

<!-- omit in toc -->

# GitHub Copilot Instructions for Go Lambda Applications

**Language Note**: This document is written in Japanese, but all generated code and comments must be in English.

## Project Overview

このリポジトリは AWS Lambda 上で動作する Go アプリケーションを含む。

| ディレクトリ/ファイル | 役割・説明                       |
| --------------------- | -------------------------------- |
| cmd/                  | Lambda 関数メイン処理            |
| internal/             | 内部パッケージ・ビジネスロジック |
| pkg/                  | 外部公開可能なライブラリコード   |

## Coding Standards

### Production and Test Code Separation

- テスト専用の依存注入関数・ラッパー・テスト用ロジックは本番 Go コード(main.go 等)に追加しない
- テスト性が必要な場合はインターフェース設計や構造体埋め込みを使い、テスト専用コードは \*\_test.go ファイルにのみ記述する
- 本番コードは AWS Lambda 標準構造・バリデーション・エラー処理・ドキュメント・可読性・保守性に集中する
- テストヘルパー・モック・テスト専用ロジックはすべてテストファイルに分離する

## Naming Conventions

| コンポーネント      | 規則                | 例                                      |
| ------------------- | ------------------- | --------------------------------------- |
| パッケージ名        | 小文字              | infrastructure, repository, service     |
| 関数名(公開)        | PascalCase          | NewLambdaConfig, ProcessEvent, GetUser  |
| 関数名(内部)        | camelCase           | validateInput, processCloudWatchEvent   |
| 変数名              | camelCase           | lambdaConfig, eventSource, userID       |
| 定数名              | PascalCase          | DefaultTimeout, MaxRetryCount           |
| インターフェース名  | PascalCase + er     | UserRepository, EventProcessor, Logger  |
| 構造体名            | PascalCase          | LambdaConfig, CloudWatchEvent, SlackMsg |
| Lambda 関数ファイル | snake_case          | main.go, event_handler.go               |
| AWS リソース名      | kebab-case + prefix | go-lambda-${stage}-${function}          |

## Go Standards

以下の内容は golangci-lint,go vet で指摘される項目以外の内容を記載する。

- Go ファイルの宣言順序を遵守: const -> var -> type (interface → struct) -> func (constructor → methods → helpers)

### Lambda Function Examples

```go
// Package main provides a lambda function that processes SNS events
// and forwards CloudWatch alarm notifications to CloudWatch Logs.
package main

import (
  // standard libraries
)

var (
	// LambdaConfig holds Lambda configuration and dependencies
	lambdaConfig *infrastructure.LambdaConfig
	// AWS region for CloudWatch Logs operations
	region string
)

const (
  // other constants
)

// Reporter ........
// nolint: wrapcheck, unused
func Reporter(ctx context.Context) error { // noinspection
	log := lambdaConfig.Log
  // some code to process the event
	return nil
}

// Main function initializes the Lambda handler and configuration
// nolint: unused
func main() {
	// Initialize Lambda configuration
	lambdaConfig = infrastructure.NewLambdaConfig()
	log := lambdaConfig.Log

	// Get required environment variables
	region = os.Getenv("AWS_REGION")

	// Validate required environment variables
	if region == "" {
		log.Panic("AWS_REGION environment variable is required")
	}

	// Start Lambda handler
	lambda.Start(Reporter)
}
```

### Package Structure Examples

...existing code...

### Error Handling Examples

...existing code...

### Testing Examples

...existing code...

## Testing and Validation

### Code Modification Guidelines

コード修正時は scripts/go/check.sh で一括検証する

```bash
# 特定ディレクトリのみ（推奨: 対象ディレクトリ）
bash scripts/go/check.sh -v -f ./cmd/cloudwatch/cloudwatch_alarm_to_sns_to_slack/

# プロジェクト全体
bash scripts/go/check.sh

# 自動修正付き
bash scripts/go/check.sh -v -f ./cmd/cloudwatch/cloudwatch_alarm_to_sns_to_slack/
```

- チェックスクリプトは以下を実施するため、個別で以下のコマンドを実行しない

  - go mod tidy（依存管理・不要パッケージ削除）
  - go fmt（自動整形）
  - go vet（静的解析・エラー検出）
  - golangci-lint（58 以上のリンターによる品質チェック）
  - go test -v（詳細な単体テスト）
  - go test -race（レースコンディション検出、CGO_ENABLED=1 時）
  - go test -cover（カバレッジ分析、80%以上必須）
  - govulncheck・ハードコード秘密検出（セキュリティチェック）
  - ベンチマークテスト（詳細モード時）

- すべての検証に合格してからコードの修正を完了とする
- Lambda 関数・パッケージはテストカバレッジ 80%以上を維持する
  - ただし、テストカバレッジが 80%を超えることが無理な場合もあるので、その点は一部許容する

### Manual Testing Requirements

- 本番コードは必ず \*\_test.go ファイルで単体テストを実施する
- アサーション・モックは testify を利用する
- テストヘルパー・モックは本番コードに追加しない
- テスト性はインターフェース設計で担保し、テスト専用フックは本番コードに追加しない
- CI/CD で全テスト・カバレッジ・レースチェックを実施する
- テスト失敗は必ず修正してからマージする
- テスト戦略・カバレッジは README 等に記載する

## Security Guidelines

### Go Security Best Practices

- 機密情報は AWS Secrets Manager や Parameter Store を利用する
- ソースコードに秘密情報をハードコーディングしない
- 設定は環境変数で管理する
- エラーメッセージに機密情報を含めない
- go mod tidy で不要な依存を削除する
- context でタイムアウト・キャンセルを管理する
- 共有データは適切に同期し、go test -race で検証する

### AWS Lambda Infrastructure Security

- Lambda 実行ロールは最小権限で設計する
- Lambda 環境変数で設定を管理する
- VPC 設定はネットワーク分離を考慮する
- CloudWatch ログは適切なレベルで有効化する
- Lambda 監視用 CloudWatch アラームを設定する

## MCP Tools

- Go の作業では Serena を既定で使用する（コードベース解析・参照関係・シンボル単位の安全な編集）。
- MCP の目的・利用サーバ一覧・運用ルールの正本は `.github/instructions/general.instructions.md` の「MCP Tools」を参照すること。
