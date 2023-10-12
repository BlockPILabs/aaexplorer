package task

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aatransactioninfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactiondecode"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactionreceiptdecode"
	"log"
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
		var i = 0
		for {
			aaInfos, err := client.AaTransactionInfo.Query().Order(aatransactioninfo.ByTime(sql.OrderAsc())).Limit(1).Offset(i).All(context.Background())
			if err != nil {
				break
			}
			if len(aaInfos) == 0 {
				break
			}
			i += 1
			aaInfo := aaInfos[0]
			txHash := aaInfo.ID
			txInfos, err := client.TransactionDecode.Query().Where(transactiondecode.IDEQ(txHash)).All(context.Background())
			if len(txInfos) == 0 {
				continue
			}
			txInfo := txInfos[0]
			copyTxProperties(aaInfo, txInfo)
			receiptInfos, err := client.TransactionReceiptDecode.Query().Where(transactionreceiptdecode.IDEQ(txHash)).All(context.Background())
			if len(receiptInfos) == 0 {
				client.AaTransactionInfo.Update().
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
					SetValue(txInfo.V).
					SetChainID(txInfo.ChainID).
					SetType(txInfo.Type).
					SetMaxFeePerGas(*txInfo.MaxFeePerGas).
					SetMaxPriorityFeePerGas(*txInfo.MaxPriorityFeePerGas).
					SetAccessList(txInfo.AccessList).
					SetMethod(txInfo.Method).
					SetStatus("0").Where(aatransactioninfo.IDEQ(aaInfo.ID)).Exec(context.Background())
				log.Printf("aa-transaction sync success part, hash:%s, network:%s", txHash, network)
				continue
			}
			receiptInfo := receiptInfos[0]
			copyReceiptProperties(aaInfo, receiptInfo)
			client.AaTransactionInfo.Update().
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
				SetValue(txInfo.V).
				SetChainID(txInfo.ChainID).
				SetType(txInfo.Type).
				SetMaxFeePerGas(*txInfo.MaxFeePerGas).
				SetMaxPriorityFeePerGas(*txInfo.MaxPriorityFeePerGas).
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
			log.Printf("aa-transaction sync success, hash:%s, network:%s", txHash, network)
		}

	}
}

func copyReceiptProperties(aaInfo *ent.AaTransactionInfo, receiptInfo *ent.TransactionReceiptDecode) {
	aaInfo.ContractAddress = receiptInfo.ContractAddress
	aaInfo.CumulativeGasUsed = receiptInfo.CumulativeGasUsed
	aaInfo.EffectiveGasPrice = receiptInfo.EffectiveGasPrice
	aaInfo.GasUsed = receiptInfo.GasUsed
	aaInfo.Logs = receiptInfo.Logs
	aaInfo.LogsBloom = receiptInfo.LogsBloom
	aaInfo.Status = receiptInfo.Status
}

func copyTxProperties(aaInfo *ent.AaTransactionInfo, txInfo *ent.TransactionDecode) {
	aaInfo.Nonce = txInfo.Nonce
	aaInfo.TransactionIndex = txInfo.TransactionIndex
	aaInfo.FromAddr = txInfo.FromAddr
	aaInfo.ToAddr = txInfo.ToAddr
	aaInfo.Value = txInfo.Value
	aaInfo.GasPrice = txInfo.GasPrice
	aaInfo.Gas = txInfo.Gas
	aaInfo.Input = txInfo.Input
	aaInfo.R = txInfo.R
	aaInfo.S = txInfo.S
	aaInfo.V = txInfo.V
	aaInfo.ChainID = txInfo.ChainID
	aaInfo.Type = txInfo.Type
	aaInfo.MaxFeePerGas = txInfo.MaxFeePerGas
	aaInfo.MaxPriorityFeePerGas = txInfo.MaxPriorityFeePerGas
	aaInfo.AccessList = txInfo.AccessList
	aaInfo.Method = txInfo.Method
}
