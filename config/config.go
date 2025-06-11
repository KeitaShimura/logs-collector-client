// internal/config/config.go（クライアント用）
package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config は、環境変数から読み込まれるアプリケーション設定を保持する構造体
type Config struct {
	GRPCEndpoint  string `env:"GRPC_ENDPOINT"  envDefault:"localhost:50051"`
	RESTEndpoint  string `env:"REST_ENDPOINT"  envDefault:"http://localhost:8080"`
	DefaultLimit  int    `env:"DEFAULT_LIMIT"  envDefault:"10"`
	DefaultOffset int    `env:"DEFAULT_OFFSET" envDefault:"0"`
}

// LoadConfig は、環境変数を読み込んで Config を生成する
func LoadConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	return &cfg, nil
}
