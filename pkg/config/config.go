package config

import (
	"context"
	"github.com/TheonAegor/go-framework/pkg/config"
	"github.com/TheonAegor/go-framework/pkg/config/envConfig"
	"log"
)

var (
	DefaultConfig Config
)

type Config struct {
	ApiId   int32  `json:"api_id" env:"api_id"`
	ApiHash string `json:"api_hash" env:"api_hash"`
}

func init() {
	ctx := context.Background()
	cfg := Config{}

	opts := append([]config.Config(nil),
		config.NewConfig(config.Struct(&cfg)),
		envConfig.NewConfig(config.Struct(&cfg)),
	)

	if err := config.Load(ctx, opts...); err != nil {
		log.Fatal(ctx, "failed to load config: %+v", err)
	}

	DefaultConfig = cfg
}
