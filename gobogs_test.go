package main

import (
	"context"
	"github.com/TheonAegor/go-framework/pkg/config"
	"github.com/TheonAegor/go-framework/pkg/config/envConfig"
	"github.com/TheonAegor/go-framework/pkg/logger"
	"github.com/TheonAegor/gobog/internal"
	"github.com/TheonAegor/gobog/pkg/tg"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	apiId := os.Getenv("api_id")
	apiIdInt, err := strconv.ParseInt(apiId, 10, 32)
	if err != nil {
		t.Fatal()
	}

	ctx := context.Background()
	cfg := internal.Config{
		ApiHash: os.Getenv("api_hash"),
		ApiId:   int32(apiIdInt),
	}

	opts := append([]config.Config(nil),
		config.NewConfig(config.Struct(&cfg)),
		envConfig.NewConfig(config.Struct(&cfg)),
	)

	if err := config.Load(ctx, opts...); err != nil {
		log.Fatal(ctx, "failed to load config: %+v", err)
	}

	log := logger.NewLogger(&logger.LoggerConfig{LogLevel: cfg.LogLevel})

	tgClientOptions := []tg.Option{
		tg.WithLog(log),
	}

	tgClient, err := tg.NewTgClient(cfg.ApiId, cfg.ApiHash, tgClientOptions...)
	if err != nil {
		t.Fatal("cannot create tg client", "error", err)
	}

	handlerOptions := []internal.Option{
		internal.WithLog(log),
	}

	h := internal.NewHandler(tgClient, cfg, handlerOptions...)

	suite := godog.TestSuite{
		ScenarioInitializer: func(scenarioContext *godog.ScenarioContext) {
			InitializeScenario(scenarioContext, h)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
