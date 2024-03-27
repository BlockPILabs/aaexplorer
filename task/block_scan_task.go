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
	"github.com/BlockPILabs/aaexplorer/internal/utils"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/BlockPILabs/aaexplorer/third/schedule"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
	"runtime"
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
		blockScanTaskChain = make(chan *ent.Network, config.Task.GetBlockScanThreads()*3/2)
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
	for i := 0; i < config.Task.GetBlockScanThreads(); i++ {
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
				Order(blocksync.ByID(sql.OrderDesc())).
				Limit(10).
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
			wg.SetLimit(runtime.NumCPU())

			blockDataDecodes := ent.BlockDataDecodes{}
			blockDataDecodeCreates := []*ent.BlockDataDecodeCreate{}

			blockDataTransactionDecodes := ent.TransactionDecodes{}
			blockDataTransactionDecodeCreates := []*ent.TransactionDecodeCreate{}

			blockDataTransactionReceiptDecodes := ent.TransactionReceiptDecodes{}
			blockDataTransactionReceiptDecodeCreates := []*ent.TransactionReceiptDecodeCreate{}

			//
			//blocksMap := map[int64]*parserBlock{}
			//transactionMap := map[string]*parserTransaction{}

			results := make([]*vo.BlockScanNetworkBlockDoResult, len(blockSyncs))
			for i, blockSync := range blockSyncs {
				(func(i int, blockSync *ent.BlockSync) {
					wg.Go(func() error {
						ret, err := blockScanNetworkBlockDo(ctx, bc, blockSync, logger)
						if err == nil {
							results[i] = ret

							blockDataDecode, blockDataDecodeCreate, transactionDecodes, transactionDecodeCreates, transactionReceiptDecodes, transactionReceiptDecodeCreates := parseBlockScanNetworkBlockDoResult(ctx, tx, client, ret)
							blockDataDecodes = append(blockDataDecodes, blockDataDecode)
							blockDataDecodeCreates = append(blockDataDecodeCreates, blockDataDecodeCreate)

							blockDataTransactionDecodes = append(blockDataTransactionDecodes, transactionDecodes...)
							blockDataTransactionDecodeCreates = append(blockDataTransactionDecodeCreates, transactionDecodeCreates...)

							blockDataTransactionReceiptDecodes = append(blockDataTransactionReceiptDecodes, transactionReceiptDecodes...)
							blockDataTransactionReceiptDecodeCreates = append(blockDataTransactionReceiptDecodeCreates, transactionReceiptDecodeCreates...)

						}
						return nil
					})
				})(i, blockSync)
			}
			wg.Wait()

			err = tx.BlockDataDecode.CreateBulk(blockDataDecodeCreates...).OnConflict(sql.DoNothing()).Exec(ctx)
			if err != nil {
				return err
			}
			err = tx.TransactionDecode.CreateBulk(blockDataTransactionDecodeCreates...).OnConflict(sql.DoNothing()).Exec(ctx)
			if err != nil {
				return err
			}
			err = tx.TransactionReceiptDecode.CreateBulk(blockDataTransactionReceiptDecodeCreates...).OnConflict(sql.DoNothing()).Exec(ctx)
			if err != nil {
				return err
			}

			fmt.Sprintln(results)
			return nil
		})
		if err != nil {
			time.Sleep(time.Second / 2)
			continue
		}
	}
}

func blockScanNetworkBlockDo(ctx context.Context, bc *ethclient.Client, blockSync *ent.BlockSync, logger log.Logger) (*vo.BlockScanNetworkBlockDoResult, error) {

	ret := &vo.BlockScanNetworkBlockDoResult{
		Block:    &vo.BlockWithBlockByNumber{},
		Receipts: []*vo.BlockWithGetBlockReceipt{},
	}

	batchCall := []rpc.BatchElem{
		{
			Method: "eth_getBlockByNumber",
			Args: []interface{}{
				rpc.BlockNumber(blockSync.ID).String(),
				true,
			},
			Result: ret.Block,
			Error:  nil,
		},
		{
			Method: "eth_getBlockReceipts",
			Args: []interface{}{
				rpc.BlockNumber(blockSync.ID).String(),
			},
			Result: &ret.Receipts,
			Error:  nil,
		},
	}
	err := bc.Client().BatchCall(batchCall)
	if err != nil {
		logger.Error("error in batch call", "err", err)
		return nil, err
	}

	if len(ret.Block.Hash) < 5 {
		return nil, errors.New("block not found")
	}

	if len(ret.Block.Transactions) != len(ret.Receipts) {
		return nil, errors.New("receipts fail")
	}
	return ret, nil
}

func parseBlockScanNetworkBlockDoResult(ctx context.Context, networkTx *ent.Client, tx *ent.Client, ret *vo.BlockScanNetworkBlockDoResult) (
	blockDataDecode *ent.BlockDataDecode,
	blockDataDecodeCreate *ent.BlockDataDecodeCreate,

	transactionDecodes ent.TransactionDecodes,
	transactionDecodeCreates []*ent.TransactionDecodeCreate,

	transactionReceiptDecodes ent.TransactionReceiptDecodes,
	transactionReceiptDecodeCreates []*ent.TransactionReceiptDecodeCreate,
) {
	// parse block data
	timestamp := time.Unix(utils.DecodeDecimal(ret.Block.Timestamp).IntPart(), 0)
	blockNumber := utils.DecodeDecimal(ret.Block.Number)
	uncles := &pgtype.TextArray{}
	uncles.Set(ret.Block.Uncles)
	blockDataDecode = &ent.BlockDataDecode{
		ID:               blockNumber.IntPart(),
		Time:             timestamp,
		CreateTime:       timestamp,
		Hash:             ret.Block.Hash,
		ParentHash:       ret.Block.ParentHash,
		Nonce:            utils.DecodeDecimal(ret.Block.Nonce),
		Sha3Uncles:       ret.Block.Sha3Uncles,
		LogsBloom:        ret.Block.LogsBloom,
		TransactionsRoot: ret.Block.TransactionsRoot,
		StateRoot:        ret.Block.StateRoot,
		ReceiptsRoot:     ret.Block.ReceiptsRoot,
		Miner:            ret.Block.Miner,
		MixHash:          ret.Block.MixHash,
		Difficulty:       utils.DecodeDecimal(ret.Block.Difficulty),
		TotalDifficulty:  utils.DecodeDecimal(ret.Block.TotalDifficulty),
		ExtraData:        ret.Block.ExtraData,
		Size:             utils.DecodeDecimal(ret.Block.Size),
		GasLimit:         utils.DecodeDecimal(ret.Block.GasLimit),
		GasUsed:          utils.DecodeDecimal(ret.Block.GasUsed),
		Timestamp:        utils.DecodeDecimal(ret.Block.Timestamp),
		TransactionCount: decimal.NewFromInt(int64(len(ret.Block.Transactions))),
		Uncles:           uncles,
		BaseFeePerGas:    utils.DecodeDecimal(ret.Block.BaseFeePerGas),
	}
	blockDataDecodeCreate = networkTx.BlockDataDecode.Create().
		SetID(blockDataDecode.ID).
		SetTime(blockDataDecode.Time).
		SetCreateTime(blockDataDecode.CreateTime).
		SetHash(blockDataDecode.Hash).
		SetParentHash(blockDataDecode.ParentHash).
		SetNonce(blockDataDecode.Nonce).
		SetSha3Uncles(blockDataDecode.Sha3Uncles).
		SetLogsBloom(blockDataDecode.LogsBloom).
		SetTransactionsRoot(blockDataDecode.TransactionsRoot).
		SetStateRoot(blockDataDecode.StateRoot).
		SetReceiptsRoot(blockDataDecode.ReceiptsRoot).
		SetMiner(blockDataDecode.Miner).
		SetMixHash(blockDataDecode.MixHash).
		SetDifficulty(blockDataDecode.Difficulty).
		SetTotalDifficulty(blockDataDecode.TotalDifficulty).
		SetExtraData(blockDataDecode.ExtraData).
		SetSize(blockDataDecode.Size).
		SetGasLimit(blockDataDecode.GasLimit).
		SetGasUsed(blockDataDecode.GasUsed).
		SetTimestamp(blockDataDecode.Timestamp).
		SetTransactionCount(blockDataDecode.TransactionCount).
		SetUncles(blockDataDecode.Uncles).
		SetBaseFeePerGas(blockDataDecode.BaseFeePerGas)

	for i, transaction := range ret.Block.Transactions {
		accessList := &pgtype.JSONB{}
		accessList.Set(transaction.AccessList)
		maxFeePerGas := utils.DecodeDecimal(transaction.MaxFeePerGas)
		maxPriorityFeePerGas := utils.DecodeDecimal(transaction.MaxPriorityFeePerGas)
		transactionDecode := &ent.TransactionDecode{
			ID:                   transaction.Hash,
			Time:                 timestamp,
			CreateTime:           timestamp,
			BlockHash:            blockDataDecode.Hash,
			BlockNumber:          blockDataDecode.ID,
			Nonce:                utils.DecodeDecimal(transaction.Nonce),
			TransactionIndex:     utils.DecodeDecimal(transaction.TransactionIndex),
			FromAddr:             transaction.From,
			ToAddr:               transaction.To,
			Value:                utils.DecodeDecimal(transaction.Value),
			GasPrice:             utils.DecodeDecimal(transaction.GasPrice),
			Gas:                  utils.DecodeDecimal(transaction.Gas),
			Input:                transaction.Input,
			R:                    transaction.R,
			S:                    transaction.S,
			V:                    utils.DecodeDecimal(transaction.V),
			ChainID:              utils.DecodeDecimal(transaction.ChainId).IntPart(),
			Type:                 transaction.Type,
			MaxFeePerGas:         &maxFeePerGas,
			MaxPriorityFeePerGas: &maxPriorityFeePerGas,
			AccessList:           accessList,
			Method:               "",
		}
		transactionDecodes = append(transactionDecodes, transactionDecode)
		transactionDecodeCreate := tx.TransactionDecode.Create().
			SetID(transactionDecode.ID).
			SetTime(transactionDecode.Time).
			SetCreateTime(transactionDecode.CreateTime).
			SetBlockHash(transactionDecode.BlockHash).
			SetBlockNumber(transactionDecode.BlockNumber).
			SetNonce(transactionDecode.Nonce).
			SetTransactionIndex(transactionDecode.TransactionIndex).
			SetFromAddr(transactionDecode.FromAddr).
			SetToAddr(transactionDecode.ToAddr).
			SetValue(transactionDecode.Value).
			SetGasPrice(transactionDecode.GasPrice).
			SetGas(transactionDecode.Gas).
			SetInput(transactionDecode.Input).
			SetR(transactionDecode.R).
			SetS(transactionDecode.S).
			SetV(transactionDecode.V).
			SetChainID(transactionDecode.ChainID).
			SetType(transactionDecode.Type).
			SetMaxFeePerGas(maxFeePerGas).
			SetMaxPriorityFeePerGas(maxPriorityFeePerGas).
			SetAccessList(transactionDecode.AccessList).
			SetMethod(transactionDecode.Method)
		transactionDecodeCreates = append(transactionDecodeCreates, transactionDecodeCreate)
		receipt := ret.Receipts[i]

		transactionReceiptDecode := &ent.TransactionReceiptDecode{
			ID:                receipt.TransactionHash,
			Time:              timestamp,
			CreateTime:        timestamp,
			BlockHash:         blockDataDecode.Hash,
			BlockNumber:       blockDataDecode.ID,
			ContractAddress:   receipt.ContractAddress,
			CumulativeGasUsed: utils.DecodeDecimal(receipt.CumulativeGasUsed).IntPart(),
			EffectiveGasPrice: receipt.EffectiveGasPrice,
			FromAddr:          receipt.From,
			GasUsed:           utils.DecodeDecimal(receipt.GasUsed),
			Logs:              string(receipt.Logs),
			LogsBloom:         receipt.LogsBloom,
			Status:            receipt.Status,
			ToAddr:            transactionDecode.ToAddr,
			TransactionIndex:  receipt.TransactionIndex,
			Type:              receipt.Type,
		}
		transactionReceiptDecodes = append(transactionReceiptDecodes, transactionReceiptDecode)

		transactionReceiptDecodeCreate := tx.TransactionReceiptDecode.Create().
			SetID(transactionReceiptDecode.ID).
			SetTime(transactionReceiptDecode.Time).
			SetCreateTime(transactionReceiptDecode.CreateTime).
			SetBlockHash(transactionReceiptDecode.BlockHash).
			SetBlockNumber(transactionReceiptDecode.BlockNumber).
			SetContractAddress(transactionReceiptDecode.ContractAddress).
			SetCumulativeGasUsed(transactionReceiptDecode.CumulativeGasUsed).
			SetEffectiveGasPrice(transactionReceiptDecode.EffectiveGasPrice).
			SetFromAddr(transactionReceiptDecode.FromAddr).
			SetGasUsed(transactionReceiptDecode.GasUsed).
			SetLogs(transactionReceiptDecode.Logs).
			SetLogsBloom(transactionReceiptDecode.LogsBloom).
			SetStatus(transactionReceiptDecode.Status).
			SetToAddr(transactionReceiptDecode.ToAddr).
			SetTransactionIndex(transactionReceiptDecode.TransactionIndex).
			SetType(transactionReceiptDecode.Type)
		transactionReceiptDecodeCreates = append(transactionReceiptDecodeCreates, transactionReceiptDecodeCreate)

	}

	return
}
