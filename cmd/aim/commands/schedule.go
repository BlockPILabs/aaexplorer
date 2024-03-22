package commands

import (
	"fmt"
	"github.com/BlockPILabs/aaexplorer/third/schedule"
	"github.com/spf13/cobra"
)

func NewScheduleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "Run schedule",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("start")
			return nil
		},
	}
	cmd.AddCommand(schedule.Commands()...)
	return cmd
}
