# アーキテクチャ設計書

## システム概要
### 技術スタック
- バックエンド: Go(latest)
- Webフレームワーク: atreugo(latest)
- データベース: postgreSQL(latest)
- API形式: REST
- 開発手法: クリーンアーキテクチャ + TDD

### 非機能要件
- レスポンスタイム: 100ms以内
- スケーラビリティ: 水平スケーリング対応
- 可用性: 99.9%以上

## アーキテクチャ詳細
### レイヤー構成
1. プレゼンテーション層 (interface/handler)
   - HTTPリクエスト/レスポンスの処理
   - バリデーション
   - DTOの変換

2. ユースケース層 (usecase)
   - ビジネスロジックの実装
   - トランザクション管理
   - ドメインサービスの利用

3. ドメイン層 (domain)
   - エンティティ定義
   - ドメインロジック
   - リポジトリインターフェース

4. インフラストラクチャ層 (infrastructure)
   - データベース実装
   - 外部APIクライアント
   - キャッシュ実装

## ディレクトリ構造
```
.
├── cmd/
│   └── api/                    # APIサーバーのエントリーポイント
│       └── main.go
├── internal/
│   ├── domain/                 # ドメインモデル
│   │   ├── entity/            # エンティティ定義
│   │   ├── repository/        # リポジトリインターフェース
│   │   └── service/           # ドメインサービス
│   ├── usecase/               # ユースケース実装
│   ├── infrastructure/        # インフラストラクチャ層
│   │   ├── persistence/       # データベース実装
│   │   └── auth/             # 認証関連実装
│   └── interface/
│       ├── handler/           # HTTPハンドラー
│       └── middleware/        # ミドルウェア
├── pkg/                       # 外部公開可能なパッケージ
│   ├── logger/               # ログユーティリティ
│   └── validator/            # バリデーションユーティリティ
├── docs/
│   ├── api/                  # API仕様書
│   ├── architecture/         # アーキテクチャドキュメント
│   └── deployment/           # デプロイメントガイド
├── tests/
│   ├── e2e/                 # E2Eテスト
│   └── integration/         # 統合テスト
├── deployments/             # デプロイメント関連ファイル
│   ├── docker/
│   └── kubernetes/
├── scripts/                 # ユーティリティスクリプト
├── .github/
│   └── workflows/          # GitHub Actions設定
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```
## 依存関係の方向
- 外層から内層への依存のみ許可
- 依存性注入を活用
- インターフェースによる疎結合

## エラーハンドリング
- カスタムエラー型の定義
- エラーラッピング
- 適切なHTTPステータスコードの使用

## セキュリティ
- 入力バリデーション (SQLインジェクション対策、XSS対策)
- CORS設定
- レート制限
- JWTによる認証
- HTTPSによる通信

## 生成コード品質
* Goのコーディング規約（`gofmt`, `golint`準拠）
    * `gofmt`によるコード整形を必須とする。
    * `golint`による静的解析を行い、指摘事項を修正する。
* Go style guide準拠
    * Effective Go: [https://go.dev/doc/effective_go](https://go.dev/doc/effective_go)
    * Code Review Comments: [https://go.dev/wiki/CodeReviewComments](https://go.dev/wiki/CodeReviewComments)
* エラーハンドリングの方針
    * エラーは可能な限り詳細な情報を付与して返す。
