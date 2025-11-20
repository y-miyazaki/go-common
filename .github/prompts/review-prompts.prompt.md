# Review Instructions for \*.prompt.md Files

## Objective

`.github/prompts/review-*.prompt.md`ファイルの実用性・統一性・詳細度確保。

## Review Checklist

### 1. File Structure (G: General)

- [ ] **G-01: Title**: レビュー対象が明確なタイトル（例: "Review Terraform Code"）
- [ ] **G-02: Purpose**: レビュー目的・スコープ記載
- [ ] **G-03: ID-Based Format**: 全チェック項目に ID 付与（例: G-01:, SEC-01:）

### 2. Category Structure (CAT: Category)

必須カテゴリ（順序統一）:

```markdown
1. General (G)
2. Code Quality (CODE) / Module (M) / Variables (V) など言語固有
3. Functionality (FUNC) / Error Handling (ERR)
4. Security (SEC)
5. Performance (PERF)
6. Testing (TEST)
7. Documentation (DOC)
8. Dependencies (DEP)
```

チェック項目:

- [ ] **CAT-01**: カテゴリ順序が統一されている
- [ ] **CAT-02**: カテゴリ ID が統一（G, CODE, FUNC, ERR, SEC, PERF, TEST, DOC, DEP）
- [ ] **CAT-03**: 番号が連番（G-01, G-02, ...）

### 3. Category Coverage (COV: Coverage)

言語別必須カテゴリ:

**Terraform** (20 カテゴリ):

- G, M, V, O, T, S, TAG, E, VERS, N, CI, P, STATE, COMP, COST, TEST, MIG, PERF, DEP, DATA

**Script** (10 カテゴリ):

- G, CODE, FUNC, ERR, SEC, PERF, TEST, DOC, DEP, LOG

**Go** (10 カテゴリ):

- G, CODE, FUNC, ERR, SEC, PERF, TEST, ARCH, DOC, DEP

**GitHub Actions** (6 カテゴリ):

- G, SEC, TOOL, ERR, PERF, BP

- [ ] **COV-01**: 対象言語の全カテゴリ存在
- [ ] **COV-02**: カテゴリ番号重複なし
- [ ] **COV-03**: 各カテゴリに最低 3 項目以上

### 4. Item Structure (ITEM: Item)

各チェック項目の必須要素:

```markdown
- **CAT-NN: Item Title**
  - Problem: 問題・リスク説明
  - Impact: 影響範囲・重要度
  - Recommendation: 具体的推奨事項
```

チェック項目:

- [ ] **ITEM-01**: 全項目に ID 付与（CAT-NN 形式）
- [ ] **ITEM-02**: Problem/Impact/Recommendation 3 要素記載
- [ ] **ITEM-03**: 具体的・実用的な推奨事項
- [ ] **ITEM-04**: コード例記載（該当する場合）

### 5. Content Quality (QUAL: Quality)

- [ ] **QUAL-01: Conciseness**: 体言止め使用、冗長表現削除
- [ ] **QUAL-02: Actionable**: レビュー実施可能な情報量
- [ ] **QUAL-03: Specific**: 具体的チェック項目（抽象的表現回避）
- [ ] **QUAL-04: No Duplication**: カテゴリ間で重複なし

### 6. Balance (BAL: Balance)

適切なバランス:

- [ ] **BAL-01: Not Too Short**: 過度圧縮回避（75 行未満は要確認）
- [ ] **BAL-02: Not Too Long**: 冗長回避（400 行超は要確認）
- [ ] **BAL-03: Detail Level**: 全 20 カテゴリ詳細保持（Terraform の場合）
- [ ] **BAL-04: Optimal Range**: 150-200 行が理想範囲（言語により変動）

現在の行数:

- review-terraform.prompt.md: 162 行 ✅
- review-script.prompt.md: 150 行 ✅
- review-go.prompt.md: 164 行 ✅
- review-github-actions-workflow.prompt.md: 88 行 ✅（カテゴリ少ないため適切）

### 7. Consistency Across Files (CONS: Consistency)

- [ ] **CONS-01: ID Format**: 全ファイルで ID 形式統一（CAT-NN:）
- [ ] **CONS-02: Category Order**: 共通カテゴリの順序統一（G→ERR→SEC→PERF→TEST→DOC→DEP）
- [ ] **CONS-03: Structure**: Problem/Impact/Recommendation 構造統一
- [ ] **CONS-04: Naming**: 同種カテゴリ名統一（例: General=G, Security=SEC）

### 8. Usability (USE: Usability)

- [ ] **USE-01: Reviewable**: 記載内容でレビュー実施可能
- [ ] **USE-02: Clear Criteria**: 判断基準が明確
- [ ] **USE-03: Examples**: 具体例・パターン記載
- [ ] **USE-04: Tools**: 使用ツール・コマンド記載

### 9. Completeness (COMP: Completeness)

- [ ] **COMP-01: All Categories**: 対象言語の全カテゴリカバー
- [ ] **COMP-02: No Missing Items**: 重要チェック項目漏れなし
- [ ] **COMP-03: Cross-reference**: 対応 instructions.md 参照整合性
- [ ] **COMP-04: Tool Coverage**: 検証ツール全て言及

## Validation Process

### 1. カテゴリ構成確認

```bash
# 全promptファイルのカテゴリ抽出
for f in /workspace/.github/prompts/review-*.prompt.md; do
  echo "=== $(basename $f) ==="
  grep -E '^\*\*[A-Z]+-[0-9]+:' "$f" | sed 's/:.*$//' | cut -d'-' -f1 | sort -u
  echo
done
```

期待結果: 各ファイルで定義されたカテゴリセット表示

### 2. ID 重複確認

```bash
# 各ファイルのID重複チェック
for f in /workspace/.github/prompts/review-*.prompt.md; do
  echo "=== $(basename $f) ==="
  grep -oE '[A-Z]+-[0-9]+:' "$f" | sort | uniq -d
  echo "---"
done
```

期待結果: 重複 ID 表示なし（空出力）

### 3. 行数バランス確認

```bash
wc -l /workspace/.github/prompts/review-*.prompt.md
```

期待範囲:

- Terraform: 150-200 行（20 カテゴリ）
- Script/Go: 150-170 行（10 カテゴリ）
- GitHub Actions: 80-100 行（6 カテゴリ）

### 4. 構造統一確認

```bash
# Problem/Impact/Recommendation構造確認
grep -c "Problem:" /workspace/.github/prompts/review-terraform.prompt.md
grep -c "Impact:" /workspace/.github/prompts/review-terraform.prompt.md
grep -c "Recommendation:" /workspace/.github/prompts/review-terraform.prompt.md
```

期待結果: 3 つの出力が同数（全項目で 3 要素記載）

## Common Issues and Fixes

### Issue 1: 過度圧縮版（75 行）

**Problem**: "2-20. その他カテゴリ（簡略版）"1 行要約のみ
**Fix**: バランス版へ修正、全 20 カテゴリ詳細保持（150-200 行）

### Issue 2: リスト形式（ID なし）

**Problem**: 単純リスト形式、ID 付与なし
**Example**:

```markdown
- Check workflow syntax
- Verify permissions
```

**Fix**: ID 付き形式へ変更
**Example**:

```markdown
- **G-01: Workflow Syntax**
  - Problem: ...
  - Impact: ...
  - Recommendation: ...
```

### Issue 3: カテゴリ番号重複

**Problem**: 同じ番号が複数カテゴリで使用
**Example**: 4 番が 2 つ（ERR-04 と SEC-04 が両方 4 番）

**Fix**: 連番修正（1,2,3,4,5,6...）

### Issue 4: カテゴリ順序不統一

**Problem**: ファイル毎にカテゴリ順序が異なる
**Fix**: 統一順序適用（G→CODE→FUNC→ERR→SEC→PERF→TEST→DOC→DEP）

## Common Category Definitions

全 prompt ファイルで共通のカテゴリ:

| Category ID | Full Name      | Description        |
| ----------- | -------------- | ------------------ |
| G           | General        | 一般的品質・構成   |
| CODE        | Code Quality   | コード品質・可読性 |
| FUNC        | Functionality  | 機能性・設計       |
| ERR         | Error Handling | エラーハンドリング |
| SEC         | Security       | セキュリティ       |
| PERF        | Performance    | パフォーマンス     |
| TEST        | Testing        | テスト             |
| DOC         | Documentation  | ドキュメント       |
| DEP         | Dependencies   | 依存関係           |

## Final Verification

全チェック完了後:

1. **統一性**: 全 prompt ファイルで ID 形式・構造統一確認
2. **実用性**: 記載内容でレビュー実施可能確認
3. **完全性**: 対象言語の全カテゴリカバー確認
4. **バランス**: 過度圧縮・冗長でない適切な詳細度確認

## Reference Files

最良の参考例:

- **review-terraform.prompt.md** (162 行): 全 20 カテゴリ詳細、バランス版
- **review-go.prompt.md** (164 行): 10 カテゴリ、Problem/Impact/Recommendation 構造
- **review-script.prompt.md** (150 行): カテゴリ番号修正済み、統一構造

## Version History

- **v1 (過度圧縮版)**: 75 行、"その他カテゴリ（簡略版）"1 行のみ → レビュー不可能
- **v2 (バランス版)**: 150-164 行、全カテゴリ詳細保持 → 現行推奨バージョン
