package commands

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	aimos "github.com/BlockPILabs/aa-scan/internal/os"
	"github.com/BlockPILabs/aa-scan/parser"
	"github.com/BlockPILabs/aa-scan/task"
	"github.com/procyon-projects/chrono"
	"github.com/spf13/cobra"
	"log"
	"time"
)

// ScanCmd ...
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan block",
	Run: func(cmd *cobra.Command, args []string) {
		taskScheduler := chrono.NewDefaultTaskScheduler()
		// db start
		err := entity.Start(config)
		if err != nil {
			return
		}
		_, err = taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
			parser.ScanBlock()
		}, 5*time.Second)

		if err == nil {
			log.Print("Task: scan block has been scheduled successfully.")
		}

		aimos.TrapSignal(logger, func() {})

		task.InitTask()

		// Run forever.
		select {}
	},
}
