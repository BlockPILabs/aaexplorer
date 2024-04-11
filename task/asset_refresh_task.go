package task

import (
	"context"
	constConfig "github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaasset"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaassetdetail"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/token"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/chenzhijie/go-web3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"math/big"
	"time"
)

const Abi = "[{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

func InitAssetRefreshTask() {
	go AssetRefreshTask(context.Background())
	hourScheduler := chrono.NewDefaultTaskScheduler()
	_, err := hourScheduler.ScheduleWithCron(func(ctx context.Context) {
		AssetRefreshTask(ctx)
	}, "0 35 0 1/6 * *")

	if err == nil {
		logger.Info("AssetRefreshTask has been scheduled")
	}
}

func AssetRefreshTask(ctx context.Context) {
	cli, err := entity.Client(ctx)
	if err != nil {
		logger.Error("AssetRefreshTask err, ", "msg", err)
		return
	}
	networks, err := cli.Network.Query().All(ctx)
	if err != nil {
		return
	}
	if len(networks) == 0 {
		return
	}

	for _, net := range networks {
		network := net.ID
		client, err := entity.Client(ctx, network)
		if err != nil {
			continue
		}
		lastTime := time.Now().UnixMilli() - constConfig.WhaleTxDay*24*3600*1000
		aas, err := client.AaAsset.Query().Where(aaasset.LastTimeLT(lastTime)).All(ctx)
		if err != nil {
			logger.Error("AssetRefreshTask query asset err ", "msg", err)
			continue
		}
		if len(aas) == 0 {
			continue
		}
		w3, err := web3.NewWeb3(net.HTTPRPC)
		if err != nil {
			logger.Error("AssetRefreshTask newWeb3 err ", "msg", err)
			continue
		}
		w3.Eth.SetChainId(net.ChainID)
		tokens, err := client.Token.Query().All(ctx)
		if len(tokens) == 0 {
			continue
		}
		blockNum, err := w3.Eth.GetBlockNumber()
		if err != nil {
			logger.Error("AssetRefreshTask blockNum err ", "msg", err)
			continue
		}
		oneSize := len(aas) / 10
		var allArrs [][]*ent.AaAsset
		for i := 0; i <= 10; i++ {
			var oneArr []*ent.AaAsset
			allArrs = append(allArrs, oneArr)
		}
		logger.Info("AssetRefreshTask-arrSize ", "size", len(allArrs))
		start := 0
		for idx, aa := range aas {
			r := idx % oneSize
			if r == 0 && idx != 0 {
				logger.Info("AssetRefreshTask-oneSize ", "size", len(allArrs[start]), "start", start)
				go doRefresh(ctx, client, tokens, w3, blockNum, network, allArrs[start])
				start = start + 1
			}
			allArrs[start] = append(allArrs[start], aa)
		}
	}
}

func doRefresh(ctx context.Context, client *ent.Client, tokens []*ent.Token, w3 *web3.Web3, blockNum uint64, network string, assets []*ent.AaAsset) {
	if len(assets) == 0 {
		return
	}
	logger.Info("AssetRefreshTask-doRefresh start.")
	for _, aa := range assets {
		userAddress := aa.ID
		totalValue := decimal.Zero
		for _, token := range tokens {
			contractAddress := token.ContractAddress
			if len(contractAddress) == 0 {
				continue
			}
			decimals := GetDecimals(ctx, contractAddress, w3)
			balance := GetBalance(ctx, contractAddress, userAddress, decimals, w3)
			assetValue := balance.Mul(token.TokenPrice)
			addOrUpdateAssetDetail(ctx, contractAddress, userAddress, assetValue, client, network, token.Symbol, balance)
			totalValue = totalValue.Add(assetValue)
		}

		balance, err := w3.Eth.GetBalance(common.HexToAddress(userAddress), big.NewInt(int64(blockNum)))
		if err != nil {
			logger.Error("AssetRefreshTask get balance err ", "user", userAddress, "network", network, "msg", err)
			continue
		}
		nativeBalance := decimal.NewFromBigInt(balance, 0).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(constConfig.DefaultDecimals)))
		nativeTokens, err := client.Token.Query().Where(token.TypeEQ("base"), token.NetworkEQ(network)).Limit(1).All(ctx)
		nativePrice := decimal.Zero
		if len(nativeTokens) > 0 {
			nativePrice = nativeTokens[0].TokenPrice
		}
		nativeValue := nativeBalance.Mul(nativePrice)
		totalValue = totalValue.Add(nativeValue)
		client.AaAsset.Update().SetAssetValue(totalValue).SetLastTime(time.Now().UnixMilli()).SetBalance(nativeBalance).Where(aaasset.IDEqualFold(userAddress)).Exec(ctx)
		if nativeBalance.Cmp(decimal.Zero) > 0 {
			logger.Info("AssetRefreshTask update balance success, ", "userAddress", aa.ID, "network", network)
		} else {
			logger.Info("AssetRefreshTask update balance success empty, ", "userAddress", aa.ID, "network", network)
		}
	}

}

func addOrUpdateAssetDetail(ctx context.Context, contractAddress string, userAddress string, value decimal.Decimal, client *ent.Client, network string, symbol string, amount decimal.Decimal) {
	details, err := client.AaAssetDetail.Query().Where(aaassetdetail.UserAddressEqualFold(userAddress), aaassetdetail.ContractAddressEqualFold(contractAddress)).All(ctx)
	if err != nil {
		return
	}
	if len(details) > 0 {
		detail := details[0]
		client.AaAssetDetail.Update().SetAssetValue(value).SetLastTime(time.Now().UnixMilli()).Where(aaassetdetail.IDEQ(detail.ID)).Exec(ctx)
	} else {
		detail := client.AaAssetDetail.Create().SetAssetValue(value).SetLastTime(time.Now().UnixMilli()).SetCreateTime(time.Now()).SetUpdateTime(time.Now()).SetNetwork(network).SetContractAddress(contractAddress).SetSymbol(symbol).SetUserAddress(userAddress).SetAssetAmount(amount).SetIsNative("false")
		_, err := detail.Save(ctx)
		if err != nil {
			logger.Error("addOrUpdateAssetDetail err ", "symbol", symbol, "msg", err)
		}
	}
}

func GetBalance(ctx context.Context, tokenAddress string, userAddress string, decimals int64, web3 *web3.Web3) decimal.Decimal {

	contract, err := web3.Eth.NewContract(Abi, tokenAddress)

	res, err := contract.Call("balanceOf", common.HexToAddress(userAddress))
	if err != nil {
		log.Context(ctx).Error("GetBalance balanceOf err ", "msg", err)
		return decimal.Zero
	}
	value, ok := res.(*big.Int)
	if ok {
		balance := decimal.NewFromBigInt(value, 0).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(decimals)))
		return balance
	}
	return decimal.Zero
}

func GetDecimals(ctx context.Context, tokenAddress string, web3 *web3.Web3) int64 {

	contract, err := web3.Eth.NewContract(Abi, tokenAddress)

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
