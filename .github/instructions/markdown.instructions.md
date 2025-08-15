---
applyTo: "**/*.md"
---

<!-- omit in toc -->

# GitHub Copilot Instructions for Markdown Documentation

**Language Note**: This document is written in Japanese, but all generated code and comments must be in English.

<!-- omit in toc -->

## Table of Contents

- [Markdown Guidelines](#markdown-guidelines)
  - [Consistency with Existing Documentation](#consistency-with-existing-documentation)
  - [Workspace-aware Documentation](#workspace-aware-documentation)
- [Overview](#overview)
- [Markdown Guidelines](#markdown-guidelines-1)
  - [Consistency with Existing Documentation](#consistency-with-existing-documentation-1)
  - [Standardized Chapter Structure](#standardized-chapter-structure)
  - [Workspace-aware Documentation](#workspace-aware-documentation-1)
  - [Markdown Standards](#markdown-standards)
  - [Markdown Templates](#markdown-templates)
    - [README.md Template](#readmemd-template)
    - [Technical Documentation Template](#technical-documentation-template)
- [Program Language Template](#program-language-template)
  - [Content Guidelines](#content-guidelines)
  - [Documentation Revision Process](#documentation-revision-process)
  - [Security Best Practices](#security-best-practices)
- [Reference Resources](#reference-resources)

## Markdown Guidelines

### Consistency with Existing Documentation

- **ドキュメント構成**: ワークスペース内の既存 Markdown ファイルを必ず確認し、構成やフォーマットの統一性を保つ

### Workspace-aware Documentation

- **設定例**: サンプルを示す際は、プロジェクトディレクトリ内の実際の設定パターンを使用する
- 長い URL は参照型リンクを使う
- 重要な注意事項は引用やアドモンションで記載する
- 複雑なドキュメントパターンには使用例を記載する

## Overview

git clone https://github.com/owner/repo.git

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

## Markdown Guidelines

### Consistency with Existing Documentation

- 既存 Markdown ファイルの構成・フォーマットを必ず確認し、統一性を保つ
- ファイル名はプロジェクトの命名規則に従う
- セクション順・見出し階層は類似ドキュメントと揃える
- 内部リンク・参照は一貫性を保つ

### Standardized Chapter Structure

1. **言語固有の章**: 必要に応じて以下の章を追加可能

- Documentation and Comments（ドキュメントとコメント）
- Error Handling（エラー処理）

2. **第３項目（###）以下**: 各言語・フレームワークの特性に応じて調整可能

3. **保持要素**:

- Table of Contents: 英語で維持
- Reference Resources: 英語で維持
- 本文コンテンツ: 日本語で記載（Reference Resources 除く）

4. **テンプレート参照**: copilot-instructions-scripts.md を統一標準のテンプレートとして使用

5. **第３項目（###）レベルの統一**:

- Testing and Validation 章の統一構造:
  - Code Modification Guidelines（共通）
  - Validation Requirements（共通）
  - Manual Testing Requirements（Go 固有）
  - Script Options（Go 固有）
- すべて英語名で統一

### Workspace-aware Documentation

- 実際のディレクトリ構成を参照してドキュメントを作成する
- 設定例は実際のプロジェクトパターンを使う

### Markdown Standards

- 見出し階層・画像 alt テキストを正しく使う
- ドキュメントファイル名は kebab-case を使う
- コードブロックは言語タグ付きで記載する
- 表は列揃えで記載する
- 長い URL は参照型リンクを使う
- 重要事項は引用・アドモンションで記載する
- 複雑なドキュメントには使用例を記載する

### Markdown Templates

#### README.md Template

...existing code...

#### Technical Documentation Template

...existing code...

## Program Language Template

```
---
applyTo: "**/*.go"
---

<!-- omit in toc -->
# GitHub Copilot Instructions for {Name}
**Language Note**: This document is written in Japanese, but all generated code and comments must be in English.

<!-- omit in toc -->
## Table of Contents

## Project Overview
## Coding Standards
## Naming Conventions
## Testing and Validation
### Code Modification Guidelines
## Security Guidelines
## Reference Resources

```

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

### Security Best Practices

- 機密情報はドキュメントに記載しない
- 外部リンクは信頼できる HTTPS のみ使用する
- コード例はセキュリティベストプラクティス（環境変数利用・秘密情報ハードコーディング禁止）を守る
- 機能の権限・アクセス制御は必ず記載する
