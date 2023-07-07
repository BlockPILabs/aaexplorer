package parser

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity/schema"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const HandleOpsSign = "0x1fad948c"
const UserOperationEventSign = "0x49628fd1471006c1482da88028e9ce4dbb080b815c9b0344d39e5a8e6ec1419f"
const LogTransferEventSign = "0xe6497e3ee548a3372136af2fcb0696db31fc6cf20260707645068bd3fe97f3c4"

type UserOperationEvent struct {
	OpsHash       string
	Sender        string
	Paymaster     string
	Nonce         int64
	Success       int
	ActualGasCost int64
	ActualGasUsed int64
}

func ScanBlock() {
	fmt.Println("start execute task")
	client, err := ethclient.Dial("https://patient-crimson-slug.matic.discover.quiknode.pro/4cb47dc694ccf2998581548feed08af5477aa84b/")
	if err != nil {
		log.Fatal(err)
		return
	}
	record, err := schema.GetBlockScanRecordsByNetwork(config.Polygon)
	if err != nil {
		log.Fatal(err)
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

	block, err := client.BlockByNumber(context.Background(), big.NewInt(nextBlockNumber))
	if err != nil {
		log.Fatal(err)
		return
	}
	doParse(block, client)
	fmt.Println("parse end")

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
			continue
		}

		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			continue
		}
		//from, err := client.TransactionSender(context.Background(), tx, common.Hash{}, receipt.TransactionIndex)
		//if err != nil {
		//	log.Fatal(err)
		//	continue
		//}
		//fmt.Println(from)
		userOpsInfos, transactionInfo := parseUserOps(input, tx, receipt, block.Time())

		if userOpsInfos != nil {
			for _, userOps := range *userOpsInfos {
				blockUserOpsInfos = append(blockUserOpsInfos, userOps)
			}
		}

		if transactionInfo != nil {
			blockTransactionInfos = append(blockTransactionInfos, *transactionInfo)
		}

	}
	schema.BulkInsertUserOpsInfo(blockUserOpsInfos)
	schema.BulkInsertTransactions(blockTransactionInfos)
}

func parseUserOps(input string, tx *types.Transaction, receipt *types.Receipt, txTime uint64) (*[]schema.UserOpsInfo, *schema.TransactionInfo) {

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

	gasFee := tx.GasFeeCap()
	fmt.Println(gasFee)
	var userOpsInfos []schema.UserOpsInfo
	fee := (float64(receipt.GasUsed) / 1e18) * float64(receipt.EffectiveGasPrice.Int64()) / 1e18
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

		target := parseCallData(callData)

		userOpsInfo := schema.UserOpsInfo{
			TxHash:               tx.Hash().String(),
			BlockNumber:          receipt.BlockNumber.Int64(),
			Network:              config.Polygon,
			Sender:               sender,
			Target:               target,
			Bundler:              "",
			EntryPoint:           tx.To().String(),
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
			userOpsInfo.Fee = float64(opsVal.ActualGasUsed) / 1e18
		}
		opsTxValue, opsTxValueOk := opsValMap[sender]
		if opsTxValueOk {
			userOpsInfo.TxValue = opsTxValue
		}

		userOpsInfos = append(userOpsInfos, userOpsInfo)
	}

	transactionInfo := schema.TransactionInfo{
		TxHash:      tx.Hash().String(),
		BlockNumber: receipt.BlockNumber.Int64(),
		Network:     config.Polygon,
		Bundler:     "",
		EntryPoint:  tx.To().String(),
		UserOpsNum:  int64(arrNumInt),
		Fee:         fee,
		TxValue:     float64(tx.Value().Uint64()) / 1e18,
		GasPrice:    tx.GasPrice().String(),
		GasLimit:    int64(receipt.GasUsed),
		Status:      int(receipt.Status),
		TxTime:      int64(txTime),
		Beneficiary: beneficary,
		CreateTime:  time.Now(),
	}

	return &userOpsInfos, &transactionInfo
}

func parseLogs(logs []*types.Log) (map[string]UserOperationEvent, map[string]float64) {

	events := make(map[string]UserOperationEvent)
	opsTxValMap := make(map[string]float64)
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
			txValue := float64(hexToDecimal(truncateString(data, 64)).Int64()) / 1e18
			opsTxValMap[hexToAddress(topics[2].String())] = txValue
		}

	}
	return events, opsTxValMap
}

func getAddr(initCode string, paymasterAndData string) (string, string) {
	var factoryAddr string
	var paymaster string

	if len(initCode) > 0 {
		factoryAddr = "0x" + truncateString(initCode, 40)
	}

	if len(paymasterAndData) > 0 {
		paymaster = "0x" + truncateString(paymasterAndData, 40)
	}
	return factoryAddr, paymaster
}

func parseCallData(callData string) string {
	target := strings.ToLower(hexToAddress(substring(callData, 8, 8+64)))
	return target
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
