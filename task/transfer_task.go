package task

import (
	"context"
	"encoding/json"
	"math"
	"strings"
	"time"

	constConfig "github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aatransactioninfo"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/token"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/transactiondecode"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/transactionreceiptdecode"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/transfertransaction"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/utils"
	"github.com/BlockPILabs/aaexplorer/task/aa"
	"github.com/chenzhijie/go-web3"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
)

const TokenAbi = "[{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
const SimpleTransferEventSign = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

func InitTransferTask(ctx context.Context) {
	//go TransferTaskOld(ctx)
	mevScheduler := chrono.NewDefaultTaskScheduler()
	_, err := mevScheduler.ScheduleWithCron(func(ctx context.Context) {
		TransferTaskNew(ctx)
	}, "0/30 * * * * *")

	if err == nil {
		logger.Info("TransferTask has been scheduled")
	}

}

func TransferTaskNew(ctx context.Context) {
	logger.Info("TransferTaskNew start.")
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
		logger.Info("TransferTaskNew network.", "network", network)
		//if network != "ethereum" {
		//	continue
		//}
		client, err := entity.Client(ctx, network)
		if err != nil {
			continue
		}
		w3, err := web3.NewWeb3(net.HTTPRPC)
		if err != nil {
			logger.Error("TransferTaskNew newWeb3 err ", "network", network, "msg", err)
			continue
		}
		w3.Eth.SetChainId(net.ChainID)

		transferTxs, err := client.TransferTransaction.Query().Order(ent.Desc(transfertransaction.FieldBlockNumber)).Limit(1).All(ctx)
		maxReceipts, err := client.TransactionReceiptDecode.Query().Order(ent.Desc(aatransactioninfo.FieldBlockNumber)).Limit(1).All(ctx)
		lastBlockNum := int64(104884700)
		maxBlockNum := int64(0)
		if len(maxReceipts) > 0 {
			maxBlockNum = maxReceipts[0].BlockNumber
		}
		if len(transferTxs) > 0 {
			lastBlockNum = transferTxs[0].BlockNumber
		}
		if err != nil {
			continue
		}
		aaAccounts, err := client.AaAccountData.Query().All(ctx)
		var accountMap = make(map[string]*ent.AaAccountData)
		if len(aaAccounts) > 0 {
			for _, aaAccount := range aaAccounts {
				accountMap[aaAccount.ID] = aaAccount
			}
		}
		tokenAlls, err := client.TokenAll.Query().All(ctx)
		var tokenMap = make(map[string]*ent.TokenAll)
		if len(tokenAlls) > 0 {
			for _, tokenAll := range tokenAlls {
				if len(tokenAll.ContractAddress) == 0 {
					continue
				}
				tokenMap[tokenAll.ContractAddress] = tokenAll
			}
		}
		logger.Info("TransferTaskNew get receipts ", "lastBlockNum", lastBlockNum, "maxBlock", maxBlockNum, "network", network)
		for {
			s0 := time.Now().UnixMilli()
			allReceipts, err := client.TransactionReceiptDecode.Query().Where(transactionreceiptdecode.BlockNumberGTE(lastBlockNum), transactionreceiptdecode.BlockNumberLT(lastBlockNum+20)).Order(ent.Asc(transactionreceiptdecode.FieldBlockNumber)).All(ctx)
			logger.Info("TransferTaskNew get receipts ", "size", len(allReceipts), "start", lastBlockNum, "end", lastBlockNum+20, "network", network)
			if err != nil {
				break
			}
			transferTxs, err := client.TransferTransaction.Query().Where(transfertransaction.BlockNumberGTE(lastBlockNum), transfertransaction.BlockNumberLT(lastBlockNum+20)).All(ctx)
			var transferTxMap = make(map[string]*ent.TransferTransaction)
			if len(transferTxs) > 0 {
				for _, transferTx := range transferTxs {
					transferTxMap[transferTx.TxHash] = transferTx
				}
			}
			lastBlockNum = lastBlockNum + 21
			if lastBlockNum > maxBlockNum {
				break
			}
			if len(allReceipts) == 0 {
				continue
			}
			var transferTxss []*ent.TransferTransactionCreate
			for _, receipt := range allReceipts {
				logs := receipt.Logs
				if len(logs) <= 2 {
					//continue
				}
				if receipt.Status == "0x0" {
					continue
				}
				var typeLogs []*aa.Log
				err := json.Unmarshal([]byte(logs), &typeLogs)
				if err != nil {
					continue
				}
				if len(typeLogs) == 0 {
					//handleNativeTransfer(client, ctx, receipt)
					continue
				}
				for _, log := range typeLogs {
					topics := log.Topics
					if len(topics) < 3 {
						continue
					}
					address := log.Address
					sign := topics[0]
					data := log.Data
					if len(data) <= 2 {
						continue
					}
					if sign == SimpleTransferEventSign {
						from := utils.HexToAddress(topics[1])
						to := utils.HexToAddress(topics[2])
						//s2 := time.Now().UnixMilli()
						//aaDatas, err := client.AaAccountData.Query().Where(aaaccountdata.IDIn(from, to)).All(ctx)
						//e2 := time.Now().UnixMilli()

						//logger.Info("TransferTaskNew complete step1 ", "spent", e2-s2)
						fromData := accountMap[from]
						toData := accountMap[to]

						if fromData == nil && toData == nil {
							logger.Info("TransferTaskNew data is null ", "hash", receipt.ID, "from", from, "to", to, "network", network)
							continue
						}

						val := hexToDecimal(substring(data, 0, 64*1))
						var tokenAll = tokenMap[strings.ToLower(address)]
						if tokenAll == nil {
							tokenAll = addToken(ctx, client, address, w3, network)
						}
						if tokenAll == nil {
							continue
						}
						tokenMap[strings.ToLower(address)] = tokenAll

						existTx := transferTxMap[receipt.ID]
						if existTx != nil {
							continue
						}

						//s4 := time.Now().UnixMilli()
						//count, err := client.TransferTransaction.Query().Where(transfertransaction.TxHashEQ(receipt.ID)).Count(ctx)
						//e4 := time.Now().UnixMilli()
						//logger.Info("TransferTaskNew complete step3 ", "spent", e4-s4)
						//if count > 0 {
						//	continue
						//}

						decimals := tokenAll.Decimals
						amount := decimal.NewFromBigInt(val, 0).DivRound(decimal.NewFromFloat(math.Pow10(int(decimals))), int32(decimals))
						tx := client.TransferTransaction.Create().SetTime(receipt.Time).SetCreateTime(time.Now()).SetTxHash(receipt.ID).SetGasPrice(decimal.Zero).
							SetGas(receipt.GasUsed).SetValue(decimal.Zero).SetTransferValue(amount).SetFromAddr(from).SetToAddr(to).
							SetTransactionIndex(receipt.TransactionIndex).SetBlockNumber(receipt.BlockNumber).SetBlockHash(receipt.BlockHash).
							SetTokenAddress(tokenAll.ContractAddress).SetTokenSymbol(tokenAll.Symbol).SetTokenURL(tokenAll.ImageURL)

						transferTxss = append(transferTxss, tx)
						//_, err = tx.Save(ctx)
						//if err == nil {
						//	logger.Info("TransferTaskNew add tx success ", "txHash", receipt.ID)
						//}

						//logger.Info("TransferTaskNew complete once ", "spent", e1-s1)
					}

				}
			}
			if len(transferTxss) > 0 {
				s5 := time.Now().UnixMilli()
				_, err := client.TransferTransaction.CreateBulk(transferTxss[:]...).Save(ctx)
				e5 := time.Now().UnixMilli()
				logger.Info("TransferTaskNew complete step4 ", "spent", e5-s5, "network", network)
				e0 := time.Now().UnixMilli()
				if err == nil {
					logger.Info("TransferTaskNew complete all ", "spent", e0-s0, "network", network)
				}
			}

		}

	}
}

func TransferTaskOld(ctx context.Context) {
	logger.Info("TransferTaskOld start.")
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
		if network != "ethereum" {
			continue
		}
		client, err := entity.Client(ctx, network)
		if err != nil {
			continue
		}

		lastBlockNum := int64(19380056)
		maxBlockNum := int64(20267076)
		if err != nil {
			continue
		}
		aaAccounts, err := client.AaAccountData.Query().All(ctx)
		var accountMap = make(map[string]*ent.AaAccountData)
		if len(aaAccounts) > 0 {
			for _, aaAccount := range aaAccounts {
				accountMap[aaAccount.ID] = aaAccount
			}
		}
		tokenAlls, err := client.TokenAll.Query().All(ctx)
		var tokenMap = make(map[string]*ent.TokenAll)
		if len(tokenAlls) > 0 {
			for _, tokenAll := range tokenAlls {
				if len(tokenAll.ContractAddress) == 0 {
					continue
				}
				tokenMap[tokenAll.ContractAddress] = tokenAll
			}
		}
		logger.Info("TransferTaskOld blockNum ", "lastBlockNum", lastBlockNum, "maxBlock", maxBlockNum)
		for {
			s0 := time.Now().UnixMilli()
			allReceipts, err := client.TransactionReceiptDecode.Query().Where(transactionreceiptdecode.BlockNumberGTE(lastBlockNum), transactionreceiptdecode.BlockNumberLT(lastBlockNum+20)).Order(ent.Asc(transactionreceiptdecode.FieldBlockNumber)).All(ctx)
			e0 := time.Now().UnixMilli()
			logger.Info("TransferTaskOld get receipts ", "size", len(allReceipts), "spent", e0-s0)
			if err != nil {
				break
			}
			transferTxs, err := client.TransferTransaction.Query().Where(transfertransaction.BlockNumberGTE(lastBlockNum), transfertransaction.BlockNumberLT(lastBlockNum+20)).All(ctx)
			var transferTxMap = make(map[string]*ent.TransferTransaction)
			if len(transferTxs) > 0 {
				for _, transferTx := range transferTxs {
					transferTxMap[transferTx.TxHash] = transferTx
				}
			}
			lastBlockNum = lastBlockNum + 21
			if lastBlockNum > maxBlockNum {
				break
			}
			if len(allReceipts) == 0 {
				continue
			}
			var transferTxss []*ent.TransferTransactionCreate
			for _, receipt := range allReceipts {
				logs := receipt.Logs
				if len(logs) <= 2 {
					continue
				}
				if receipt.Status == "0x0" {
					continue
				}
				var typeLogs []*aa.Log
				err := json.Unmarshal([]byte(logs), &typeLogs)
				if err != nil {
					continue
				}
				if len(typeLogs) == 0 {
					//handleNativeTransfer(client, ctx, receipt)
					continue
				}
				for _, log := range typeLogs {
					topics := log.Topics
					if len(topics) < 3 {
						continue
					}
					address := log.Address
					sign := topics[0]
					data := log.Data
					if len(data) <= 2 {
						continue
					}

					if sign == SimpleTransferEventSign {
						from := utils.HexToAddress(topics[1])
						to := utils.HexToAddress(topics[2])
						fromData := accountMap[from]
						toData := accountMap[to]

						if fromData == nil && toData == nil {
							continue
						}

						val := hexToDecimal(substring(data, 0, 64*1))
						//curTokenAlls, err := client.TokenAll.Query().Where(tokenall.ContractAddressEqualFold(address)).All(ctx)
						//if err != nil {
						//	continue
						//}
						var tokenAll = tokenMap[strings.ToLower(address)]
						if tokenAll == nil {
							//tokenAll = addToken(ctx, client, address, w3, network)
							continue
						}

						existTx := transferTxMap[receipt.ID]
						if existTx != nil {
							continue
						}
						//count, _ := client.TransferTransaction.Query().Where(transfertransaction.TxHashEQ(receipt.ID)).Count(ctx)
						//if count > 0 {
						//	continue
						//}

						decimals := tokenAll.Decimals
						amount := decimal.NewFromBigInt(val, 0).DivRound(decimal.NewFromFloat(math.Pow10(int(decimals))), int32(decimals))
						tx := client.TransferTransaction.Create().SetTime(receipt.Time).SetCreateTime(time.Now()).SetTxHash(receipt.ID).SetGasPrice(decimal.Zero).
							SetGas(receipt.GasUsed).SetValue(decimal.Zero).SetTransferValue(amount).SetFromAddr(from).SetToAddr(to).
							SetTransactionIndex(receipt.TransactionIndex).SetBlockNumber(receipt.BlockNumber).SetBlockHash(receipt.BlockHash).
							SetTokenAddress(tokenAll.ContractAddress).SetTokenSymbol(tokenAll.Symbol).SetTokenURL(tokenAll.ImageURL)

						transferTxss = append(transferTxss, tx)
					}

				}
			}

			if len(transferTxss) > 0 {
				_, err := client.TransferTransaction.CreateBulk(transferTxss[:]...).Save(ctx)
				e0 := time.Now().UnixMilli()
				if err == nil {
					logger.Info("TransferTaskOld batch complete all ", "spent", e0-s0)
				}
			}

		}

	}
}

func addToken(ctx context.Context, client *ent.Client, address string, w3 *web3.Web3, network string) *ent.TokenAll {
	if len(address) == 0 {
		return nil
	}
	tokens, err := client.Token.Query().Where(token.ContractAddressEqualFold(address)).All(ctx)
	imageUrl := ""
	if len(tokens) > 0 {
		imageUrl = tokens[0].ImageURL
	}

	symbol := GetTokenSymbol(ctx, address, TokenAbi, w3)
	decimals := GetTokenDecimals(ctx, address, TokenAbi, w3)
	tokenName := GetTokenName(ctx, address, TokenAbi, w3)
	address = strings.ToLower(address)
	tokenAll := client.TokenAll.Create().SetTokenPrice(decimal.Zero).SetNetwork(network).SetLastTime(time.Now().UnixMilli()).
		SetSymbol(symbol).SetDecimals(decimals).SetFullName(tokenName).SetImageURL(imageUrl).SetType("").SetMarketRank(1).
		SetCreateTime(time.Now()).SetUpdateTime(time.Now()).SetContractAddress(address)
	res, err := tokenAll.Save(ctx)
	if err == nil {
		logger.Info("TransferTask add token success ", "tokenAll", tokenAll)
		return res
	}
	return nil
}

func handleNativeTransfer(client *ent.Client, ctx context.Context, receipt *ent.TransactionReceiptDecode) {

	txs, err := client.TransactionDecode.Query().Where(transactiondecode.IDEQ(receipt.ID)).All(ctx)
	if err != nil {
		return
	}
	if len(txs) == 0 {
		return
	}
	input := txs[0].Input
	if input != "0x" {
		return
	}
	tokens, err := client.Token.Query().Where(token.TypeEQ("native")).All(ctx)
	if len(tokens) == 0 {
		return
	}
	val := txs[0].Value.DivRound(decimal.NewFromFloat(math.Pow10(18)), 18)
	tx := client.TransferTransaction.Create().SetTime(txs[0].Time).SetCreateTime(time.Now()).SetTxHash(receipt.ID).SetGasPrice(txs[0].GasPrice).
		SetGas(receipt.GasUsed).SetValue(val).SetTransferValue(val).SetFromAddr(txs[0].FromAddr).SetToAddr(txs[0].ToAddr).
		SetTransactionIndex(receipt.TransactionIndex).SetBlockNumber(receipt.BlockNumber).SetBlockHash(receipt.BlockHash).
		SetTokenAddress(tokens[0].ContractAddress).SetTokenSymbol(tokens[0].Symbol).SetTokenURL(tokens[0].ImageURL)
	_, err = tx.Save(ctx)
	if err == nil {
		logger.Info("TransferTask add native tx success ", "txHash", receipt.ID)
	}
}

func GetTokenDecimals(ctx context.Context, tokenAddress string, abi string, web3 *web3.Web3) int64 {

	contract, err := web3.Eth.NewContract(abi, tokenAddress)

	res, err := contract.Call("decimals")
	if err != nil {
		log.Context(ctx).Error("GetDecimals decimals err ", "msg", err)
		return constConfig.DefaultDecimals
	}
	value, ok := res.(uint8)
	if ok {
		val := int64(value)
		return val
	}
	return constConfig.DefaultDecimals
}

func GetTokenName(ctx context.Context, tokenAddress string, abi string, web3 *web3.Web3) string {

	contract, err := web3.Eth.NewContract(abi, tokenAddress)

	res, err := contract.Call("name")
	if err != nil {
		log.Context(ctx).Error("GetTokenName name err ", "msg", err)
		return ""
	}
	value, ok := res.(string)
	if ok {
		return value
	}
	return ""
}

func GetTokenSymbol(ctx context.Context, tokenAddress string, abi string, web3 *web3.Web3) string {

	contract, err := web3.Eth.NewContract(abi, tokenAddress)

	res, err := contract.Call("symbol")
	if err != nil {
		log.Context(ctx).Error("GetTokenSymbol symbol err ", "msg", err)
		return ""
	}
	value, ok := res.(string)
	if ok {
		return value
	}
	return ""
}
