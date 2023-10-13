package task

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aatransactioninfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactiondecode"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactionreceiptdecode"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

func AATransactionFix() {
	cli, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	records, err := cli.Network.Query().All(context.Background())
	if len(records) == 0 {
		return
	}
	for _, record := range records {
		network := record.ID
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		now := time.Now()
		for {
			aaInfos, err := client.AaTransactionInfo.Query().Where(aatransactioninfo.NonceIsNil()).Order(aatransactioninfo.ByTime(sql.OrderAsc())).Limit(200).All(context.Background())
			if err != nil {
				continue
			}
			if len(aaInfos) == 0 {
				break
			}
			lastAaInfo := aaInfos[len(aaInfos)-1]
			if lastAaInfo.Time.Compare(now) > 0 {
				break
			}
			hashes, aaMap := getHashes(aaInfos)
			txInfos, err := client.TransactionDecode.Query().Where(transactiondecode.IDIn(hashes[:]...)).All(context.Background())
			if len(txInfos) == 0 {
				continue
			}
			txMap := getTxMap(txInfos)
			receiptInfos, err := client.TransactionReceiptDecode.Query().Where(transactionreceiptdecode.IDIn(hashes[:]...)).All(context.Background())
			receiptMap := getReceiptMap(receiptInfos)

			for _, hash := range hashes {

				txInfo := txMap[hash]
				aaInfo := aaMap[hash]
				var receiptInfo *ent.TransactionReceiptDecode
				if receiptMap == nil {
					receiptInfo = nil
				} else {
					receiptInfo = receiptMap[hash]
				}
				maxFeePerGas := txInfo.MaxFeePerGas
				if maxFeePerGas == nil {
					maxFeePerGas = &decimal.Zero
				}
				maxPriorityFeePerGas := txInfo.MaxPriorityFeePerGas
				if maxPriorityFeePerGas == nil {
					maxPriorityFeePerGas = &decimal.Zero
				}
				if receiptInfo == nil {
					err = client.AaTransactionInfo.Update().
						SetNonce(txInfo.Nonce).
						SetTransactionIndex(txInfo.TransactionIndex).
						SetFromAddr(txInfo.FromAddr).
						SetToAddr(txInfo.ToAddr).
						SetValue(txInfo.Value).
						SetGasPrice(txInfo.GasPrice).
						SetGas(txInfo.Gas).
						SetInput(txInfo.Input).
						SetR(txInfo.R).
						SetS(txInfo.S).
						SetV(txInfo.V).
						SetChainID(txInfo.ChainID).
						SetType(txInfo.Type).
						SetMaxFeePerGas(*maxFeePerGas).
						SetMaxPriorityFeePerGas(*maxPriorityFeePerGas).
						SetAccessList(txInfo.AccessList).
						SetMethod(txInfo.Method).
						SetStatus("0").Where(aatransactioninfo.IDEQ(aaInfo.ID)).Exec(context.Background())
					log.Printf("aa-transaction sync success part, hash:%s, network:%s", hash, network)
				} else {
					err = client.AaTransactionInfo.Update().
						SetNonce(txInfo.Nonce).
						SetTransactionIndex(txInfo.TransactionIndex).
						SetFromAddr(txInfo.FromAddr).
						SetToAddr(txInfo.ToAddr).
						SetValue(txInfo.Value).
						SetGasPrice(txInfo.GasPrice).
						SetGas(txInfo.Gas).
						SetInput(txInfo.Input).
						SetR(txInfo.R).
						SetS(txInfo.S).
						SetV(txInfo.V).
						SetChainID(txInfo.ChainID).
						SetType(txInfo.Type).
						SetMaxFeePerGas(*maxFeePerGas).
						SetMaxPriorityFeePerGas(*maxPriorityFeePerGas).
						SetAccessList(txInfo.AccessList).
						SetMethod(txInfo.Method).
						SetContractAddress(receiptInfo.ContractAddress).
						SetCumulativeGasUsed(receiptInfo.CumulativeGasUsed).
						SetEffectiveGasPrice(receiptInfo.EffectiveGasPrice).
						SetGasUsed(receiptInfo.GasUsed).
						SetLogs(receiptInfo.Logs).
						SetLogsBloom(receiptInfo.LogsBloom).
						SetStatus(receiptInfo.Status).
						Where(aatransactioninfo.IDEQ(aaInfo.ID)).Exec(context.Background())
					log.Printf("aa-transaction sync success, hash:%s, network:%s", hash, network)
				}
			}

		}

	}
}

func getTxMap(infos []*ent.TransactionDecode) map[string]*ent.TransactionDecode {
	txMap := make(map[string]*ent.TransactionDecode)
	if len(infos) == 0 {
		return txMap
	}
	for _, info := range infos {
		txMap[info.ID] = info
	}
	return txMap
}

func getHashes(infos []*ent.AaTransactionInfo) ([]string, map[string]*ent.AaTransactionInfo) {
	var hashArray []string
	hashMap := make(map[string]*ent.AaTransactionInfo)
	for _, info := range infos {
		hashArray = append(hashArray, info.ID)
		hashMap[info.ID] = info
	}

	return hashArray, hashMap
}

func copyReceiptProperties(aaInfo *ent.AaTransactionInfo, receiptInfo *ent.TransactionReceiptDecode) {
	aaInfo.ContractAddress = &receiptInfo.ContractAddress
	aaInfo.CumulativeGasUsed = &receiptInfo.CumulativeGasUsed
	aaInfo.EffectiveGasPrice = &receiptInfo.EffectiveGasPrice
	aaInfo.GasUsed = &receiptInfo.GasUsed
	aaInfo.Logs = &receiptInfo.Logs
	aaInfo.LogsBloom = &receiptInfo.LogsBloom
	aaInfo.Status = &receiptInfo.Status
}

func copyTxProperties(aaInfo *ent.AaTransactionInfo, txInfo *ent.TransactionDecode) {
	aaInfo.Nonce = &txInfo.Nonce
	aaInfo.TransactionIndex = &txInfo.TransactionIndex
	aaInfo.FromAddr = &txInfo.FromAddr
	aaInfo.ToAddr = &txInfo.ToAddr
	aaInfo.Value = &txInfo.Value
	aaInfo.GasPrice = &txInfo.GasPrice
	aaInfo.Gas = &txInfo.Gas
	aaInfo.Input = &txInfo.Input
	aaInfo.R = &txInfo.R
	aaInfo.S = &txInfo.S
	aaInfo.V = &txInfo.V
	aaInfo.ChainID = &txInfo.ChainID
	aaInfo.Type = &txInfo.Type
	aaInfo.MaxFeePerGas = txInfo.MaxFeePerGas
	aaInfo.MaxPriorityFeePerGas = txInfo.MaxPriorityFeePerGas
	aaInfo.AccessList = txInfo.AccessList
	aaInfo.Method = &txInfo.Method
}
