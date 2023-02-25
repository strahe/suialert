package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/strahe/suialert/build"
)

func (c *command) initVersionCmd() {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of saas",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(build.UserVersion())
		},
	}
	c.setNodeFlags(cmd)
	c.root.AddCommand(cmd)
}
