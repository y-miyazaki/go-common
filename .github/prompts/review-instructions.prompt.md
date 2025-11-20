# Review Instructions for \*.instructions.md Files

## Objective

`.github/instructions/*.instructions.md`ファイルの品質・統一性・実用性確保。

## Review Checklist

### 1. File Structure (G: General)

- [ ] **G-01: Front Matter**: `applyTo`と`description`が正確に定義されている
- [ ] **G-02: Language Policy**: "**言語ポリシー**: ドキュメント日本語、コード・コメント英語"記載
- [ ] **G-03: Title**: ファイル用途が明確なタイトル

### 2. Chapter Structure (STRUCT: Structure)

必須章構成（順序厳守）:

```markdown
## Standards

## Guidelines

## Testing and Validation

## Security Guidelines
```

チェック項目:

- [ ] **STRUCT-01**: 4 つの必須章が全て存在
- [ ] **STRUCT-02**: 章順序が統一（Standards → Guidelines → Testing and Validation → Security Guidelines）
- [ ] **STRUCT-03**: 見出しレベル適切（## 主要章、### サブセクション）

### 3. Standards Chapter (STD: Standards)

- [ ] **STD-01: Naming Conventions**: 表形式でコンポーネント別命名規則記載
- [ ] **STD-02: Tool Standards**: 対象ツール固有の標準規約記載（例: `terraform fmt`, `go fmt`）
- [ ] **STD-03: Consistency**: 他ファイルと記載レベル統一

### 4. Guidelines Chapter (GUIDE: Guidelines)

必須サブセクション:

- [ ] **GUIDE-01: Documentation and Comments**: コメント・ドキュメント規約
- [ ] **GUIDE-02: Code Modification Guidelines**: 修正時手順・検証方法
- [ ] **GUIDE-03: Tool Usage**: MCP Tool Usage または言語固有ガイドライン
- [ ] **GUIDE-04: Error Handling**: エラーハンドリングパターン（該当言語のみ）

### 5. Testing and Validation Chapter (TEST: Testing)

必須サブセクション:

- [ ] **TEST-01: Validation Commands**: 必須検証コマンド記載（実行例付き）
- [ ] **TEST-02: Command Count**: 検証コマンド 3 項目以上（他ファイルと同レベル）
- [ ] **TEST-03: Code Block**: 実行例がコードブロック形式（```bash）
- [ ] **TEST-04: Validation Items**: 検証項目リスト記載

検証コマンド例:

- **script**: `bash -n`, `shellcheck`, `validate_all_scripts.sh`
- **go**: `go fmt`, `go vet`, `golangci-lint`, `go test`, `govulncheck`（8 項目）
- **terraform**: `terraform fmt`, `terraform validate`, `tflint`, `trivy config`
- **github-actions**: `actionlint`, `ghalint run`, `disable-checkout-persist-credentials`, `ghatm`
- **markdown**: `markdownlint`, `markdown-link-check`
- **dac**: YAML 構文チェック、図生成テスト

### 6. Security Guidelines Chapter (SEC: Security)

- [ ] **SEC-01: Security Items**: セキュリティ必須項目記載
- [ ] **SEC-02: Secrets Management**: 機密情報管理方法記載
- [ ] **SEC-03: Best Practices**: 具体的セキュリティ対策記載
- [ ] **SEC-04: Examples**: YAML/コード例あり（該当する場合）

### 7. Content Quality (QUAL: Quality)

- [ ] **QUAL-01: Conciseness**: 体言止め使用、冗長表現削除
- [ ] **QUAL-02: Practical Examples**: 実用的な例・コード片記載
- [ ] **QUAL-03: No Redundancy**: 重複内容なし
- [ ] **QUAL-04: Token Efficiency**: 不要な大規模コード例削除済み

### 8. Consistency Across Files (CONS: Consistency)

- [ ] **CONS-01: Chapter Order**: 全ファイルで章順序統一
- [ ] **CONS-02: Section Names**: 同種セクション名統一（例: "Code Modification Guidelines"）
- [ ] **CONS-03: Detail Level**: 記載詳細度が他ファイルと同等
- [ ] **CONS-04: Format**: 表・リスト形式統一

### 9. Completeness (COMP: Completeness)

- [ ] **COMP-01: All Required Sections**: 必須セクション全て存在
- [ ] **COMP-02: No Missing Commands**: 検証コマンド全て記載
- [ ] **COMP-03: Tool Coverage**: aqua.yaml の関連ツール全てカバー
- [ ] **COMP-04: Real Commands**: 実行可能なコマンド記載

## Validation Process

### 1. 章構成確認

```bash
# 全ファイルの章構成抽出
for f in /workspace/.github/instructions/*.instructions.md; do
  echo "=== $(basename $f) ==="
  grep -E '^## ' "$f"
  echo
done
```

期待結果: 全ファイルで 4 章統一（Standards/Guidelines/Testing and Validation/Security Guidelines）

### 2. 行数バランス確認

```bash
wc -l /workspace/.github/instructions/*.instructions.md
```

期待範囲:

- 最小: 70 行程度（terraform: 73 行）
- 最大: 230 行程度（go: 222 行、特殊ケース）
- 標準: 100-180 行

### 3. 検証コマンド網羅性確認

各ファイルの"Testing and Validation"章で検証コマンド数確認:

- 最低 3 項目以上
- 実行例付き
- コードブロック形式

### 4. セキュリティガイドライン確認

全ファイルで"Security Guidelines"章存在確認:

```bash
grep -l "## Security Guidelines" /workspace/.github/instructions/*.instructions.md | wc -l
```

期待: 6 ファイル全て

## Common Issues and Fixes

### Issue 1: 章順序不統一

**Problem**: Testing and Validation が Guidelines 内にある
**Fix**: 独立章として抽出、Security Guidelines の前に配置

### Issue 2: 検証コマンド不足

**Problem**: 検証コマンドが 1-2 項目のみ
**Fix**: aqua.yaml 確認、関連ツール全て追加（最低 3 項目）

### Issue 3: Security Guidelines 章なし

**Problem**: セキュリティ章が存在しない
**Fix**: 機密情報管理・ベストプラクティス記載の章追加

### Issue 4: 記載レベル不統一

**Problem**: 他ファイルより詳細度が低い
**Fix**: 他ファイル参照、同等の詳細度に拡充

## Final Verification

全チェック完了後:

1. **統一性**: 全 6 ファイルで章構成・順序統一確認
2. **実用性**: 各検証コマンドが実行可能確認
3. **完全性**: 必須セクション全て存在確認
4. **バランス**: 行数・詳細度が他ファイルと同等確認

## Reference Files

最良の参考例:

- **go.instructions.md** (222 行): 最も詳細、Testing 章 8 項目
- **github-actions-workflow.instructions.md** (180 行): 拡充後の良例
- **script.instructions.md** (106 行): 標準的なバランス
