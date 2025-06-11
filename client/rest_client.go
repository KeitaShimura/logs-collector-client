package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/KeitaShimura/logs-collector-client/model"
)

// ErrUnexpectedHTTPStatus は、想定外の HTTP ステータスが返された場合のエラー
var ErrUnexpectedHTTPStatus = errors.New("unexpected HTTP status")

// RESTClient は、ログ送信・取得を行う REST API クライアント
type RESTClient struct {
	Endpoint string // REST API のエンドポイント（例: http://localhost:8080）
}

// NewRESTClient は、指定されたエンドポイントで RESTClient を初期化する
func NewRESTClient(endpoint string) *RESTClient {
	return &RESTClient{Endpoint: endpoint}
}

// sendLogRequest は POST /api/logs に送信するリクエストボディの構造体
// Protobuf 仕様に合わせて log フィールドでネストされる
type sendLogRequest struct {
	Log *model.Log `json:"log"`
}

// SendLog はログデータを REST API に POST で送信する
func (c *RESTClient) SendLog(ctx context.Context, log *model.Log) error {
	// リクエストボディ構造に変換
	bodyStruct := sendLogRequest{Log: log}

	// JSON にシリアライズ
	body, err := json.Marshal(bodyStruct)
	if err != nil {
		return fmt.Errorf("failed to marshal log: %w", err)
	}

	// POST リクエスト作成
	url := c.Endpoint + "/api/logs"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// ヘッダー設定
	req.Header.Set("Content-Type", "application/json")

	// リクエスト送信
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer res.Body.Close()

	// ステータスコード確認
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %s", ErrUnexpectedHTTPStatus, res.Status)
	}

	return nil
}

// GetLogs は指定された条件に基づいてログを取得する
// クエリパラメータとして service, level, limit, offset を使用する
func (c *RESTClient) GetLogs(ctx context.Context, service, level string, limit, offset int32) ([]*model.Log, error) {
	// クエリパラメータ構築
	queryParams := url.Values{}
	queryParams.Set("service", service)
	queryParams.Set("level", level)
	queryParams.Set("limit", strconv.Itoa(int(limit)))
	queryParams.Set("offset", strconv.Itoa(int(offset)))

	// リクエスト URL を組み立て
	reqURL := fmt.Sprintf("%s/api/logs?%s", c.Endpoint, queryParams.Encode())

	// GET リクエスト作成
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// リクエスト送信
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer res.Body.Close()

	// レスポンスをデコードしてログ配列に変換
	var logs []*model.Log
	if err := json.NewDecoder(res.Body).Decode(&logs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return logs, nil
}
