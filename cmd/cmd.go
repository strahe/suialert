package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/strahe/suialert/build"
	"go.uber.org/zap"
)

const (
	optionDebug        = "debug"
	optionRpcEndpoints = "rpc"
)

const (
	DevNetRpcUrl  = "wss://fullnode.devnet.sui.io"
	TestnetRpcUrl = "wss://fullnode.testnet.sui.io"
)

type command struct {
	root    *cobra.Command
	vp      *viper.Viper
	cfgFile string
}

type option func(*command)

func newCommand(opts ...option) (c *command, err error) {

	c = &command{
		root: &cobra.Command{
			Use:           strings.ToLower(build.AppName),
			Short:         "Simple And Fast Object Storage",
			SilenceErrors: true,
			SilenceUsage:  true,
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				return c.initConfig()
			},
		},
		vp: viper.New(),
	}

	for _, o := range opts {
		o(c)
	}

	c.initGlobalFlags()
	c.initRunCmd()
	c.initMigrateCmd()
	c.initVersionCmd()

	return c, nil
}

func (c *command) initConfig() error {

	// info default
	configName := "config"
	if c.cfgFile != "" {
		c.vp.SetConfigFile(c.cfgFile)
	} else {
		c.vp.SetConfigName(configName)
	}
	c.vp.AddConfigPath(".")
	c.vp.SetEnvPrefix(strings.ToLower(build.AppName))
	c.vp.AutomaticEnv()
	c.vp.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	c.vp.SetConfigType("toml")
	if err := c.vp.ReadInConfig(); err != nil {
		var e viper.ConfigFileNotFoundError
		if !errors.As(err, &e) {
			zap.S().Errorf("failed to read config file: %s", err)
			return err
		}
		zap.S().Warnf("config file not found: %s", c.cfgFile)
	}

	var logger *zap.Logger
	var err error
	if c.vp.GetBool("debug") {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		return err
	}

	defer logger.Sync() // nolint: errcheck
	zap.ReplaceGlobals(logger)
	return nil
}

func Execute() (err error) {
	c, err := newCommand()
	if err != nil {
		return err
	}
	return c.Execute()
}

func (c *command) Execute() (err error) {
	return c.root.Execute()
}

func (c *command) initGlobalFlags() {
	globalFlags := c.root.PersistentFlags()
	globalFlags.StringVarP(&c.cfgFile, "config", "c", "", fmt.Sprintf("config file (default is $HOME/.%s.toml)", "config"))
}
