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
	"github.com/shopspring/decimal"
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
	Data    string
	Topics  []string
	Address string
}

func ScanBlock() {
	client, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	records, err := client.BlockScanRecord.Query().All(context.Background())
	if len(records) == 0 {
		return
	}
	for _, record := range records {
		network := record.Network
		last := record.LastBlockNumber
		cli, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		fmt.Println(network)
		//1.get max block num by network
		//2.
		blockData, err := cli.BlockDataDecode.Query().Order(ent.Desc(blockdatadecode.FieldNumber)).Limit(1).All(context.Background())
		if err != nil {
			continue
		}
		if len(blockData) == 0 {
			continue
		}
		for i := last + 1; i <= blockData[0].Number; i++ {
			//do biz
			transactions, err := cli.TransactionDecode.Query().Where(transactiondecode.BlockNumberEQ(decimal.NewFromInt(i))).All(context.Background())
			if err != nil {
				continue
			}
			if len(transactions) == 0 {
				continue
			}
			receipts, err := cli.TransactionReceiptDecode.Query().Where(transactionreceiptdecode.BlockNumberEQ(decimal.NewFromInt(int64(i)))).All(context.Background())
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
		var receiptLogs []*ReceiptLog
		err := json.Unmarshal(logBytes, receiptLogs)
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
			fmt.Println(topic0)
			if topic0 == Transfer {
				if sender != config.ZeroAddress {
					accountMap[sender] = 1
				}
				if target != config.ZeroAddress {
					accountMap[target] = 1
				}
				if change {
					tokenMap[address] = 1
					saveTrace(network, address, config.AddressTypeToken, receipt, cli)
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
		SetTxHash(receipt.TransactionHash).
		SetBlockNumber(receipt.BlockNumber.CoefficientInt64()).
		SetLastChangeTime(time.Now()).
		SetAddress(address).
		SetAddressType(addressType)
	trace.Save(context.Background())
}

func checkChange(logs []*ReceiptLog) bool {

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