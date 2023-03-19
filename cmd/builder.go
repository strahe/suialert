package cmd

import (
	"context"
	"fmt"

	"github.com/strahe/suialert/rule"
	"github.com/strahe/suialert/service"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/strahe/suialert/bots"
	"github.com/strahe/suialert/bots/discord"
	"github.com/strahe/suialert/client"
	"github.com/strahe/suialert/config"
	"github.com/strahe/suialert/handlers"
	"github.com/strahe/suialert/model"
	"github.com/strahe/suialert/processors"
	"go.uber.org/fx"
)

func (c *command) Config() (*config.Config, error) {
	cfg := config.DefaultConfig
	if err := c.vp.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func NewEngine(lc fx.Lifecycle, ruleService *service.RuleService) (*rule.Engine, error) {
	eng, err := rule.NewEngine(ruleService)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return eng.LoadRules(ctx)
		},
	})
	return eng, nil
}

func NewBot(lc fx.Lifecycle, cfg *config.Config, userService *service.UserService, ruleService *service.RuleService) (bots.Bot, error) {
	bot, err := discord.NewDiscord(cfg.Bots.Discord, userService, ruleService)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return bot.Run(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return bot.Close(ctx)
		},
	})
	return bot, nil
}

func NewHandler(lc fx.Lifecycle, bot bots.Bot, db *gorm.DB, eng *rule.Engine) *handlers.SubHandler {
	hd := handlers.NewSubHandler(bot, db, eng)
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

func NewRuleService(db *gorm.DB) *service.RuleService {
	return service.NewRuleService(db)
}

func NewUserService(db *gorm.DB) *service.UserService {
	return service.NewUserService(db)
}

func NewDB(lc fx.Lifecycle, cfg *config.Config) (*gorm.DB, error) {
	var dia gorm.Dialector
	switch cfg.Database.Driver {
	case "sqlite3":
		dia = sqlite.Open(cfg.Database.DSN)
	case "mysql":
		dia = mysql.Open(cfg.Database.DSN)
	case "postgres":
		dia = postgres.Open(cfg.Database.DSN)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}
	gc := &gorm.Config{
		CreateBatchSize: 100,
	}
	if cfg.Debug {
		gc.Logger = logger.Default.LogMode(logger.Warn) // logger.Default.LogMode(logger.Info)
	} else {
		gc.Logger = logger.Default.LogMode(logger.Error)
	}
	db, err := gorm.Open(dia, gc)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return model.Migration(db)
		},
		OnStop: func(context.Context) error {
			d, err := db.DB()
			if err != nil {
				return err
			}
			// maybe not necessary
			return d.Close()
		},
	})

	return db, nil
}
