package commands

import (
	"encoding/json"
	"fmt"
	"github.com/BlockPILabs/aa-scan/version"
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

// VersionCmd ...
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info",
	Run: func(cmd *cobra.Command, args []string) {
		aimVersion := version.Version
		if version.GitCommitHash != "" {
			aimVersion += "+" + version.GitCommitHash
		}

		if verbose {
			values, _ := json.MarshalIndent(struct {
				Version string `json:"version"`
			}{
				Version: aimVersion,
			}, "", "  ")
			fmt.Println(string(values))
		} else {
			fmt.Println(aimVersion)
		}
	},
}

func init() {
	VersionCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show protocol and library versions")
}
