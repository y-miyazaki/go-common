---
mode: "agent"
description: "Go言語のコード品質・セキュリティ・パフォーマンス・ベストプラクティス準拠レビュー"
---

# Go Review Prompt

- Go 言語のベストプラクティスに精通したエキスパート
- Go アプリケーションのコード品質・セキュリティ・パフォーマンス・業界標準準拠レビュー
- scripts/go/check.sh を利用した自動化検証を前提とした高度なレビュー
- Serena MCP を利用する
- レビューコメントは日本語

## Review Guidelines (ID Based)

> 構造化チェックリスト。レビュー出力時は ID + ✅/❌ を利用。英語コメント統一は国際的な可読性と自動解析(検索/静的解析ツール)精度向上が目的。

### 1. Global / Base (G)

- G-01 Go 構文・型安全性・コンパイル妥当性 (go vet / golangci-lint 合格)
- G-02 パッケージ構成・import 文適切利用 (不要 import/循環依存/未使用変数検出)
- G-03 機密情報ハードコーディング禁止 (API キー/認証情報/DB 接続文字列など)
- G-04 Go モジュール依存管理適切 (go.mod/go.sum 整合性・脆弱性チェック)
- G-05 Error handling 明示的実装 (err != nil チェック・適切なエラーラップ)
- G-06 context.Context 適切利用 (タイムアウト・キャンセル処理)
- G-07 Goroutine・Channel 安全利用 (data race 無・適切な同期)
- G-08 関数シグネチャ設計適切 (引数・戻り値・エラー処理)
- G-09 標準ライブラリ活用 (適切なパッケージ選択・効率的な利用)
- G-10 ログ出力適切レベル (本番用レベル・機密情報非出力)

### 2. Code Standards (CODE)

- CODE-01 命名規則準拠 (snake_case/camelCase/PascalCase 適材適所)
- CODE-02 関数サイズ適切 (単一責任・可読性・50 行以下推奨)
- CODE-03 複雑度適切 (サイクロマティック複雑度・ネスト制限)
- CODE-04 DRY 原則準拠 (重複ロジック排除・共通化)
- CODE-05 インターフェース適切設計 (依存性注入・テスタビリティ)
- CODE-06 構造体適切設計 (JSON タグ・validation・埋め込み)
- CODE-07 定数・変数適切利用 (magic number 排除・const 利用)
- CODE-08 型アサーション・型変換安全 (nil チェック・パニック回避)
- CODE-09 defer 適切利用 (リソース解放・エラー時クリーンアップ)
- CODE-10 slice・map 適切操作 (nil チェック・境界チェック・容量管理)
- CODE-11 Go ファイルの宣言順序を遵守: const -> var -> type (interface → struct) -> func (constructor → methods → helpers)

### 3. Function Design (FUNC)

- FUNC-01 関数分割適切 (単一責任・適切なサイズ・可読性)
- FUNC-02 引数設計適切 (個数制限・構造体活用・デフォルト値)
- FUNC-03 戻り値設計適切 (named return・multiple return・error 位置)
- FUNC-04 純粋関数推奨 (副作用最小化・テスタビリティ)
- FUNC-05 レシーバー設計適切 (pointer/value・命名・一貫性)
- FUNC-06 メソッドセット設計 (interface 実装・組み込み・拡張性)
- FUNC-07 初期化関数適切 (init 関数・コンストラクタ・依存関係)
- FUNC-08 高次関数活用 (closure・callback・functional programming)
- FUNC-09 ジェネリクス適切利用 (型パラメータ・制約・パフォーマンス)
- FUNC-10 関数ドキュメント充実 (用途・引数・戻り値・例外・例)

### 4. Error Handling (ERR)

- ERR-01 エラー処理必須実装 (全ての error 戻り値チェック)
- ERR-02 エラーラップ適切 (pkg/errors または fmt.Errorf で context 追加)
- ERR-03 カスタムエラー適切定義 (errors.New・typed error・sentinel error)
- ERR-04 パニック回避・復旧 (recover 利用・graceful degradation)
- ERR-05 ログエラー情報適切 (スタックトレース・context・重要度)
- ERR-06 上位層エラー伝播適切 (error chain・原因保持)
- ERR-07 エラーハンドリング戦略 (fail-fast・graceful degradation・retry)
- ERR-08 外部依存エラー処理 (timeout・connection・service unavailable)
- ERR-09 バリデーションエラー (入力値検証・ビジネスルール・フォーマット)
- ERR-10 エラーメッセージセキュリティ (機密情報非露出・ユーザーフレンドリー)

### 5. Testing (TEST)

- TEST-01 単体テスト存在・充足性 (全関数・80%以上カバレッジ)
- TEST-02 テーブル駆動テスト利用 (複数ケース・境界値・異常系)
- TEST-03 testify 利用 (assert・require・mock・suite)
- TEST-04 モック適切利用 (外部依存・AWS SDK・HTTP client)
- TEST-05 テストヘルパー分離 (\*\_test.go ファイル・本番コード非混入)
- TEST-06 ベンチマークテスト (パフォーマンス重要箇所・メモリ使用量)
- TEST-07 競合状態テスト (go test -race・goroutine 安全性)
- TEST-08 統合テスト分離 (build tag・環境分離・E2E)
- TEST-09 テストデータ管理 (testdata/・golden file・fixture)
- TEST-10 テスト実行効率 (並列実行・キャッシュ・CI 最適化)

### 6. Security (SEC)

- SEC-01 機密情報管理適切 (環境変数・外部サービス・暗号化保存)
- SEC-02 入力値検証実装 (JSON validation・SQL injection 対策・XSS 対策)
- SEC-03 出力値サニタイズ (ログ・レスポンス・エラーメッセージ)
- SEC-04 暗号化適切実装 (TLS・AES・ハッシュ関数・salt)
- SEC-05 認証・認可実装 (JWT・OAuth・token validation・session)
- SEC-06 レート制限・DOS 対策 (request limit・throttling・circuit breaker)
- SEC-07 依存関係脆弱性管理 (govulncheck・定期更新・audit)
- SEC-08 ログセキュリティ (機密情報マスク・適切なレベル・retention)
- SEC-09 安全なデフォルト値 (設定・権限・暗号化・通信)
- SEC-10 セキュアコーディング (OWASP 準拠・脆弱性パターン回避)

### 7. Performance (PERF)

- PERF-01 メモリ使用量最適化 (slice capacity・map pre-allocation・GC pressure)
- PERF-02 CPU 使用量最適化 (アルゴリズム効率・並列処理・プロファイリング)
- PERF-03 I/O 最適化 (buffering・batching・connection pooling)
- PERF-04 データ構造選択適切 (slice/array/map/channel・用途別最適化)
- PERF-05 ガベージコレクション配慮 (allocation 削減・object pooling)
- PERF-06 文字列処理最適化 (strings.Builder・buffer・避けるべきパターン)
- PERF-07 並列処理最適化 (goroutine 数制限・worker pool・load balancing)
- PERF-08 キャッシュ戦略 (in-memory・LRU・TTL・invalidation)
- PERF-09 プロファイリング活用 (pprof・benchmark・bottleneck 特定)
- PERF-10 Hot path 最適化 (critical section・lock-free・atomic)

### 8. Architecture (ARCH)

- ARCH-01 レイヤー分離適切 (presentation/business/data・明確な責務)
- ARCH-02 依存性注入実装 (interface・DI container・testability)
- ARCH-03 ドメイン駆動設計準拠 (entity・value object・aggregate)
- ARCH-04 SOLID 原則準拠 (単一責任・開放閉鎖・依存性逆転)
- ARCH-05 パッケージ構成適切 (internal/・pkg/・cmd/・明確な責務)
- ARCH-06 設定管理統一 (config struct・環境別・validation)
- ARCH-07 ログ管理統一 (structured logging・correlation ID・level)
- ARCH-08 エラー管理統一 (error interface・wrapping・handling)
- ARCH-09 外部連携抽象化 (adapter pattern・port/adapter・circuit breaker)
- ARCH-10 モジュール設計 (cohesion・coupling・reusability)

### 9. Documentation (DOC)

- DOC-01 パッケージドキュメント存在 (package comment・用途説明)
- DOC-02 公開関数ドキュメント (godoc・引数説明・戻り値・例外)
- DOC-03 複雑ロジックコメント (アルゴリズム・ビジネスルール・TODO)
- DOC-04 構造体フィールドコメント (JSON tag・validation・用途)
- DOC-05 定数・変数説明 (用途・値の意味・単位・制約)
- DOC-06 英語コメント統一 (国際化・自動解析・検索効率)
- DOC-07 README.md 整備 (セットアップ・使用法・設定・FAQ)
- DOC-08 API 仕様書 (OpenAPI・request/response・example)
- DOC-09 運用ドキュメント (deploy・monitoring・troubleshooting)
- DOC-10 変更履歴管理 (CHANGELOG・breaking change・migration)

### 10. Dependencies (DEP)

- DEP-01 go.mod 適切管理 (semantic versioning・互換性・最新版)
- DEP-02 go.sum 整合性 (checksum verification・改竄検出)
- DEP-03 不要依存削除 (go mod tidy・unused package・bloat 削減)
- DEP-04 直接依存明示 (indirect 回避・explicit require)
- DEP-05 依存関係更新戦略 (定期更新・breaking change 対応・test)
- DEP-06 vendor 管理 (必要時のみ・サイズ制限・security)
- DEP-07 標準ライブラリ優先 (external dependency 削減・security・performance)
- DEP-08 AWS SDK バージョン管理 (v1/v2 選択・feature・deprecation)
- DEP-09 開発依存分離 (build tag・test only・tool)
- DEP-10 ライセンス互換性 (license check・legal compliance・audit)

### 11. Build & Deploy (BUILD)

- BUILD-01 go build 設定適切 (GOOS/GOARCH・ldflags・build constraint)
- BUILD-02 バイナリサイズ最適化 (-s -w・asset embedding・圧縮)
- BUILD-03 クロスコンパイル対応 (複数プラットフォーム・依存関係)
- BUILD-04 環境別ビルド (dev/staging/prod・feature flag・config)
- BUILD-05 CI/CD パイプライン (test→build→deploy・rollback・automation)
- BUILD-06 静的解析統合 (golangci-lint・govulncheck・security scan)
- BUILD-07 テスト自動化 (unit/integration/e2e・coverage・race)
- BUILD-08 アーティファクト管理 (versioning・repository・retention)
- BUILD-09 リリース戦略 (semantic versioning・changelog・migration)
- BUILD-10 品質ゲート (coverage threshold・lint pass・security check)

### 12. Configuration (CONF)

- CONF-01 環境変数設計 (naming convention・type safety・validation)
- CONF-02 設定構造体統一 (yaml/json・default value・required field)
- CONF-03 機密情報分離 (外部管理・暗号化・アクセス制御)
- CONF-04 環境別設定 (dev/staging/prod・inheritance・override)
- CONF-05 設定検証 (startup validation・type check・range check)
- CONF-06 設定変更影響 (hot reload・graceful restart・backward compatibility)
- CONF-07 設定管理戦略 (centralized・versioning・audit trail)
- CONF-08 ログ設定 (level・format・destination・rotation)
- CONF-09 feature flag 管理 (実験・段階リリース・緊急無効化)
- CONF-10 設定ドキュメント (説明・例・制約・migration)

### 13. Monitoring & Observability (MON)

- MON-01 ログ出力適切 (structured logging・correlation ID・searchable)
- MON-02 メトリクス収集 (business metric・technical metric・SLI)
- MON-03 トレーシング実装 (distributed tracing・span・annotation)
- MON-04 ヘルスチェック実装 (readiness・liveness・dependency check)
- MON-05 エラー監視 (error rate・alert threshold・escalation)
- MON-06 パフォーマンス監視 (latency・throughput・resource usage)
- MON-07 ビジネスメトリクス (conversion・user behavior・feature usage)
- MON-08 セキュリティ監視 (access log・anomaly detection・audit trail)
- MON-09 ダッシュボード (operation・business・real-time・historical)
- MON-10 アラート設計 (SLO・escalation・on-call・runbook)

### 14. External Integrations (EXT)

- EXT-01 HTTP クライアント設計 (timeout・retry・connection pool・TLS)
- EXT-02 データベース統合 (connection pool・transaction・prepared statement)
- EXT-03 メッセージキュー統合 (producer/consumer・error handling・backpressure)
- EXT-04 外部 API 統合 (circuit breaker・fallback・rate limiting)
- EXT-05 ファイルシステム操作 (path handling・permission・cleanup)
- EXT-06 ネットワーク通信 (protocol selection・security・monitoring)
- EXT-07 認証プロバイダー連携 (OAuth・SAML・LDAP・token 管理)
- EXT-08 ストレージサービス (object storage・caching・versioning)
- EXT-09 監視システム連携 (metrics export・log forwarding・alerting)
- EXT-10 外部ライブラリ管理 (dependency injection・version compatibility)

### 15. Data Processing (DATA)

- DATA-01 JSON 処理最適化 (streaming・validation・schema・encoding)
- DATA-02 CSV/XML 処理 (parser selection・memory efficiency・encoding)
- DATA-03 バイナリデータ処理 (encoding・compression・streaming)
- DATA-04 大容量データ処理 (chunk・stream・pagination・backpressure)
- DATA-05 データベース操作 (connection pool・transaction・prepared statement)
- DATA-06 キャッシュ戦略 (TTL・eviction・consistency・warming)
- DATA-07 データ変換・ETL (schema mapping・validation・error handling)
- DATA-08 ファイル処理 (multipart・streaming・temporary file・cleanup)
- DATA-09 暗号化・復号化 (key management・algorithm・performance)
- DATA-10 データ検証 (schema・business rule・sanitization)

### 16. Concurrency (CONC)

- CONC-01 goroutine 安全管理 (leak 防止・panic recovery・context)
- CONC-02 channel 適切設計 (buffered/unbuffered・close・select)
- CONC-03 mutex 適切利用 (deadlock 回避・granularity・RWMutex)
- CONC-04 sync package 利用 (WaitGroup・Once・Pool・atomic)
- CONC-05 競合状態回避 (data race・shared state・immutable)
- CONC-06 並列処理パターン (worker pool・pipeline・fan-out/fan-in)
- CONC-07 context 伝播 (timeout・cancellation・value・deadline)
- CONC-08 エラー処理・集約 (errgroup・error channel・partial failure)
- CONC-09 リソース制限 (goroutine limit・semaphore・rate limiting)
- CONC-10 テスト・デバッグ (race detector・deadlock detection・profiling)

### 17. HTTP & API (HTTP)

- HTTP-01 HTTP client 設定 (timeout・retry・connection pool・TLS)
- HTTP-02 リクエスト処理 (validation・sanitization・rate limiting)
- HTTP-03 レスポンス設計 (status code・header・body・cache)
- HTTP-04 CORS 設定 (origin・method・header・credential)
- HTTP-05 認証・認可 (JWT・OAuth・API key・session)
- HTTP-06 エラーレスポンス (RFC 7807・consistent format・logging)
- HTTP-07 ミドルウェア設計 (logging・recovery・compression・auth)
- HTTP-08 API versioning (URL/header/content type・backward compatibility)
- HTTP-09 OpenAPI 仕様 (documentation・validation・code generation)
- HTTP-10 テスト (httptest・mock server・contract testing)

### 18. Database (DB)

- DB-01 Connection pool 設定 (max connection・idle timeout・lifetime)
- DB-02 トランザクション管理 (ACID・isolation level・rollback)
- DB-03 プリペアードステートメント (SQL injection 対策・performance)
- DB-04 ORM 適切利用 (GORM・query optimization・N+1 problem)
- DB-05 スキーマ管理 (migration・versioning・rollback)
- DB-06 インデックス最適化 (query plan・covering index・composite)
- DB-07 データ型適切選択 (storage・performance・constraint)
- DB-08 デッドロック対策 (lock ordering・timeout・retry)
- DB-09 ヘルスチェック (ping・circuit breaker・graceful degradation)
- DB-10 監視・メトリクス (slow query・connection・error rate)

### 19. Maintenance (MAINT)

- MAINT-01 コードフォーマット統一 (gofmt・goimports・consistent style)
- MAINT-02 リファクタリング戦略 (incremental・test coverage・backward compatibility)
- MAINT-03 技術的負債管理 (TODO・FIXME・deprecation・metrics)
- MAINT-04 コードレビュー (checklist・automation・knowledge sharing)
- MAINT-05 継続的改善 (metric tracking・retrospective・automation)
- MAINT-06 ドキュメント更新 (code change・API・runbook)
- MAINT-07 依存関係更新 (security patch・feature・compatibility)
- MAINT-08 テスト保守 (flaky test・test data・coverage)
- MAINT-09 監視調整 (alert tuning・dashboard・SLO)
- MAINT-10 知識共有 (onboarding・best practice・troubleshooting)

### 20. Migration & Compatibility (MIG)

- MIG-01 Go バージョン互換性 (language feature・standard library・build)
- MIG-02 ライブラリマイグレーション (breaking change・deprecation・alternative)
- MIG-03 API バージョニング (backward compatibility・deprecation・sunset)
- MIG-04 データマイグレーション (schema change・data transformation・rollback)
- MIG-05 設定マイグレーション (environment variable・config file・default)
- MIG-06 モジュール分割・統合 (package restructure・import path・interface)
- MIG-07 プラットフォーム移行 (OS・architecture・runtime・container)
- MIG-08 アーキテクチャ変更 (monolith↔microservice・event-driven・sync/async)
- MIG-09 開発ツール移行 (build system・CI/CD・monitoring・testing)
- MIG-10 ロールバック戦略 (version control・feature flag・graceful degradation)

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
# Go Review Result

## Issues

None ✅ (No issues detected across all checklist items)
```

### ❌ Issues Found Example

```markdown
# Go Review Result

## Issues

1. ERR-01 エラー処理未実装

   - Problem: func processData() 内で os.Open() のエラーを無視
   - Impact: ファイル操作失敗時にパニック・予期しない動作
   - Recommendation: if err != nil { return fmt.Errorf("failed to open file: %w", err) }

2. SEC-02 入力値検証不足

   - Problem: JSON unmarshaling 後のバリデーションが未実装
   - Impact: 不正データによる SQL injection・XSS 脆弱性
   - Recommendation: validator パッケージまたは手動検証の追加

3. PERF-02 非効率なループ処理

   - Problem: nested loop 内で毎回 DB 接続を作成
   - Impact: パフォーマンス低下・Connection pool 枯渇
   - Recommendation: DB 接続を外側で作成し再利用

4. TEST-01 単体テスト不足

   - Problem: coverage 45% (目標 80%未満)
   - Impact: バグ発見率低下・リファクタリング困難
   - Recommendation: 主要関数とエラーケースのテスト追加

5. DOC-02 関数ドキュメント欠落
   - Problem: 公開関数 ProcessData() に godoc コメント無し
   - Impact: API 利用者の理解困難・保守性低下
   - Recommendation: // ProcessData processes user data and returns formatted results
```

### Compact Table Variant (Optional)

```markdown
| ID      | Status | Note                     |
| ------- | ------ | ------------------------ |
| ERR-01  | ❌     | Error handling missing   |
| SEC-02  | ❌     | Input validation missing |
| PERF-02 | ❌     | Inefficient loop         |
| TEST-01 | ❌     | Test coverage low        |
| DOC-02  | ❌     | Function doc missing     |
```

---
