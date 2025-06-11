package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/KeitaShimura/logs-collector-client/model"
	pb "github.com/KeitaShimura/logs-collector-protos/go/logs/v1"
)

// GRPCClient は gRPC 経由でログの送信・取得を行うクライアント
type GRPCClient struct {
	conn   *grpc.ClientConn    // gRPC接続
	client pb.LogServiceClient // gRPCクライアント（LogService）
}

// NewGRPCClient は指定されたエンドポイントに接続する GRPCClient を作成する
func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	client := pb.NewLogServiceClient(conn)

	return &GRPCClient{conn: conn, client: client}, nil
}

// Close は gRPC 接続をクローズする
func (c *GRPCClient) Close() {
	c.conn.Close()
}

// SendLog はログを gRPC API 経由で送信する
func (c *GRPCClient) SendLog(ctx context.Context, log *model.Log) error {
	// 文字列の timestamp を protobuf の Timestamp 型に変換
	timestamp, err := parseTimestamp(log.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}

	// gRPC のリクエストを構築
	req := &pb.SendLogRequest{
		Log: &pb.Log{
			Id:        log.ID,
			TraceId:   log.TraceID,
			Timestamp: timestamp,
			Level:     log.Level,
			Service:   log.Service,
			Message:   log.Message,
			Metadata:  log.Metadata,
		},
	}

	// リクエスト送信
	if _, err := c.client.SendLog(ctx, req); err != nil {
		return fmt.Errorf("failed to send log via gRPC: %w", err)
	}

	return nil
}

// GetLogs は指定された条件でログを gRPC API 経由で取得する
func (c *GRPCClient) GetLogs(ctx context.Context, service, level string, limit, offset int32) ([]*model.Log, error) {
	// リクエスト構築
	req := &pb.GetLogsRequest{
		Service:   StringPtr(service),
		Level:     StringPtr(level),
		Limit:     limit,
		Offset:    offset,
		StartTime: nil,
		EndTime:   nil,
	}

	// リクエスト送信
	resp, err := c.client.GetLogs(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs via gRPC: %w", err)
	}

	// 結果を model.Log にマッピング
	protoLogs := resp.GetLogs()
	logs := make([]*model.Log, 0, len(protoLogs))

	for _, protoLog := range protoLogs {
		logs = append(logs, &model.Log{
			ID:        protoLog.GetId(),
			TraceID:   protoLog.GetTraceId(),
			Timestamp: formatTimestamp(protoLog.GetTimestamp()),
			Level:     protoLog.GetLevel(),
			Service:   protoLog.GetService(),
			Message:   protoLog.GetMessage(),
			Metadata:  protoLog.GetMetadata(),
		})
	}

	return logs, nil
}

// StringPtr は string 値を *string に変換するヘルパー関数
func StringPtr(s string) *string {
	return &s
}

// parseTimestamp は RFC3339 フォーマットの文字列を protobuf の Timestamp に変換する
func parseTimestamp(ts string) (*timestamppb.Timestamp, error) {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp format: %w", err)
	}

	return timestamppb.New(t), nil
}

// formatTimestamp は protobuf の Timestamp を RFC3339 文字列に変換する
func formatTimestamp(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}

	return ts.AsTime().Format(time.RFC3339)
}
