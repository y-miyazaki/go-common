---
applyTo: "**"
---

<!-- omit in toc -->

# GitHub Copilot Base Instructions

**Language Note**: This document is written in Japanese, but all generated code and comments must be in English.

## Git Command Guidelines

- コミットメッセージは以下のフォーマットに従う
  - コミットメッセージは英語で記載
  - markdown 形式で記載
  - 先頭行には、#(h1)をつけて、<概要>
  - ２行目以降には、リスト形式で詳細を記載

## Copilot fixed code Guidelines

- コード修正後は必ずコマンド動作検証を行う
- 修正完了報告は全て完了してから報告する
  - 一部の場合は最後に"まだ作業中のため次の指示が必要です"と太字で記載する
  - 残っている作業内容はリスト化して表示する
- 複数のコード修正時でもユーザの毎回の確認は必要ない
- 統一性
  - 修正内容が他のコードにも適用すべき場合は、grep で他ファイルも検索し、必要なら全体を修正する
  - 共通ライブラリの関数内容は全て把握した上で修正する
  - 修正完了後は何を検索しても問題なかったかを grep レベルで記載する
  - grep だけでは問題がありそうな場合は、コード自体も中身も実際に読み込む
- エラー修正は Copilot が自律的に実行する
- コマンド動作検証
  - すべてワンライナーで実行し、複数回の動作確認はユーザー確認不要
  - dry-run オプションは使用しない
  - && で全てワンライナー対応
  - 複数結果確認は || でワンライナー対応
  - ターミナルが終了するため、コマンド実行時に `set -e`は利用しない。
- コマンド動作実行
  - 出力ファイル名はデフォルトから変更しないこと。そのため output オプションは指定しない
  - 出力ファイルは毎回名前変更しないこと。比較が必要な場合のみ許可

## General Principles

- プロジェクト全体で統一性を保つ
  - 修正時は同ディレクトリの既存コードを参照
  - コードは一貫性・明瞭性を持って記述する
- 変数・関数・構造体・リソース名の記載方法は言語ルールに乗っ取り、統一する
- 固有ライブラリを利用する際には事前にコンテキストとして利用できるよう把握する
- すべての関数・リソースは明確な説明を含める
- DRY 原則（重複排除）を守る
- 明示的なコードを優先する
- 非自明なロジックには必ずコメントを付与する

## Documentation and Comments

- すべてのファイル/スクリプト/モジュールは目的を記載したヘッダーを含める
- すべての関数は目的・引数説明を含める
- 複雑なロジックにはインラインコメントを付与する
- コメント・ドキュメントは英語で記載すること（全言語共通）

## Error Handling

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

## Development Workflow

- 新規作業は feature/xxx または hotfix/xxx ブランチで行う
- 必ずプルリクエストを作成し、コードレビューを受ける
- マージ前に最低 1 回レビューを受ける
- マージ前にベースブランチと同期する
- 必要に応じてドキュメント・使用例も更新する

## MCP Tools

本プロジェクトでは Model Context Protocol (MCP) 対応ツールを活用する。
目的: コードベース解析、ドキュメント検索、AWS 操作支援、変更の一貫性向上。

### 利用可能な MCP サーバ

#### serena (Code Analysis & Safe Editing)

**プロジェクト初期化時の必須操作:**

```
/mcp__serena__initial_instructions
```

**主要な使用パターン:**

1. **プロジェクト構造の理解:**

   ```
   mcp_serena_list_dir with relative_path="." and recursive=true
   mcp_serena_get_symbols_overview with relative_path="pkg/service/user_service.go"
   ```

2. **シンボル検索と編集:**

   ```
   mcp_serena_find_symbol with name_path="UserService" and include_body=true
   mcp_serena_replace_symbol_body with name_path="UserService/GetUser"
   ```

3. **安全な一括変更:**
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

### 実践的な使用ワークフロー

#### Go 開発での典型的フロー

1. **初期設定**: `mcp_serena_activate_project` で対象プロジェクトをアクティベート
2. **構造把握**: `mcp_serena_list_dir` でディレクトリ構造確認
3. **コード理解**: `mcp_serena_get_symbols_overview` でファイル概要取得
4. **詳細分析**: `mcp_serena_find_symbol` で特定のシンボル詳細確認
5. **安全な編集**: `mcp_serena_replace_symbol_body` で変更実行
6. **影響確認**: `mcp_serena_find_referencing_symbols` で参照元チェック

#### AWS 操作での典型的フロー

1. **権限確認**: `aws sts get-caller-identity`
2. **リソース調査**: `mcp_awslabs_aws-a_suggest_aws_commands` で適切なコマンド提案
3. **安全実行**: 提案されたコマンドを段階的に実行
4. **ドキュメント参照**: 不明点は `aws-knowledge-mcp-server` で公式ドキュメント確認

### 注意事項とベストプラクティス

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
