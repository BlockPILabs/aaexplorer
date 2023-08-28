package parser

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/schema"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"log"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const HandleOpsSign = "0x1fad948c"
const UserOperationEventSign = "0x49628fd1471006c1482da88028e9ce4dbb080b815c9b0344d39e5a8e6ec1419f"
const LogTransferEventSign = "0xe6497e3ee548a3372136af2fcb0696db31fc6cf20260707645068bd3fe97f3c4"
const TransferEventSign = "0xe6497e3ee548a3372136af2fcb0696db31fc6cf20260707645068bd3fe97f3c4"
const AccountDeploySign = "0xd51a9c61267aa6196961883ecf5ff2da6619c37dac0fa92122513fb32c032d2d"

const ExecuteSign = "0xb61d27f6"
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

func ScanBlock() {

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

}

func getCurrentTimestampMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func doParse(block *types.Block, client *ethclient.Client) {

	transactions := block.Transactions()

	var blockUserOpsInfos []schema.UserOpsInfo
	var blockTransactionInfos []schema.TransactionInfo

	for _, tx := range transactions {
		input := "0x" + common.Bytes2Hex(tx.Data())
		if len(input) <= 10 {
			continue
		}
		sign := input[:10]
		input = input[10:]
		if sign != HandleOpsSign {
			addMevTx(tx, client)
			continue
		}

		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			continue
		}
		from := getFrom(tx, client)

		userOpsInfos, transactionInfo := parseUserOps(input, tx, receipt, block.Time(), from, block.BaseFee())

		if userOpsInfos != nil {
			for _, userOps := range *userOpsInfos {
				blockUserOpsInfos = append(blockUserOpsInfos, userOps)
			}
		}

		if transactionInfo != nil {
			blockTransactionInfos = append(blockTransactionInfos, *transactionInfo)
		}

	}
	insertUserOpsInfo(blockUserOpsInfos)
	insertTransactions(blockTransactionInfos)
}

func addMevTx(tx *types.Transaction, client *ethclient.Client) {
	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return
	}
	logs := receipt.Logs
	events, _ := parseLogs(logs)
	if len(events) == 0 {
		return
	}
	cli, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	for key, event := range events {
		sender := substring(key, 0, 42)
		nonce, err := strconv.Atoi(substringFromIndex(key, 42))
		if err != nil {
			continue
		}
		fmt.Println(event)
		fmt.Println(sender)
		fmt.Println(nonce)
		cli.AAUserOpsCalldata.Create().SetSender(sender).SetNetwork("").SetTxHash(tx.Hash().String()).SetTarget(event.Target).SetBlockNumber(receipt.BlockNumber.Int64()).Save(context.Background())
	}

}

func getFrom(tx *types.Transaction, client *ethclient.Client) string {
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Printf("get networkId err, %s\n", err)
		return ""
	}
	signer := types.LatestSignerForChainID(chainID)
	from, err := types.Sender(signer, tx)
	if err != nil {
		return ""
	}
	return strings.ToLower(from.String())
}

func insertTransactions(infos []schema.TransactionInfo) {
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
		log.Println(err)
	}
}

func insertUserOpsInfo(infos []schema.UserOpsInfo) {
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
		log.Println(err)
	}
}

func parseUserOps(input string, tx *types.Transaction, receipt *types.Receipt, txTime uint64, from string, baseFee *big.Int) (*[]schema.UserOpsInfo, *schema.TransactionInfo) {

	start := truncateString(input, 64)
	startIdx := hexToDecimal(start)
	beneficary := hexToAddress(substring(input, 64*1, 64*2))
	startInt, err := strconv.Atoi(startIdx.String())
	if err != nil {
		return nil, nil
	}
	arrNum := hexToDecimalInt(substring(input, startInt*2, startInt*2+64))
	if arrNum == nil {
		return nil, nil
	}

	var userOpsInfos []schema.UserOpsInfo
	fee := decimal.NewFromInt(int64(receipt.GasUsed)).Mul(DivRav(receipt.EffectiveGasPrice.Int64()))
	logs := receipt.Logs

	events, opsValMap := parseLogs(logs)

	arrNumInt := *arrNum
	for i := 1; i <= arrNumInt; i++ {
		offset := hexToDecimalInt(substring(input, 64*(2+i), 64*(i+3)))
		if offset == nil {
			continue
		}
		offsetInt := *offset
		realData := substringFromIndex(input, 64*3+offsetInt*2)
		sender := hexToAddress(substring(realData, 0, 64*1))
		nonce := hexToDecimalInt(substring(realData, 64*1, 64*2))
		initCodeOffset := *hexToDecimalInt(substring(realData, 64*2, 64*3))
		callDataOffset := *hexToDecimalInt(substring(realData, 64*3, 64*4))
		callGasLimit := hexToDecimal(substring(realData, 64*4, 64*5))
		verificationGasLimit := hexToDecimal(substring(realData, 64*5, 64*6))
		preVerificationGas := hexToDecimal(substring(realData, 64*6, 64*7))
		maxFeePerGas := hexToDecimal(substring(realData, 64*7, 64*8))
		maxPriorityFeePerGas := hexToDecimal(substring(realData, 64*8, 64*9))
		paymasterAndDataOffset := *hexToDecimalInt(substring(realData, 64*9, 64*10))
		signatureOffset := *hexToDecimalInt(substring(realData, 64*10, 64*11))
		initCodeLen := *hexToDecimalInt(substring(realData, initCodeOffset*2, initCodeOffset*2+64))
		callDataLen := *hexToDecimalInt(substring(realData, callDataOffset*2, callDataOffset*2+64))
		paymasterAndDataLen := *hexToDecimalInt(substring(realData, paymasterAndDataOffset*2, paymasterAndDataOffset*2+64))
		signatureLen := *hexToDecimalInt(substring(realData, signatureOffset*2, signatureOffset*2+64))
		initCode := substring(realData, initCodeOffset*2+64, initCodeOffset*2+64+initCodeLen*2)
		callData := substring(realData, callDataOffset*2+64, callDataOffset*2+64+callDataLen*2)
		paymasterAndData := substring(realData, paymasterAndDataOffset*2+64, paymasterAndDataOffset*2+64+paymasterAndDataLen*2)
		signature := substring(realData, signatureOffset*2+64, signatureOffset*2+64+signatureLen*2)

		factoryAddr, paymaster := getAddr(initCode, paymasterAndData)

		callDetails := parseCallData(callData)
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
			TxTimeFormat:         formatTimestamp(int64(txTime)),
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
			userOpsInfo.Fee = DivRav(opsVal.ActualGasCost)
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
		TxValue:      DivRav(tx.Value().Int64()),
		GasPrice:     tx.GasPrice().String(),
		GasLimit:     int64(receipt.GasUsed),
		Status:       int(receipt.Status),
		TxTime:       int64(txTime),
		TxTimeFormat: formatTimestamp(int64(txTime)),
		Beneficiary:  beneficary,
		CreateTime:   time.Now(),
	}

	return &userOpsInfos, &transactionInfo
}

func getUserOpsCalldatas(details []*CallDetail, tx *types.Transaction, receipt *types.Receipt, txTime uint64, userOpsHash, sender string) []*ent.UserOpsCalldataCreate {
	if len(details) == 0 {
		return nil
	}
	client, err := entity.Client(context.Background())
	if err != nil {
		log.Printf("UserOpsCalldata err, %s\n", err)
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

func DivRav(data int64) decimal.Decimal {
	return decimal.NewFromInt(data).DivRound(decimal.NewFromFloat(math.Pow10(18)), 20)
}

func formatTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	formatted := t.Format("2006-01-02 15:04:05")
	return formatted
}

func parseLogs(logs []*types.Log) (map[string]UserOperationEvent, map[string]decimal.Decimal) {

	events := make(map[string]UserOperationEvent)
	opsTxValMap := make(map[string]decimal.Decimal)
	for _, log := range logs {
		topics := log.Topics
		if len(topics) < 1 {
			continue
		}
		sign := topics[0].String()
		dataStr := hex.EncodeToString(log.Data)
		if len(dataStr) <= 2 {
			continue
		}
		data := hex.EncodeToString(log.Data)
		if sign == UserOperationEventSign {

			event := UserOperationEvent{
				OpsHash:       topics[1].String(),
				Sender:        strings.ToLower(hexToAddress(topics[2].String())),
				Paymaster:     strings.ToLower(hexToAddress(topics[3].String())),
				Nonce:         hexToDecimal(substring(data, 0, 64*1)).Int64(),
				Success:       *hexToDecimalInt(substring(data, 64*1, 64*2)),
				ActualGasCost: hexToDecimal(substring(data, 64*2, 64*3)).Int64(),
				ActualGasUsed: hexToDecimal(substring(data, 64*3, 64*4)).Int64(),
			}
			events[event.Sender+strconv.Itoa(int(event.Nonce))] = event
		} else if sign == LogTransferEventSign {
			txValue := DivRav(hexToDecimal(truncateString(data, 64)).Int64())
			opsTxValMap[hexToAddress(topics[2].String())] = txValue
		}

	}
	return events, opsTxValMap
}

func getAddr(initCode string, paymasterAndData string) (string, string) {
	var factoryAddr string
	var paymaster string

	if len(initCode) > 0 {
		factory := truncateString(initCode, 40)
		if isAddress(factory) {
			factoryAddr = "0x" + truncateString(initCode, 40)
		} else {
			factoryAddr = ""
		}

	}

	if len(paymasterAndData) > 0 {
		paymaster = truncateString(paymasterAndData, 40)
		if isAddress(paymaster) {
			paymaster = "0x" + paymaster
		} else {
			paymaster = ""
		}
	}
	return factoryAddr, paymaster
}

func isAddress(address string) bool {
	if len(address) != 40 {
		return false
	}
	for i := 0; i < 15; i++ {
		if address[i] != '0' {
			return true
		}
	}
	return false
}

func parseCallData(callData string) []*CallDetail {
	if len(callData) < 8 {
		return nil
	}
	sign := substring(callData, 0, 8)
	paramData := substringFromIndex(callData, 8)
	var callDetails []*CallDetail
	switch sign {
	case ExecuteSign:
		callDetails = parseExecute(paramData)
		break
	case ExecuteBatchSign:
		callDetails = parseExecuteBatch(paramData)
		break
	case ExecuteBatchCallSign:
		callDetails = parseExecuteBatchCall(paramData)
		break
	}
	return callDetails
}

func parseExecuteBatchCall(paramData string) []*CallDetail {
	offset1 := hexToDecimalInt(substring(paramData, 0, 64*1))
	offset2 := hexToDecimalInt(substring(paramData, 64*1, 64*2))
	offset3 := hexToDecimalInt(substring(paramData, 64*2, 64*3))
	num1 := hexToDecimalInt(substring(paramData, *offset1*2, *offset1*2+64*1))
	var callDetails []*CallDetail
	for i := 1; i <= *num1; i++ {
		target := hexToAddress(substring(paramData, *offset1*2+64*i, *offset1*2+64*(i+1)))
		value := DivRav(hexToDecimal(substring(paramData, *offset2*2+64*i, *offset2*2+64*(i+1))).Int64())
		data := substring(paramData, *offset3*2+64*i, *offset3*2+64*(i+1))
		callDetails = append(callDetails, &CallDetail{
			target: target,
			value:  value,
			data:   data,
		})
	}
	return callDetails
}

func parseExecuteBatch(paramData string) []*CallDetail {
	offset1 := hexToDecimalInt(substring(paramData, 0, 64*1))
	offset2 := hexToDecimalInt(substring(paramData, 64*1, 64*2))
	offset3 := hexToDecimalInt(substring(paramData, 64*2, 64*3))
	num1 := hexToDecimalInt(substring(paramData, *offset1*2, *offset1*2+64*1))
	//num2 := hexToDecimalInt(substring(paramData, *offset2*2, *offset2*2+64*1))
	//num3 := hexToDecimalInt(substring(paramData, *offset3*2, *offset3*2+64*1))
	var callDetails []*CallDetail
	for i := 1; i <= *num1; i++ {
		target := hexToAddress(substring(paramData, *offset1*2+64*i, *offset1*2+64*(i+1)))
		value := DivRav(hexToDecimal(substring(paramData, *offset2*2+64*i, *offset2*2+64*(i+1))).Int64())
		data := substring(paramData, *offset3*2+64*i, *offset3*2+64*(i+1))
		callDetails = append(callDetails, &CallDetail{
			target: target,
			value:  value,
			data:   data,
		})
	}
	return callDetails
}

func parseExecute(paramData string) []*CallDetail {
	target := strings.ToLower(hexToAddress(substring(paramData, 0, 64*1)))
	value := DivRav(hexToDecimal(substring(paramData, 64*1, 64*2)).Int64())
	offset := hexToDecimalInt(substring(paramData, 64*2, 64*3))
	len := hexToDecimalInt(substring(paramData, *offset*2, *offset*2+64*1))
	data := substring(paramData, *offset*2+64*1, *offset*2+64*1+*len*2)
	var details []*CallDetail
	details = append(details, &CallDetail{
		target: target,
		value:  value,
		data:   data,
	})

	return details
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length]
}

func hexToDecimal(hexStr string) *big.Int {
	hexStr = strings.TrimPrefix(hexStr, "0x")

	decimal := new(big.Int)
	_, success := decimal.SetString(hexStr, 16)
	if !success {
		return nil
	}

	return decimal
}

func hexToDecimalInt(hexStr string) *int {
	hexStr = strings.TrimPrefix(hexStr, "0x")

	decimal := new(big.Int)
	_, success := decimal.SetString(hexStr, 16)
	if !success {
		return nil
	}
	res, err := strconv.Atoi(decimal.String())
	if err != nil {
		return nil
	}
	return &res
}

func hexToAddress(hexStr string) string {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	address := strings.ToLower(common.HexToAddress(hexStr).String())
	return address
}

func substring(input string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(input) {
		end = len(input)
	}

	return input[start:end]
}

func substringFromIndex(input string, index int) string {
	if index < 0 || index >= len(input) {
		return ""
	}
	return input[index:]
}
