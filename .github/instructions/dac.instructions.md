---
applyTo: "**/*.yaml,**/aws_architecture_diagram*.{yaml,yml,png}"
description: "AI Assistant Instructions for Diagram as Code (DAC)"
---

# AI Assistant Instructions for Diagram as Code (DAC)

**言語ポリシー**: ドキュメントは日本語、コード・コメントは英語。

このリポジトリでは AWS Architecture Diagram を Diagram as Code (DAC) 形式で管理します。

| File/Pattern                          | Purpose / Description        |
| ------------------------------------- | ---------------------------- |
| aws*architecture_diagram*\*.yaml      | DAC ソースファイル（環境別） |
| aws*architecture_diagram*\*.png       | 生成された構成図（環境別）   |
| scripts/terraform/aws*architecture*\* | DAC 関連の自動生成スクリプト |

## Standards

### Naming Conventions

| Component       | Rule                        | Example                            |
| --------------- | --------------------------- | ---------------------------------- |
| YAML ファイル名 | snake_case + env suffix     | aws_architecture_diagram_prd.yaml  |
| PNG ファイル名  | snake_case + env suffix     | aws_architecture_diagram_prd.png   |
| Resource ID     | PascalCase                  | VPC, ECSBackend, RDSCluster        |
| Stack/Group ID  | PascalCase + "Stack" suffix | EdgeServicesStack, DataLayerStack  |
| Resource Title  | Human-readable with spaces  | "ECS Backend", "Private Subnet 1a" |
| Link Labels     | Short descriptive text      | "HTTPS", "SQL", "ETL"              |

### Resource Type Mapping (Terraform → DAC)

Terraform コードから DAC リソースタイプへのマッピング：

| Terraform Resource                   | DAC Resource Type                         | Notes                                   |
| ------------------------------------ | ----------------------------------------- | --------------------------------------- |
| aws_lb (application)                 | AWS::ElasticLoadBalancingV2::LoadBalancer | Title に "(ALB)" を追加                 |
| aws_lb (network)                     | AWS::ElasticLoadBalancingV2::LoadBalancer | Title に "(NLB)" を追加                 |
| aws_ecs_service                      | AWS::ECS::Service                         | クラスタ名を Title に含める             |
| aws_rds_cluster (aurora-postgresql)  | AWS::RDS::DBCluster                       | Title に "Aurora PostgreSQL" など明記   |
| aws_rds_cluster (aurora-mysql)       | AWS::RDS::DBCluster                       | Title に "Aurora MySQL" など明記        |
| aws_redshift_cluster                 | AWS::Redshift::Cluster                    | -                                       |
| aws_lambda_function                  | AWS::Lambda::Function                     | VPC 配置の有無を Title に含める（任意） |
| aws_api_gateway_rest_api             | AWS::ApiGateway::RestApi                  | Backend/External を Title に明記        |
| aws_cloudfront_distribution          | AWS::CloudFront::Distribution             | -                                       |
| aws_wafv2_web_acl                    | AWS::WAFv2::WebACL                        | "WAF" または "WAF v2" を Title に       |
| aws_route53_zone                     | AWS::Route53::HostedZone                  | -                                       |
| aws_cognito_user_pool                | AWS::Cognito::UserPool                    | -                                       |
| aws_glue_job                         | AWS::Glue::Job                            | Title に "ETL" を含める（任意）         |
| aws_athena_workgroup                 | AWS::Athena::WorkGroup                    | -                                       |
| aws*quicksight*\*                    | AWS::QuickSight::Dashboard                | -                                       |
| aws_s3_bucket                        | AWS::S3::Bucket                           | バケット目的を Title に含める           |
| aws*ses*\*                           | AWS::SES::ConfigurationSet                | -                                       |
| aws_kinesis_firehose_delivery_stream | AWS::KinesisFirehose::DeliveryStream      | -                                       |
| aws_transfer_server                  | AWS::Transfer::Server                     | Title に "SFTP" など明記                |
| aws_nat_gateway                      | AWS::EC2::NatGateway                      | -                                       |
| aws_subnet (public)                  | AWS::EC2::Subnet                          | Title に "Public" を含める              |
| aws_subnet (private)                 | AWS::EC2::Subnet                          | Title に "Private" + 用途を含める       |
| aws_ec2_instance (bastion)           | AWS::EC2::Instance                        | Title に "Bastion" を含める             |

### DAC Best Practices

1. **階層構造の原則**

   - Canvas → Cloud → Region → VPC → Subnet → Resources の順に階層化
   - Public/Private/Data 等のサブネット層を明確に分離
   - VerticalStack/HorizontalStack を適切に使用して配置を制御

2. **リンクの描画規則**

   - ユーザーからのトラフィックフローは North → South
   - 水平方向の通信は East ↔ West
   - データベースへのクエリは orthogonal タイプを使用
   - 関連性の弱い接続は LineStyle: dashed を使用

3. **ラベルの付与**

   - 重要なトラフィックには Labels で通信プロトコルや目的を明記
   - 複数のリンクが同一リソースから出る場合は TargetLeft/TargetRight で重複回避

4. **タイトルの命名**

   - 環境名（prd/stg/dev）は Region の Title に含める
   - リソースの役割が明確になるよう具体的な Title を設定
   - AZ 情報は Subnet の Title に含める（例: "Private App 1a"）

   ### Route53 と WAF の可視化ルール

   - Route53 はグローバルサービスとして扱い、Region 配下ではなく Global セクションに記載すること。
   - `User` ノード（訪問者）から直接 `Route53` へ線を引かないこと。ユーザーは DNS を解決した後、直接 CloudFront や API Gateway などの実際の配信先へ接続する表現とする。
   - ドメイン名（例: example.com や CDN の alternate domain）は Route53 ノードではなく、実際にそのドメインが割り当てられているリソース（CloudFront, API Gateway, Load Balancer など）の Title に記載すること。

   - WAF の配置ルール:
     - WAF は、その WAF を適用するリソースの「右横」に配置する。配置方法は 2 通り許容する：
       1. 右側に単一の WAF を置き、複数リソースから破線で結ぶ（従来方式）。
       2. リソースごとに WAF ノードを右横に個別に配置する方式（推奨）。単一 WAF から線を出すと図が見にくくなるため、各リソースの横に WAF を置くと読みやすくなるケースが多い。
     - WAF ノードから線を出す表記は避ける（WAF -> リソース）。必ずリソース -> WAF の向きで破線を描くこと。
     - 同一 WAF を複数のリソースが参照する場合は、どちらの方式でも可だが、視認性を優先し、必要に応じてリソースごとの WAF ノードを使うこと。

## Guidelines

### Documentation and Comments

- YAML ファイル内にコメントで構成の意図を記載する
- 複雑なリンク構造には説明コメントを追加する
- 環境固有の設定は TODO コメントで明記する

### Code Modification Guidelines

DAC ファイル修正時は以下の手順で検証する：

```bash
# 1. YAML 構文チェック
yamllint aws_architecture_diagram_prd.yaml

# 2. DAC から PNG を生成
awsdac -d aws_architecture_diagram_prd.yaml -o aws_architecture_diagram_prd.png

# 3. 生成された画像を確認
file aws_architecture_diagram_prd.png
```

### MCP Tool Usage (awsdac-mcp-server)

awsdac-mcp-server の MCP ツールを使用して DAC を生成・編集する：

```bash
# 1. DAC フォーマット情報の取得
mcp_awsdac-mcp-se_getDiagramAsCodeFormat

# 2. YAML から PNG を生成（ファイル保存）
mcp_awsdac-mcp-se_generateDiagramToFile with yamlContent="..." and outputFilePath="/workspace/diagram.png"

# 3. Base64 形式で取得（表示可能なクライアント向け）
mcp_awsdac-mcp-se_generateDiagram with yamlContent="..." and outputFormat="png"
```

### Terraform コードから DAC への変換手順

1. **リソース収集**

   - `grep_search` で Terraform ファイルから主要リソースを抽出
   - tfvars ファイルで環境固有の設定値を確認
   - load_balancer_type 等の詳細設定を必ず確認

2. **構造設計**

   - VPC 階層: Public/Private App/Private Data の 3 層構成
   - 各 AZ を HorizontalStack でグループ化
   - サービス種別ごとに VerticalStack で配置

3. **リンク設計**

   - User → Route53 → CloudFront → API Gateway → ALB → ECS → RDS の主要フロー
   - East-West トラフィック（ECS ↔ S3, Glue ↔ Redshift 等）
   - 管理アクセス（Bastion 等）

4. **検証**
   - 生成された PNG で視覚的にレイアウトを確認
   - リンクの交差や重複がないか確認
   - すべてのリソースが適切な親子関係にあるか確認

### Common Pitfalls

- **ALB と NLB の混同**: tfvars で `load_balancer_type` を必ず確認
- **VPC 内外の配置ミス**: Lambda や API Gateway の VPC 配置有無を確認
- **リンク先の存在確認**: Links セクションで存在しないリソース ID を参照しない
- **Stack の Source/Target 指定**: VerticalStack/HorizontalStack はリンクの Source/Target に使用不可

### Version Control

- YAML ファイルと PNG ファイルは両方とも Git で管理
- 環境ごとに別ファイルとして管理（prd/stg/dev）
- 大きな構成変更は PR でレビューを受ける
- コミットメッセージには変更内容を明記（例: "Add Kinesis Firehose to data flow"）

## Testing and Validation

### Validation Checklist

- [ ] YAML 構文エラーがないこと
- [ ] すべての Resources が Canvas から到達可能
- [ ] Links の Source/Target がすべて存在するリソース
- [ ] Title が人間にとって理解しやすい
- [ ] 環境名（prd/stg/dev）が Region Title に含まれている
- [ ] load_balancer_type が正しく反映されている（ALB/NLB）
- [ ] VPC/Subnet の階層構造が正しい

### Generation Test

```bash
# Test diagram generation
awsdac -d aws_architecture_diagram_prd.yaml -o test_output.png && file test_output.png && rm -f test_output.png
```

## Security Guidelines

- 構成図に機密情報（IP アドレス、アカウント ID 等）を含めない
- 公開リポジトリにプッシュする前に機密情報の有無を確認
- Title には一般的な名称を使用（具体的なドメイン名等は避ける）
