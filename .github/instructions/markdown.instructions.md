---
applyTo: "**/*.md"
---

# GitHub Copilot Instructions for Markdown Documentation

**Language Note**: This document is written in Japanese, but all generated code and comments must be in English.

## Overview

このリポジトリは技術ドキュメント、README、手順書等の Markdown ファイルを含みます。

| Section       | Description                                                          |
| ------------- | -------------------------------------------------------------------- |
| Prerequisites | List all required dependencies, minimum versions, installation links |
| Installation  | Step-by-step installation instructions                               |
| Usage         | Basic usage examples and common scenarios                            |
| Configuration | Configuration options and environment variables                      |
| API Reference | Link to detailed API documentation or key endpoints                  |
| Examples      | Real-world examples and use cases                                    |
| Contributing  | Link to CONTRIBUTING.md with contribution guidelines                 |
| License       | License information and link to LICENSE file                         |

## Standards

### Markdown Guidelines

### Consistency with Existing Documentation

- **ドキュメント構成**: ワークスペース内の既存 Markdown ファイルを必ず確認し、構成やフォーマットの統一性を保つ
- ファイル名はプロジェクトの命名規則に従う
- セクション順・見出し階層は類似ドキュメントと揃える
- 内部リンク・参照は一貫性を保つ

### Workspace-aware Documentation

- 実際のディレクトリ構成を参照してドキュメントを作成する
- 設定例は実際のプロジェクトパターンを使う
- 長い URL は参照型リンクを使う
- 重要な注意事項は引用やアドモンションで記載する
- 複雑なドキュメントパターンには使用例を記載する

### Standardized Chapter Structure

1. **言語固有の章**: 必要に応じて以下の章を追加可能

- Documentation and Comments（ドキュメントとコメント）
- Error Handling（エラー処理）

2. **第３項目（###）以下**: 各言語・フレームワークの特性に応じて調整可能

3. **テンプレート参照**: copilot-instructions-scripts.md を統一標準のテンプレートとして使用

4. **第３項目（###）レベルの統一**:

- Testing and Validation 章の統一構造:
  - Code Modification Guidelines（共通）
  - Validation Requirements（共通）
  - Manual Testing Requirements（Go 固有）
  - Script Options（Go 固有）
- すべて英語名で統一

### Markdown Standards

- 見出し階層・画像 alt テキストを正しく使う
- ドキュメントファイル名は kebab-case を使う
- コードブロックは言語タグ付きで記載する
- 表は列揃えで記載する
- 長い URL は参照型リンクを使う
- 重要事項は引用・アドモンションで記載する
- 複雑なドキュメントには使用例を記載する
- リスト・数字付きリスト
  - タイトルの後は 2 スペースと改行を入れ、その後に説明を記載
  - インデントを考慮し、リストの次以降には 4 スペースを文章・コードブロックの先頭に付ける

### Markdown Templates

### README.md Template

```
<!-- omit in toc -->
# Project Title
<!-- omit in toc -->
## Table of Contents

## Overview

### Directory Structure

| Directory/File | Description                          |
| -------------- | ------------------------------------ |
| `dir1/`        | Description of dir1                  |

## Installation

## Local Development Environment

### Required

### Setting

### Create Local Development Environment

## Commands

## Troubleshooting

### Common Issues

### Getting Help

## Note

```

### Technical Documentation Template

````markdown
<!-- omit in toc -->

# Technical Component Name

**Language Note**: This document is written in Japanese, but all generated code and comments must be in English.

## Overview

Brief description of the component purpose and scope.

## Architecture

- System design overview
- Component relationships
- Data flow diagrams (if applicable)

## Installation

### Prerequisites

- List required dependencies
- Minimum versions
- Installation links

### Setup Instructions

Step-by-step installation process.

## Configuration

### Environment Variables

| Variable   | Description      | Default         | Required |
| ---------- | ---------------- | --------------- | -------- |
| `VAR_NAME` | Variable purpose | `default_value` | Yes/No   |

### Configuration Files

Configuration file formats and examples.

## API Reference

### Endpoints

- `GET /api/endpoint` - Endpoint description
- `POST /api/resource` - Resource creation

### Request/Response Examples

```json
{
  "example": "request_body"
}
```
````

## Examples

### Basic Usage

```bash
# Command example
command --option value
```

### Advanced Scenarios

Real-world use cases and implementations.

## Troubleshooting

### Common Issues

- **Issue**: Description and solution
- **Error**: Specific error message and fix

## Contributing

Link to CONTRIBUTING.md and development guidelines.

## License

License information and link to LICENSE file.

`````

### Program Language Template

````markdown
---
applyTo: "**/*.go"
---

<!-- omit in toc -->
# Project Title
<!-- omit in toc -->
## Table of Contents

## Overview
## Coding Standards
## Naming Conventions
## Testing and Validation
### Code Modification Guidelines
## Security Guidelines
## Reference Resources
`````

## Guidelines

### Content Guidelines

- 変更前に類似ドキュメントを確認し、トーン・構成・フォーマットを揃える
- プロジェクト用語は一貫性を保つ
- バージョン・機能は現状に合わせる

### Documentation Revision Process

1. **構造分析**: 既存の copilot-instructions-\*.md ファイルの構造を確認
2. **統一化**: 第２項目（##）レベルを標準構造に合わせる
3. **第３項目統一**: Testing and Validation 章の###レベルを英語名で統一
4. **重複削除**: 混乱したセクション・重複コンテンツを整理
5. **言語統一**: 本文を日本語に翻訳（TOC・Reference Resources は英語維持）
6. **検証**: grep_search や read_file で構造と内容を確認

## Security Guidelines

**詳細な security guidelines は `.github/instructions/general.instructions.md` を参照。**

### Markdown Specific Security

- 外部リンクは信頼できる HTTPS のみ使用する
- コード例はセキュリティベストプラクティス（環境変数利用・秘密情報ハードコーディング禁止）を守る
- 機能の権限・アクセス制御は必ず記載する

## MCP Tools

**詳細な MCP Tools の設定は `.github/instructions/general.instructions.md` を参照。**

ドキュメント作業での主な活用：

### aws-knowledge-mcp-server (公式ドキュメント参照)

**技術ドキュメント作成時:**

```bash
# AWS サービスの最新情報確認
mcp_aws-knowledge_aws___search_documentation with search_phrase="Lambda function URLs CloudFormation"

# 詳細設定手順の参照
mcp_aws-knowledge_aws___read_documentation with url="https://docs.aws.amazon.com/lambda/latest/dg/lambda-urls.html"

# 関連ドキュメントの発見
mcp_aws-knowledge_aws___recommend with url="https://docs.aws.amazon.com/lambda/latest/dg/"
```

### context7 (ライブラリドキュメント管理)

**README・技術文書での活用:**

```bash
# フレームワーク使用例の確認
mcp_context7_resolve-library-id with libraryName="next.js"
mcp_context7_get-library-docs with context7CompatibleLibraryID="/vercel/next.js" and topic="deployment"

# ライブラリのベストプラクティス取得
mcp_context7_resolve-library-id with libraryName="terraform"
mcp_context7_get-library-docs with context7CompatibleLibraryID="/hashicorp/terraform" and topic="best practices"
```

**ドキュメント品質向上での使用パターン:**

- **技術仕様確認**: 公式ドキュメントから最新の仕様・制限事項を確認
- **コード例作成**: ライブラリドキュメントから適切なサンプルコードを取得
- **ベストプラクティス反映**: 最新の推奨設定・セキュリティ要件を文書に反映
