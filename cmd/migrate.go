package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/strahe/suialert/config"
	"github.com/strahe/suialert/storage"
)

func (c *command) initMigrateCmd() {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Manage the schema version installed in a database.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()

			cfg := config.DefaultConfig
			if err := c.vp.Unmarshal(&cfg); err != nil {
				return fmt.Errorf("error reading config: %w", err)
			}

			db, err := storage.NewDatabase(cfg.Database.Postgres)
			if err != nil {
				return fmt.Errorf("error creating to database: %w", err)
			}

			return db.MigrateSchema(ctx)
		},
	}
	c.setNodeFlags(cmd)
	c.root.AddCommand(cmd)
}
