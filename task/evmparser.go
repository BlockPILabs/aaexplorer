package task

import (
	"bytes"
	"context"
	"encoding/json"
	"entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aablocksync"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/blockdatadecode"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactiondecode"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactionreceiptdecode"
	"github.com/BlockPILabs/aa-scan/internal/entity/schema"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/service"
	"github.com/BlockPILabs/aa-scan/internal/utils"
	"github.com/BlockPILabs/aa-scan/task/aa"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgtype"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
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
const ExecuteCall = "0x9e5d4c49"
const ExecuteBatchSign = "0x47e1da2a"
const ExecuteBatchCallSign = "0x912ccaa3"

type UserOperationEvent struct {
	OpsHash       string
	Sender        string
	Paymaster     string
	Nonce         int64
	Success       int
	ActualGasCost int64
	ActualGasUsed int64
	Target        string
	Factory       string
}

type CallDetail struct {
	target string
	value  decimal.Decimal
	data   string
}

type _evmParser struct {
	logger          log.Logger
	config          *config.Config
	startBlock      map[string]int64
	abi             abi.ABI
	handleOpsMethod *abi.Method
}

type parserBlock struct {
	block       *ent.BlockDataDecode
	transitions []*parserTransaction
	userOpInfo  *ent.AaBlockInfo
}
type parserTransaction struct {
	transaction     *ent.TransactionDecode
	receipt         *ent.TransactionReceiptDecode
	userOpInfo      *ent.AaBlockInfo
	userops         []*ent.AAUserOpsInfo
	logs            []*aa.Log
	userOpsCalldata []*ent.AAUserOpsCalldata
}

func InitEvmParse(config *config.Config, logger log.Logger) error {
	logger = logger.With("task", "evmparser")
	dayScheduler := chrono.NewDefaultTaskScheduler()
	t := _evmParser{
		logger:     logger,
		config:     config,
		startBlock: map[string]int64{},
	}

	for network, blockNumber := range t.config.EvmParser.StartBlock {
		t.startBlock[network] = blockNumber
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

	for fiend {
		fiend = false
		networks, err := service.NetworkService.GetNetworks(context.Background())
		if err != nil {
			t.logger.Error("network find error", "err", err)
			return
		}

		wg := &sync.WaitGroup{}
		for _, network := range networks {
			if _, ok := t.startBlock[network.ID]; !ok {
				t.startBlock[network.ID] = 0
			}
			wg.Add(1)
			ctx := log.WithContext(ctx, t.logger.With("network", network.ID))
			if t.ScanBlockByNetwork(ctx, network, wg) {
				fiend = true
			}
		}
		wg.Wait()
	}

	/*
		startMillis := getCurrentTimestampMillis()
		client, err := ethclient.Dial("https://patient-crimson-slug.matic.discover.quiknode.pro/4cb47dc694ccf2998581548feed08af5477aa84b/")
		if err != nil {
			log.Printf("RPC client err, %s\n", err)
			return
		}
		record, err := schema.GetBlockScanRecordsByNetwork(config.Polygon)
		if err != nil {
			return
		}
		var nextBlockNumber int64
		if record == nil {
			latestBlockNumber, err := client.BlockNumber(context.Background())
			if err != nil {
				log.Printf("evmparser--ScanBlock--err, %s\n", err)
				return
			}
			nextBlockNumber = int64(latestBlockNumber)
		} else {
			nextBlockNumber = record.LastBlockNumber + 1
		}
		nextBlockNumber = 44787202
		block, err := client.BlockByNumber(context.Background(), big.NewInt(nextBlockNumber))
		if err != nil {
			return
		}
		doParse(block, client)

		if record == nil {
			scanRecord := &schema.BlockScanRecord{
				Network:         config.Polygon,
				LastBlockNumber: nextBlockNumber,
				LastScanTime:    time.Now(),
				UpdateTime:      time.Now(),
			}
			schema.InsertBlockScanRecord(scanRecord)
		} else {
			record.LastBlockNumber = nextBlockNumber
			record.LastScanTime = time.Now()
			record.UpdateTime = time.Now()
			schema.UpdateBlockScanRecordByID(record.ID, record)
		}

		endMillis := getCurrentTimestampMillis()
		log.Printf("Block parse end, [%d], [%d]\n", nextBlockNumber, endMillis-startMillis)

	*/

}
func (t *_evmParser) ScanBlockByNetwork(ctx context.Context, network *ent.Network, wg *sync.WaitGroup) (fiend bool) {
	defer func() {
		if !fiend {
			wg.Done()
		}
	}()
	log.Context(ctx).Info("start block", "net", network)
	client, err := entity.Client(ctx, network.ID)
	if err != nil {
		log.Context(ctx).Error("network db client", "err", err)
		return false
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		log.Context(ctx).Error("network db client tx", "err", err)
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
		log.Context(ctx).Error("find AaBlockSync  tx", "err", err)
		return false
	}
	if len(aaBlockSyncs) < 1 {
		log.Context(ctx).Debug("not find AaBlockSync")
		return false
	}

	blockIds := make([]int64, len(aaBlockSyncs))
	for i, blockSync := range aaBlockSyncs {
		blockIds[i] = blockSync.ID
		t.startBlock[network.ID] = blockSync.ID
	}

	ctx = log.WithContext(context.Background(), log.Context(ctx))
	go func() {
		defer func() {
			tx.Commit()
			wg.Done()
		}()
		blockDataDecodes, transactionDecodes, receiptDecodes, blocksMap, transactionMap, err := t.getParseData(ctx, client, blockIds...)
		_ = (blockDataDecodes)
		_ = (transactionDecodes)
		_ = (receiptDecodes)
		_ = (blocksMap)
		_ = (transactionMap)
		if err != nil {
			log.Context(ctx).Error("get parse data error", "err", err)
		}
		for _, block := range blocksMap {
			t.doParse(network, block)
		}

	}()

	return true
}

func (t *_evmParser) getParseData(ctx context.Context, client *ent.Client, blockIds ...int64) (
	blockDataDecodes []*ent.BlockDataDecode,
	transactionDecodes []*ent.TransactionDecode,
	transactionReceiptDecodes []*ent.TransactionReceiptDecode,
	blocksMap map[int64]*parserBlock,
	transactionMap map[string]*parserTransaction,
	retErr error,
) {

	timeoutCtx, _ := context.WithTimeout(ctx, time.Minute)

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
		}
		return err
	})
	retErr = g.Wait()

	blocksMap = map[int64]*parserBlock{}
	transactionMap = map[string]*parserTransaction{}
	for _, blockDataDecode := range blockDataDecodes {
		blocksMap[blockDataDecode.ID] = &parserBlock{
			block: blockDataDecode,
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

func (t *_evmParser) doParse(network *ent.Network, block *parserBlock) {

	parserTransactions := block.transitions

	var blockUserOpsInfos []schema.UserOpsInfo
	var blockTransactionInfos []schema.TransactionInfo

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

		userOpsInfos, transactionInfo := t.parseUserOps(network, block, parserTx)

		if userOpsInfos != nil {
			for _, userOps := range *userOpsInfos {
				blockUserOpsInfos = append(blockUserOpsInfos, userOps)
			}
		}

		if transactionInfo != nil {
			blockTransactionInfos = append(blockTransactionInfos, *transactionInfo)
		}

	}
	t.insertUserOpsInfo(blockUserOpsInfos)
	t.insertTransactions(blockTransactionInfos)
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

func (t *_evmParser) insertTransactions(infos []schema.TransactionInfo) {
	if len(infos) == 0 {
		return
	}
	client, err := entity.Client(context.Background())
	if err != nil {
		return
	}

	var transactionInfoCreates []*ent.TransactionInfoCreate
	for _, tx := range infos {
		txCreate := client.TransactionInfo.Create().
			SetFee(tx.Fee).
			SetBundler(tx.Bundler).
			SetNetwork(tx.Network).
			SetStatus(tx.Status).
			SetTxTime(tx.TxTime).
			SetTxHash(tx.TxHash).
			SetTxTimeFormat(tx.TxTimeFormat).
			SetEntryPoint(tx.EntryPoint).
			SetBlockNumber(tx.BlockNumber).
			SetUserOpsNum(tx.UserOpsNum).
			SetBeneficiary(tx.Beneficiary).
			SetGasLimit(tx.GasLimit).
			SetGasPrice(tx.GasPrice).
			SetTxValue(tx.TxValue)

		transactionInfoCreates = append(transactionInfoCreates, txCreate)
	}
	_, err = client.TransactionInfo.CreateBulk(transactionInfoCreates...).Save(context.Background())
	if err != nil {
		t.logger.Info("", "err", err)
	}
}

func (t *_evmParser) insertUserOpsInfo(infos []schema.UserOpsInfo) {
	client, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	var userOpsInfoCreates []*ent.UserOpsInfoCreate
	for _, ops := range infos {
		userOpsCreate := client.UserOpsInfo.Create().
			SetUserOperationHash(ops.UserOperationHash).
			SetFactory(ops.Factory).
			SetPaymaster(ops.Paymaster).
			SetNetwork(ops.Network).
			SetBundler(ops.Bundler).
			SetActualGasCost(ops.ActualGasCost).
			SetActualGasUsed(ops.ActualGasUsed).
			SetBlockNumber(ops.BlockNumber).
			SetCalldata(ops.Calldata).
			SetCallGasLimit(ops.CallGasLimit).
			SetEntryPoint(ops.EntryPoint).
			SetFee(ops.Fee).
			SetInitCode(ops.InitCode).
			SetMaxFeePerGas(ops.MaxFeePerGas).
			SetMaxPriorityFeePerGas(ops.MaxPriorityFeePerGas).
			SetNonce(ops.Nonce).
			SetPaymasterAndData(ops.PaymasterAndData).
			SetPreVerificationGas(ops.PreVerificationGas).
			SetVerificationGasLimit(ops.VerificationGasLimit).
			SetTxValue(ops.TxValue).
			SetTxTimeFormat(ops.TxTimeFormat).
			SetTxTime(ops.TxTime).
			SetTxHash(ops.TxHash).
			SetTarget(ops.Target).
			SetStatus(ops.Status).
			SetSource(ops.Source).
			SetSignature(ops.Signature).
			SetSender(ops.Sender)
		userOpsInfoCreates = append(userOpsInfoCreates, userOpsCreate)
	}
	_, err = client.UserOpsInfo.CreateBulk(userOpsInfoCreates...).Save(context.Background())
	if err != nil {
		t.logger.Info("", "err", err)
	}
}

func (t *_evmParser) parseUserOps(network *ent.Network, block *parserBlock, parserTx *parserTransaction) (*[]schema.UserOpsInfo, *schema.TransactionInfo) {

	data, err := hexutil.Decode(parserTx.transaction.Input)
	if err != nil {
		return nil, nil
	}

	unpack, err := t.handleOpsMethod.Inputs.UnpackValues(data[4:])
	if err != nil {
		return nil, nil
	}
	if len(unpack) < 2 {
		return nil, nil
	}

	beneficiary := parserTx.transaction.FromAddr
	if beneficiaryAddr, ok := unpack[1].(common.Address); ok {
		beneficiary = beneficiaryAddr.Hex()
	}

	opsBytes, _ := json.Marshal(unpack[0])
	var ops []*aa.UserOperation
	_ = json.Unmarshal(opsBytes, &ops)
	err = json.Unmarshal([]byte(parserTx.receipt.Logs), &parserTx.logs)
	if err != nil {
		return nil, nil
	}

	events, opsValMap := t.parseLogs(parserTx.logs)
	fmt.Println(events)
	fmt.Println(opsValMap)

	var userOpsInfos []ent.AAUserOpsInfo
	for _, op := range ops {
		callDetails := t.parseCallData(hexutil.Encode(op.CallData))
		var target = ""
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
		userOpsInfo := ent.AAUserOpsInfo{
			ID:                   userOpHash.Hex(),
			Time:                 parserTx.transaction.Time,
			TxHash:               parserTx.transaction.ID,
			BlockNumber:          parserTx.transaction.BlockNumber,
			Network:              network.ID,
			Sender:               op.Sender.Hex(),
			Target:               target,
			Targets:              &pgTarges,
			TxValue:              parserTx.transaction.Value,
			Fee:                  parserTx.transaction.Gas,
			Bundler:              parserTx.transaction.FromAddr,
			EntryPoint:           parserTx.transaction.ToAddr,
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
			TxTime:               parserTx.transaction.Time,
			InitCode:             hexutil.Encode(op.InitCode),
			Status:               0,
			Source:               "",
			ActualGasCost:        0,
			ActualGasUsed:        0,
			CreateTime:           time.Now(),
			UsdAmount:            decimal.Decimal{},
		}

		userOpsInfos = append(userOpsInfos, userOpsInfo)
		fmt.Println(userOpsInfos)
	}

	fmt.Println(beneficiary)
	//
	//fmt.Println(ret[0])
	//
	//ops, ok := (ret[0]).([]*aa.UserOperation)
	//if !ok {
	//	return nil, nil
	//}
	//
	//for _, op := range ops {
	//	fmt.Println(op)
	//}
	//
	//fmt.Println(ops)
	/*
		start := utils.TruncateString(input, 64)
		startIdx := utils.HexToDecimal(start)
		beneficary := utils.HexToAddress(utils.Substring(input, 64*1, 64*2))
		startInt, err := strconv.Atoi(startIdx.String())
		if err != nil {
			return nil, nil
		}
		arrNum := utils.HexToDecimalInt(utils.Substring(input, startInt*2, startInt*2+64))
		if arrNum == nil {
			return nil, nil
		}
		var userOpsInfos []schema.UserOpsInfo
		fee := decimal.NewFromInt(int64(receipt.GasUsed)).Mul(utils.DivRav(receipt.EffectiveGasPrice.Int64()))
		logs := receipt.Logs

		events, opsValMap := t.parseLogs(logs)

		arrNumInt := *arrNum
		for i := 1; i <= arrNumInt; i++ {
			offset := utils.HexToDecimalInt(utils.Substring(input, 64*(2+i), 64*(i+3)))
			if offset == nil {
				continue
			}
			offsetInt := *offset
			realData := utils.SubstringFromIndex(input, 64*3+offsetInt*2)
			sender := utils.HexToAddress(utils.Substring(realData, 0, 64*1))
			nonce := utils.HexToDecimalInt(utils.Substring(realData, 64*1, 64*2))
			initCodeOffset := *utils.HexToDecimalInt(utils.Substring(realData, 64*2, 64*3))
			callDataOffset := *utils.HexToDecimalInt(utils.Substring(realData, 64*3, 64*4))
			callGasLimit := utils.HexToDecimal(utils.Substring(realData, 64*4, 64*5))
			verificationGasLimit := utils.HexToDecimal(utils.Substring(realData, 64*5, 64*6))
			preVerificationGas := utils.HexToDecimal(utils.Substring(realData, 64*6, 64*7))
			maxFeePerGas := utils.HexToDecimal(utils.Substring(realData, 64*7, 64*8))
			maxPriorityFeePerGas := utils.HexToDecimal(utils.Substring(realData, 64*8, 64*9))
			paymasterAndDataOffset := *utils.HexToDecimalInt(utils.Substring(realData, 64*9, 64*10))
			signatureOffset := *utils.HexToDecimalInt(utils.Substring(realData, 64*10, 64*11))
			initCodeLen := *utils.HexToDecimalInt(utils.Substring(realData, initCodeOffset*2, initCodeOffset*2+64))
			callDataLen := *utils.HexToDecimalInt(utils.Substring(realData, callDataOffset*2, callDataOffset*2+64))
			paymasterAndDataLen := *utils.HexToDecimalInt(utils.Substring(realData, paymasterAndDataOffset*2, paymasterAndDataOffset*2+64))
			signatureLen := *utils.HexToDecimalInt(utils.Substring(realData, signatureOffset*2, signatureOffset*2+64))
			initCode := utils.Substring(realData, initCodeOffset*2+64, initCodeOffset*2+64+initCodeLen*2)
			callData := utils.Substring(realData, callDataOffset*2+64, callDataOffset*2+64+callDataLen*2)
			paymasterAndData := utils.Substring(realData, paymasterAndDataOffset*2+64, paymasterAndDataOffset*2+64+paymasterAndDataLen*2)
			signature := utils.Substring(realData, signatureOffset*2+64, signatureOffset*2+64+signatureLen*2)

			factoryAddr, paymaster := t.getAddr(initCode, paymasterAndData)

			callDetails := t.parseCallData(callData)
			var target = ""
			if callDetails != nil && len(callDetails) > 0 {
				target = callDetails[0].target
			}

			userOpsInfo := schema.UserOpsInfo{
				TxHash:               tx.Hash().String(),
				BlockNumber:          receipt.BlockNumber.Int64(),
				Network:              config.Polygon,
				Sender:               sender,
				Target:               target,
				Bundler:              from,
				EntryPoint:           strings.ToLower(tx.To().String()),
				Factory:              factoryAddr,
				Paymaster:            paymaster,
				PaymasterAndData:     "0x" + paymasterAndData,
				Calldata:             "0x" + callData,
				Nonce:                int64(*nonce),
				CallGasLimit:         callGasLimit.Int64(),
				PreVerificationGas:   preVerificationGas.Int64(),
				VerificationGasLimit: verificationGasLimit.Int64(),
				MaxFeePerGas:         maxFeePerGas.Int64(),
				MaxPriorityFeePerGas: maxPriorityFeePerGas.Int64(),
				TxTime:               int64(txTime),
				TxTimeFormat:         utils.FormatTimestamp(int64(txTime)),
				InitCode:             "0x" + initCode,
				Source:               "",
				Signature:            "0x" + signature,
			}
			opsVal, ok := events[sender+strconv.Itoa(int(*nonce))]
			if ok {
				userOpsInfo.ActualGasUsed = opsVal.ActualGasUsed
				userOpsInfo.UserOperationHash = opsVal.OpsHash
				userOpsInfo.ActualGasCost = opsVal.ActualGasCost
				userOpsInfo.Status = opsVal.Success
				userOpsInfo.Fee = utils.DivRav(opsVal.ActualGasCost)
			}
			opsTxValue, opsTxValueOk := opsValMap[sender]
			if opsTxValueOk {
				userOpsInfo.TxValue = opsTxValue
			}

			//userOpsCalldatas := getUserOpsCalldatas(callDetails, tx, receipt, txTime, userOpsInfo.UserOperationHash, sender)

			userOpsInfos = append(userOpsInfos, userOpsInfo)
		}

		transactionInfo := schema.TransactionInfo{
			TxHash:       tx.Hash().String(),
			BlockNumber:  receipt.BlockNumber.Int64(),
			Network:      config.Polygon,
			Bundler:      from,
			EntryPoint:   strings.ToLower(tx.To().String()),
			UserOpsNum:   int64(arrNumInt),
			Fee:          fee,
			TxValue:      utils.DivRav(tx.Value().Int64()),
			GasPrice:     tx.GasPrice().String(),
			GasLimit:     int64(receipt.GasUsed),
			Status:       int(receipt.Status),
			TxTime:       int64(txTime),
			TxTimeFormat: utils.FormatTimestamp(int64(txTime)),
			Beneficiary:  beneficary,
			CreateTime:   time.Now(),
		}

		return &userOpsInfos, &transactionInfo
	*/

	return nil, nil
}

func (t *_evmParser) getUserOpsCalldatas(details []*CallDetail, tx *types.Transaction, receipt *types.Receipt, txTime uint64, userOpsHash, sender string) []*ent.UserOpsCalldataCreate {
	if len(details) == 0 {
		return nil
	}
	client, err := entity.Client(context.Background())
	if err != nil {
		t.logger.Info("UserOpsCalldata", "err", err)
	}
	var userOpsCalldatas []*ent.UserOpsCalldataCreate
	for _, detail := range details {
		opsCalldata := client.UserOpsCalldata.Create().
			SetUserOpsHash(userOpsHash).
			SetTxHash(tx.Hash().String()).
			SetTxValue(detail.value).
			SetBlockNumber(receipt.BlockNumber.Int64()).
			SetTxTime(int64(txTime)).
			SetSender(sender).
			SetSource("").
			SetTarget(detail.target).
			SetCalldata(detail.data)

		userOpsCalldatas = append(userOpsCalldatas, opsCalldata)
	}
	return userOpsCalldatas
}

func (t *_evmParser) parseLogs(logs []*aa.Log) (map[string]UserOperationEvent, map[string]decimal.Decimal) {

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
				Success:       *utils.HexToDecimalInt(utils.Substring(data, 64*1, 64*2)),
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

func (t *_evmParser) getAddr(initCode string, paymasterAndData string) (string, string) {
	var factoryAddr string
	var paymaster string

	if len(initCode) > 0 {
		factory := utils.TruncateString(initCode, 40)
		if utils.IsAddress(factory) {
			factoryAddr = "0x" + utils.TruncateString(initCode, 40)
		} else {
			factoryAddr = ""
		}

	}

	if len(paymasterAndData) > 0 {
		paymaster = utils.TruncateString(paymasterAndData, 40)
		if utils.IsAddress(paymaster) {
			paymaster = "0x" + paymaster
		} else {
			paymaster = ""
		}
	}
	return factoryAddr, paymaster
}

func (t *_evmParser) parseCallData(callData string) []*CallDetail {
	if len(callData) < 8 {
		return nil
	}
	sign := utils.Substring(callData, 0, 10)
	paramData := utils.SubstringFromIndex(callData, 10)
	var callDetails []*CallDetail
	switch sign {
	case ExecuteSign:
		callDetails = t.parseExecute(paramData)
		break
	case ExecuteCall:
		callDetails = t.parseExecute(paramData)
		break
	case ExecuteBatchSign:
		callDetails = t.parseExecuteBatch(paramData)
		break
	case ExecuteBatchCallSign:
		callDetails = t.parseExecuteBatchCall(paramData)
		break
	}
	return callDetails
}

func (t *_evmParser) parseExecuteBatchCall(paramData string) []*CallDetail {
	offset1 := utils.HexToDecimalInt(utils.Substring(paramData, 0, 64*1))
	offset2 := utils.HexToDecimalInt(utils.Substring(paramData, 64*1, 64*2))
	offset3 := utils.HexToDecimalInt(utils.Substring(paramData, 64*2, 64*3))
	num1 := utils.HexToDecimalInt(utils.Substring(paramData, *offset1*2, *offset1*2+64*1))
	var callDetails []*CallDetail
	for i := 1; i <= *num1; i++ {
		target := utils.HexToAddress(utils.Substring(paramData, *offset1*2+64*i, *offset1*2+64*(i+1)))
		value := utils.DivRav(utils.HexToDecimal(utils.Substring(paramData, *offset2*2+64*i, *offset2*2+64*(i+1))).Int64())
		data := utils.Substring(paramData, *offset3*2+64*i, *offset3*2+64*(i+1))
		callDetails = append(callDetails, &CallDetail{
			target: target,
			value:  value,
			data:   data,
		})
	}
	return callDetails
}

func (t *_evmParser) parseExecuteBatch(paramData string) []*CallDetail {
	offset1 := utils.HexToDecimalInt(utils.Substring(paramData, 0, 64*1))
	offset2 := utils.HexToDecimalInt(utils.Substring(paramData, 64*1, 64*2))
	offset3 := utils.HexToDecimalInt(utils.Substring(paramData, 64*2, 64*3))
	num1 := utils.HexToDecimalInt(utils.Substring(paramData, *offset1*2, *offset1*2+64*1))
	//num2 := utils.HexToDecimalInt(utils.Substring(paramData, *offset2*2, *offset2*2+64*1))
	//num3 := utils.HexToDecimalInt(utils.Substring(paramData, *offset3*2, *offset3*2+64*1))
	var callDetails []*CallDetail
	for i := 1; i <= *num1; i++ {
		target := utils.HexToAddress(utils.Substring(paramData, *offset1*2+64*i, *offset1*2+64*(i+1)))
		value := utils.DivRav(utils.HexToDecimal(utils.Substring(paramData, *offset2*2+64*i, *offset2*2+64*(i+1))).Int64())
		data := utils.Substring(paramData, *offset3*2+64*i, *offset3*2+64*(i+1))
		callDetails = append(callDetails, &CallDetail{
			target: target,
			value:  value,
			data:   data,
		})
	}
	return callDetails
}

func (t *_evmParser) parseExecute(paramData string) []*CallDetail {
	target := strings.ToLower(utils.HexToAddress(utils.Substring(paramData, 0, 64*1)))
	value := utils.DivRav(utils.HexToDecimal(utils.Substring(paramData, 64*1, 64*2)).Int64())
	offset := utils.HexToDecimalInt(utils.Substring(paramData, 64*2, 64*3))
	len := utils.HexToDecimalInt(utils.Substring(paramData, *offset*2, *offset*2+64*1))
	data := utils.Substring(paramData, *offset*2+64*1, *offset*2+64*1+*len*2)
	var details []*CallDetail
	details = append(details, &CallDetail{
		target: target,
		value:  value,
		data:   data,
	})

	return details
}
