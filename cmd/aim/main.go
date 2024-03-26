package main

import (
	"github.com/BlockPILabs/aaexplorer/cmd/aim/commands"
	"github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/config/cli"
	"os"
	"path/filepath"
)

func main() {
	rootCmd := commands.RootCmd
	rootCmd.AddCommand(
		commands.VersionCmd,
		commands.ScanCmd,
		commands.MigrateCmd,
	)
	// Create & start
	rootCmd.AddCommand(commands.NewStartCmd())
	cli.AddCommand(rootCmd, commands.NewExecCmd())
	cmd := cli.PrepareBaseCmd(rootCmd, "AIM", os.ExpandEnv(filepath.Join("$HOME", config.DefaultHomeDir)))
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
