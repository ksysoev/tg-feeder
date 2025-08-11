package cmd

import (
	"context"
	"fmt"

	"github.com/ksysoev/tg-feeder/pkg/bot"
	"github.com/ksysoev/tg-feeder/pkg/core"
	"github.com/ksysoev/tg-feeder/pkg/prov/someapi"
	"github.com/ksysoev/tg-feeder/pkg/repo/user"
	"github.com/redis/go-redis/v9"
)

// RunCommand initializes the logger, loads configuration, creates the core and API services,
// and starts the API service. It returns an error if any step fails.
func RunCommand(ctx context.Context, flags *cmdFlags) error {
	if err := initLogger(flags); err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}

	cfg, err := loadConfig(flags)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
	})

	someAPI := someapi.New(cfg.Provider.SomeAPI)
	userRepo := user.New(rdb)
	svc := core.New(userRepo, someAPI)

	tgBot, err := bot.New(&cfg.Bot, svc)
	if err != nil {
		return fmt.Errorf("failed to create API service: %w", err)
	}

	err = tgBot.Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to run API service: %w", err)
	}

	return nil
}
