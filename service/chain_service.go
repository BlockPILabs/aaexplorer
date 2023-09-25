package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/assetchangetrace"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/blockdatadecode"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/blockscanrecord"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactiondecode"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactionreceiptdecode"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"strings"
	"time"
)

// 1.scan block, get active account, contract
// 2.add trace table
// 3.mev detect
const SWAP3 = "0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"
const SWAP2 = "0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"
const Transfer = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

type ReceiptLog struct {
	Data    string   `json:"data"`
	Topics  []string `json:"topics"`
	Address string   `json:"address"`
}

func ScanBlock() {
	//go scanBlock()
}

func scanBlock() {
	for {
		doScanBlock()
		time.Sleep(20 * time.Millisecond)
	}

}

func doScanBlock() {
	cli, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	networks, err := cli.Network.Query().All(context.Background())
	if len(networks) == 0 {
		return
	}
	for _, net := range networks {
		network := net.ID
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		records, err := client.BlockScanRecord.Query().Where(blockscanrecord.NetworkEQ(network)).All(context.Background())
		if len(records) == 0 {
			continue
		}
		record := records[0]
		last := record.LastBlockNumber
		cli, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		//1.get max block num by network
		//2.
		blockData, err := cli.BlockDataDecode.Query().Order(ent.Desc(blockdatadecode.FieldID)).Limit(1).All(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		if len(blockData) == 0 {
			continue
		}
		for i := last + 1; i <= blockData[0].ID; i++ {
			//do biz
			//task.MEVTask(i, network)
			transactions, err := cli.TransactionDecode.Query().Where(transactiondecode.BlockNumberEQ(i)).All(context.Background())
			if err != nil {
				continue
			}
			if len(transactions) == 0 {
				continue
			}
			receipts, err := cli.TransactionReceiptDecode.Query().Where(transactionreceiptdecode.BlockNumberEQ(i)).All(context.Background())
			if err != nil {
				continue
			}
			parseReceipt(network, receipts, cli)
			client.BlockScanRecord.Update().Where(blockscanrecord.IDEQ(record.ID)).SetLastBlockNumber(i).Exec(context.Background())
		}
	}

}

func parseReceipt(network string, receipts []*ent.TransactionReceiptDecode, cli *ent.Client) {
	for _, receipt := range receipts {
		logStr := receipt.Logs
		if len(logStr) == 0 {
			continue
		}
		logBytes := []byte(logStr)
		var receiptLogs []ReceiptLog
		err := json.Unmarshal(logBytes, &receiptLogs)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(receiptLogs) == 0 {
			continue
		}
		change := checkChange(receiptLogs)
		accountMap := make(map[string]int)
		tokenMap := make(map[string]int)
		for _, log := range receiptLogs {
			topics := log.Topics
			address := log.Address
			if len(topics) < 3 {
				continue
			}
			topic0 := topics[0]
			sender := hexToAddress(topics[1])
			target := hexToAddress(topics[2])
			if topic0 == Transfer {
				if sender != config.ZeroAddress {
					accountMap[sender] = 1
				}
				if target != config.ZeroAddress {
					accountMap[target] = 1
				}
				if change {
					tokenMap[address] = 1
					//saveTrace(network, address, config.AddressTypeToken, receipt, cli)
				}
			}
		}

		if len(accountMap) > 0 {
			for address, _ := range accountMap {
				traces, err := cli.AssetChangeTrace.Query().Where(assetchangetrace.AddressEqualFold(address), assetchangetrace.AddressTypeEQ(config.AddressTypeAccount)).All(context.Background())
				if err != nil {
					continue
				}
				if len(traces) == 0 {
					saveTrace(network, address, config.AddressTypeAccount, receipt, cli)
				} else {
					cli.AssetChangeTrace.Update().SetSyncFlag(0).Where(assetchangetrace.AddressEqualFold(address), assetchangetrace.AddressTypeEQ(config.AddressTypeAccount)).Exec(context.Background())
				}

			}
		}

		if len(tokenMap) > 0 {
			for address, _ := range tokenMap {
				traces, err := cli.AssetChangeTrace.Query().Where(assetchangetrace.AddressEqualFold(address), assetchangetrace.AddressTypeEQ(config.AddressTypeToken)).All(context.Background())
				if err != nil {
					continue
				}
				if len(traces) == 0 {
					saveTrace(network, address, config.AddressTypeToken, receipt, cli)
				} else {
					cli.AssetChangeTrace.Update().SetSyncFlag(0).Where(assetchangetrace.AddressEqualFold(address), assetchangetrace.AddressTypeEQ(config.AddressTypeToken)).Exec(context.Background())
				}

			}
		}

	}
}

func saveTrace(network string, address string, addressType int, receipt *ent.TransactionReceiptDecode, cli *ent.Client) {
	trace := cli.AssetChangeTrace.Create().
		SetNetwork(network).
		SetSyncFlag(0).
		SetTxHash(receipt.ID).
		SetBlockNumber(receipt.BlockNumber).
		SetLastChangeTime(0).
		SetAddress(address).
		SetAddressType(addressType)
	_, err := trace.Save(context.Background())
	if err != nil {
		fmt.Println(err)
	}

}

func checkChange(logs []ReceiptLog) bool {

	if len(logs) == 0 {
		return false
	}
	for _, log := range logs {
		topics := log.Topics
		if len(topics) == 0 {
			continue
		}
		topic0 := topics[0]
		if SWAP2 == topic0 || SWAP3 == topic0 {
			return true
		}
	}
	return false
}

func hexToAddress(hexStr string) string {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	address := strings.ToLower(common.HexToAddress(hexStr).String())
	return address
}
