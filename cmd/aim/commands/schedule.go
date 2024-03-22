package commands

import (
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/memo"
	"github.com/BlockPILabs/aaexplorer/third/moralis"
	"github.com/BlockPILabs/aaexplorer/third/schedule"
	"github.com/spf13/cobra"
)

func NewExecCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec",
		Short: "execute schedule task",
		//RunE: func(cmd *cobra.Command, args []string) error {
		//	fmt.Println("start")
		//	return nil
		//},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// db start
			err := entity.Start(logger.With("lib", "ent"), config)
			if err != nil {
				logger.Error("error in db start", "err", err)
				return err
			}

			err = memo.Start(logger.With("lib", "memo"), config)
			if err != nil {
				logger.Error("error in memo start", "err", err)
				return err
			}
			moralis.SetConfig(config)
			return nil
		},
	}
	cmd.AddCommand(schedule.Commands()...)
	return cmd
}
