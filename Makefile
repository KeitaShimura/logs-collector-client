ALL_PACKAGES := ./...         # 全てのGoパッケージ
CMD_PACKAGES := ./cmd/main.go # メイン実行ファイル

.PHONY: all init format lint lint-fix test-fast test cover grpc-send grpc-get rest-send rest-get

# すべての主要なタスクを順に実行
all: format lint test

# プロジェクトの初期セットアップ（依存関係の整備）
init:
	go mod tidy

# コードフォーマットとインポート整理
format:
	go fmt ${ALL_PACKAGES}
	gci write -s standard -s default -s "prefix(github.com/KeitaShimura)" $(shell find . -name '*.go')
	gofumpt -w .

# Lint チェック（静的解析）
lint:
	golangci-lint run

# Lint の自動修正
lint-fix:
	golangci-lint run --fix

# 簡易テスト（レース検出・カバレッジなし、詳細出力あり）
test-fast:
	go test -v $(ALL_PACKAGES)

# テスト実行（詳細出力 + カバレッジ + レース検出）
test:
	go test -v -race -cover $(ALL_PACKAGES)

# カバレッジ付きテストの実行
cover:
	mkdir -p coverage
	go test -cover $(ALL_PACKAGES) -coverprofile=coverage/cover.out

# gRPCでログ送信
grpc-send:
	go run cmd/main.go grpc-send

# gRPCでログ取得
grpc-get:
	go run cmd/main.go grpc-get

# RESTでログ送信
rest-send:
	go run cmd/main.go rest-send

# RESTでログ取得
rest-get:
	go run cmd/main.go rest-get