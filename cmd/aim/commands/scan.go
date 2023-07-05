package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

// ScanCmd ...
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan block",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scan block")
	},
}
