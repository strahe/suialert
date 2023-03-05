package cmd

import (
	"context"

	"github.com/strahe/suialert/bots"
	"github.com/strahe/suialert/bots/discord"
	"github.com/strahe/suialert/client"
	"github.com/strahe/suialert/config"
	"github.com/strahe/suialert/handlers"
	"github.com/strahe/suialert/model"
	"github.com/strahe/suialert/processors"
	"github.com/strahe/suialert/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error
	if cfg.Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	return logger, err
}

func (c *command) Config() (*config.Config, error) {
	cfg := config.DefaultConfig
	if err := c.vp.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func NewStorage(lc fx.Lifecycle, cfg *config.Config) (model.Storage, error) {
	db, err := storage.NewDatabase(cfg.Database.Postgres)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return db.Connect(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return db.Close(ctx)
		},
	})
	return db, nil
}

func NewBot(lc fx.Lifecycle, cfg *config.Config) (bots.Bot, error) {
	bot, err := discord.NewDiscord(cfg.Bots.Discord)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return bot.Run()
		},
		OnStop: func(context.Context) error {
			return bot.Close()
		},
	})
	return bot, nil
}

func NewHandler(lc fx.Lifecycle, bot bots.Bot, store model.Storage) *handlers.SubHandler {
	hd := handlers.NewSubHandler(bot, store)
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return hd.Close()
		},
	})
	return hd
}

func NewProcessor(lc fx.Lifecycle, cfg *config.Config, rpcClient *client.Client, hd *handlers.SubHandler) (*processors.Processor, error) {
	p, err := processors.NewProcessor(cfg, rpcClient, hd)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return p.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return p.Close(ctx)
		},
	})
	return p, nil
}

func NewPRCClient(lc fx.Lifecycle, cfg *config.Config, hd *handlers.SubHandler) (*client.Client, error) {
	c, closer, err := client.NewClient(context.TODO(), cfg.Sui.Endpoint, hd)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			closer()
			return nil
		},
	})
	return c, nil
}
