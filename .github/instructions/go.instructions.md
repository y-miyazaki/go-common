---
applyTo: "**/*.go"
description: "AI Assistant Instructions for Go Development"
---

# AI Assistant Instructions for Go Development

**言語ポリシー**: ドキュメントは日本語、コード・コメントは英語。

Go 共通ライブラリとサンプルアプリケーションを含むプロジェクト。

| Directory/File | Purpose / Description          |
| -------------- | ------------------------------ |
| example/       | サンプルアプリ（Gin、S3 等）   |
| pkg/           | 共通ライブラリ・ユーティリティ |
| cmd/           | 実行可能コマンドツール         |
| internal/      | 内部パッケージ                 |

## Standards

- GoDoc スタイルコメント: 全パッケージ、公開関数/メソッド/構造体
- コメント・ドキュメントは英語記載
- パッケージ名: 小文字、単一単語
- ファイル名: snake_case
- 変数/関数: camelCase、公開: PascalCase
- エラーは明示的処理（`_`無視禁止）
- エラーラッピング: `fmt.Errorf` + `%w`

## Guidelines

### Code Organization

- パッケージ機能別分離
- interface は使用側で定義（Dependency Inversion）
- 循環参照回避

### Error Handling

- error 無視禁止
- context.Context でキャンセル・タイムアウト処理
- 構造化ログ使用

### Performance

- goroutine リーク防止
- channel close 責任明確化
- 高頻度操作でメモリプール検討

### Go-Specific MCP Patterns

- Struct/Interface 理解: `mcp_serena_find_symbol`
- 依存関係確認: `mcp_serena_find_referencing_symbols`
- テスト整合性: relative_path 指定で検索

### Project Structure (go-common)

#### Repository Layout

go-common 構造: .github(workflows/instructions), pkg(infrastructure/repository/service/handler/utils), example(gin/mysql/postgres/s3), scripts(go/terraform/lib), coverage

#### Editing Guidelines

- **pkg/**: 品質・テスト重視
- **example/**: 理解しやすさ重視
- **scripts/**: 安全性・エラーハンドリング重視
- **.github/**: 一貫性・メンテナンス性重視

#### Initial Onboarding

1. serena: `activate_project` → `onboarding`
2. 構造把握: `list_dir` recursive=true
3. 主要ファイル: go.mod, Makefile, README.md, .github/workflows/

### Coding Standards

#### Production and Test Separation

- テスト専用コードは本番コードに追加禁止
- テスト性はインターフェース設計で確保
- テストヘルパー・モックは`*_test.go`に分離

### Naming Conventions

| コンポーネント | 規則            | 例                 |
| -------------- | --------------- | ------------------ |
| パッケージ名   | 小文字          | infrastructure     |
| 関数(公開)     | PascalCase      | NewConfig, GetUser |
| 関数(内部)     | camelCase       | validateInput      |
| 定数           | PascalCase      | DefaultTimeout     |
| Interface      | PascalCase + er | UserRepository     |
| ファイル名     | snake_case      | event_handler.go   |

### Go Standards

- ファイル宣言順序: const → var → type (interface → struct) → func

### Lambda Function

- パッケージコメント: 機能概要記載
- 環境変数: 必須変数バリデーション実装
- エラーハンドリング: context.Context でタイムアウト処理
- ログ出力: 構造化ログ使用

### Testing

#### Test Naming

- ファイル名: `_test.go`付与（例: `user_service_test.go`）
- 関数名: `TestFunctionName_Scenario`（例: `TestUserService_GetUser`）

#### testify Usage

- Assert: `testify/assert`、Mock: `testify/mock`、Suite: `testify/suite`
- AAA Pattern: Arrange-Act-Assert で構造化

#### Mock Implementation

- インターフェースベース依存注入前提
- モックはテストファイル内定義
- モックメソッドはインターフェースと完全一致

#### Table-Driven Tests

- テーブル形式で複数ケース定義
- 正常系・異常系網羅
- 明確なケース名

#### Test Structure

- AAA Pattern 使用
- テスト独立実行可能
- データ準備明確分離
- モック期待値はテスト開始時設定

#### Error Testing

- エラー種別ごと個別ケース
- エラーメッセージ内容検証
- エラー型チェック

#### Test Helpers

- 共通セットアップヘルパー化
- 複雑モック設定ヘルパー化

#### Coverage

- 目標: 80%以上
- 全公開関数/メソッドテスト
- `go test -cover`確認

#### Integration Test

- 別ファイル: `*_integration_test.go`
- ビルドタグ: `// +build integration`
- 実依存関係使用

#### Benchmark

- 関数名: `Benchmark`始まり
- `testing.B`使用

#### Test Organization

- テストファイル同ディレクトリ配置
- テスト専用: `package_test`形式

#### Common Patterns

- Constructor: `TestNewXxx`、Method: `TestXxx_MethodName`
- Validation: `TestXxx_ValidateXxx`、Edge: `TestXxx_EdgeCase`

#### Test Data

- 定数またはヘルパー定義
- 独立性保持

## Testing and Validation

### Code Modification Guidelines

```bash
# 特定ディレクトリ検証（推奨）
bash /workspace/scripts/go/check.sh -f ./example/gin1/

# プロジェクト全体
bash /workspace/scripts/go/check.sh

# 自動修正付き
bash /workspace/scripts/go/check.sh -f ./example/gin1/ --fix
```

検証内容: go mod tidy, go fmt, go vet, golangci-lint, go test (-v/-race/-cover), govulncheck

### Validation Requirements

- 全検証合格後にコード修正完了
- テストカバレッジ 80%以上目標
- テスト失敗は修正後コミット

### Manual Testing

- 本番コードは`*_test.go`で単体テスト実施
- testify 使用、テストヘルパー・モックは本番コード追加禁止

## Security Guidelines

- エラーメッセージに機密情報含めない
- 不要依存削除（`go mod tidy`）
- タイムアウト・キャンセル管理（`context`）
- 共有データ同期と競合検証（`go test -race`）
- 外部入力バリデーション

## MCP Tools

**詳細は `.github/copilot-instructions.md` 参照。**

### Go 開発特有パターン

- **serena**: 構造体・関数理解、メソッド編集、インターフェース実装確認
- **AWS 開発**: Go SDK 設定・ベストプラクティス確認
- **依存関係**: gin/gorm ドキュメント取得
