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

#### Test File Naming Convention

- テストファイルは対象ファイル名に `_test.go` を付与
- 例: `user_service.go` → `user_service_test.go`

#### Test Function Naming Convention

- テスト関数は `Test` で始まる PascalCase
- テスト対象の関数/メソッド名を基に命名
- 例: `TestNewUserService`, `TestUserService_GetUser`, `TestUserService_CreateUser_Error`

#### testify Usage Guidelines

- アサーションには `testify/assert` を使用
- モックには `testify/mock` を使用
- スイートテストには `testify/suite` を使用

```go
// Good: testify/assert を使用した明確なアサーション
func TestUserService_GetUser(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)
    userID := "123"
    expectedUser := &User{ID: userID, Name: "John"}

    mockRepo.On("GetUser", userID).Return(expectedUser, nil)

    // Act
    user, err := service.GetUser(userID)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, expectedUser.ID, user.ID)
    assert.Equal(t, expectedUser.Name, user.Name)
    mockRepo.AssertExpectations(t)
}

// Bad: 標準 testing パッケージのみを使用
func TestUserService_GetUser_Bad(t *testing.T) {
    // 冗長で読みにくいアサーション
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if user == nil {
        t.Error("Expected user, got nil")
    }
    // ... 多くの if 文が必要
}
```

#### Mock Implementation Guidelines

- インターフェースベースの依存注入を前提としたモック実装
- モックはテストファイル内に定義
- モックメソッドは対象インターフェースと完全に一致させる

```go
// Mock implementation example
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetUser(id string) (*User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *User) error {
    args := m.Called(user)
    return args.Error(0)
}
```

#### Table-Driven Tests

- 複数のテストケースをテーブル形式で定義
- 正常系・異常系の両方を網羅的にテスト
- テストケースごとに明確な名前を付与

```go
func TestUserService_ValidateUser(t *testing.T) {
    tests := []struct {
        name        string
        user        *User
        expectError bool
        errorMsg    string
    }{
        {
            name:        "valid user",
            user:        &User{ID: "123", Name: "John", Email: "john@example.com"},
            expectError: false,
        },
        {
            name:        "empty name",
            user:        &User{ID: "123", Name: "", Email: "john@example.com"},
            expectError: true,
            errorMsg:    "name is required",
        },
        {
            name:        "invalid email",
            user:        &User{ID: "123", Name: "John", Email: "invalid-email"},
            expectError: true,
            errorMsg:    "invalid email format",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            service := NewUserService(nil) // No repository needed for validation
            err := service.ValidateUser(tt.user)

            if tt.expectError {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errorMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

#### Test Structure Best Practices

- **Arrange-Act-Assert (AAA)** パターンを使用
- テストは独立して実行可能
- テストデータの準備は明確に分離
- モックの期待値設定はテスト開始時に完了

```go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange: テストデータの準備とモックの設定
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)

    inputUser := &User{Name: "John", Email: "john@example.com"}
    expectedUser := &User{ID: "generated-id", Name: "John", Email: "john@example.com"}

    mockRepo.On("CreateUser", mock.MatchedBy(func(u *User) bool {
        return u.Name == inputUser.Name && u.Email == inputUser.Email
    })).Return(nil).Run(func(args mock.Arguments) {
        user := args.Get(0).(*User)
        user.ID = "generated-id" // Simulate ID generation
    })

    // Act: テスト対象の実行
    createdUser, err := service.CreateUser(inputUser)

    // Assert: 結果の検証
    assert.NoError(t, err)
    assert.NotNil(t, createdUser)
    assert.Equal(t, expectedUser.ID, createdUser.ID)
    assert.Equal(t, expectedUser.Name, createdUser.Name)
    assert.Equal(t, expectedUser.Email, createdUser.Email)
    mockRepo.AssertExpectations(t)
}
```

#### Error Testing Guidelines

- エラーの種類ごとに個別のテストケースを作成
- エラーメッセージの内容も検証
- 期待されるエラー型をチェック

```go
func TestUserService_GetUser_Error(t *testing.T) {
    tests := []struct {
        name          string
        userID        string
        mockError     error
        expectedError string
    }{
        {
            name:          "user not found",
            userID:        "nonexistent",
            mockError:     ErrUserNotFound,
            expectedError: "user not found",
        },
        {
            name:          "database error",
            userID:        "123",
            mockError:     errors.New("database connection failed"),
            expectedError: "database connection failed",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(MockUserRepository)
            service := NewUserService(mockRepo)

            mockRepo.On("GetUser", tt.userID).Return(nil, tt.mockError)

            user, err := service.GetUser(tt.userID)

            assert.Error(t, err)
            assert.Nil(t, user)
            assert.Contains(t, err.Error(), tt.expectedError)
            mockRepo.AssertExpectations(t)
        })
    }
}
```

#### Test Helper Functions

- 共通のテストセットアップはヘルパー関数として分離
- ヘルパー関数はテストファイル内で定義
- 複雑なモック設定はヘルパー関数化

```go
// Test helper functions
func createMockUserRepository() *MockUserRepository {
    return new(MockUserRepository)
}

func createTestUserService(repo *MockUserRepository) *UserService {
    return NewUserService(repo)
}

func setupValidUser() *User {
    return &User{
        ID:    "123",
        Name:  "John Doe",
        Email: "john@example.com",
    }
}

// Usage in tests
func TestUserService_UpdateUser(t *testing.T) {
    mockRepo := createMockUserRepository()
    service := createTestUserService(mockRepo)
    user := setupValidUser()

    // ... test implementation
}
```

#### Coverage Guidelines

- 目標カバレッジ: 80%以上
- すべての公開関数/メソッドをテスト
- エラーケースも含めて網羅的にテスト
- カバレッジレポートは `go test -cover` で確認

```bash
# カバレッジレポート生成
go test -cover ./pkg/...

# HTMLレポート生成
go test -coverprofile=coverage.out ./pkg/...
go tool cover -html=coverage.out -o coverage.html
```

#### Integration Test Guidelines

- 統合テストは別ファイル（`*_integration_test.go`）に分離
- ビルドタグ `// +build integration` を使用
- 実際の依存関係（DB、外部 API）を使用

```go
// +build integration

package service_test

import (
    "testing"
    // ... imports
)

func TestUserService_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    // Setup real database connection
    db := setupTestDatabase()
    defer db.Close()

    repo := repository.NewUserRepository(db)
    service := NewUserService(repo)

    // ... integration test implementation
}
```

#### Benchmark Test Guidelines

- ベンチマークテストは `Benchmark` で始まる関数名
- パフォーマンス測定対象の関数をテスト
- `testing.B` を使用

```go
func BenchmarkUserService_GetUser(b *testing.B) {
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)

    userID := "123"
    expectedUser := &User{ID: userID, Name: "John"}

    mockRepo.On("GetUser", userID).Return(expectedUser, nil)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = service.GetUser(userID)
    }
}
```

#### Test File Organization

- テストファイルは対象パッケージと同じディレクトリに配置
- テスト専用パッケージは `package_test` 形式を使用
- 例:
  ```
  pkg/user/
  ├── user.go
  ├── user_test.go          # 同じパッケージ
  └── user_integration_test.go  # 統合テスト
  ```

#### Common Test Patterns

- **Constructor Tests**: `TestNewXxx`, `TestNewXxx_Error`
- **Method Tests**: `TestXxx_MethodName`, `TestXxx_MethodName_Error`
- **Validation Tests**: `TestXxx_ValidateXxx`
- **Edge Case Tests**: `TestXxx_EdgeCase`, `TestXxx_BoundaryCondition`

#### Test Data Management

- テストデータは定数またはヘルパー関数で定義
- ハードコードされた値は避け、意味のある名前を付与
- テストデータの変更が他のテストに影響しないよう独立させる

```go
const (
    testUserID    = "test-user-123"
    testUserName  = "Test User"
    testUserEmail = "test@example.com"
)

func createTestUser() *User {
    return &User{
        ID:    testUserID,
        Name:  testUserName,
        Email: testUserEmail,
    }
}
```

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
