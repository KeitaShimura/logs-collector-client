# logs-collector-client

このリポジトリは、[logs-collector-api](https://github.com/KeitaShimura/logs-collector-api) に対してログの送信・取得を行う **動作確認用 CLI クライアント** です。gRPC / REST 両方に対応しています。

## 使用方法

### ログ送信・取得（gRPC）

```bash
make grpc-send  # ログを gRPC 経由で送信
make grpc-get   # ログを gRPC 経由で取得
```

### ログ送信・取得（REST）

```bash
make rest-send  # ログを REST 経由で送信
make rest-get   # ログを REST 経由で取得
```

## 環境変数（`.env`）

| 変数名           | 説明               | デフォルト値            |
| ---------------- | ------------------ | ----------------------- |
| `GRPC_ENDPOINT`  | gRPC の接続先      | `localhost:50051`       |
| `REST_ENDPOINT`  | REST API の接続先  | `http://localhost:8080` |
| `DEFAULT_LIMIT`  | ログ取得件数の上限 | `10`                    |
| `DEFAULT_OFFSET` | ログ取得の開始位置 | `0`                     |

## ディレクトリ構成

```
logs-collector-client/
├── .github/
│   └── workflows/
│       ├── ci.yaml
│       └── release.yaml
├── .gitignore
├── .golangci.yaml
├── .goreleaser.yaml
├── Makefile
├── README.md
├── go.mod
├── go.sum
├── cmd/
│   └── main.go
└── internal/
    ├── client/
    │   ├── client.go
    │   ├── grpc_client.go
    │   └── rest_client.go
    ├── config/
    │   └── config.go
    ├── logger/
    │   ├── logger.go
    │   └── logger_test.go
    └── model/
        └── log.go
```

## 対応 API

### REST API

| メソッド | パス      | 説明         | 主なクエリ/ボディ                      |
| -------- | --------- | ------------ | -------------------------------------- |
| POST     | /api/logs | ログ送信     | body: { log: Log }                     |
| GET      | /api/logs | ログ一覧取得 | service, level, limit, offset (クエリ) |

- **POST /api/logs**

  - ログデータ（JSON, `{ log: ... }`）を送信
  - 成功時: 200 OK

- **GET /api/logs**
  - クエリパラメータでサービス名・レベル・件数・オフセット指定
  - レスポンス: ログ配列（JSON）

### gRPC API

| サービス名         | メソッド | 説明         | 主なリクエスト/レスポンス        |
| ------------------ | -------- | ------------ | -------------------------------- |
| logs.v1.LogService | SendLog  | ログ送信     | SendLogRequest / SendLogResponse |
| logs.v1.LogService | GetLogs  | ログ一覧取得 | GetLogsRequest / GetLogsResponse |

- **SendLog (logs.v1.LogService)**

  - ログデータを送信
  - リクエスト: `SendLogRequest { log: Log }`
  - レスポンス: `SendLogResponse`

- **GetLogs (logs.v1.LogService)**
  - サービス名・レベル・件数・オフセット等でログを取得
  - リクエスト: `GetLogsRequest { service, level, limit, offset, ... }`
  - レスポンス: `GetLogsResponse { logs: [Log] }`

---
