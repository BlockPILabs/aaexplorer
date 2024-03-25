package task

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"errors"
	"fmt"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/blocksync"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/BlockPILabs/aaexplorer/third/schedule"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

var blockScanTaskChain = make(chan *ent.Network, 10)

func init() {
	schedule.Add("block_sync", func(ctx context.Context) {
		BlockSyncRun(ctx)
	}).ScheduleWithCron("*/1 * * * * *")
	schedule.Add("scan_block", func(ctx context.Context) {
		startBlockScanRun()
		BlockScanRun(ctx)
	}).ScheduleWithCron("*/1 * * * * *")
	schedule.Add("scan_block_test", func(ctx context.Context) {
		startBlockScanRun()
		BlockScanRun(ctx)
		select {}
	})
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

var startBlockScanRunOnce = sync.Once{}

func startBlockScanRun() {
	startBlockScanRunOnce.Do(func() {
		blockScanTaskChain = make(chan *ent.Network, config.Task.BlockScanThreads*3/2)
		go startBlockScanNetworkDo()
	})
}

var blockScanRunNetworks ent.Networks

func BlockScanRun(ctx context.Context) {
	logger.Debug("start", "nets", config.Task.Networks)

	tx, err := entity.Client(ctx)
	if err != nil {
		log.Context(ctx).Error("db connect error", "err", err)
		return
	}

	blockScanRunNetworks, err = getTaskNetworks(ctx, tx)
	if err != nil {
		log.Context(ctx).Error("error in getTaskNetworks", "err", err)
		return
	}
}

func startBlockScanNetworkDo() {
	for i := 0; i < config.Task.BlockScanThreads+2; i++ {
		(func(i int) {
			go blockScanNetworkDo(i)
		})(i)
	}
	go blockScanNetworkFanout()
}

func blockScanNetworkFanout() {
	for {
		blockScanRunNetworks := blockScanRunNetworks[:]
		for _, network := range blockScanRunNetworks {
			blockScanTaskChain <- network
		}
		time.Sleep(time.Second / 100)
	}
}

func blockScanNetworkDo(i int) {
	for network := range blockScanTaskChain {
		ctx := context.Background()
		logger := logger.With("network", network.ID, "networkName", network.Name, "threads", i)
		client, err := entity.Client(ctx, network.ID)
		if err != nil {
			logger.Debug("error in network client", "err", err)
			continue
		}

		err = entity.WithTx(ctx, client, func(tx *ent.Client) error {
			blockSyncs, err := tx.BlockSync.
				Query().
				ForUpdate(sql.WithLockAction(sql.SkipLocked)).
				Where(
					//blocksync.Scanned(false),
					blocksync.Scanned(true),
				).
				Limit(1000).
				All(ctx)
			if err != nil {
				logger.Error("error in block sync query", "err", err)
				return err
			}

			if len(blockSyncs) < 1 {
				logger.Error("not found in block sync query", "err", err)
				return errors.New("not found")
			}

			bc, err := ethclient.DialContext(ctx, network.HTTPRPC)
			if err != nil {
				return err
			}

			wg := errgroup.Group{}
			wg.SetLimit(config.Task.BlockScanThreads + 1)

			for i, blockSync := range blockSyncs {
				(func(i int, blockSync *ent.BlockSync) {
					wg.Go(func() error {
						blockScanNetworkBlockDo(ctx, bc, blockSync, logger)
						return nil
					})
				})(i, blockSync)
			}
			wg.Wait()

			fmt.Sprintln(blockSyncs)
			return nil
		})
		if err != nil {
			time.Sleep(time.Second / 2)
			continue
		}
	}
}

func blockScanNetworkBlockDo(ctx context.Context, bc *ethclient.Client, blockSync *ent.BlockSync, logger log.Logger) {

	block := &vo.BlockWithBlockByNumber{}
	receipts := types.Receipts{}

	batchCall := []rpc.BatchElem{
		{
			Method: "eth_getBlockByNumber",
			Args: []interface{}{
				rpc.BlockNumber(blockSync.ID).String(),
				true,
			},
			Result: block,
			Error:  nil,
		},
		{
			Method: "eth_getBlockReceipts",
			Args: []interface{}{
				rpc.BlockNumber(blockSync.ID).String(),
			},
			Result: &receipts,
			Error:  nil,
		},
	}
	err := bc.Client().BatchCall(batchCall)
	if err != nil {
		logger.Error("error in batch call", "err", err)
		return
	}

	fmt.Sprintln(block)
}
