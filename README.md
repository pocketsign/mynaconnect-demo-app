# PocketSign MynaConnect のデモアプリケーション

## 概要

このデモアプリケーションは、PocketSign MynaConnectを使用して、本人確認情報と自己情報の取得を実装したWebアプリケーションです。

## 技術スタック

- **Go 1.24**: バックエンド言語
- **ConnectRPC**: PocketSign MynaConnectとの通信に使用
  - 本デモアプリではConnectRPCを使用したPocketSign MynaConnectが提供するAPIクライアントを使用しますが、一般的なHTTPのJSON POST APIを使用した通信も可能です。詳しくはドキュメントをご覧ください。

## 環境変数

| 変数名                       | デフォルト値 | 説明          |
|---------------------------|--------|-------------|
| `MYNA_CONNECT_API_TOKEN`  | (必須)   | テナントAPIトークン |
| `MYNA_CONNECT_SERVICE_ID` | (必須)   | APIサービスID   |

## セットアップ

### 事前準備

PocketSign Platformでテナント及び「自己情報取得API」のAPIサービスを作成してください。
また、テナントのAPIトークンを発行し、手元に控えてください。

APIサービスでは以下の設定を行なってください。
- マイナコネクトからのコールバック先として許可するURIのプレフィックス: `http://localhost:3000/`
- 取得する本人確認の種類: `基本4情報`
- マイナポータルAPI上でのカード読み取り前に本人確認情報を検証を実施: `する`

### 1. 環境変数の設定

```bash
export MYNA_CONNECT_API_TOKEN="your_tenant_api_token_here" # PocketSign Platformから取得したテナントAPIトークンを設定
export MYNA_CONNECT_SERVICE_ID="your_api_service_id_here"  # PocketSign Platformで作成した自己情報取得APIのAPIサービスIDを設定
```

### 2. 依存関係のインストール

```bash
go mod download
```

### 3. アプリケーションの起動

```bash
go run .
```

### 4. ブラウザでアクセス

```
http://localhost:3000
```

## API フロー

### 1. セッション作成

- `/start` エンドポイントで `CreateSelfPersonalDataRequestSession` を呼び出し
- 今日の日付を照会条件として設定
- セッションIDを保存
- PocketSign MynaConnectにリダイレクト
  - ポケットサイン開発環境ではPocketSign Platformに事前登録してある「本人確認情報テストデータ」の選択画面が表示されます

### 2. コールバック処理

- `/callback` エンドポイントでPocketSign MynaConnectからのコールバックを受信
- `GetSelfPersonalDataRequestStatus` でステータスを確認
- ステータスに応じて適切なページにリダイレクト：
  - `SUCCESS` → `/result`（結果表示）
  - `NEED_TO_VERIFY_USER` → `/verify`（本人確認）
  - `ERROR`/`EXPIRED`/`PENDING` → エラーページ

### 3. 本人確認

- `/verify` エンドポイントで `GetUserIdentity` を呼び出し
- 本人確認情報をJSON形式で表示
- `/verify-identity` エンドポイントで `SubmitUserIdentityVerificationResult` を呼び出し
- PocketSign MynaConnectにリダイレクト
  - ポケットサイン開発環境ではPocketSign Platformに事前登録してある「自己情報取得APIテストデータ」の選択画面が表示されます


### 4. 結果表示

- `/result` エンドポイントで `GetSelfPersonalDataRequestResult` を呼び出し
- 生データとパース済みデータを表示

## ファイル構成

```
demoapp/
├── main.go             # メイン関数とグローバル変数定義
├── handlers.go         # HTTPハンドラー関数
├── session.go          # セッション管理関数
├── go.mod              # Go モジュール定義
├── go.sum              # 依存関係のハッシュ
├── README.md           # このファイル
└── templates/          # HTMLテンプレートと静的ファイル
```

## 注意事項

- **これはデモアプリケーションです**。本番環境では適切なセキュリティ対策を実装してください
- セッション管理、エラーハンドリング、認証などは簡略化されています
