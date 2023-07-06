package main

import (
	"github.com/BlockPILabs/aa-scan/cmd/aim/commands"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/config/cli"
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
	cmd := cli.PrepareBaseCmd(rootCmd, "AIM", os.ExpandEnv(filepath.Join("$HOME", config.DefaultHomeDir)))
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
