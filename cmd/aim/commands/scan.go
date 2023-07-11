package commands

import (
	"context"
	aimos "github.com/BlockPILabs/aa-scan/internal/os"
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

		_, err := taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
			//parser.ScanBlock()
			logger.Info("scan block end")
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
