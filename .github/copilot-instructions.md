# GitHub Copilot Instructions

**Language Note**: This document is written in Japanese, but all generated code and comments must be in English.

## Core Principles

- 修正完了報告は全て完了してから報告する
  - 残っている作業内容はリスト化して表示する
- 各作業開始時に対象ファイルの instructions を読み、要点をメモとして明示する
- **Path-Specific Instructions の遵守**
  - 対象ファイルの instructions ファイルが存在する場合は、**作業前に必ず明示的に読み込み・確認する**
  - 自動適用されていても、コンテキストとして確実に把握するため `read_file` ツールで内容を確認する
  - instructions ファイルで指定された検証コマンド・コーディング規約を必ず遵守する
- **統一性の維持**
  - 修正内容が他のコードにも適用すべき場合は、grep で他ファイルも検索し、必要なら全体を修正する
- 作業中は常に「残タスクリスト」を維持し、進捗に応じて更新する

## General Standards

### Git Command Guidelines

- コミットメッセージは以下のフォーマットに従う
  - コミットメッセージは英語で記載
  - markdown 形式で記載
  - 先頭行には、#(h1)をつけて、<概要>
  - ２行目以降には、リスト形式で詳細を記載

### Copilot Fixed Code Guidelines

- コード修正後は必ずコマンド動作検証を行う
- 修正完了報告は全て完了してから報告する
  - 残っている作業内容はリスト化して表示する
- 統一性
  - 修正内容が他のコードにも適用すべき場合は、grep で他ファイルも検索し、必要なら全体を修正する
  - **修正実施時の必須手順**:
    1. 修正対象のパターンを grep で検索
    2. 同様のパターンが他ファイルにも存在するか確認
    3. 存在する場合は全てのファイルに一括適用
    4. 修正完了後は使用した grep コマンドと検索結果を報告に記載
  - 共通ライブラリの関数内容は全て把握した上で修正する
  - 修正完了後は何を検索しても問題なかったかを grep レベルで記載する
  - grep だけでは問題がありそうな場合は、コード自体も中身も実際に読み込む
- エラー修正は Copilot が自律的に実行する
- コマンド動作検証
  - dry-run オプションは使用しない
  - `||` or `&&` で全てワンライナー対応
  - ターミナルが終了するため、コマンド実行時に `set -e`は利用しない
- コマンド動作実行
  - 出力ファイル名はデフォルトから変更しないこと。
- 長時間の対話（10 ターン以上または 10,000 トークン相当）では、途中で再度ガイドラインを読み直し、遵守状況を報告する

#### Final Checklist

作業完了報告の直前に、以下のチェック項目を Markdown のチェックボックス形式で出力すること。

- [ ] Path-Specific Instructions を再確認した
- [ ] grep による統一性確認を実施し、結果を報告した
- [ ] 追加で確認が必要なファイルや関数を読み込んだ
- [ ] コマンド動作検証は、`{language}.instructions.md` にある `Code Modification Guidelines` に従って実施した
- [ ] 残作業リストを更新し、未完了項目がないことを確認した
- [ ] 逸脱や懸念点がある場合は振り返りコメントを残した

### General Principles

- プロジェクト全体で統一性を保つ
  - 修正時は同ディレクトリの既存コードを参照
  - コードは一貫性・明瞭性を持って記述する
- 変数・関数・構造体・リソース名の記載方法は言語ルールに乗っ取り、統一する
- 固有ライブラリを利用する際には事前にコンテキストとして利用できるよう把握する
- すべての関数・リソースは明確な説明を含める
- DRY 原則（重複排除）を守る
- 非自明なロジックには必ずコメントを付与する

## Language-Specific Guidelines

**重要**: 各言語の作業時は、必ず対応する instructions ファイルの内容を確認・適用すること。

## Path-Specific Instructions

**重要**: このプロジェクトでは `.github/instructions/*.instructions.md` ファイルでパス固有のカスタム指示を定義しています。

### How Path-Specific Instructions Work

各 instructions ファイルには frontmatter で `applyTo` パターンが定義されており、指定されたファイルを編集する際に自動的に適用されます:

```markdown
---
applyTo: "**/*.go"
description: "AI Assistant Instructions for Go Development"
---
```

### Available Instructions Files

| File                                                        | Applies To                                               | Description                                            |
| ----------------------------------------------------------- | -------------------------------------------------------- | ------------------------------------------------------ |
| `/workspace/.github/instructions/go.instructions.md`        | \*_/_.go                                                 | AI Assistant Instructions for Go Development           |
| `/workspace/.github/instructions/markdown.instructions.md`  | \*_/_.md                                                 | AI Assistant Instructions for Markdown Documentation   |
| `/workspace/.github/instructions/scripts.instructions.md`   | **/\*.sh,scripts/**                                      | AI Assistant Instructions for Shell Scripts            |
| `/workspace/.github/instructions/terraform.instructions.md` | **/\*.tf,**/_.tfvars,\*\*/_.hcl                          | AI Assistant Instructions for Terraform                |
| `/workspace/.github/instructions/workflows.instructions.md` | **/.github/workflows/\*.yaml,**/.github/workflows/\*.yml | AI Assistant Instructions for GitHub Actions Workflows |

### Usage Rules

1. **自動適用**: 対象ファイル編集時に自動的に適用されます
2. **明示的確認**: 自動適用されていても、作業前に必ず内容を確認すること
3. **検証の徹底**: instructions ファイルで指定された検証コマンドを必ず実行すること

## Guidelines

### Go Language

あなたは Go の専門家です。

**必須参照**:

- ファイルパス: `.github/instructions/go.instructions.md`
- Go ファイル編集時は必ずこのファイルの内容を確認すること
- **自動適用**: `applyTo: "**/*.go"` で自動適用されますが、作業前に必ず明示的に確認すること

### Terraform Language

あなたは Terraform の専門家です。

**必須参照**:

- ファイルパス: `.github/instructions/terraform.instructions.md`
- Terraform ファイル (_.tf, _.tfvars, \*.hcl) 編集時は必ずこのファイルの内容を確認すること
- **自動適用**: `applyTo: "**/*.tf,**/*.tfvars,**/*.hcl"` で自動適用されますが、作業前に必ず明示的に確認すること
- 検証コマンド、命名規則、セキュリティガイドラインが記載されています

### Shell Scripts

あなたは Shell スクリプトの専門家です。

**必須参照**:

- ファイルパス: `.github/instructions/scripts.instructions.md`
- Shell スクリプト編集時は必ずこのファイルの内容を確認すること
- **自動適用**: `applyTo: "**/*.sh,scripts/**"` で自動適用されますが、作業前に必ず明示的に確認すること

### Markdown

あなたは Github Markdown の専門家です。

**必須参照**:

- ファイルパス: `.github/instructions/markdown.instructions.md`
- Markdown ファイル編集時は必ずこのファイルの内容を確認すること
- **自動適用**: `applyTo: "**/*.md"` で自動適用されますが、作業前に必ず明示的に確認すること

### GitHub Actions

あなたは GitHub Actions の専門家です。

**必須参照**:

- ファイルパス: `.github/instructions/workflows.instructions.md`
- GitHub Actions ワークフロー編集時は必ずこのファイルの内容を確認すること
- **自動適用**: `applyTo: "**/.github/workflows/*.yaml,**/.github/workflows/*.yml"` で自動適用されますが、作業前に必ず明示的に確認すること

### Documentation and Comments

- すべてのファイル/スクリプト/モジュールは目的を記載したヘッダーを含める
- すべての関数は目的・引数説明を含める
- 複雑なロジックにはインラインコメントを付与する
- コメント・ドキュメントは英語で記載すること（全言語共通）

### Error Handling

- 必ずエラーを検知し、適切に処理する
- 具体的かつ行動可能なエラーメッセージを出す
- デバッグしやすいよう十分な情報をログ出力する
- 機密情報はエラーメッセージに含めない

## Testing and Validation

- コード修正後は必ず構文・静的解析チェックを行う
- 可能な限り機能・統合テストを実施する
- テスト結果は成功/失敗を明示する
- テスト失敗時は修正案を提示する

## Security Guidelines

- 秘密情報はハードコーディングせず、シークレット管理や環境変数を利用する
- 権限は最小限にする
- 外部入力は必ずバリデーションする
- 機密データは保存・転送時に暗号化する
- 機密操作はログに記録するが、機密データ自体は記録しない

### Development Workflow

- 新規作業は feature/xxx または hotfix/xxx ブランチで行う
- 必ずプルリクエストを作成し、コードレビューを受ける
- マージ前に最低 1 回レビューを受ける
- マージ前にベースブランチと同期する
- 必要に応じてドキュメント・使用例も更新する

## MCP Tools

本プロジェクトでは Model Context Protocol (MCP) 対応ツールを活用する。
目的: コードベース解析、ドキュメント検索、AWS 操作支援、変更の一貫性向上。

### Available MCP Servers

#### serena (Code Analysis & Safe Editing)

**Required step during project initialization:**

```
/mcp__serena__initial_instructions
```

**Common usage patterns:**

1. **Understand project structure:**

   ```
   mcp_serena_list_dir with relative_path="." and recursive=true
   mcp_serena_get_symbols_overview with relative_path="pkg/service/user_service.go"
   ```

2. **Find and edit symbols:**

   ```
   mcp_serena_find_symbol with name_path="UserService" and include_body=true
   mcp_serena_replace_symbol_body with name_path="UserService/GetUser"
   ```

3. **Safe bulk changes:**
   ```
   mcp_serena_find_referencing_symbols with name_path="GetUser" and relative_path="pkg/service/"
   ```

#### awslabs.aws-api-mcp-server (AWS CLI Operations)

**基本操作パターン:**

1. **リソース確認:**

   ```
   mcp_awslabs_aws-a_call_aws with cli_command="aws sts get-caller-identity"
   mcp_awslabs_aws-a_call_aws with cli_command="aws s3 ls --region us-east-1"
   ```

2. **不明なコマンドの提案:**

   ```
   mcp_awslabs_aws-a_suggest_aws_commands with query="List all running EC2 instances in us-east-1"
   ```

3. **安全な実行:**
   - 必ずリージョンを明示的に指定
   - 重要な操作は事前に dry-run で確認
   - 出力は最小限に制限（max_results を活用）

#### aws-knowledge-mcp-server (Documentation Reference)

**ドキュメント検索パターン:**

1. **サービス固有の検索:**

   ```
   mcp_aws-knowledge_aws___search_documentation with search_phrase="S3 bucket versioning policy"
   ```

2. **詳細ドキュメント取得:**

   ```
   mcp_aws-knowledge_aws___read_documentation with url="https://docs.aws.amazon.com/s3/latest/userguide/Versioning.html"
   ```

3. **関連情報の発見:**
   ```
   mcp_aws-knowledge_aws___recommend with url="https://docs.aws.amazon.com/s3/latest/userguide/"
   ```

#### context7 (Library Documentation)

**ライブラリ情報取得パターン:**

1. **ライブラリ ID の解決:**

   ```
   mcp_context7_resolve-library-id with libraryName="gin-gonic"
   ```

2. **ドキュメント取得:**
   ```
   mcp_context7_get-library-docs with context7CompatibleLibraryID="/gin-gonic/gin" and topic="routing"
   ```

### Practical Usage Workflows

#### Typical Go Development Flow

1. **初期設定**: `mcp_serena_activate_project` で対象プロジェクトをアクティベート
2. **構造把握**: `mcp_serena_list_dir` でディレクトリ構造確認
3. **コード理解**: `mcp_serena_get_symbols_overview` でファイル概要取得
4. **詳細分析**: `mcp_serena_find_symbol` で特定のシンボル詳細確認
5. **安全な編集**: `mcp_serena_replace_symbol_body` で変更実行
6. **影響確認**: `mcp_serena_find_referencing_symbols` で参照元チェック

#### Typical AWS Operations Flow

1. **権限確認**: `aws sts get-caller-identity`
2. **リソース調査**: `mcp_awslabs_aws-a_suggest_aws_commands` で適切なコマンド提案
3. **安全実行**: 提案されたコマンドを段階的に実行
4. **ドキュメント参照**: 不明点は `aws-knowledge-mcp-server` で公式ドキュメント確認

### Tool Selection Decision Tree

```
タスク種別による使用ツール選択:

コード理解・編集タスク
├─ Goコード → serena (優先) → 標準tools (フォールバック)
├─ 設定ファイル → 標準tools → serena (大規模変更時)
└─ ドキュメント → 標準tools

インフラ・AWS操作
├─ AWSリソース確認 → aws-api-mcp-server
├─ AWS公式ドキュメント → aws-knowledge-mcp-server
└─ Terraform → aws-api-mcp-server + aws-knowledge-mcp-server

ライブラリ・フレームワーク情報
├─ Go依存関係 → context7
├─ 新しいライブラリ調査 → context7 + aws-knowledge-mcp-server
└─ 使用例・ベストプラクティス → context7
```

### Tool Usage Priority Rules

#### 1. serena Usage Rules

- **使用すべき場面**:
  - Go 言語のコード構造理解・編集
  - 大規模なリファクタリング
  - 関数・構造体の影響範囲確認
- **使用を避ける場面**:
  - 単純なファイル読み取り
  - 設定ファイルの簡単な編集
  - 1-2 行の修正

#### 2. AWS MCP Usage Rules

- **aws-api-mcp-server**:
  - 必ず`--region`を明示的に指定
  - 本番環境では慎重に実行
  - 大量データは`max_results`で制限
- **aws-knowledge-mcp-server**:
  - 公式ドキュメントが必要な場合のみ
  - 設計・アーキテクチャ決定時に活用

#### 3. Fallback Strategy

1. **Primary Tool 失敗時**: 標準 tools に切り替え
2. **情報不足時**: 複数ツールを組み合わせ
3. **パフォーマンス問題時**: より軽量なアプローチに変更

### Tool Combination Patterns

#### New Feature Development Flow

```
1. context7でライブラリ調査
2. serenaでコード構造理解
3. serenaで実装・テスト
4. aws-api-mcp-serverで動作確認（AWS関連時）
```

#### Troubleshooting Flow

```
1. serenaでエラー箇所特定
2. aws-knowledge-mcp-serverで公式情報確認
3. aws-api-mcp-serverで実際の状態確認
4. 標準toolsで修正実行
```

1. **serena 使用時:**

   - 大きなファイルは `get_symbols_overview` → `find_symbol` の順で段階的に読み取り
   - `include_body=true` は必要な場合のみ使用
   - 編集前に必ず `find_referencing_symbols` で影響範囲を確認

2. **AWS MCP 使用時:**

   - リージョンは必ず明示的に指定
   - 本番環境操作時は特に慎重に
   - `max_results` で出力量を制御

3. **全般:**
   - MCP ツールは補助的に使用し、最終的な判断は人間が行う
   - エラー発生時はツールの制限として受け入れ、代替手段を検討
