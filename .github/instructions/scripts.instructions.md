---
applyTo: "**/*.sh,scripts/**"
---

<!-- omit in toc -->

# GitHub Copilot Instructions for Shell Scripts

**Language Note**: This document is written in Japanese, but all generated code and comments must be in English.

## Project Overview

このリポジトリは自動化・インフラ管理用のシェルスクリプトを含みます。スクリプトは以下のように構成されています：

| ディレクトリ/ファイル           | 役割・説明                                                |
| ------------------------------- | --------------------------------------------------------- |
| scripts/terraform/              | Terraform 操作・AWS リソース管理用の自動化スクリプト      |
| scripts/go/                     | Go 言語プロジェクトのビルド・テスト・デプロイ用スクリプト |
| scripts/lib/                    | 共通ライブラリ・ユーティリティ関数                        |
| scripts/validate_all_scripts.sh | 全スクリプトの品質チェック・検証ツール                    |

## Coding Standards

## Naming Conventions

| コンポーネント       | 規則                | 例                                       |
| -------------------- | ------------------- | ---------------------------------------- |
| スクリプトファイル名 | snake_case          | validate_scripts.sh, deploy_terraform.sh |
| 関数名               | snake_case          | show_usage, log_message, check_status    |
| 変数名               | snake_case          | script_name, log_level, error_count      |
| 定数名               | UPPER_SNAKE_CASE    | DEFAULT_TIMEOUT, MAX_RETRY_COUNT         |
| 環境変数             | UPPER_SNAKE_CASE    | AWS_REGION, LOG_LEVEL, DRY_RUN           |
| 一時ファイル         | snake_case + suffix | temp_output.tmp, backup_file.bak         |

### Shell Script Standards

```bash
#!/bin/bash
#######################################
# Description: Description of script purpose and functionality
# Usage: ./script_name.sh [options] <required_arg>
#   options:
#     -h, --help     Display this help message
#     -v, --verbose  Enable verbose output
# Design Rules:
#   - Use clear and descriptive names for scripts, functions, and variables
#######################################

# Error handling: exit on error, unset variable, or failed pipeline
set -euo pipefail

# Get script directory for library loading
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export SCRIPT_DIR

# Load all-in-one library
# shellcheck source=../lib/all.sh
source "${SCRIPT_DIR}/../lib/all.sh"

#######################################
# Global variables and default values
#######################################

#######################################
# Display usage information
#######################################
show_usage() {
    cat << EOF
Usage: $(basename "$0") [options] <required_arg>

Description of script purpose and functionality.

Options:
  -h, --help     Display this help message
  -v, --verbose  Enable verbose output
  -d, --dry-run  Show what would be done without executing

Examples:
  $(basename "$0") -v production
  $(basename "$0") --dry-run staging

EOF
    exit 0
}

#######################################
# Parse command line arguments
#######################################
function parse_arguments {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h | --help)
                show_usage
                ;;
            -*)
                error_exit "Unknown option: $1"
                ;;
            *)
                error_exit "Unexpected argument: $1"
                ;;
        esac
    done
}

#######################################
# other functions...
#######################################


#######################################
# Main execution function
#######################################
main() {
    # Script implementation
    echo "Script execution starts"
    # Parse command line arguments
    parse_arguments "$@"

    # Validate dependencies
    validate_dependencies "some_command"

    # other functions

    echo_section "some comment completed"
}

# Script entry point
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
```

## Documentation and Comments

- すべてのスクリプトは目的を記載したヘッダーを含める
- すべての関数は詳細な説明を含める
- コメント・ドキュメントは日本語で記載する

### Help Function Standards

- シンプルなスクリプト: 基本オプション(-h, -v, -d 等)付きの show_usage 関数を使う
- 複雑なスクリプト: 多数オプション・カテゴリ対応の show_usage_complex パターンを使う
- フォーマット基準:
  - Usage: $(basename "$0") [options] で開始
  - スクリプト目的を明確に記載
  - オプションは 2 スペース揃えで整列
  - 代表的な使用例を記載
  - 複雑なスクリプトはカテゴリ一覧や詳細ヘルプを含める
  - 必ず exit 0 で終了

## Error Handling

- すべてのスクリプトは適切なエラー処理を実装する
- エラーメッセージは具体的かつ行動可能にする
- ログは INFO/WARN/ERROR などレベル分けして記録する

  ```bash
  # ログ関数例
  function log {
    local level=$1
    local message=$2
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] [$level] $message"
  }

  log "INFO" "Starting process"
  log "ERROR" "Failed to connect to database"
  ```

## Testing and Validation

### Code Modification Guidelines

コード修正時は以下コマンドで一括検証する：

```bash
# 全スクリプトの検証（推奨）
bash /workspace/scripts/validate_all_scripts.sh -v -f
```

検証内容：

- 再帰的なスクリプト検出
- Shebang チェック
- 実行権限チェック
- bash 構文チェック
- shellcheck 静的解析
- 関数/複雑度分析

### Validation Requirements

- すべてのスクリプトが "All validations passed" となることを確認する
- 警告・エラーが出た場合は shellcheck/bash の推奨に従い修正する
- 実行可能スクリプトはサンプル実行または dry-run で検証する

## Security Guidelines

**詳細な security guidelines は `.github/instructions/general.instructions.md` を参照。**

### Shell Script Specific Security

- 入力値は必ずバリデーションし、コマンドインジェクションを防ぐ
- 一時ファイルは安全に管理する
- エラーは安全に処理し、情報漏洩を防ぐ
- 外部データは必ず検証・サニタイズする

## MCP Tools

**詳細な MCP Tools の設定は `.github/instructions/general.instructions.md` を参照。**

スクリプト作業での主な活用：

- `awslabs.aws-api-mcp-server`: AWS CLI の提案・実行（明示的なリージョン指定・最小スコープ運用）
- `context7`: コンテキスト情報の管理・操作支援
