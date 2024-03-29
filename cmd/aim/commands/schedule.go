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
	AddNetworkFlag(cmd)
	cmd.AddCommand(schedule.Commands()...)
	return cmd
}

func AddNetworkFlag(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		cmd.PersistentFlags().StringSlice("task.networks", config.Task.Networks, "--task.networks neta --task.networks netb")
	}
}
