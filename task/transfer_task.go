package task

import (
	"context"
	"encoding/json"
	constConfig "github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/token"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/tokenall"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/transactiondecode"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/transactionreceiptdecode"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/transfertransaction"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/utils"
	"github.com/BlockPILabs/aaexplorer/task/aa"
	"github.com/chenzhijie/go-web3"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"math"
	"strings"
	"time"
)

const TokenAbi = "[{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

func InitTransferTask(ctx context.Context) {
	mevScheduler := chrono.NewDefaultTaskScheduler()
	_, err := mevScheduler.ScheduleWithCron(func(ctx context.Context) {
		TransferTask(ctx)
	}, "0/30 * * * * *")

	if err == nil {
		logger.Info("whaleHourStatistic has been scheduled")
	}

}

func TransferTask(ctx context.Context) {
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
		w3, err := web3.NewWeb3(net.HTTPRPC)
		if err != nil {
			logger.Error("TransferTask newWeb3 err ", "msg", err)
			continue
		}
		w3.Eth.SetChainId(net.ChainID)

		transferTxs, err := client.TransferTransaction.Query().Order(ent.Desc(transfertransaction.FieldBlockNumber)).Limit(1).All(ctx)
		lastBlockNum := int64(13920457)
		if err != nil {
			continue
		}
		if len(transferTxs) > 0 {
			lastBlockNum = transferTxs[0].BlockNumber
		}
		for {

			allReceipts, err := client.TransactionReceiptDecode.Query().Where(transactionreceiptdecode.BlockNumberGTE(lastBlockNum), transactionreceiptdecode.BlockNumberLT(lastBlockNum+10)).Order(ent.Asc(transactionreceiptdecode.FieldBlockNumber)).All(ctx)
			if err != nil {
				break
			}
			if len(allReceipts) == 0 {
				break
			}
			lastBlockNum = allReceipts[len(allReceipts)-1].BlockNumber + 1

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
					handleNativeTransfer(client, ctx, receipt)
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
					if sign == LogTransferEventSign {
						from := utils.HexToAddress(topics[1])
						to := utils.HexToAddress(topics[2])
						val := hexToDecimal(substring(data, 0, 64*1))
						curTokenAlls, err := client.TokenAll.Query().Where(tokenall.ContractAddressEqualFold(address)).All(ctx)
						if err != nil {
							continue
						}
						var tokenAll *ent.TokenAll
						if len(curTokenAlls) == 0 {
							tokenAll = addToken(ctx, client, address, w3, network)
						} else {
							tokenAll = curTokenAlls[0]
						}

						if tokenAll == nil {
							continue
						}

						count, err := client.TransferTransaction.Query().Where(transfertransaction.TxHashEQ(receipt.ID)).Count(ctx)
						if count > 0 {
							continue
						}

						decimals := tokenAll.Decimals
						amount := decimal.NewFromBigInt(val, 0).DivRound(decimal.NewFromFloat(math.Pow10(int(decimals))), int32(decimals))
						tx := client.TransferTransaction.Create().SetTime(receipt.Time).SetCreateTime(time.Now()).SetTxHash(receipt.ID).SetGasPrice(decimal.Zero).
							SetGas(receipt.GasUsed).SetValue(decimal.Zero).SetTransferValue(amount).SetFromAddr(from).SetToAddr(to).
							SetTransactionIndex(receipt.TransactionIndex).SetBlockNumber(receipt.BlockNumber).SetBlockHash(receipt.BlockHash).
							SetTokenAddress(tokenAll.ContractAddress).SetTokenSymbol(tokenAll.Symbol).SetTokenURL(tokenAll.ImageURL)

						_, err = tx.Save(ctx)
						if err == nil {
							logger.Info("TransferTask add tx success ", "txHash", receipt.ID)
						}
					}

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

	time.Sleep(time.Millisecond * 100)
	symbol := GetTokenSymbol(ctx, address, TokenAbi, w3)
	time.Sleep(time.Millisecond * 100)
	decimals := GetTokenDecimals(ctx, address, TokenAbi, w3)
	time.Sleep(time.Millisecond * 100)
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
