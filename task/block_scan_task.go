package task

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/third/schedule"
)

func init() {
	schedule.Add("scan_block", func(ctx context.Context) {
		BlockScanRun(ctx)
	})
	logger.Debug("BlockScanRun has been scheduled")
}

func BlockScanRun(ctx context.Context) {
	logger.Debug("start", "nets", config.Task.Networks)
}
