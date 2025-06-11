package client

import (
	"context"

	"github.com/KeitaShimura/logs-collector-client/internal/model"
)

// Client はログの送信および取得を行うためのインターフェース
type Client interface {
	SendLog(ctx context.Context, log *model.Log) error
	GetLogs(ctx context.Context, service string, level string, limit int, offset int) ([]model.Log, error)
}
