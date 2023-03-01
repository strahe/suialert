package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/strahe/suialert/storage"

	"github.com/strahe/suialert/bots/discord"

	"github.com/strahe/suialert/handlers"

	"github.com/strahe/suialert/client"

	"github.com/spf13/cobra"
	"github.com/strahe/suialert/build"
	"github.com/strahe/suialert/config"
	"github.com/strahe/suialert/processors"
	"go.uber.org/zap"
)

func (c *command) initRunCmd() {

	cmd := &cobra.Command{
		Use:   "run",
		Short: fmt.Sprintf("Start a %s node process", build.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			zap.S().Infof("Starting %s process", build.AppName)
			ctx, done := context.WithCancel(cmd.Context())
			defer done()

			cfg := config.DefaultConfig
			if err := c.vp.Unmarshal(&cfg); err != nil {
				return fmt.Errorf("error reading config: %w", err)
			}
			// todo:
			db, err := storage.NewDatabase(cfg.Database.Postgres)
			if err != nil {
				return fmt.Errorf("error creating database: %w", err)
			}
			if err := db.Connect(ctx); err != nil {
				return fmt.Errorf("error connecting to database: %w", err)
			}
			defer db.Close() //nolint:errcheck

			// todo: check if the bot is enabled
			bot, err := discord.New(cfg.Bots.Discord)
			if err != nil {
				return fmt.Errorf("error creating discord bot: %w", err)
			}
			if err := bot.Run(); err != nil {
				return fmt.Errorf("error starting discord bot: %w", err)
			}

			defer func() {
				if err := bot.Close(); err != nil {
					zap.S().Errorf("error closeing discord bot: %v", err)
				}
			}()

			hd := handlers.NewEthSubHandler(bot, db)
			defer func() {
				if err := hd.Close(); err != nil {
					zap.S().Errorf("error closing handlers: %v", err)
				}
			}()

			rpcClient, closer, err := client.NewClient(ctx, c.vp.GetString(optionRpcEndpoints), hd)
			if err != nil {
				return fmt.Errorf("failed to create rpc client: %w", err)
			}
			defer func() {
				zap.S().Infof("closing RPC client")
				closer()
			}()

			p, err := processors.NewProcessor(&cfg, rpcClient, hd)
			if err != nil {
				return err
			}
			defer func() {
				zap.S().Infof("closing processor")
				if err := p.Close(); err != nil {
					zap.S().Errorf("error closing processor: %v", err)
				}
			}()
			sigChan := make(chan os.Signal, 2)
			signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
			go func() {
				if err := p.Run(ctx); err != nil {
					zap.S().Errorf("error: %v", err)
					sigChan <- syscall.SIGTERM
				}
			}()

			<-sigChan
			zap.S().Infof("Shutting down %s process", build.AppName)
			return nil
		},
	}
	c.setNodeFlags(cmd)
	c.root.AddCommand(cmd)
}

func (c *command) setNodeFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(optionDebug, false, "enable debug model")
	c.vp.BindPFlag(optionDebug, cmd.Flag(optionDebug))

	cmd.Flags().String(optionRpcEndpoints, DevNetRpcUrl, "rpc endpoints of sui to connect to")
	c.vp.BindPFlag(optionRpcEndpoints, cmd.Flag(optionRpcEndpoints))
}
