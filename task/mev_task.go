package task

import (
	"context"
	"fmt"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aauseropsinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactionreceipt"
)

func MEVTask() {

	client, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	var blockMumber = 11234455
	failedReceipts, err := client.TransactionReceipt.Query().Where(transactionreceipt.StatusEQ("0X0"), transactionreceipt.BlockNumEQ(int64(blockMumber))).All(context.Background())
	if err != nil {
		return
	}
	if len(failedReceipts) == 0 {
		return
	}
	var failedHashes []string
	for _, receipt := range failedReceipts {
		failedHashes = append(failedHashes, receipt.TransactionHash)
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
			successReceipts, err := client.TransactionReceipt.Query().
				Where(transactionreceipt.TransactionHashEQ(same.TxHash), transactionreceipt.StatusEQ("0x1")).All(context.Background())
			if err != nil {
				continue
			}
			if len(successReceipts) == 0 {
				continue
			}

			successOps, err := client.AAUserOpsInfo.Query().Where(aauseropsinfo.TxHashEQ(successReceipts[0].TransactionHash)).All(context.Background())
			if err != nil {
				continue
			}
			if len(successOps) == 0 {
				continue
			}
			res := compareOps(successOps, failedMap[same.TxHash])
			if res {
				mevResMap[successOps[0].TxHash] = txHash
			}
		}
	}

	if len(mevResMap) == 0 {
		return
	}
	fmt.Printf("mev check exist, %s", mapToString(mevResMap))

	//for key, value := range mevResMap {
	//	client.MevInfo.Create().SetNetwork().set
	//}

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
