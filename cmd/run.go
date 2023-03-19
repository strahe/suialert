package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/strahe/suialert/build"
	"github.com/strahe/suialert/processors"
	"go.uber.org/fx"
)

func (c *command) initRunCmd() {

	cmd := &cobra.Command{
		Use:   "run",
		Short: fmt.Sprintf("Start a %s node process", build.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := fx.New(
				fx.Provide(c.Config),
				fx.Provide(NewDB),
				fx.Provide(NewRuleService),
				fx.Provide(NewUserService),
				fx.Provide(NewPRCClient),
				fx.Provide(NewProcessor),
				fx.Provide(NewHandler),
				fx.Provide(NewBot),
				fx.Provide(NewEngine),
				fx.Invoke(func(cfg *processors.Processor) {}),
			)
			app.Run()
			return nil
		},
	}
	c.root.AddCommand(cmd)
}
