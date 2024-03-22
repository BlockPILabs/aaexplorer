package task

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/third/schedule"
)

func BlockScanStart(ctx context.Context, cfg *config.Config, logger log.Logger) {
	schedule.Add("scan_block", func(ctx context.Context) {
		BlockScanRun(ctx)
	})
	logger.Info("BlockScanRun has been scheduled")
}

func BlockScanRun(ctx context.Context) {

}
