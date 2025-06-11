package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/KeitaShimura/logs-collector-client/client"
	"github.com/KeitaShimura/logs-collector-client/config"
	"github.com/KeitaShimura/logs-collector-client/logger"
	"github.com/KeitaShimura/logs-collector-client/model"
)

// 共通エラー定義
var (
	ErrInvalidAction = errors.New("invalid action")
	ErrIntOverflow   = errors.New("value overflows int32")
)

// os.Args の最低必要引数数（コマンド + アクション）
const minArgs = 2

func main() {
	// run() の返り値（ステータスコード）を exit code として返す
	os.Exit(run())
}

// run は CLI のメイン処理。引数に応じて対応する処理関数を呼び出す
func run() int {
	// 一時的なINFOレベルロガーを初期化
	logger := logger.NewLogger(logger.WithLevel(logger.LevelInfo))
	ctx := context.Background()

	// 引数数チェック
	if len(os.Args) < minArgs {
		logger.Error("usage: go run cmd/main.go [grpc-send|grpc-get|rest-send|rest-get]", nil)

		return 1
	}

	action := os.Args[1]

	// 入力されたアクションに応じた処理へルーティング
	switch action {
	case "grpc-send":
		return runGRPCSend(ctx, logger)
	case "grpc-get":
		return runGRPCGet(ctx, logger)
	case "rest-send":
		return runRESTSend(ctx, logger)
	case "rest-get":
		return runRESTGet(ctx, logger)
	default:
		// 不正なアクションが指定された場合のエラーハンドリング
		logger.Error("unknown action", fmt.Errorf("%w: %s", ErrInvalidAction, action))

		return 1
	}
}

// runGRPCSend は gRPC API を通じてログを送信する
func runGRPCSend(ctx context.Context, logger logger.Logger) int {
	// 環境変数から設定情報を読み込む
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("failed to load config", err)

		return 1
	}

	// gRPC クライアントを初期化
	client, err := client.NewGRPCClient(cfg.GRPCEndpoint)
	if err != nil {
		logger.Error("failed to connect to gRPC", err)

		return 1
	}
	defer client.Close()

	// テスト用ログを生成
	log := &model.Log{
		ID:        uuid.NewString(),
		TraceID:   uuid.NewString(),
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   "test-service",
		Level:     "INFO",
		Message:   "Hello, log world!",
		Metadata:  map[string]string{"env": "dev"},
	}

	// gRPC API へログ送信を試みる
	if err := client.SendLog(ctx, log); err != nil {
		logger.Error("SendLog failed", err,
			"id", log.ID,
			"trace_id", log.TraceID,
			"timestamp", log.Timestamp,
			"service", log.Service,
			"level", log.Level,
			"message", log.Message,
			"metadata", log.Metadata,
		)

		return 1
	}

	// 成功時は構造化ログで出力
	logger.Info("SendLog succeeded",
		"id", log.ID,
		"trace_id", log.TraceID,
		"timestamp", log.Timestamp,
		"service", log.Service,
		"level", log.Level,
		"message", log.Message,
		"metadata", log.Metadata,
	)

	return 0
}

// runGRPCGet は gRPC API からログ一覧を取得し、ログ出力する
func runGRPCGet(ctx context.Context, logger logger.Logger) int {
	// 環境変数から設定情報を読み込む
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("failed to load config", err)

		return 1
	}

	// gRPC クライアントを初期化
	client, err := client.NewGRPCClient(cfg.GRPCEndpoint)
	if err != nil {
		logger.Error("failed to connect to gRPC", err)

		return 1
	}
	defer client.Close()

	// DefaultLimit を int32 に変換（オーバーフローがないか安全にチェック）
	limit, err := safeIntToInt32(cfg.DefaultLimit)
	if err != nil {
		logger.Error("invalid limit", err)

		return 1
	}

	// DefaultOffset を int32 に変換（オーバーフローがないか安全にチェック）
	offset, err := safeIntToInt32(cfg.DefaultOffset)
	if err != nil {
		logger.Error("invalid offset", err)

		return 1
	}

	// ログ取得（サービス・レベルでフィルタリング）
	logs, err := client.GetLogs(ctx, "test-service", "INFO", limit, offset)
	if err != nil {
		logger.Error("GetLogs failed", err)

		return 1
	}

	// 結果を構造化ログで出力
	logger.Info("GetLogs succeeded", "count", len(logs))

	for _, log := range logs {
		logger.Info("Log entry", "id", log.ID, "message", log.Message)
	}

	return 0
}

// runRESTSend は REST API を用いてログを送信する
func runRESTSend(ctx context.Context, logger logger.Logger) int {
	// 環境変数から設定情報を読み込む
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("failed to load config", err)

		return 1
	}

	// REST クライアント初期化
	client := client.NewRESTClient(cfg.RESTEndpoint)

	log := &model.Log{
		ID:        uuid.NewString(),
		TraceID:   uuid.NewString(),
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   "test-service",
		Level:     "INFO",
		Message:   "Hello from REST!",
		Metadata:  map[string]string{"env": "dev"},
	}

	// REST API へログを送信
	if err := client.SendLog(ctx, log); err != nil {
		logger.Error("SendLog (REST) failed", err,
			"id", log.ID,
			"trace_id", log.TraceID,
			"timestamp", log.Timestamp,
			"service", log.Service,
			"level", log.Level,
			"message", log.Message,
			"metadata", log.Metadata,
		)

		return 1
	}

	// 成功時は構造化ログで出力
	logger.Info("SendLog (REST) succeeded",
		"id", log.ID,
		"trace_id", log.TraceID,
		"timestamp", log.Timestamp,
		"service", log.Service,
		"level", log.Level,
		"message", log.Message,
		"metadata", log.Metadata,
	)

	return 0
}

// runRESTGet は REST API からログ一覧を取得し、出力する
func runRESTGet(ctx context.Context, logger logger.Logger) int {
	// 環境変数から設定情報を読み込む
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("failed to load config", err)

		return 1
	}

	// REST クライアント初期化
	client := client.NewRESTClient(cfg.RESTEndpoint)

	// DefaultLimit を int32 に変換（オーバーフローがないか安全にチェック）
	limit, err := safeIntToInt32(cfg.DefaultLimit)
	if err != nil {
		logger.Error("invalid limit", err)

		return 1
	}

	// DefaultOffset を int32 に変換（オーバーフローがないか安全にチェック）
	offset, err := safeIntToInt32(cfg.DefaultOffset)
	if err != nil {
		logger.Error("invalid offset", err)

		return 1
	}

	// ログ取得（サービス・レベルでフィルタリング）
	logs, err := client.GetLogs(ctx, "test-service", "INFO", limit, offset)
	if err != nil {
		logger.Error("GetLogs (REST) failed", err)

		return 1
	}

	// 結果を構造化ログで出力
	logger.Info("GetLogs (REST) succeeded", "count", len(logs))

	for _, log := range logs {
		logger.Info("Log entry", "id", log.ID, "message", log.Message)
	}

	return 0
}

// safeIntToInt32 は int 値を int32 に安全に変換する関数
func safeIntToInt32(n int) (int32, error) {
	if n > int(^int32(0)) || n < int(^int32(0)+1) {
		return 0, fmt.Errorf("%w: %d", ErrIntOverflow, n)
	}

	return int32(n), nil
}
