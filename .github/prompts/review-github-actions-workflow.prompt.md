````prompt
---
mode: "agent"
description: "GitHub Actions Workflows の正確性・セキュリティ・保守性・ベストプラクティス準拠レビュー"
tools: ["github", "context7", "fetch_webpage"]
---

# GitHub Actions Workflow Review Prompt

- GitHub Actions およびCI/CDのベストプラクティスに精通したエキスパート
- GitHub Actions Workflows の正確性・セキュリティ・保守性・業界標準準拠レビュー
- 使用 MCP List
  - github
  - context7
  - fetch_webpage
- レビューコメントは日本語

## Review Guidelines (ID Based)

> 構造化チェックリスト。レビュー出力時は ID + ✅/❌ を利用。英語コメント統一は国際的な可読性と自動解析(検索/静的解析ツール)精度向上が目的。

### 1. Global / Base (G)

- G-01 YAML 構文妥当性 (actionlint 合格・インデント・引用符)
- G-02 ワークフロー命名規則統一 (name フィールド明確・一貫性)
- G-03 シークレット・認証情報ハードコーディング禁止 (平文トークン/パスワードなし)
- G-04 Actions バージョン管理適切 (commit SHA 利用・コメントでバージョン併記)
  - **必須**: uses で `@v1.2.3` のようなタグ指定は避け、commit SHA を使用
  - **推奨**: `# v1.2.3` のようにコメントでバージョン情報を併記
- G-05 環境変数命名規則 (UPPER_SNAKE_CASE・明確な用途)
- G-06 job 名・step 名一貫性 (プロジェクト全体で統一・わかりやすさ)
- G-07 不要な step なし (未使用・重複・デッドコード)
- G-08 コメント英語統一 (国際的共有と自動解析精度向上のため)
- G-09 条件分岐適切 (if 文・論理演算子・式評価)
- G-10 エラーハンドリング適切 (失敗時動作・continue-on-error・timeout)

### 2. Trigger Configuration (TRIG)

- TRIG-01 トリガーイベント適切 (push/pull_request/workflow_dispatch など)
- TRIG-02 ブランチフィルタ適切 (main/develop/feature/* など)
- TRIG-03 path フィルタ最適化 (変更検出範囲・不要実行回避)
- TRIG-04 schedule 設定適切 (cron 式・タイムゾーン・頻度)
- TRIG-05 workflow_dispatch 入力定義 (type/required/default/description)
- TRIG-06 workflow_call 再利用性 (汎用的設計・入力/出力明確)
- TRIG-07 不要なトリガー重複なし (同一条件の複数定義回避)
- TRIG-08 pull_request イベント適切 (types: opened/synchronize/reopened)
- TRIG-09 concurrency 設定適切 (group/cancel-in-progress)
- TRIG-10 repository_dispatch カスタムイベント適切

### 3. Permissions (PERM)

- PERM-01 最小権限原則遵守 (必要最小限の権限のみ付与)
- PERM-02 job レベル権限明示 (workflow レベルより優先)
- PERM-03 OIDC トークン発行権限 (id-token: write・必要時のみ)
- PERM-04 contents 権限適切 (read/write・用途明確)
- PERM-05 pull-requests 権限適切 (write・PR コメント時のみ)
- PERM-06 issues 権限適切 (write・issue 操作時のみ)
- PERM-07 packages 権限適切 (write・パッケージ公開時のみ)
- PERM-08 deployments 権限適切 (write・デプロイ記録時のみ)
- PERM-09 不要な広範囲権限なし (write-all 回避)
- PERM-10 fork PR からの権限昇格回避 (pull_request_target 注意)

### 4. Reusable Workflows (REUSE)

- REUSE-01 再利用可能ワークフロー設計 (workflow_call・汎用性)
- REUSE-02 入力定義完全性 (inputs: type/required/description)
- REUSE-03 出力定義適切 (outputs: value/description)
- REUSE-04 secrets 伝播適切 (inherit または明示的指定)
- REUSE-05 呼び出し元との責務分離 (特定プロジェクト依存回避)
- REUSE-06 バージョン管理戦略 (タグ/ブランチ指定・破壊的変更管理)
- REUSE-07 デフォルト値適切 (合理的・環境非依存)
- REUSE-08 エラー伝播適切 (失敗時の呼び出し元への通知)
- REUSE-09 ドキュメント整備 (使用例・入力/出力説明)
- REUSE-10 テスト実施 (全 caller workflow で動作確認)

### 5. Environment Variables (ENV)

- ENV-01 環境変数定義一箇所集約 (Setup Parameters step で集中定義)
  - **必須**: `echo "VAR=value" >> $GITHUB_ENV` で定義し後続 step で参照
  - **禁止**: 複数 step での同一変数再定義・冗長なフォールバック式
- ENV-02 変数命名規則統一 (UPPER_SNAKE_CASE・明確な用途)
- ENV-03 デフォルト値適切 (合理的・安全なフォールバック)
- ENV-04 機密情報環境変数化 (secrets 利用・平文回避)
- ENV-05 workflow レベル vs job レベル適切 (スコープ明確化)
- ENV-06 環境変数参照統一 (`${{ env.VAR }}` または `${VAR}`)
- ENV-07 不要な環境変数なし (未使用・重複削除)
- ENV-08 GITHUB_ENV への書き込み適切 (grouping・順序)
- ENV-09 GITHUB_OUTPUT への書き込み適切 (step outputs)
- ENV-10 環境変数衝突回避 (GitHub 予約変数との重複なし)

### 6. Job Configuration (JOB)

- JOB-01 job 命名明確 (役割・目的わかりやすさ)
- JOB-02 runs-on 適切 (ubuntu-latest/windows-latest など)
- JOB-03 timeout-minutes 設定 (デフォルト 360 分・適切な制限)
- JOB-04 strategy matrix 活用 (複数環境・バージョンテスト)
- JOB-05 needs 依存関係適切 (順序制御・並列実行最適化)
- JOB-06 if 条件適切 (job レベル実行制御)
- JOB-07 container 利用適切 (Docker イメージ・サービスコンテナ)
- JOB-08 outputs 定義適切 (後続 job への値渡し)
- JOB-09 環境指定適切 (environment: dev/staging/prod)
- JOB-10 concurrency 設定適切 (同時実行制御・キャンセル)

### 7. Step Design (STEP)

- STEP-01 step 命名明確 (name フィールド・処理内容明示)
- STEP-02 単一責任原則 (1 step = 1 責務)
- STEP-03 条件分岐適切 (if・dependabot 除外など)
- STEP-04 continue-on-error 適切利用 (失敗許容箇所のみ)
- STEP-05 timeout-minutes 設定 (長時間実行 step に制限)
- STEP-06 working-directory 適切 (コマンド実行ディレクトリ)
- STEP-07 shell 明示 (bash/pwsh/python など)
- STEP-08 run コマンド可読性 (複数行・適切なインデント)
- STEP-09 uses Action バージョン固定 (commit SHA 利用)
- STEP-10 with パラメータ完全性 (必須項目・適切な値)

### 8. Security (SEC)

- SEC-01 secrets 管理適切 (GitHub Secrets 利用・平文回避)
- SEC-02 OIDC 認証優先 (AWS/Azure/GCP 連携時)
- SEC-03 トークン権限最小化 (GITHUB_TOKEN・least privilege)
- SEC-04 サードパーティ Action 検証 (信頼性・メンテナンス状況)
- SEC-05 script injection 対策 (ユーザー入力サニタイズ)
- SEC-06 pull_request_target 慎重利用 (fork からの権限昇格リスク)
- SEC-07 環境変数マスキング (add-mask・機密情報非表示)
- SEC-08 artifact セキュリティ (retention・公開範囲)
- SEC-09 依存関係スキャン (Dependabot・vulnerability check)
- SEC-10 コード実行検証 (信頼できないコードの実行回避)

### 9. Caching (CACHE)

- CACHE-01 キャッシュ戦略適切 (actions/cache・ビルド高速化)
- CACHE-02 キャッシュキー設計 (バージョン・依存関係ハッシュ)
- CACHE-03 restore-keys フォールバック (部分一致キャッシュ利用)
- CACHE-04 キャッシュパス適切 (必要ファイルのみ・容量最適化)
- CACHE-05 キャッシュ無効化戦略 (バージョン suffix・強制更新)
- CACHE-06 並列ジョブキャッシュ競合回避 (キーの一意性)
- CACHE-07 キャッシュサイズ制限遵守 (10GB 上限・圧縮)
- CACHE-08 キャッシュヒット率監視 (効果測定・最適化)
- CACHE-09 language 固有キャッシュ (setup-go/setup-node cache)
- CACHE-10 不要キャッシュ削除 (古いキャッシュクリーンアップ)

### 10. Artifacts (ART)

- ART-01 artifact アップロード適切 (actions/upload-artifact)
- ART-02 artifact ダウンロード適切 (actions/download-artifact)
- ART-03 artifact 命名明確 (environment/component 含む)
- ART-04 retention-days 設定適切 (保持期間・ストレージコスト)
- ART-05 artifact パス適切 (必要ファイルのみ・除外パターン)
- ART-06 圧縮戦略 (zip・サイズ最適化)
- ART-07 job 間 artifact 共有適切 (needs・データ受け渡し)
- ART-08 機密情報含有回避 (secrets/credentials 非含有)
- ART-09 artifact 命名衝突回避 (複数 job での一意性)
- ART-10 不要 artifact 削除 (ストレージ管理・クリーンアップ)

### 11. Error Handling (ERR)

- ERR-01 失敗時処理定義 (if: failure()・クリーンアップ)
- ERR-02 常時実行処理定義 (if: always()・Summary 生成)
- ERR-03 エラーメッセージ明確 (失敗原因・対処方法)
- ERR-04 PR コメント通知 (失敗時・実行ログリンク)
- ERR-05 リトライ戦略 (一時的失敗の再試行)
- ERR-06 タイムアウト処理 (timeout-minutes・適切な制限)
- ERR-07 部分的失敗許容 (continue-on-error・適切な利用)
- ERR-08 デバッグ情報出力 (失敗時の詳細ログ)
- ERR-09 依存 job 失敗時の処理 (needs・if 条件)
- ERR-10 ロールバック戦略 (デプロイ失敗時の復旧)

### 12. Testing & Validation (TEST)

- TEST-01 静的解析実行 (actionlint・workflow 構文検証)
- TEST-02 ドライラン実行 (本番前の動作確認)
- TEST-03 複数環境テスト (dev/staging/prod・matrix)
- TEST-04 エッジケーステスト (空入力・境界値・エラー)
- TEST-05 並列実行テスト (race condition・データ競合)
- TEST-06 タイムアウトテスト (長時間実行の制限確認)
- TEST-07 権限テスト (最小権限での動作確認)
- TEST-08 fork PR テスト (secrets 利用不可シナリオ)
- TEST-09 再利用ワークフローテスト (全 caller での動作確認)
- TEST-10 破壊的変更影響確認 (既存ワークフロー互換性)

### 13. Performance (PERF)

- PERF-01 並列実行最適化 (needs 依存関係・matrix 活用)
- PERF-02 キャッシュ活用 (ビルド時間短縮・依存関係)
- PERF-03 不要 step スキップ (if 条件・path フィルタ)
- PERF-04 Docker layer キャッシュ (buildx・cache-from)
- PERF-05 依存関係最小化 (不要パッケージ削除)
- PERF-06 並列テスト実行 (test splitting・matrix)
- PERF-07 アーティファクト転送最適化 (圧縮・サイズ削減)
- PERF-08 API レート制限考慮 (GitHub API・外部サービス)
- PERF-09 実行時間監視 (ボトルネック特定・改善)
- PERF-10 コスト最適化 (実行時間・ストレージ・転送量)

### 14. Monitoring & Observability (MON)

- MON-01 Summary 生成 (GITHUB_STEP_SUMMARY・実行結果)
- MON-02 PR コメント通知 (成功/失敗・詳細情報)
- MON-03 ログ出力適切 (詳細レベル・構造化)
- MON-04 メトリクス収集 (実行時間・成功率・コスト)
- MON-05 失敗通知適切 (Slack/email・重要度別)
- MON-06 デバッグモード実装 (詳細ログ・環境変数制御)
- MON-07 実行履歴追跡 (workflow runs・トレンド分析)
- MON-08 リソース使用状況監視 (ストレージ・API 制限)
- MON-09 依存サービス監視 (外部 API・third-party actions)
- MON-10 SLO/SLI 定義 (CI/CD パフォーマンス目標)

### 15. Integration (INT)

- INT-01 Reviewdog 統合 (PR コメント・linter 結果)
- INT-02 Codecov 統合 (カバレッジレポート・OIDC)
- INT-03 AWS 統合 (configure-aws-credentials・OIDC)
- INT-04 Docker 統合 (build/push・registry 認証)
- INT-05 Terraform 統合 (plan/apply・state 管理)
- INT-06 データベース統合 (migration・テストデータ)
- INT-07 通知サービス統合 (Slack/Teams/Discord)
- INT-08 外部 API 連携 (webhook・status check)
- INT-09 パッケージレジストリ (npm/Maven/PyPI)
- INT-10 モニタリングツール (Datadog/NewRelic/Sentry)

### 16. Documentation (DOC)

- DOC-01 ワークフロー目的明記 (name・コメント・README)
- DOC-02 入力パラメータ説明 (description・required・type)
- DOC-03 出力説明 (description・用途・型)
- DOC-04 使用例提供 (caller workflow・実行手順)
- DOC-05 トラブルシューティング (FAQ・よくあるエラー)
- DOC-06 権限要件文書化 (必要な permissions・理由)
- DOC-07 環境変数一覧 (名前・用途・デフォルト値)
- DOC-08 依存関係明記 (外部 actions・サービス)
- DOC-09 変更履歴管理 (CHANGELOG・breaking changes)
- DOC-10 セキュリティガイド (secrets 管理・ベストプラクティス)

### 17. Maintenance (MAINT)

- MAINT-01 定期メンテナンス計画 (依存更新・非推奨対応)
- MAINT-02 Actions バージョン更新 (Dependabot・手動確認)
- MAINT-03 非推奨機能置換 (deprecated syntax・commands)
- MAINT-04 技術的負債管理 (TODO・FIXME・リファクタリング)
- MAINT-05 ワークフロー統廃合 (重複削除・統合)
- MAINT-06 コードレビュー (変更時の品質確保)
- MAINT-07 破壊的変更管理 (影響範囲確認・段階的展開)
- MAINT-08 パフォーマンス改善 (定期的な最適化)
- MAINT-09 セキュリティ更新 (脆弱性対応・迅速なパッチ)
- MAINT-10 ドキュメント更新 (変更に追従・正確性維持)

### 18. Cost Optimization (COST)

- COST-01 実行時間最適化 (並列化・キャッシュ・不要処理削除)
- COST-02 ストレージ最適化 (artifact retention・cache サイズ)
- COST-03 マシンサイズ適切 (runs-on・リソース要件)
- COST-04 不要実行削減 (path フィルタ・if 条件)
- COST-05 Self-hosted runner 検討 (高頻度実行・コスト削減)
- COST-06 Large runner 利用判断 (必要性・cost/benefit)
- COST-07 API 呼び出し最適化 (レート制限・キャッシュ)
- COST-08 並列ジョブ数制限 (max-parallel・リソース管理)
- COST-09 スケジュール実行最適化 (頻度・必要性)
- COST-10 コスト監視・分析 (usage report・最適化機会)

### 19. Compliance & Policy (COMP)

- COMP-01 組織ポリシー準拠 (命名規則・承認フロー)
- COMP-02 セキュリティポリシー (secrets 管理・権限)
- COMP-03 ブランチ保護ルール遵守 (必須チェック・承認)
- COMP-04 監査ログ記録 (重要操作・変更履歴)
- COMP-05 コンプライアンスチェック (GDPR・SOC2 など)
- COMP-06 ライセンス管理 (依存ライブラリ・Actions)
- COMP-07 データ保持ポリシー (artifact/logs retention)
- COMP-08 アクセス制御 (environment protection・reviewers)
- COMP-09 変更管理プロセス (承認・レビュー・テスト)
- COMP-10 インシデント対応 (失敗時の手順・escalation)

### 20. Advanced Patterns (ADV)

- ADV-01 Dynamic matrix 生成 (条件付き matrix・API 取得)
- ADV-02 Composite actions 活用 (再利用可能 step グループ)
- ADV-03 Custom actions 開発 (JavaScript/Docker/composite)
- ADV-04 Webhook トリガー (repository_dispatch・外部連携)
- ADV-05 Manual approval (environment protection rules)
- ADV-06 Blue-Green deployment (環境切り替え・ロールバック)
- ADV-07 Canary deployment (段階的リリース)
- ADV-08 Feature flags 統合 (条件付き機能有効化)
- ADV-09 Multi-repo orchestration (複数リポジトリ連携)
- ADV-10 Infrastructure as Code (workflow で管理)

## Output Format

- レビュー結果リスト形式
- 指摘事項は簡潔な説明と推奨修正案
- Checks
  - チェック項目を全て表示する
  - 問題がある場合はチェックを外す
  - 問題がない場合はチェックを入れる
- Issues
  - 問題があるもののみ表示
  - 修正が必要のないものをリストアップしない

## Example Output

視覚的に整理された出力例。アイコン凡例: ✅=Pass / ❌=Fail / ⚠=Needs Attention / ⏭=N/A

### ✅ All Pass Example

```markdown
# GitHub Actions Workflow Review Result

## Issues

None ✅ (No issues detected across all checklist items)
```

### ❌ Issues Found Example

```markdown
# GitHub Actions Workflow Review Result

## Issues

1. G-04 Actions バージョン管理不適切

   - Problem: `uses: actions/checkout@v4` のようなタグ指定を使用
   - Impact: セキュリティリスク・予期しない動作変更
   - Recommendation: commit SHA を使用 `uses: actions/checkout@08c6903cd8c0fde910a37f88322edcb5dd907a8 # v5.0.0`

2. PERM-01 過剰な権限付与

   - Problem: `permissions: write-all` を使用
   - Impact: セキュリティリスク・最小権限原則違反
   - Recommendation: 必要な権限のみ明示 `permissions: { contents: read, pull-requests: write }`

3. ENV-01 環境変数定義の重複

   - Problem: 複数の step で同じ環境変数を再定義
   - Impact: 可読性低下・保守性悪化・意図しない値の衝突
   - Recommendation: Setup Parameters step で一括定義し `$GITHUB_ENV` に書き出す

4. CACHE-02 キャッシュキー設計不適切

   - Problem: `key: deps` のような固定キーを使用
   - Impact: 依存関係変更時にキャッシュが更新されない
   - Recommendation: `key: deps-${{ hashFiles('**/go.sum') }}` のようにハッシュ値を含める

5. SEC-06 pull_request_target の危険な使用
   - Problem: fork からのコード実行に pull_request_target を使用
   - Impact: fork から secrets へのアクセス・権限昇格リスク
   - Recommendation: pull_request イベントを使用するか、environment protection を設定
```

### Compact Table Variant (Optional)

```markdown
| ID      | Status | Note                                |
| ------- | ------ | ----------------------------------- |
| G-04    | ❌     | Action version pinning missing      |
| PERM-01 | ❌     | Excessive permissions               |
| ENV-01  | ❌     | Duplicate environment variable defs |
| CACHE-02| ❌     | Cache key design inappropriate      |
| SEC-06  | ❌     | Dangerous pull_request_target usage |
```

---

````
