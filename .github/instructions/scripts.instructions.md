---
applyTo: "**/*.sh,scripts/**"
---

<!-- omit in toc -->

# GitHub Copilot Instructions for Shell Scripts

**Language Note**: This document is written in Japanese, but all generated code and comments must be in English.

<!-- omit in toc -->

## Table of Contents

- [Project Overview](#project-overview)
- [Coding Standards](#coding-standards)
- [Naming Conventions](#naming-conventions)
  - [Shell Script Standards](#shell-script-standards)
- [Documentation and Comments](#documentation-and-comments)
  - [Help Function Standards](#help-function-standards)
- [Error Handling](#error-handling)
- [Testing and Validation](#testing-and-validation)
  - [Code Modification Guidelines](#code-modification-guidelines)
  - [Validation Requirements](#validation-requirements)
- [Security Guidelines](#security-guidelines)
  - [Shell Script Security Best Practices](#shell-script-security-best-practices)
  - [AWS Environment Security](#aws-environment-security)

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
set -euo pipefail # Error handling: exit on error, unset variable, or failed pipeline

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

コード修正時は /workspace/scripts/validate_all_scripts.sh -v -f で一括検証すること。

```bash
bash /workspace/scripts/validate_all_scripts.sh -v -f
```

- チェックスクリプトは以下を実施するため、個別で以下のコマンドを実行しない
  - 再帰的なスクリプト検出
  - Shebang
  - 実行権限
  - bash 構文
  - shellcheck
  - 関数/複雑度分析
- すべてのスクリプトが All validations passed となる
- 警告・エラーが出た場合は shellcheck/bash の推奨に従い修正する
- -v オプションで詳細表示
- 警告・エラーが残る場合は shellcheck/bash の推奨に従い修正する

### Validation Requirements

- コード修正後は上記統合検証スクリプトを必ず実行する
- 実行可能スクリプトはサンプル実行または dry-run で検証する
- テスト結果は明示し、失敗時は修正案を提示する

## Security Guidelines

### Shell Script Security Best Practices

- 秘密情報は AWS Secrets Manager や Parameter Store を利用する
- スクリプトのファイル権限を適切に設定する
- 入力値は必ずバリデーションし、コマンドインジェクションを防ぐ
- 一時ファイルは安全に管理する
- エラーは安全に処理し、情報漏洩を防ぐ
- 機密操作はログに記録するが、機密データ自体は記録しない
- 必要最小限の権限でスクリプトを実行する
- 外部データは必ず検証・サニタイズする

### AWS Environment Security

- スクリプト実行は最小限の IAM 権限で行う
- 一時認証情報やインスタンスプロファイルを活用する
- すべての操作はログ記録する
- ネットワークアクセスは必要最小限に制限する
- データ保存・転送時は暗号化を有効化する
