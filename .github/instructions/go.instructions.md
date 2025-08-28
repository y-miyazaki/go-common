---
applyTo: "**/*.go"
---

<!-- omit in toc -->

# GitHub Copilot Instructions for Go

**Language Note**: This document is written in Japanese, but all generated code and comments must be in English.

## Project Overview

このリポジトリは Go で構築された共通ライブラリとサンプルアプリケーションを含む。

| ディレクトリ/ファイル | 役割・説明                               |
| --------------------- | ---------------------------------------- |
| example/              | サンプルアプリケーション（Gin、S3 等）   |
| pkg/                  | 共通ライブラリ・ユーティリティ           |
| cmd/                  | 実行可能なコマンドツール（存在する場合） |
| internal/             | 内部パッケージ（存在する場合）           |

## Coding Standards

## Coding Standards

### Production and Test Code Separation

- テスト専用の依存注入関数・ラッパー・テスト用ロジックは本番コードに追加しない
- テスト性が必要な場合はインターフェース設計や構造体埋め込みを使い、テスト専用コードは `*_test.go` ファイルにのみ記述する
- 本番コードは可読性・保守性・エラー処理・ドキュメントに集中する
- テストヘルパー・モック・テスト専用ロジックはすべてテストファイルに分離する

## Naming Conventions

| コンポーネント     | 規則            | 例                                     |
| ------------------ | --------------- | -------------------------------------- |
| パッケージ名       | 小文字          | infrastructure, repository, service    |
| 関数名(公開)       | PascalCase      | NewConfig, ProcessEvent, GetUser       |
| 関数名(内部)       | camelCase       | validateInput, processEvent            |
| 変数名             | camelCase       | config, eventSource, userID            |
| 定数名             | PascalCase      | DefaultTimeout, MaxRetryCount          |
| インターフェース名 | PascalCase + er | UserRepository, EventProcessor, Logger |
| 構造体名           | PascalCase      | Config, Event, User                    |
| ファイル名         | snake_case      | main.go, event_handler.go              |

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

コード修正時は `/workspace/scripts/go/check.sh` で一括検証する

```bash
# 特定ディレクトリのみ（推奨: 対象ディレクトリ）
bash /workspace/scripts/go/check.sh -f ./example/gin1/

# プロジェクト全体
bash /workspace/scripts/go/check.sh

# 自動修正付き
bash /workspace/scripts/go/check.sh -f ./example/gin1/
```

検証内容：

- `go mod tidy`（依存管理・不要パッケージ削除）
- `go fmt`（自動整形）
- `go vet`（静的解析・エラー検出）
- `golangci-lint`（複数リンターによる品質チェック）
- `go test -v`（詳細な単体テスト）
- `go test -race`（レースコンディション検出、CGO_ENABLED=1 時）
- `go test -cover`（カバレッジ分析、80%以上推奨）
- `govulncheck`・ハードコード秘密検出（セキュリティチェック）
- ベンチマークテスト（詳細モード時）

### Validation Requirements

- すべての検証に合格してからコードの修正を完了とする
- テストカバレッジは可能な限り高く維持する（目標 80%以上）
- テスト失敗は必ず修正してからコミットする

### Manual Testing Requirements

- 本番コードは必ず `*_test.go` ファイルで単体テストを実施する
- アサーション・モックは `testify` を利用する
- テストヘルパー・モックは本番コードに追加しない
- テスト性はインターフェース設計で担保し、テスト専用フックは本番コードに追加しない

## Security Guidelines

**詳細な security guidelines は `.github/instructions/general.instructions.md` を参照。**

### Go Specific Security Best Practices

- エラーメッセージに機密情報を含めない
- `go mod tidy` で不要な依存を削除する
- `context` でタイムアウト・キャンセルを管理する
- 共有データは適切に同期し、`go test -race` で検証する
- 外部入力は必ずバリデーションする

## MCP Tools

**詳細な MCP Tools の設定は `.github/instructions/general.instructions.md` を参照。**

Go 開発での主な活用：

- `serena`: コードベース解析・安全な編集支援。プロジェクト初期化時には `/mcp__serena__initial_instructions` を実行
- `awslabs.aws-api-mcp-server`: AWS CLI コマンド提案・実行
- `context7`: コンテキスト情報の管理・操作支援
