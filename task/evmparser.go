package task

import (
	"bytes"
	"context"
	"encoding/json"
	"entgo.io/ent/dialect/sql"
	"errors"
	"fmt"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aaaccountdata"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aablockinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aablocksync"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aatransactioninfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aauseropscalldata"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aauseropsinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/account"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/blockdatadecode"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactiondecode"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactionreceiptdecode"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/service"
	"github.com/BlockPILabs/aa-scan/internal/utils"
	ser "github.com/BlockPILabs/aa-scan/service"
	"github.com/BlockPILabs/aa-scan/task/aa"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgtype"
	"github.com/panjf2000/ants/v2"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/maps"
	"golang.org/x/sync/errgroup"
	"math"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"
)

const HandleOpsSign = "0x1fad948c"
const UserOperationEventSign = "0x49628fd1471006c1482da88028e9ce4dbb080b815c9b0344d39e5a8e6ec1419f"
const LogTransferEventSign = "0xe6497e3ee548a3372136af2fcb0696db31fc6cf20260707645068bd3fe97f3c4"
const TransferEventSign = "0xe6497e3ee548a3372136af2fcb0696db31fc6cf20260707645068bd3fe97f3c4"
const AccountDeploySign = "0xd51a9c61267aa6196961883ecf5ff2da6619c37dac0fa92122513fb32c032d2d"

const ExecuteSign = "0xb61d27f6"
const ExecuteSign1 = "0x51945447"
const ExecuteCall = "0x9e5d4c49"
const ExecuteBatchSign = "0x47e1da2a"
const ExecuteBatchCallSign = "0x912ccaa3"
const EmptyMethod = "00000000"

func InitEvmParse(ctx context.Context, config *config.Config, logger log.Logger) error {
	logger = logger.With("task", "evmparser")
	dayScheduler := chrono.NewDefaultTaskScheduler()
	t := _evmParser{
		logger:     logger,
		config:     config,
		startBlock: map[string]int64{},
	}

	for network, blockNumber := range t.config.EvmParser.StartBlock {
		t.startBlock[network] = blockNumber
		if t.startBlock[network] == -1 {
			t.startBlock[network] = 0
			client, err := entity.Client(ctx, network)
			if err != nil {
				log.Context(ctx).Warn("client error", "err", err, "network", network)
				continue
			}
			latestBlock, err := dao.AaBlockDao.GetLatestBlock(ctx, client)
			if err != nil {
				log.Context(ctx).Warn("GetLatestBlock error", "err", err, "network", network)
				continue
			}
			t.startBlock[network] = latestBlock.ID - int64(math.Max(float64(config.EvmParser.Multi*config.EvmParser.Batch), 10))

		}

		log.Context(ctx).Info("start block", "blockNumber", t.startBlock[network])
	}

	jsonAbi, err := abi.JSON(bytes.NewBufferString(t.config.EvmParser.GetAbi()))
	if err != nil {
		logger.Error("abi parse error", "err", err)
		return err
	}

	t.abi = jsonAbi
	t.handleOpsMethod, err = jsonAbi.MethodById(hexutil.MustDecode(HandleOpsSign))
	if err != nil {
		logger.Error("abi method parse error", "err", err)
		return err
	}

	_, err = dayScheduler.ScheduleWithCron(func(ctx context.Context) {
		t.ScanBlock(log.WithContext(ctx, logger.With("action", "ScanBlock")))
	}, "*/10 * * * * *")
	if err != nil {
		logger.Error("Schedule error", "err", err)

		return err
	}

	return err
}

func (t *_evmParser) ScanBlock(ctx context.Context) {
	fiend := true
	logger := log.Context(ctx)
	logger.Debug("scan start")
	pool, err := ants.NewPool(t.config.EvmParser.Batch)
	if err != nil {
		logger.Error("network find error", "err", err)
		return
	}
	defer pool.Release()
	wg := &sync.WaitGroup{}
	for fiend {
		fiend = false
		networks, err := service.NetworkService.GetNetworks(context.Background())
		if err != nil {
			logger.Error("network find error", "err", err)
			return
		}
		logger.Debug("waiting scans", "num", pool.Running())
		for _, network := range networks {
			if _, ok := t.startBlock[network.ID]; !ok {
				t.startBlock[network.ID] = 0
			}
			wg.Add(1)
			ctx := log.WithContext(context.Background(), logger.With("network", network.ID))
			if t.ScanBlockByNetwork(ctx, network, wg, pool) {
				fiend = true
			}
		}
	}
	wg.Wait()

	logger.Info("scan complete")

	// reset start
	networks, err := service.NetworkService.GetNetworks(context.Background())
	for _, network := range networks {
		t.startBlock[network.ID] = 0
	}

}
func (t *_evmParser) ScanBlockByNetwork(ctx context.Context, network *ent.Network, wg *sync.WaitGroup, pool *ants.Pool) (fiend bool) {
	defer func() {
		if !fiend {
			wg.Done()
		}
	}()

	logger := log.Context(ctx)
	logger.Info("start block", "network", network.ID)
	client, err := entity.NetworkClient(ctx, network)
	if err != nil {
		logger.Error("network db client", "err", err)
		return false
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		logger.Error("network db client tx", "err", err)
		return false
	}

	defer func() {
		if !fiend {
			err := tx.Rollback()
			if err != nil {
				t.logger.Error("roll back error", "err", err)
			}
		}
	}()

	aaBlockSyncs, err := tx.AaBlockSync.Query().
		Where(
			aablocksync.Scanned(false),
			aablocksync.BlockScannedNotNil(),
			aablocksync.BlockScanned(true),
			aablocksync.TxScannedNotNil(),
			aablocksync.TxScanned(true),
			aablocksync.TxrScannedNotNil(),
			aablocksync.TxrScanned(true),
			aablocksync.IDGT(t.startBlock[network.ID]),
		).
		ForUpdate(sql.WithLockAction(sql.SkipLocked)).
		Order(
			aablocksync.ByID(
				sql.OrderAsc(),
			),
		).
		Limit(t.config.EvmParser.Multi).
		All(ctx)
	if err != nil {
		logger.Error("find AaBlockSync  tx", "err", err)
		return false
	}
	logger.Debug("queried aaBlockSyncs", "count", len(aaBlockSyncs))
	if len(aaBlockSyncs) < 1 {
		logger.Debug("not find AaBlockSync")
		return false
	}
	fiend = true

	blockIds := make([]int64, len(aaBlockSyncs))
	for i, blockSync := range aaBlockSyncs {
		blockIds[i] = blockSync.ID
		t.startBlock[network.ID] = blockSync.ID
	}

	err = pool.Submit(func() {

		ctx = log.WithContext(context.Background(), logger)
		defer func() {
			err := tx.Commit()
			if err != nil {

				log.Context(ctx).Warn("commit error", "err", err)
			}
			wg.Done()
		}()
		client, err = entity.NetworkClient(ctx, network)
		logger.Debug("pre parse , get parse data")
		blockDataDecodes, transactionDecodes, receiptDecodes, blocksMap, transactionMap, err := t.getParseData(ctx, client, blockIds...)
		_ = (blockDataDecodes)
		_ = (transactionDecodes)
		_ = (receiptDecodes)
		_ = (blocksMap)
		_ = (transactionMap)

		logger.Debug("complete get parse data", "err", err,
			"blockDataDecodes", len(blockDataDecodes),
			"transactionDecodes", len(transactionDecodes),
			"receiptDecodes", len(receiptDecodes),
			"blocksMap", len(blocksMap),
			"transactionMap", len(transactionMap),
		)
		if err != nil {
			logger.Error("get parse data error", "err", err)
			return
		}

		start := time.Now()
		defer func() {
			logger.Debug("block parse", "blockIds", blockIds, "count", len(blockIds), "duration", time.Now().Sub(start).Round(time.Millisecond))
		}()

		var aaUserOpsInfos ent.AAUserOpsInfos
		var aaTransactionInfos ent.AaTransactionInfos
		var userOpsInfoCalldatas ent.AAUserOpsCalldataSlice
		var aaBlockInfos ent.AaBlockInfos
		var setBlockSyncedId []int64
		var aaAccountDataMap = map[string]*ent.AaAccountData{}
		for _, block := range blocksMap {
			t.doParse(ctx, client, network, block)

			setBlockSyncedId = append(setBlockSyncedId, block.block.ID)

			if block.userOpInfo == nil || block.userOpInfo.UseropCount < 1 {
				continue
			}
			aaBlockInfos = append(aaBlockInfos, block.userOpInfo)

			for _, transition := range block.transitions {
				if len(transition.userops) > 0 {
					aaUserOpsInfos = append(aaUserOpsInfos, transition.userops...)
				}

				if len(transition.userOpsCalldata) > 0 {
					userOpsInfoCalldatas = append(userOpsInfoCalldatas, transition.userOpsCalldata...)
				}

				if transition.userOpInfo != nil {
					aaTransactionInfos = append(aaTransactionInfos, transition.userOpInfo)
				}
			}

			accountDataSlice := block.AaAccountDataSlice()
			for i, data := range accountDataSlice {
				if accountData, ok := aaAccountDataMap[data.ID]; ok {
					if len(data.AaType) > 0 && len(accountData.AaType) < 1 {
						accountData.AaType = data.AaType
					}
					if len(data.Factory) > 0 && len(accountData.Factory) < 1 {
						accountData.Factory = data.Factory
						accountData.FactoryTime = data.FactoryTime
					}
				} else {
					aaAccountDataMap[data.ID] = accountDataSlice[i]
				}
			}
		}

		t.insertUserOpsInfo(ctx, client, network, aaUserOpsInfos)
		t.insertTransactions(ctx, client, network, aaTransactionInfos)
		t.insertBlockInfos(ctx, client, network, aaBlockInfos)
		t.insertuserOpsInfoCalldatas(ctx, client, network, userOpsInfoCalldatas)
		t.insertAccounts(ctx, client, network, aaAccountDataMap)
		t.insertAaAccounts(ctx, client, network, aaAccountDataMap)

		// set sync status
		if len(setBlockSyncedId) > 0 {
			affected, err := tx.AaBlockSync.Update().
				Where(
					aablocksync.IDIn(setBlockSyncedId...),
				).
				SetScanned(true).
				SetUpdateTime(time.Now()).Save(ctx)
			if err != nil {
				logger.Warn("set block sync status error", "err", err)
			} else {
				logger.Info("set block scanned", "ids", setBlockSyncedId, "num", affected)
			}
		}

	})
	if err != nil {
		logger.Warn("block scanned error", "err", err)
		return false
	}
	logger.Debug("block success")
	return fiend
}

func (t *_evmParser) getParseData(ctx context.Context, client *ent.Client, blockIds ...int64) (
	blockDataDecodes []*ent.BlockDataDecode,
	transactionDecodes []*ent.TransactionDecode,
	transactionReceiptDecodes []*ent.TransactionReceiptDecode,
	blocksMap map[int64]*parserBlock,
	transactionMap map[string]*parserTransaction,
	retErr error,
) {

	//timeoutCtx, _ := context.WithTimeout(ctx, time.Minute)
	timeoutCtx := ctx

	g, _ := errgroup.WithContext(timeoutCtx)
	g.Go(func() error {

		var err error
		blockDataDecodes, err = client.BlockDataDecode.Query().
			Where(
				blockdatadecode.IDIn(
					blockIds...,
				),
			).All(timeoutCtx)
		if err != nil {
			log.Context(ctx).Error("not find BlockDataDecode", "err", err)
		} else {
			log.Context(ctx).Debug("find BlockDataDecode", "count", len(blockDataDecodes))
		}
		return err
	})

	g.Go(func() error {
		var err error

		transactionDecodes, err = client.TransactionDecode.Query().
			Where(
				transactiondecode.BlockNumberIn(blockIds...),
			).All(timeoutCtx)
		if err != nil {
			log.Context(ctx).Error("not find TransactionDecode", "err", err)
		} else {
			log.Context(ctx).Debug("find transactionDecodes", "count", len(transactionDecodes))
		}
		return err
	})

	g.Go(func() error {
		var err error
		transactionReceiptDecodes, err = client.TransactionReceiptDecode.Query().
			Where(
				transactionreceiptdecode.BlockNumberIn(blockIds...),
			).All(timeoutCtx)
		if err != nil {
			log.Context(ctx).Error("not find TransactionReceiptDecode", "err", err)
		} else {
			log.Context(ctx).Debug("find transactionReceiptDecodes", "count", len(transactionReceiptDecodes))
		}
		return err
	})
	retErr = g.Wait()

	blocksMap = map[int64]*parserBlock{}
	transactionMap = map[string]*parserTransaction{}
	for _, blockDataDecode := range blockDataDecodes {
		blocksMap[blockDataDecode.ID] = &parserBlock{
			block:         blockDataDecode,
			transitions:   []*parserTransaction{},
			userOpInfo:    &ent.AaBlockInfo{},
			aaAccounts:    &sync.Map{},
			aaAccountsLck: &sync.Mutex{},
		}
	}
	for _, transactionDecode := range transactionDecodes {
		if b, ok := blocksMap[transactionDecode.BlockNumber]; ok {
			transactionMap[transactionDecode.ID] = &parserTransaction{
				transaction:     transactionDecode,
				receipt:         nil,
				userOpInfo:      nil,
				userops:         nil,
				userOpsCalldata: nil,
			}
			b.transitions = append(b.transitions, transactionMap[transactionDecode.ID])
		}
	}

	for _, transactionReceiptDecode := range transactionReceiptDecodes {
		if tx, ok := transactionMap[transactionReceiptDecode.ID]; ok {
			tx.receipt = transactionReceiptDecode
		}
	}

	// filter blocks
	delKeys := []int64{}
	for blockNumber, block := range blocksMap {
		if len(block.transitions) < 1 {
			delKeys = append(delKeys, blockNumber)
			continue
		}

		for _, transition := range block.transitions {
			if transition.receipt == nil || transition.transaction == nil {
				delKeys = append(delKeys, blockNumber)
				break
			}
		}
	}

	for _, blockNumber := range delKeys {
		delete(blocksMap, blockNumber)
	}

	return
}

func (t *_evmParser) getCurrentTimestampMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (t *_evmParser) doParse(ctx context.Context, client *ent.Client, network *ent.Network, block *parserBlock) {

	ctx, logger := log.With(ctx, "blockNumber", block.block.ID)

	logger.Debug("start parse")
	parserTransactions := block.transitions

	block.userOpInfo = &ent.AaBlockInfo{
		ID:             block.block.ID,
		Time:           block.block.Time,
		Hash:           block.block.Hash,
		UseropCount:    0,
		UseropMevCount: 0,
		BundlerProfit:  decimal.Decimal{},
		CreateTime:     time.Now(),
	}
	for _, parserTx := range parserTransactions {
		tx := parserTx.transaction
		input := tx.Input
		if len(input) <= 10 {
			continue
		}
		sign := input[:10]
		input = input[10:]
		if sign != HandleOpsSign {
			continue
		}

		t.parseUserOps(ctx, client, network, block, parserTx)

		block.userOpInfo.BundlerProfit = block.userOpInfo.BundlerProfit.Add(parserTx.userOpInfo.BundlerProfit)
		block.userOpInfo.UseropCount += len(parserTx.userops)
	}
	block.userOpInfo.BundlerProfitUsd = block.userOpInfo.BundlerProfit.Mul(ser.GetNativePrice(network.ID))
}

func (t *_evmParser) getFrom(tx *types.Transaction, client *ethclient.Client) string {
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		t.logger.Info("get networkId", "err", err)
		return ""
	}
	signer := types.LatestSignerForChainID(chainID)
	from, err := types.Sender(signer, tx)
	if err != nil {
		return ""
	}
	return strings.ToLower(from.String())
}

func (t *_evmParser) insertTransactions(ctx context.Context, client *ent.Client, network *ent.Network, infos ent.AaTransactionInfos) {
	if len(infos) < 1 {
		return
	}

	var transactionInfoCreates []*ent.AaTransactionInfoCreate
	for _, tx := range infos {
		txCreate := client.AaTransactionInfo.Create().
			SetTime(tx.Time).
			SetBlockHash(tx.BlockHash).
			SetBlockNumber(tx.BlockNumber).
			SetUseropCount(tx.UseropCount).
			SetIsMev(tx.IsMev).
			SetBundlerProfit(tx.BundlerProfit).
			SetCreateTime(tx.CreateTime).
			SetBundlerProfitUsd(tx.BundlerProfitUsd).
			SetID(tx.ID)

		transactionInfoCreates = append(transactionInfoCreates, txCreate)
	}
	err := client.AaTransactionInfo.
		CreateBulk(transactionInfoCreates...).
		OnConflictColumns(aatransactioninfo.FieldTime, aatransactioninfo.FieldID).
		Update(func(upsert *ent.AaTransactionInfoUpsert) {
			upsert.UpdateTime().
				UpdateBlockHash().
				UpdateBlockNumber().
				UpdateUseropCount().
				UpdateIsMev().
				UpdateBundlerProfit().
				UpdateBundlerProfitUsd()
		}).
		Exec(context.Background())
	if err != nil {
		log.Context(ctx).Info("insert AaTransactionInfo error", "err", err)
	}
}
func (t *_evmParser) insertBlockInfos(ctx context.Context, client *ent.Client, network *ent.Network, infos ent.AaBlockInfos) {
	if len(infos) < 1 {
		return
	}

	var transactionInfoCreates []*ent.AaBlockInfoCreate
	for _, tx := range infos {
		txCreate := client.AaBlockInfo.Create().
			SetTime(tx.Time).
			SetHash(tx.Hash).
			SetUseropCount(tx.UseropCount).
			SetUseropMevCount(tx.UseropMevCount).
			SetBundlerProfit(tx.BundlerProfit).
			SetBundlerProfitUsd(tx.BundlerProfitUsd).
			SetCreateTime(tx.CreateTime).
			SetID(tx.ID)

		transactionInfoCreates = append(transactionInfoCreates, txCreate)
	}
	err := client.AaBlockInfo.
		CreateBulk(transactionInfoCreates...).
		OnConflictColumns(aablockinfo.FieldTime, aablockinfo.FieldID).
		Update(func(upsert *ent.AaBlockInfoUpsert) {
			upsert.UpdateTime().
				UpdateHash().
				UpdateUseropCount().
				UpdateUseropMevCount().
				UpdateBundlerProfit().
				UpdateBundlerProfitUsd()
		}).
		Exec(context.Background())
	if err != nil {
		log.Context(ctx).Info("insert AaBlockInfo error", "err", err)
	}
}
func (t *_evmParser) insertuserOpsInfoCalldatas(ctx context.Context, client *ent.Client, network *ent.Network, infos ent.AAUserOpsCalldataSlice) {
	if len(infos) == 0 {
		return
	}

	var transactionInfoCreates []*ent.AAUserOpsCalldataCreate
	for _, tx := range infos {
		txCreate := client.AAUserOpsCalldata.Create().
			SetTime(tx.Time).
			SetUserOpsHash(tx.UserOpsHash).
			SetTxHash(tx.TxHash).
			SetBlockNumber(tx.BlockNumber).
			SetNetwork(tx.Network).
			SetSender(tx.Sender).
			SetTarget(tx.Target).
			SetTxValue(tx.TxValue).
			SetSource(tx.Source).
			SetCalldata(tx.Calldata).
			SetTxTime(tx.TxTime).
			SetCreateTime(tx.CreateTime).
			SetUpdateTime(tx.UpdateTime).
			SetAaIndex(tx.AaIndex).
			SetID(tx.ID)

		transactionInfoCreates = append(transactionInfoCreates, txCreate)
	}
	err := client.AAUserOpsCalldata.
		CreateBulk(transactionInfoCreates...).
		OnConflictColumns(aauseropscalldata.FieldTime, aauseropscalldata.FieldID).
		Update(func(upsert *ent.AAUserOpsCalldataUpsert) {

			upsert.UpdateTime().
				UpdateUserOpsHash().
				UpdateTxHash().
				UpdateBlockNumber().
				UpdateNetwork().
				UpdateSender().
				UpdateTarget().
				UpdateTxValue().
				UpdateSource().
				UpdateCalldata().
				UpdateTxTime().
				UpdateUpdateTime().
				UpdateAaIndex()

		}).
		Exec(context.Background())
	if err != nil {
		log.Context(ctx).Info("insert AAUserOpsCalldata error", "err", err)
	}
}

func (t *_evmParser) insertUserOpsInfo(ctx context.Context, client *ent.Client, network *ent.Network, infos ent.AAUserOpsInfos) {
	if len(infos) < 1 {
		return
	}
	var userOpsInfoCreates []*ent.AAUserOpsInfoCreate
	for _, ops := range infos {
		userOpsCreate := client.AAUserOpsInfo.Create().
			SetTime(ops.Time).
			SetTxHash(ops.TxHash).
			SetBlockNumber(ops.BlockNumber).
			SetNetwork(ops.Network).
			SetSender(ops.Sender).
			SetTarget(ops.Target).
			SetTargets(ops.Targets).
			SetTxValue(ops.TxValue).
			SetFee(ops.Fee).
			SetBundler(ops.Bundler).
			SetEntryPoint(ops.EntryPoint).
			SetFactory(ops.Factory).
			SetPaymaster(ops.Paymaster).
			SetPaymasterAndData(ops.PaymasterAndData).
			SetSignature(ops.Signature).
			SetCalldata(ops.Calldata).
			SetCalldataContract(ops.CalldataContract).
			SetNonce(ops.Nonce).
			SetCallGasLimit(ops.CallGasLimit).
			SetPreVerificationGas(ops.PreVerificationGas).
			SetVerificationGasLimit(ops.VerificationGasLimit).
			SetMaxFeePerGas(ops.MaxFeePerGas).
			SetMaxPriorityFeePerGas(ops.MaxPriorityFeePerGas).
			SetTxTime(ops.TxTime).
			SetInitCode(ops.InitCode).
			SetStatus(ops.Status).
			SetSource(ops.Source).
			SetActualGasCost(ops.ActualGasCost).
			SetActualGasUsed(ops.ActualGasUsed).
			SetCreateTime(ops.CreateTime).
			SetUpdateTime(ops.UpdateTime).
			SetUsdAmount(*ops.UsdAmount).
			SetID(ops.ID).
			SetTargetsCount(ops.TargetsCount).
			SetAaIndex(ops.AaIndex).
			SetFeeUsd(ops.FeeUsd).
			SetTxValueUsd(ops.TxValueUsd)

		userOpsInfoCreates = append(userOpsInfoCreates, userOpsCreate)
	}
	err := client.AAUserOpsInfo.CreateBulk(userOpsInfoCreates...).
		OnConflict(
			sql.ConflictColumns(aauseropsinfo.FieldTime, aauseropsinfo.FieldTxHash, aauseropsinfo.FieldID),
		).
		Update(func(upsert *ent.AAUserOpsInfoUpsert) {
			upsert.UpdateTime().
				UpdateTxHash().
				UpdateBlockNumber().
				UpdateNetwork().
				UpdateSender().
				UpdateTarget().
				UpdateTargets().
				UpdateTxValue().
				UpdateFee().
				UpdateBundler().
				UpdateEntryPoint().
				UpdateFactory().
				UpdatePaymaster().
				UpdatePaymasterAndData().
				UpdateSignature().
				UpdateCalldata().
				UpdateCalldataContract().
				UpdateNonce().
				UpdateCallGasLimit().
				UpdatePreVerificationGas().
				UpdateVerificationGasLimit().
				UpdateMaxFeePerGas().
				UpdateMaxPriorityFeePerGas().
				UpdateTxTime().
				UpdateInitCode().
				UpdateStatus().
				UpdateSource().
				UpdateActualGasCost().
				UpdateActualGasUsed().
				UpdateUpdateTime().
				UpdateUsdAmount().
				UpdateAaIndex().
				UpdateTargetsCount().
				UpdateFeeUsd().
				UpdateTxValueUsd()
		}).Exec(context.Background())
	if err != nil {
		log.Context(ctx).Info("insert AAUserOpsInfo error", "err", err)
	}
}

func (t *_evmParser) insertAccounts(ctx context.Context, client *ent.Client, network *ent.Network, dataMap map[string]*ent.AaAccountData) {
	if len(dataMap) < 1 {
		return
	}
	keys := maps.Keys(dataMap)

	accounts := client.Account.
		Query().
		Where(
			account.IDIn(keys...),
		).AllX(ctx)

	accountsMap := map[string]*ent.Account{}
	for i, a := range accounts {
		accountsMap[a.ID] = accounts[i]
	}
	var insertAccounts []*ent.AccountCreate

	// find insert
	for id, aaAccount := range dataMap {
		if _, ok := accountsMap[id]; !ok {
			emptyArray := pgtype.TextArray{}
			emptyArray.Set([]string{})
			create := client.Account.Create().
				SetID(aaAccount.ID).
				SetAbi("").
				SetLabel(&emptyArray).
				SetTag(&emptyArray).
				SetUpdateTime(time.Now())

			if len(aaAccount.AaType) > 0 {
				tags := []string{aaAccount.AaType}
				textArray := &pgtype.TextArray{}
				err := textArray.Set(tags)
				if err == nil {
					create.SetTag(textArray)
				}
			}
			insertAccounts = append(insertAccounts, create)
		}
	}

	// find update
	for id, a := range accountsMap {
		if aaAccount, ok := dataMap[id]; ok {
			upd := client.Account.UpdateOneID(id)
			needUpd := false

			if len(aaAccount.AaType) > 0 {
				var tags []string
				if a.Tag != nil && len(a.Tag.Elements) > 0 {
					_ = a.Tag.AssignTo(&tags)
				}

				tagContains := false
				for _, tag := range tags {
					if tag == aaAccount.AaType {
						tagContains = true
					}
				}

				if !tagContains {
					tags = append(tags, aaAccount.AaType)
					textArray := &pgtype.TextArray{}
					err := textArray.Set(tags)
					if err == nil {
						upd.SetTag(textArray)
						needUpd = true
					}
				}
			}

			if needUpd {
				_ = upd.Exec(ctx)
			}

		}
	}
	if len(insertAccounts) > 0 {
		err := client.Account.
			CreateBulk(insertAccounts...).
			OnConflictColumns(account.FieldID).
			Update(func(upsert *ent.AccountUpsert) {
				upsert.UpdateUpdateTime()
			}).Exec(ctx)
		if err != nil {
			log.Context(ctx).Error("account create error", "err", err)
		}
	}

	//accountsMap := make()
	//fmt.Println(accounts)
}
func (t *_evmParser) insertAaAccounts(ctx context.Context, client *ent.Client, network *ent.Network, dataMap map[string]*ent.AaAccountData) {
	if len(dataMap) < 1 {
		return
	}
	keys := maps.Keys(dataMap)

	accounts := client.AaAccountData.
		Query().
		Where(
			aaaccountdata.IDIn(keys...),
		).AllX(ctx)

	accountsMap := map[string]*ent.AaAccountData{}
	for i, a := range accounts {
		accountsMap[a.ID] = accounts[i]
	}
	var insertAccounts []*ent.AaAccountDataCreate

	factoryMap := map[string]*ent.AaAccountData{}
	paymasterMap := map[string]*ent.AaAccountData{}
	bundlerMap := map[string]*ent.AaAccountData{}

	// find insert
	for id, aaAccount := range dataMap {
		switch aaAccount.AaType {
		case config.AaAccountTypeFactory:
			factoryMap[aaAccount.ID] = dataMap[id]
		case config.AaAccountTypePaymaster:
			paymasterMap[aaAccount.ID] = dataMap[id]
		case config.AaAccountTypeBundler:
			bundlerMap[aaAccount.ID] = dataMap[id]
		}
		if _, ok := accountsMap[id]; !ok {
			create := client.AaAccountData.Create().
				SetID(aaAccount.ID).
				SetAaType(aaAccount.AaType).
				SetFactory(aaAccount.Factory).
				SetFactoryTime(aaAccount.FactoryTime).
				SetUpdateTime(time.Now())
			insertAccounts = append(insertAccounts, create)
		}
	}

	// find update
	for id, a := range accountsMap {
		if aaAccount, ok := dataMap[id]; ok {
			upd := client.AaAccountData.UpdateOneID(id)
			needUpd := false

			if len(a.AaType) < 1 && len(aaAccount.AaType) > 0 {
				upd.SetAaType(aaAccount.AaType)
				needUpd = true
			}

			if len(aaAccount.Factory) > 0 {
				upd.SetFactory(aaAccount.Factory)
				upd.SetFactoryTime(aaAccount.FactoryTime)
				needUpd = true
			}

			if needUpd {
				_ = upd.Exec(ctx)
			}

		}
	}
	if len(insertAccounts) > 0 {
		err := client.AaAccountData.CreateBulk(insertAccounts...).
			OnConflictColumns(aaaccountdata.FieldID).
			Update(func(upsert *ent.AaAccountDataUpsert) {
				upsert.UpdateUpdateTime()
			}).Exec(ctx)
		if err != nil {
			log.Context(ctx).Error("account create error", "err", err)
		}
	}

}

func (t *_evmParser) parseUserOps(ctx context.Context, client *ent.Client, network *ent.Network, block *parserBlock, parserTx *parserTransaction) error {
	ctx, logger := log.With(ctx, "transaction", parserTx.transaction.ID)
	logger.Debug("start parse transaction")
	data, err := hexutil.Decode(parserTx.transaction.Input)
	if err != nil {
		logger.Warn("decode input error", "err", err)
		return err
	}

	unpack, err := t.handleOpsMethod.Inputs.UnpackValues(data[4:])
	if err != nil {
		logger.Warn("abi unpack input error", "err", err)
		return err
	}
	if len(unpack) < 2 {
		logger.Warn("abi unpack data error", "err", err)
		return errors.New("abi unpack error")
	}
	logger.Debug("abi packed")
	//beneficiary := parserTx.transaction.FromAddr
	//if beneficiaryAddr, ok := unpack[1].(common.Address); ok {
	//	beneficiary = strings.ToLower(beneficiaryAddr.Hex())
	//}

	opsBytes, _ := json.Marshal(unpack[0])
	var ops []*aa.UserOperation
	_ = json.Unmarshal(opsBytes, &ops)
	err = json.Unmarshal([]byte(parserTx.receipt.Logs), &parserTx.logs)
	if err != nil {
		logger.Warn("abi  Unmarshal success", "err", err)
		return err
	}

	events, _ := t.parseLogs(ctx, parserTx.logs)

	parserTx.userOpInfo = &ent.AaTransactionInfo{
		ID:            parserTx.transaction.ID,
		Time:          parserTx.transaction.Time,
		BlockHash:     parserTx.transaction.BlockHash,
		BlockNumber:   parserTx.transaction.BlockNumber,
		UseropCount:   0,
		IsMev:         false,
		BundlerProfit: decimal.Decimal{},
		CreateTime:    time.Now(),
	}

	bundler := block.AaAccountData(parserTx.transaction.FromAddr)
	bundler.AaType = config.AaAccountTypeBundler

	entryPoint := block.AaAccountData(parserTx.transaction.ToAddr)
	entryPoint.AaType = config.AaAccountTypeEntryPoint

	for aaIndex, op := range ops {
		var source = ""
		var target = ""
		callDetails, source := t.parseCallData(ctx, client, network, hexutil.Encode(op.CallData))
		var targetsMap = map[string]string{}
		var targets = []string{}
		for i, callDetail := range callDetails {
			if i == 0 {
				target = callDetail.target
			}
			if _, ok := targetsMap[callDetail.target]; !ok {
				targetsMap[callDetail.target] = callDetail.target
				targets = append(targets, callDetail.target)
			}
		}
		var pgTarges = pgtype.TextArray{}
		pgTarges.Set(targets)

		userOpHash := op.GetUserOpHash(common.HexToAddress(parserTx.transaction.ToAddr), big.NewInt(network.ChainID))
		now := time.Now()
		userOpsInfo := &ent.AAUserOpsInfo{
			ID:                   userOpHash.Hex(),
			Time:                 parserTx.transaction.Time,
			TxHash:               parserTx.transaction.ID,
			BlockNumber:          parserTx.transaction.BlockNumber,
			Network:              network.ID,
			Sender:               strings.ToLower(op.Sender.Hex()),
			Target:               target,
			Targets:              &pgTarges,
			TxValue:              decimal.Zero,
			Fee:                  decimal.Zero,
			Bundler:              strings.ToLower(parserTx.transaction.FromAddr),
			EntryPoint:           strings.ToLower(parserTx.transaction.ToAddr),
			Factory:              "",
			Paymaster:            "",
			PaymasterAndData:     hexutil.Encode(op.PaymasterAndData),
			Signature:            hexutil.Encode(op.Signature),
			Calldata:             hexutil.Encode(op.CallData),
			CalldataContract:     "",
			Nonce:                op.Nonce.Int64(),
			CallGasLimit:         op.CallGasLimit.Int64(),
			PreVerificationGas:   op.PreVerificationGas.Int64(),
			VerificationGasLimit: op.VerificationGasLimit.Int64(),
			MaxFeePerGas:         op.MaxFeePerGas.Int64(),
			MaxPriorityFeePerGas: op.MaxPriorityFeePerGas.Int64(),
			TxTime:               parserTx.transaction.Time.Unix(),
			InitCode:             hexutil.Encode(op.InitCode),
			Status:               0,
			Source:               source,
			ActualGasCost:        0,
			ActualGasUsed:        0,
			CreateTime:           now,
			UpdateTime:           now,
			UsdAmount:            &decimal.Decimal{},
			AaIndex:              aaIndex + 1,
			TargetsCount:         len(callDetails),
		}
		sender := block.AaAccountData(userOpsInfo.Sender)
		sender.AaType = config.AaAccountTypeAA
		factoryAddr, paymaster := t.getAddr(ctx, userOpsInfo.InitCode, userOpsInfo.PaymasterAndData)
		userOpsInfo.Factory = strings.ToLower(factoryAddr)
		userOpsInfo.Paymaster = strings.ToLower(paymaster)

		if len(userOpsInfo.Paymaster) > 0 {
			paymaster := block.AaAccountData(userOpsInfo.Paymaster)
			paymaster.AaType = config.AaAccountTypePaymaster
		}

		if len(userOpsInfo.Factory) > 0 {
			factory := block.AaAccountData(userOpsInfo.Factory)
			factory.AaType = config.AaAccountTypeFactory
			sender.Factory = userOpsInfo.Factory
			sender.FactoryTime = userOpsInfo.Time
		}
		opsVal, ok := events[userOpsInfo.Sender+strconv.Itoa(int(userOpsInfo.Nonce))]
		if ok {
			userOpsInfo.ActualGasUsed = opsVal.ActualGasUsed
			userOpsInfo.ID = opsVal.OpsHash
			userOpsInfo.ActualGasCost = opsVal.ActualGasCost
			userOpsInfo.Status = int32(opsVal.Success)
		}

		userOpsInfo.Fee = utils.DivRav(opsVal.ActualGasCost)
		//opsTxValue, opsTxValueOk := opsValMap[userOpsInfo.Sender]
		//if opsTxValueOk {
		//	userOpsInfo.TxValue = opsTxValue
		//}

		for aaCallIndex, callDetail := range callDetails {
			id := crypto.
				Keccak256Hash(
					[]byte(fmt.Sprintf(
						"%s-%s-%s-%s-%s-%d-%d",
						userOpsInfo.TxHash,
						userOpsInfo.ID,
						userOpsInfo.Sender,
						callDetail.target,
						callDetail.data,
						callDetail.value.IntPart(),
						aaIndex+1,
					)),
				).
				Hex()
			aaUserOpsCalldata := &ent.AAUserOpsCalldata{
				ID:          id,
				Time:        userOpsInfo.Time,
				UserOpsHash: userOpsInfo.ID,
				TxHash:      userOpsInfo.TxHash,
				BlockNumber: userOpsInfo.BlockNumber,
				Network:     userOpsInfo.Network,
				Sender:      userOpsInfo.Sender,
				Target:      callDetail.target,
				TxValue:     &callDetail.value,
				Source:      callDetail.source,
				Calldata:    callDetail.data,
				TxTime:      userOpsInfo.TxTime,
				CreateTime:  now,
				UpdateTime:  now,
				AaIndex:     aaCallIndex + 1,
			}
			parserTx.userOpsCalldata = append(parserTx.userOpsCalldata, aaUserOpsCalldata)

			block.AaAccountData(aaUserOpsCalldata.Target)
			userOpsInfo.TxValue = userOpsInfo.TxValue.Add(callDetail.value)
		}

		userOpsInfo.FeeUsd = userOpsInfo.Fee.Mul(ser.GetNativePrice(network.ID))
		userOpsInfo.TxValueUsd = userOpsInfo.TxValue.Mul(ser.GetNativePrice(network.ID))

		parserTx.userOpInfo.UseropCount++
		parserTx.userOpInfo.BundlerProfit = parserTx.userOpInfo.BundlerProfit.Add(userOpsInfo.Fee)

		parserTx.userops = append(parserTx.userops, userOpsInfo)
	}

	parserTx.userOpInfo.BundlerProfit = parserTx.userOpInfo.BundlerProfit.Sub(GetReceiptGasRayDiv(parserTx.receipt))
	parserTx.userOpInfo.BundlerProfitUsd = parserTx.userOpInfo.BundlerProfit.Mul(ser.GetNativePrice(network.ID))
	logger.Info("parse success")
	return nil
}

func (t *_evmParser) parseLogs(ctx context.Context, logs []*aa.Log) (map[string]UserOperationEvent, map[string]decimal.Decimal) {

	events := make(map[string]UserOperationEvent)
	opsTxValMap := make(map[string]decimal.Decimal)
	for _, log := range logs {
		topics := log.Topics
		if len(topics) < 1 {
			continue
		}
		sign := topics[0]
		dataStr := log.Data
		if len(dataStr) <= 2 {
			continue
		}
		data := log.Data[2:]
		if sign == UserOperationEventSign {

			event := UserOperationEvent{
				OpsHash:       topics[1],
				Sender:        strings.ToLower(utils.HexToAddress(topics[2])),
				Paymaster:     strings.ToLower(utils.HexToAddress(topics[3])),
				Nonce:         utils.HexToDecimal(utils.Substring(data, 0, 64*1)).Int64(),
				Success:       utils.HexToDecimalInt(utils.Substring(data, 64*1, 64*2)),
				ActualGasCost: utils.HexToDecimal(utils.Substring(data, 64*2, 64*3)).Int64(),
				ActualGasUsed: utils.HexToDecimal(utils.Substring(data, 64*3, 64*4)).Int64(),
			}
			events[event.Sender+strconv.Itoa(int(event.Nonce))] = event
		} else if sign == LogTransferEventSign {
			txValue := utils.DivRav(utils.HexToDecimal(utils.TruncateString(data, 64)).Int64())
			opsTxValMap[utils.HexToAddress(topics[2])] = txValue
		}

	}
	return events, opsTxValMap
}

func (t *_evmParser) getAddr(ctx context.Context, initCode string, paymasterAndData string) (string, string) {
	var factoryAddr string
	var paymaster string

	if len(initCode) > 0 {
		factory := utils.TruncateString(initCode, 42)
		if common.IsHexAddress(factory) {
			factoryAddr = utils.TruncateString(initCode, 42)
		} else {
			factoryAddr = ""
		}

	}

	if len(paymasterAndData) > 0 {
		paymaster = utils.TruncateString(paymasterAndData, 42)
		if common.IsHexAddress(paymaster) {
			paymaster = paymaster
		} else {
			paymaster = ""
		}
	}
	return factoryAddr, paymaster
}

func (t *_evmParser) parseCallData(ctx context.Context, client *ent.Client, network *ent.Network, callData string) ([]*CallDetail, string) {
	if len(callData) < 8 {
		return nil, ""
	}
	sign := utils.Substring(callData, 0, 10)
	paramData := utils.SubstringFromIndex(callData, 10)
	var callDetails []*CallDetail
	var source = ""
	switch sign {

	case ExecuteCall:
		callDetails = t.parseExecute(ctx, paramData)
		source = "executeCall"
		break
	case ExecuteBatchSign:
		callDetails = t.parseExecuteBatch(ctx, paramData)
		source = "executeBatch"
		break
	case ExecuteBatchCallSign:
		callDetails = t.parseExecuteBatchCall(ctx, paramData)
		source = "executeBatchCall"
		break
	default:
		callDetails = t.parseExecute(ctx, paramData)
		source = "execute"
		break
	}

	client, _ = entity.Client(ctx, network.ID)

	for _, detail := range callDetails {
		if len(detail.data) < 8 {
			continue
		}
		detail.source = detail.data[0:8]
		if detail.source == EmptyMethod {
			detail.source = ""
			continue
		}

		functionSignature, err := service.FunctionSignatureService.GetMethodBySignature(ctx, entity.MustClient(), detail.source)
		if err == nil {
			detail.source = functionSignature.Name
			continue
		}

		accountAbi, err := service.AccountService.GetAbiByAddress(ctx, client, detail.target)
		if err != nil {
			continue
		}
		detail.source = "0x" + detail.source
		method, err := accountAbi.MethodById(hexutil.MustDecode(detail.source))
		if err != nil {
			continue
		}
		detail.source = method.Name
	}

	if len(callDetails) > 0 {
		for _, detail := range callDetails {
			if len(detail.source) > 0 {
				source = detail.source
				break
			}
		}
	}

	return callDetails, source
}

func (t *_evmParser) parseExecuteBatchCall(ctx context.Context, paramData string) []*CallDetail {
	offset1 := utils.HexToDecimalInt(utils.Substring(paramData, 0, 64*1))
	offset2 := utils.HexToDecimalInt(utils.Substring(paramData, 64*1, 64*2))
	offset3 := utils.HexToDecimalInt(utils.Substring(paramData, 64*2, 64*3))
	num1 := utils.HexToDecimalInt(utils.Substring(paramData, offset1*2, offset1*2+64*1))
	var callDetails []*CallDetail
	for i := 1; i <= num1; i++ {
		target := utils.HexToAddress(utils.Substring(paramData, offset1*2+64*i, offset1*2+64*(i+1)))
		if len(target) < 1 {
			continue
		}
		value := utils.DivRav(utils.HexToDecimal(utils.Substring(paramData, offset2*2+64*i, offset2*2+64*(i+1))).Int64())
		data := utils.Substring(paramData, offset3*2+64*i, offset3*2+64*(i+1))
		callDetails = append(callDetails, &CallDetail{
			target: target,
			value:  value,
			data:   data,
		})
	}
	return callDetails
}

func (t *_evmParser) parseExecuteBatch(ctx context.Context, paramData string) []*CallDetail {
	offset1 := utils.HexToDecimalInt(utils.Substring(paramData, 0, 64*1))
	offset2 := utils.HexToDecimalInt(utils.Substring(paramData, 64*1, 64*2))
	offset3 := utils.HexToDecimalInt(utils.Substring(paramData, 64*2, 64*3))
	num1 := utils.HexToDecimalInt(utils.Substring(paramData, offset1*2, offset1*2+64*1))
	//num2 := utils.HexToDecimalInt(utils.Substring(paramData, offset2*2, offset2*2+64*1))
	//num3 := utils.HexToDecimalInt(utils.Substring(paramData, *offset3*2, *offset3*2+64*1))
	var callDetails []*CallDetail
	for i := 1; i <= num1; i++ {
		target := utils.HexToAddress(utils.Substring(paramData, offset1*2+64*i, offset1*2+64*(i+1)))
		if len(target) < 1 {
			continue
		}
		value := utils.DivRav(utils.HexToDecimal(utils.Substring(paramData, offset2*2+64*i, offset2*2+64*(i+1))).Int64())
		data := utils.Substring(paramData, offset3*2+64*i, offset3*2+64*(i+1))
		callDetails = append(callDetails, &CallDetail{
			target: target,
			value:  value,
			data:   data,
		})
	}
	return callDetails
}

func (t *_evmParser) parseExecute(ctx context.Context, paramData string) []*CallDetail {

	target := strings.ToLower(utils.HexToAddress(utils.Substring(paramData, 0, 64*1)))
	if len(target) < 1 {
		return nil
	}
	value := utils.DivRav(utils.HexToDecimal(utils.Substring(paramData, 64*1, 64*2)).Int64())
	offset := utils.HexToDecimalInt(utils.Substring(paramData, 64*2, 64*3))
	len := utils.HexToDecimalInt(utils.Substring(paramData, offset*2, offset*2+64*1))
	data := utils.Substring(paramData, offset*2+64*1, offset*2+64*1+len*2)
	var details []*CallDetail
	details = append(details, &CallDetail{
		target: target,
		value:  value,
		data:   data,
	})

	return details
}
