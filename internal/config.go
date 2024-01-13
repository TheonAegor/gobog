package internal

import (
	"context"
	"github.com/TheonAegor/go-framework/pkg/config"
	"github.com/TheonAegor/go-framework/pkg/config/envConfig"
	"log"
)

type Config struct {
	ApiId    int32  `json:"api_id" env:"api_id"`
	ApiHash  string `json:"api_hash" env:"api_hash"`
	BotId    int64  `json:"bot_id" env:"bot_id"`
	LogLevel string `json:"log_level" env:"log_level"`
}

func InitConfig(ctx context.Context) *Config {
	cfg := Config{}

	opts := append([]config.Config(nil),
		config.NewConfig(config.Struct(&cfg)),
		envConfig.NewConfig(config.Struct(&cfg)),
	)

	if err := config.Load(ctx, opts...); err != nil {
		log.Fatal(ctx, "failed to load config: %+v", err)
	}

	return &cfg
}
