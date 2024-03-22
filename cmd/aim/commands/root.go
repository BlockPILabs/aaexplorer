package commands

import (
	"fmt"
	cfg "github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/config/cli"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/task"
	"github.com/BlockPILabs/aaexplorer/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var logger = log.L()
var config = cfg.DefaultConfig()

func init() {
	task.SetConfig(config)
	task.SetLogger(logger)
	log.SetDefaultLogger(logger)
}

// ParseConfig retrieves the default environment configuration,
// sets up the aim root and ensures that the root exists
func ParseConfig(cmd *cobra.Command) (*cfg.Config, error) {
	conf := cfg.DefaultConfig()

	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}

	var home string
	if os.Getenv("AIM_HOME") != "" {
		home = os.Getenv("AIM_HOME")
	} else {
		home, err = cmd.Flags().GetString(cli.HomeFlag)
		if err != nil {
			return nil, err
		}
	}

	conf.RootDir = home

	conf.SetRoot(conf.RootDir)
	cfg.EnsureRoot(conf.RootDir)
	if err := conf.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("error in config file: %v", err)
	}
	if warnings := conf.CheckDeprecated(); len(warnings) > 0 {
		for _, warning := range warnings {
			logger.Info("deprecated usage found in configuration file", "usage", warning)
		}
	}
	return conf, nil
}

var RootCmd = &cobra.Command{
	Use: version.Name,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if cmd.Name() == VersionCmd.Name() {
			return nil
		}
		config, err = ParseConfig(cmd)
		if err != nil {
			return err
		}

		if config.LogFormat == cfg.LogFormatJSON {
			logger = log.NewTMJSONLogger(log.NewSyncWriter(os.Stdout))
		}

		logger, err = cli.ParseLogLevel(config.LogLevel, logger, cfg.DefaultLogLevel)
		if err != nil {
			return err
		}

		if viper.GetBool(cli.TraceFlag) {
			logger = log.NewTracingLogger(logger)
		}

		logger = logger.With("module", "main")

		// Set default logger
		log.SetDefaultLogger(logger)
		task.SetLogger(logger)
		task.SetConfig(config)
		cmd.SetContext(
			log.WithContext(cmd.Context(), logger),
		)

		return nil
	},
	//Run: func(cmd *cobra.Command, args []string) {
	//	logger.Warn("warn")
	//	logger.Error("err")
	//	logger.Info("info")
	//	logger.Debug("debug")
	//},
}
