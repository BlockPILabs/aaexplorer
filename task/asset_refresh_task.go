package task

import (
	"context"
	constConfig "github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaasset"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaassetdetail"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aauseropsinfo"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/token"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/chenzhijie/go-web3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"math/big"
	"strings"
	"time"
)

const Abi = "[{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

func InitAssetRefreshTask() {
	hourScheduler := chrono.NewDefaultTaskScheduler()
	_, err := hourScheduler.ScheduleWithCron(func(ctx context.Context) {
		AssetRefreshTask(ctx)
	}, "0 35 0 1/6 * *")

	if err == nil {
		logger.Info("AssetRefreshTask has been scheduled")
	}
}

func AddDecimals(ctx context.Context) {
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
		w3, err := web3.NewWeb3(net.HTTPRPC)
		if err != nil {
			logger.Error("AddDecimals newWeb3 err ", "msg", err)
			continue
		}
		w3.Eth.SetChainId(net.ChainID)
		tokens, err := client.Token.Query().All(ctx)
		if len(tokens) == 0 {
			continue
		}

		for _, oneToken := range tokens {
			contractAddress := oneToken.ContractAddress
			if len(contractAddress) == 0 {
				continue
			}
			decimals := GetDecimals(ctx, contractAddress, w3)
			err = client.Token.Update().SetDecimals(decimals).Where(token.IDEQ(oneToken.ID)).Exec(ctx)
			if err != nil {
				logger.Error("update decimals err ", "symbol", oneToken.Symbol, "msg", err)
			} else {
				logger.Info("update decimals success ", "symbol", oneToken.Symbol, "decimals", decimals)

			}
			time.Sleep(5 * time.Second)
		}
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
		oneSize := len(aas) / 3
		if oneSize == 0 {
			oneSize = len(aas)
		}
		var allArrs [][]*ent.AaAsset
		for i := 0; i <= 3; i++ {
			var oneArr []*ent.AaAsset
			allArrs = append(allArrs, oneArr)
		}
		logger.Info("AssetRefreshTask-arrSize ", "size", len(allArrs), "oneSize", oneSize, "network", network)
		start := 0
		for idx, aa := range aas {
			r := idx % oneSize
			if r == 0 && idx != 0 {
				logger.Info("AssetRefreshTask-oneSize ", "size", len(allArrs[start]), "start", start, "idx", idx, "r", r, "network", network)
				go doRefresh(ctx, client, tokens, w3, blockNum, network, allArrs[start])
				start = start + 1
			}
			if len(allArrs)-1 < start {
				break
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
		startTime := time.Now().Second() - constConfig.WhaleTxDay*24*3600
		lastTxCount, err := client.AAUserOpsInfo.Query().Where(aauseropsinfo.TxTimeGTE(int64(startTime)), aauseropsinfo.SenderEqualFold(userAddress)).Count(ctx)
		if err != nil {
			continue
		}
		if lastTxCount == 0 {
			oldRefresh(ctx, client, userAddress, tokens, network, aa)
			continue
		}

		for _, token := range tokens {
			contractAddress := token.ContractAddress
			if len(contractAddress) == 0 {
				continue
			}
			decimals := token.Decimals
			balance := GetBalance(ctx, contractAddress, userAddress, decimals, w3)
			assetValue := balance.Mul(token.TokenPrice)
			addOrUpdateAssetDetail(ctx, contractAddress, userAddress, assetValue, client, network, token.Symbol, balance)
			totalValue = totalValue.Add(assetValue)
			//logger.Info("AssetRefreshTask-one-token ", "user", userAddress, "token", contractAddress, "network", network)
		}
		logger.Info("AssetRefreshTask-one-address ", "user", userAddress, "network", network)
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

func oldRefresh(ctx context.Context, client *ent.Client, address string, tokens []*ent.Token, network string, aa *ent.AaAsset) {
	details, err := client.AaAssetDetail.Query().Where(aaassetdetail.UserAddressEqualFold(address)).All(ctx)
	if err != nil {
		return
	}

	balance := aa.Balance
	if len(details) == 0 && balance.Cmp(decimal.Zero) == 0 {
		return
	}

	var tokenMap = make(map[string]decimal.Decimal)
	var nativePrice = decimal.Zero
	for _, one := range tokens {
		if one.Type != nil && *one.Type == "base" {
			nativePrice = one.TokenPrice
		}
		if len(one.ContractAddress) == 0 {
			continue
		}
		contractAddress := strings.ToLower(one.ContractAddress)
		tokenMap[contractAddress] = one.TokenPrice
	}

	totalValue := decimal.Zero
	if len(details) > 0 {
		for _, detail := range details {
			if detail.AssetAmount.Cmp(decimal.Zero) == 0 {
				continue
			}
			contractAddress := strings.ToLower(detail.ContractAddress)
			price := tokenMap[contractAddress]
			oneValue := price.Mul(detail.AssetAmount)
			totalValue = totalValue.Add(oneValue)
			client.AaAssetDetail.Update().SetAssetValue(oneValue).SetLastTime(time.Now().UnixMilli()).Where(aaassetdetail.IDEQ(detail.ID)).Exec(ctx)
		}
	}
	if balance.Cmp(decimal.Zero) > 0 {
		totalValue = totalValue.Add(nativePrice.Mul(balance))
	}

	err = client.AaAsset.Update().SetAssetValue(totalValue).SetLastTime(time.Now().UnixMilli()).Where(aaasset.IDEqualFold(address)).Exec(ctx)
	if err != nil {
		logger.Info("oldRefresh update asset err, ", "userAddress", address, "network", network, "msg", err)
	} else {
		logger.Info("oldRefresh update asset success, ", "userAddress", address, "network", network)
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
