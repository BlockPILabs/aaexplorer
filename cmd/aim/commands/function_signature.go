package commands

import (
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/memo"
	"github.com/BlockPILabs/aaexplorer/service"
	"github.com/spf13/cobra"
)

var FunctionSignatureCmd = &cobra.Command{
	Use:    "function_signature",
	Short:  "load solid function signature",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = entity.Start(logger.With("lib", "ent"), config)
		if err != nil {
			return
		}

		err = memo.Start(logger.With("lib", "memo"), config)
		if err != nil {
			return
		}
		service.ScanSignature(cmd.Context())
		return nil
	},
}

func init() {
	RootCmd.AddCommand(FunctionSignatureCmd)
}
