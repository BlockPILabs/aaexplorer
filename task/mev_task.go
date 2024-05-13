package task

import (
	"context"
	"encoding/json"
	"fmt"
	constConfig "github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaaccountdata"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aatransactioninfo"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aauseropsinfo"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/blockscanrecord"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/mevtransaction"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/token"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/transactiondecode"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/transactionreceiptdecode"
	"github.com/BlockPILabs/aaexplorer/task/aa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"log"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

func InitMEVTask(ctx context.Context) {
	mevScheduler := chrono.NewDefaultTaskScheduler()
	_, err := mevScheduler.ScheduleWithCron(func(ctx context.Context) {
		MEVTask(ctx)
	}, "0/30 * * * * *")

	if err == nil {
		logger.Info("whaleHourStatistic has been scheduled")
	}

}

func MEVTask0(blockNumber int64, network string) {

	client, err := entity.Client(context.Background())
	if err != nil {
		return

	}
	failedReceipts, err := client.TransactionReceiptDecode.Query().Where(transactionreceiptdecode.StatusEQ("0x0"), transactionreceiptdecode.BlockNumberEQ(blockNumber)).All(context.Background())
	if err != nil {
		return
	}
	if len(failedReceipts) == 0 {
		return
	}
	var failedHashes []string
	for _, receipt := range failedReceipts {
		failedHashes = append(failedHashes, receipt.ID)
	}
	failedOps, err := client.AAUserOpsInfo.Query().Where(aauseropsinfo.TxHashIn(failedHashes[:]...)).All(context.Background())
	if err != nil {
		return
	}
	if len(failedOps) == 0 {
		return
	}
	var failedMap = make(map[string]map[string]bool)
	for _, opsInfo := range failedOps {
		opsMap, opsMapOk := failedMap[opsInfo.TxHash]
		if !opsMapOk {
			opsMap = make(map[string]bool)
		}
		opsMap[opsInfo.Sender+":"+string(opsInfo.Nonce)] = true
		failedMap[opsInfo.TxHash] = opsMap
	}

	var mevResMap = make(map[string]string)
	for _, opsInfo := range failedOps {
		sender := opsInfo.Sender
		nonce := opsInfo.Nonce
		txHash := opsInfo.TxHash
		sameOps, err := client.AAUserOpsInfo.Query().Where(aauseropsinfo.SenderEqualFold(sender), aauseropsinfo.NonceEQ(nonce)).All(context.Background())
		if err != nil {
			continue
		}
		if len(sameOps) == 0 {
			continue
		}
		for _, same := range sameOps {
			if txHash == same.TxHash {
				continue
			}
			successReceipts, err := client.TransactionReceiptDecode.Query().
				Where(transactionreceiptdecode.ID(same.TxHash), transactionreceiptdecode.StatusEQ("0x1")).All(context.Background())
			if err != nil {
				continue
			}
			if len(successReceipts) == 0 {
				continue
			}

			successOps, err := client.AAUserOpsInfo.Query().Where(aauseropsinfo.TxHashEQ(successReceipts[0].ID)).All(context.Background())
			if err != nil {
				continue
			}
			if len(successOps) == 0 {
				continue
			}
			res := compareOps(successOps, failedMap[txHash])
			if res {
				mevResMap[successOps[0].TxHash] = txHash
			}
		}
	}

	if len(mevResMap) == 0 {
		return
	}
	fmt.Printf("mev check exist, %s", mapToString(mevResMap))

	for key, value := range mevResMap {
		_, err := client.MevInfo.Create().
			SetNetwork(network).
			SetTxHash(key).
			SetRelatedTxHash(value).
			SetBlockNumber(blockNumber).
			SetTxTime(1).
			SetTxFrom("").
			SetTxTo("").
			Save(context.Background())
		if err != nil {
			log.Println(err)
		}
	}

}

func compareOps(ops []*ent.AAUserOpsInfo, failMap map[string]bool) bool {
	if len(ops) != len(failMap) {
		return false
	}
	for _, opsInfo := range ops {
		key := opsInfo.Sender + ":" + string(opsInfo.Nonce)
		_, keyOk := failMap[key]
		if !keyOk {
			return false
		}
	}

	return true
}

func mapToString(myMap map[string]string) string {
	result := "{"
	for key, value := range myMap {
		result += fmt.Sprintf("%s: %s, ", key, value)
	}
	result = result[:len(result)-2] + "}"
	return result
}

func MEVTask(ctx context.Context) {
	cli, err := entity.Client(ctx)
	if err != nil {
		return
	}

	networks, err := cli.Network.Query().All(ctx)
	if len(networks) == 0 {
		return
	}

	for _, net := range networks {
		network := net.ID
		//if network != "ethereum" {
		//	continue
		//}
		client, err := entity.Client(ctx, network)
		if err != nil {
			continue
		}
		lastRecords, err := client.BlockScanRecord.Query().Where(blockscanrecord.TypeEQ("mev")).All(ctx)
		if err != nil {
			continue
		}
		if len(lastRecords) == 0 {
			continue
		}
		lastBlockNum := lastRecords[0].LastBlockNumber
		failedTxs, err := client.AaTransactionInfo.Query().Where(aatransactioninfo.StatusEQ("0x0"), aatransactioninfo.BlockNumberGT(lastBlockNum)).Order(ent.Asc(aatransactioninfo.FieldBlockNumber)).All(ctx)
		if len(failedTxs) == 0 {
			continue
		}
		tokens, err := client.Token.Query().Where(token.TypeEQ("base")).All(ctx)
		if len(tokens) == 0 {
			continue
		}
		tokenPrice := tokens[0].TokenPrice
		for _, tx := range failedTxs {
			blockNum := tx.BlockNumber
			txHash := tx.ID
			userOpsInfos, err := client.AAUserOpsInfo.Query().Where(aauseropsinfo.TxHashEqualFold(txHash)).All(ctx)
			if err != nil {
				continue
			}
			if len(userOpsInfos) == 0 {
				continue
			}
			blockNums := getLast3Block(blockNum)
			if blockNums == nil {
				continue
			}
			receipts, err := client.TransactionReceiptDecode.Query().Where(transactionreceiptdecode.BlockNumberIn(blockNums[:]...)).All(ctx)
			if len(receipts) == 0 {
				continue
			}

			victimReceipts, err := client.TransactionReceiptDecode.Query().Where(transactionreceiptdecode.IDEQ(txHash)).All(ctx)
			if len(victimReceipts) == 0 {
				continue
			}
			victimReceipt := victimReceipts[0]

			var userOpsArr []string
			for _, userOpsInfo := range userOpsInfos {
				nonce := userOpsInfo.Nonce
				sender := userOpsInfo.Sender
				key := sender + strconv.Itoa(int(nonce))
				userOpsArr = append(userOpsArr, key)
			}
		receiptFor:
			for _, receipt := range receipts {
				logs := receipt.Logs
				if len(logs) <= 2 {
					continue
				}
				if receipt.Status == "0x0" {
					continue
				}
				if receipt.ID == txHash {
					continue
				}
				var typeLogs []*aa.Log
				err := json.Unmarshal([]byte(logs), &typeLogs)
				if err != nil {
					continue
				}
				eventMap := parseLogs(typeLogs)
				if len(eventMap) == 0 {
					continue
				}
				for _, userOpsKey := range userOpsArr {
					event := eventMap[userOpsKey]
					if event != nil {

						oldMev, err := client.MevTransaction.Query().Where(mevtransaction.IDEQ(receipt.ID)).All(ctx)
						if len(oldMev) > 0 {
							break
						}
						successUserOps, err := client.AAUserOpsInfo.Query().Where(aauseropsinfo.TxHashEQ(receipt.ID)).All(ctx)
						if err != nil {
							continue
						}
						if len(successUserOps) == 0 {
							continue
						}
						var totalUserCost = decimal.Zero
						for _, oneUserOps := range successUserOps {
							totalUserCost = totalUserCost.Add(RayDiv(decimal.NewFromInt(oneUserOps.ActualGasCost)))
						}

						accounts, err := client.AaAccountData.Query().Where(aaaccountdata.IDEqualFold(tx.ID)).All(ctx)
						victimType := ""
						if len(accounts) > 0 {
							victimType = accounts[0].AaType
						}

						successTxs, err := client.TransactionDecode.Query().Where(transactiondecode.IDEQ(receipt.ID)).All(ctx)
						if len(successTxs) == 0 {
							continue
						}
						successTx := successTxs[0]

						attackerGas := successTx.GasPrice.DivRound(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18)), 18).Mul(receipt.GasUsed)
						victimGas := (*tx.GasPrice).DivRound(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18)), 18).Mul(victimReceipt.GasUsed)

						bundlerLossUsd := victimGas.Mul(tokenPrice)
						mevProfitUsd := totalUserCost.Sub(attackerGas).Mul(tokenPrice)

						mevTx := client.MevTransaction.Create().
							SetCreateTime(time.Now()).SetID(receipt.ID).SetTime(receipt.Time).SetFromAddr(receipt.FromAddr).
							SetToAddr(receipt.ToAddr).SetBlockNumber(receipt.BlockNumber).SetBlockHash(receipt.BlockHash).
							SetValue(decimal.Zero).SetGas(decimal.NewFromInt(receipt.CumulativeGasUsed)).SetGasPrice(*tx.GasPrice).SetTransactionIndex(*tx.TransactionIndex).SetVictim(*tx.FromAddr).
							SetAttacker(receipt.FromAddr).SetVictimTxHash(txHash).SetVictimFromAddr(*tx.FromAddr).SetVictimType(victimType).SetVictimToAddr(*tx.ToAddr).
							SetVictimBlockNumber(tx.BlockNumber).SetMevType(constConfig.MevFront).SetBundlerLoss(victimGas).SetBundlerLossUsd(bundlerLossUsd).SetMevProfit(totalUserCost.Sub(attackerGas)).
							SetMevProfitUsd(mevProfitUsd)
						logger.Info("find mev tx success ", "userHash", txHash, "mevHash", receipt.ID, "sender", userOpsKey)
						_, err = mevTx.Save(ctx)
						if err != nil {
							logger.Info("TestMEV save mev tx success ", "userHash", txHash, "mevHash", receipt.ID, "sender", userOpsKey)
						}
						break receiptFor
					}
				}

			}
			logger.Info("complete tx ", "hash", txHash)
		}

		err = client.BlockScanRecord.Update().SetLastBlockNumber(lastBlockNum).SetLastScanTime(time.Now()).Where(blockscanrecord.IDEQ(lastRecords[0].ID)).Exec(ctx)
		if err != nil {
			logger.Error("MEVTask update block num err ", "network", "blockNum", "msg", network, lastBlockNum, err)
		}
	}
}

func getLast3Block(num int64) []int64 {
	var nums []int64
	if num < 3 {
		return nil
	}
	nums = append(nums, num)
	nums = append(nums, num-1)
	nums = append(nums, num-2)

	return nums

}

func parseLogs(logs []*aa.Log) map[string]*UserOperationEvent {
	if len(logs) == 0 {
		return nil
	}
	events := make(map[string]*UserOperationEvent)
	for _, log := range logs {
		topics := log.Topics
		if len(topics) < 1 {
			continue
		}
		sign := topics[0]
		data := log.Data
		if len(data) <= 2 {
			continue
		}
		if sign == UserOperationEventSign {
			event := UserOperationEvent{
				OpsHash:       topics[1],
				Sender:        strings.ToLower(hexToAddress(topics[2])),
				Paymaster:     strings.ToLower(hexToAddress(topics[3])),
				Nonce:         hexToDecimal(substring(data, 0, 64*1)).Int64(),
				ActualGasCost: hexToDecimal(substring(data, 64*2, 64*3)).Int64(),
				ActualGasUsed: hexToDecimal(substring(data, 64*3, 64*4)).Int64(),
			}
			events[event.Sender+strconv.Itoa(int(event.Nonce))] = &event
		}

	}
	return events
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

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length]
}

func DivRav(data int64) decimal.Decimal {
	return decimal.NewFromInt(data).DivRound(decimal.NewFromFloat(math.Pow10(18)), 20)
}
