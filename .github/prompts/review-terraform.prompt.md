---
mode: "agent"
description: "Terraformコードの正確性・セキュリティ・保守性・ベストプラクティス準拠レビュー"
---

# Terraform Review Prompt

- Terraform およびインフラコードのベストプラクティスに精通したエキスパート
- Terraform コードの正確性・セキュリティ・保守性・業界標準準拠レビュー
- Terraform MCP を利用する
- レビューコメントは日本語

## Review Guidelines (ID Based)

> 構造化チェックリスト。レビュー出力時は ID + ✅/❌ を利用。英語コメント統一は国際的な可読性と自動解析(検索/静的解析ツール)精度向上が目的。

### 1. Global / Base (G)

- G-01 Terraform ファイル構文・リソース設定妥当性 (構文/属性/必須項目)
- G-02 変数 / outputs / Module の適切利用 (未使用/重複/不要定義検出)
- G-03 シークレット・認証情報ハードコーディング禁止 (平文キー/トークンなし)
- G-04 外部 Module Version ピン + 最新バージョン併記 (ローカル module は除外)
- G-05 Provider Version constraint 記載 + (レビューコメントで最新バージョンを表示)
- G-06 apply 後決定値を for_each / count のキーに使用していない
- G-07 count より for_each 推奨 (ただし有効化/無効化トグル用途の count は許容)
- G-08 Module 使用方法: 引数設定妥当 (型/期待値)・不要引数なし
- G-09 Module 出力の活用 (不要な未参照 output 無 / 必要 output 欠落無)
- G-10 tfsec 言及があれば trivy へ移行する指摘

### 2. Modules (M)

- M-01 モジュールディレクトリ内の全 tf ファイル対象化 (main/variables/outputs/versions など)
- M-02 Provider Version チェック (aws provider は最新必須ではなく妥当性判断)
- M-03 locals / variables / outputs の責務明確 (ネーミング/分離)
- M-04 重複したタグ・命名プリフィックスの統一

### 3. variables.tf (V)

- V-01 変数名 snake_case
- V-02 型が過度に map(any)/any へ逃げていない (具体化)
- V-03 デフォルト値の妥当性 (不要 default の削除 / sentinel 値回避)
- V-04 説明コメント存在 & (Required)/(Optional) 規則順守
- V-05 validation ブロックで禁止パターン (length > 0 等) 不使用
- V-06 不要/未使用変数なし

### 4. outputs.tf (O)

- O-01 各 output に description (なければ指摘)
- O-02 機密情報出力なし (ARN/ID 可, 秘密値不可)
- O-03 未参照 output 無 (参照されていないものは削除提案)

### 5. tfvars (T)

- T-01 変数名 snake_case
- T-02 シークレット/認証情報未記載 (Secret Manager / SSM Parameter へ誘導)
- T-03 環境ごとファイル分離 (dev/stg/qa/prd)
- T-04 他環境識別子 (アカウント ID / VPC ID / Subnet Group / SG ID) の混在なし
- T-05 環境名 prefix の誤混在なし (terraform.{env}.tfvars 内の他環境文字列)

### 6. Security (S)

- S-01 KMS 暗号化可能リソースで暗号化設定 (SNS/S3/Logs/StateMachines 等)
- S-02 IAM Policy 最小権限 ("\*" は必要最小限定理由提示)
- S-03 EventBridge → SNS / Step Functions などのリソースポリシーに Condition(SourceArn 等) 付与
- S-04 平文シークレット無し (token/password/API key)
- S-05 Logging 設定 (CloudTrail / CloudWatch Logs / Step Functions) 適切

### 7. Tagging (TAG)

- TAG-01 `locals.tags = merge(try(data.aws_default_tags.provider.tags, {}), var.tags == null ? {} : var.tags)` 統一
- TAG-02 Name 追加は `merge(local.tags,{ Name = "..." })` 形式
- TAG-03 不要な手動重複キー再定義なし (provider default に一致する再指定回避)

### 8. Events & Observability (E)

- E-01 EventBridge event_pattern 過剰キャッチ回避 (必要フィルタ有り)
- E-02 CloudWatch Log Group retention 設定 (デフォルト無制限回避)
- E-03 アラーム / メトリクス / Dashboard 設定整合 (命名規則, prefix)
- E-04 Step Functions ログ出力レベル適切

### 9. Versioning (VERS)

- VERS-01 `required_version` 範囲がプロジェクト標準に準拠
- VERS-02 provider version 範囲 (>= lower,< upper) の互換性表記
- VERS-03 外部 module バージョン固定 (SHA/pseudo version 回避)

### 10. Naming & Docs (N)

- N-01 ファイル命名規約 (snake_case / 明確な purpose)
- N-02 コメント英語 (国際的共有と自動解析精度向上のため) ※ 違反時指摘
- N-03 Module 冒頭ヘッダー (目的 / 概要)
- N-04 重要リソース説明コメント (ポリシー/イベントパターン)

### 11. CI & Lint (CI)

- CI-01 terraform fmt / validate / tflint / trivy / security scan 前提
- CI-02 `terraform plan` 差分が意図通り (タグ単純化などで無駄差分なし)
- CI-03 新規リソース追加は明確な要件裏付け

### 12. Patterns (P)

- P-01 dynamic blocks 過剰使用回避 (可読性確保)
- P-02 for_each のキー安定 (map/object keys 明示)
- P-03 不要な count = 0/1 トグルを複数段で連鎖させない

### 13. State & Backend (STATE)

- STATE-01 remote backend 設定 (bucket / key / region / dynamodb_table など) が暗号化 (SSE) とロック機構(DynamoDB) を利用
- STATE-02 backend 設定に資格情報/シークレット直接記載なし (環境変数 / 認証プロファイル使用)
- STATE-03 複数 workspace 運用方針が明文化 (workspace 名の衝突なし)
  - workspace は利用しない方針
- STATE-04 `terraform state` 手動操作が必要なケースはドキュメント化 (moved/import 記録)

### 14. Compliance & Policy (COMP)

- COMP-01 Organization / Security Hub / Config などガバナンスリソースは意図と整合
- COMP-02 脆弱性スキャン (trivy) 結果をパイプラインに統合 (失敗基準定義)
- COMP-03 デフォルト VPC / オープンな SG / パブリック S3 バケット等の禁止パターン検出

### 15. Cost Optimization (COST)

- COST-01 不要な高コストメトリクス/ログ保持期間 (retention) が長期化していない
- COST-02 大量リソース生成ループ (for_each 大量展開) に cost justification あり
- COST-03 無効化可能なオプション (detailed monitoring, xray, log retention) のデフォルト最小化

### 16. Testing & Validation (TEST)

- TEST-01 `terraform validate` / `tflint` / `trivy` / `terraform plan` が CI パイプラインで実行

### 17. Migration & Refactor (MIG)

- MIG-01 `moved` ブロック利用でリネーム/構造変更時のリソース再作成回避
- MIG-02 廃止(deprecated) リソース/属性が残存せず最新版へ置換
- MIG-03 一時的コメントアウトリソース (drift の原因) 不残存

### 18. Performance & Limits (PERF)

- PERF-01 大量 for_each / count による plan 実行時間過多を回避 (分割検討)
- PERF-02 Provider 呼び出し数削減 (重複 data source を locals/outputs 共有)
- PERF-03 CloudWatch イベント/アラーム過剰生成 (サービスクォータ近傍) 監視

### 19. Dependency & Ordering (DEP)

- DEP-01 `depends_on` は最小限 (暗黙依存で十分な箇所に冗長記述なし)
- DEP-02 循環参照回避 (data <-> resource 間でループなし)
- DEP-03 implicit dependency の明示化が必要なケース (e.g., bucket policy before replication) は適切に depends_on 記載

### 20. Data Sources & Imports (DATA)

- DATA-01 data source 利用は再評価 (静的値に置換できる場合は簡素化)
- DATA-02 import 予定/実施リソースは README やコメントに手順記録
- DATA-03 外部 ID / ARN 参照はハードコードせず変数化 (アカウント間で再利用性確保)
- DATA-04 不要になった data source 削除 (未参照検出)

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
# Terraform Review Result

## Issues

None ✅ (No issues detected across all checklist items)
```

### ❌ Issues Found Example

```markdown
# Terraform Review Result

## Issues

1. G-06 apply 後決定値 for_each 利用

   - Problem: for_each = aws_s3_bucket.example.tags をキーに使用
   - Impact: 計画差分不安定 / 並列適用リスク
   - Recommendation: 事前決定可能な map(var.enabled_buckets) のような静的キーへ置換

2. S-03 EventBridge → SNS Policy Condition 欠落

   - Problem: sns:Publish ポリシーに SourceArn Condition 不在
   - Impact: 他アカウント/予期せぬイベントからの Publish 余地
   - Recommendation: Condition.StringEquals { aws:SourceArn = module.eventbridge_rule.arn }

3. TAG-02 Name タグ統一違反

   - Problem: resource 内で tags = merge(local.tags,{ Name = ... }) 形式未使用
   - Recommendation: merge(local.tags,{ Name = "${var.name_prefix}example" }) へ統一

4. O-01 output description 欠落

   - Problem: outputs.tf の output "sns_topic_arn" に description 未定義
   - Recommendation: description = "SNS topic ARN for alarm notifications"

5. TEST-03 Force replacement 未説明
   - Problem: bucket versioning 変更で recreate が plan に表示も README に手順無し
   - Recommendation: MIGRATION.md に手順と影響範囲記載 + moved ブロック検討
```

### Compact Table Variant (Optional)

```markdown
| ID      | Status | Note                            |
| ------- | ------ | ------------------------------- |
| G-06    | ❌     | for_each uses runtime attribute |
| S-03    | ❌     | Missing SourceArn condition     |
| TAG-02  | ❌     | Name tag pattern mismatch       |
| O-01    | ❌     | No description on output        |
| TEST-03 | ❌     | Force replacement undocumented  |
```

---
