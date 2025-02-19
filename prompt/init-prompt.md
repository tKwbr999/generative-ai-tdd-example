# 初期化プロンプト

## 目的
このプロンプトは、プロジェクトの初期セットアップを自動化し、開発環境を整備するためのガイドラインです。

## タスク実行順序

1.  **プロジェクト構造の理解**
    *   `rule.md`の解析
        *   開発フローの理解
        *   ドキュメント構成の確認
        *   各アクターの役割の把握
    *   `architecture.md`の解析
        *   システム構成の理解
        *   ディレクトリ構造の確認
        *   技術スタックの確認

2.  **サンプル機能仕様の作成**
    *   ユーザー管理機能の仕様書作成
        *   ユースケース図 (PlantUML形式)
        *   API仕様（OpenAPI 3.0形式）
        *   データモデル (JSON Schema形式)
    *   以下のCRUD操作を含む
        *   ユーザー作成
        *   ユーザー取得
        *   ユーザー更新
        *   ユーザー削除

3.  **テストコード作成**
    *   テーブル駆動テスト手法を採用
    *   単体テスト
        *   エンティティのテスト
        *   ユースケースのテスト
        *   リポジトリのモックを使用したテスト
    *   統合テスト
        *   APIエンドポイントのテスト (httptestパッケージ)
        *   データベース操作のテスト
    *   テストカバレッジ80%以上を目標 (go test -cover)

4.  **開発環境セットアップ**
    *   必要なツールとバージョン
    *   依存パッケージの管理
    *   データベース設定の環境変数設定
        *   環境変数のファイル名は`env.[環境]`という形式にする (`env.local`, `env.dev`)
        *   URIはコード内で環境変数から作成する。
        *   `env.sample`を作成し、環境変数のサンプルを記述する。

5.  **サンプルコード生成**
    *   TDDに基づいた実装
    *   以下のコンポーネントを含む
        *   エンティティ定義
        *   リポジトリインターフェース
        *   ユースケース実装
        *   HTTPハンドラー

6.  **Docker環境構築**
    *   アプリケーションコンテナ (Dockerfile)
    *   MongoDBコンテナ (公式イメージ)
    *   開発用ネットワーク設定 (docker-compose.yml)
    *   ホットリロード対応 (air)
    *   データベースURIは環境変数`ENVIRONMENT`によって切り替えます。
    *   各種ミドルウェアなどのuser名、パスワード、そのほかセキュリティ的にgithubに乗せてはいけないものは全て環境変数にまとめるようにする。
    *   `env.sample`はgitに載せて良い。
    *   docker-compose.ymlの設定例
        ```yaml
        version: "3.9"
        services:
          app:
            build: .
            ports:
              - "8080:8080"
            environment:
              ENVIRONMENT: dev
            depends_on:
              - mongo
          mongo:
            image: mongo:6.0
            ports:
              - "27017:27017"
            volumes:
              - mongo_data:/data/db
        volumes:
          mongo_data:
        ```

## 成果物
*   プロジェクトの基本構造
*   機能仕様書
*   テストコード
*   サンプルCRUD機能の実装
*   Docker開発環境
*   README.md
