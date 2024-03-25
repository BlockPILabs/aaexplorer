package task

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/third/schedule"
	"github.com/ethereum/go-ethereum/ethclient"
	"time"
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
	log.Context(ctx).Info("start", "nets", config.Task.Networks)
	tx, err := entity.Client(ctx)
	if err != nil {
		log.Context(ctx).Error("db connect error", "err", err)
		return
	}

	networks, err := getTaskNetworks(ctx, tx)
	if err != nil {
		log.Context(ctx).Error("error in getTaskNetworks", "err", err)
		return
	}

	for _, network := range networks {
		ctx, logger := log.With(ctx, "network", network.ID, "networkName", network.Name)

		dialContext, err := ethclient.DialContext(ctx, network.HTTPRPC)
		if err != nil {
			logger.Error("error in getTaskNetworks", "err", err)
			continue
		}
		blockNumber, err := dialContext.BlockNumber(ctx)
		if err != nil {
			logger.Error("error in BlockNumber", "err", err)
			continue
		}

		networkTx, err := entity.Client(ctx, network.ID)
		if err != nil {
			logger.Error("error in network db connect", "err", err)
			return
		}

		result, err := networkTx.ExecContext(ctx, `insert into block_sync(block_num, create_time, update_time) select generate_series(max(block_num) , $1 ) , $2 , $2 from block_sync on conflict do nothing`, blockNumber, time.Now())
		if err != nil {
			logger.Error("error in block_sync generate_series", "err", err)
			continue
		}
		rowsAffected, _ := result.RowsAffected()

		logger.Info("block sync result", "rowsAffected", rowsAffected)

	}

}

func BlockScanRun(ctx context.Context) {
	logger.Debug("start", "nets", config.Task.Networks)

}
