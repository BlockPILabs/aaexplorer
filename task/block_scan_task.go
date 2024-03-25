package task

import (
	"context"
	"fmt"
	"github.com/BlockPILabs/aaexplorer/internal/dao"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/third/schedule"
)

func init() {
	schedule.Add("block_sync", func(ctx context.Context) {
		BlockSyncRun(ctx)
	}).ScheduleWithCron("*/1 * * * * *")
	schedule.Add("scan_block", func(ctx context.Context) {
		BlockScanRun(ctx)
	}).ScheduleWithCron("*/1 * * * * *")
	logger.Debug("BlockScanRun has been scheduled")
}

func BlockSyncRun(ctx context.Context) {
	logger.Info("start", "nets", config.Task.Networks)
	tx, err := entity.Client(ctx)
	if err != nil {
		return
	}
	networks, err := dao.NetworkDao.GetNetworks(ctx, tx)
	if err != nil {
		return
	}
	fmt.Sprintln(networks)
}

func BlockScanRun(ctx context.Context) {
	logger.Debug("start", "nets", config.Task.Networks)
}
