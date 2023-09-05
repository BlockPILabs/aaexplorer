package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/tokenpriceinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/userassetinfo"
	"github.com/BlockPILabs/aa-scan/third/moralis"
	"github.com/shopspring/decimal"
	"time"
)

type WalletBalance struct {
	Symbol          string
	ContractAddress string
	Amount          decimal.Decimal
	ValueUsd        decimal.Decimal
	Percent         decimal.Decimal
}

func GetWalletBalanceDetail(accountAddress string, network string) []*WalletBalance {
	client, err := entity.Client(context.Background())
	if err != nil {
		return nil
	}
	userAssetInfos, err := client.UserAssetInfo.Query().Where(userassetinfo.AccountAddressEqualFold(accountAddress)).All(context.Background())
	if err != nil {
		return nil
	}
	if len(userAssetInfos) == 0 {
		userAssetInfos = moralis.GetUserAsset(accountAddress, network)
	}
	if len(userAssetInfos) == 0 {
		return nil
	}

	var totalBalance = decimal.Zero
	var details []*WalletBalance

	for _, asset := range userAssetInfos {
		tokenPrice, err := client.TokenPriceInfo.Query().Where(tokenpriceinfo.ContractAddressEqualFold(asset.ContractAddress)).Limit(1).All(context.Background())
		if err != nil {
			continue
		}
		var detail *WalletBalance
		if tokenPrice == nil || len(tokenPrice) == 0 {
			onePrice := moralis.GetTokenPrice(asset.ContractAddress, network)
			detail = &WalletBalance{
				Symbol:          asset.Symbol,
				ContractAddress: asset.ContractAddress,
				Amount:          asset.Amount,
				ValueUsd:        onePrice.UsdPrice.Mul(asset.Amount),
			}

		} else {
			onePrice := tokenPrice[0]
			detail = &WalletBalance{
				Symbol:          asset.Symbol,
				ContractAddress: asset.ContractAddress,
				Amount:          asset.Amount,
				ValueUsd:        onePrice.TokenPrice.Mul(asset.Amount),
			}
		}

		if detail.ValueUsd.Equal(decimal.Zero) {
			continue
		}
		details = append(details, detail)
		totalBalance = totalBalance.Add(detail.ValueUsd)
	}

	for _, one := range details {
		one.Percent = one.ValueUsd.DivRound(totalBalance, 4)
	}

	return details
}

func GetNativePrice(network string) *decimal.Decimal {
	client, err := entity.Client(context.Background(), network)
	if err != nil {
		return nil
	}
	var contract string
	if network == config.EthNetwork {
		contract = config.WETH
	} else if network == config.BscNetwork {
		contract = config.WBNB
	} else if network == config.PolygonNetwork {
		contract = config.WMATIC
	}
	prices, err := client.TokenPriceInfo.Query().Where(tokenpriceinfo.ContractAddressEqualFold(contract), tokenpriceinfo.NetworkEqualFold(network)).All(context.Background())
	if err != nil {
		return nil
	}
	if len(prices) == 0 {
		token := moralis.GetTokenPrice(contract, network)
		if token == nil {
			return nil
		}
		client.TokenPriceInfo.Create().
			SetTokenPrice(token.UsdPrice).
			SetNetwork(network).
			SetContractAddress(contract).
			SetLastTime(time.Now().UnixMilli()).
			SetSymbol(token.TokenSymbol).Save(context.Background())

		return &token.UsdPrice
	}

	return &prices[0].TokenPrice
}

type WhaleOverview struct {
	TxDominance     decimal.Decimal
	Represent       decimal.Decimal
	TotalBalance    decimal.Decimal
	TotalBalanceETH decimal.Decimal
}

func WalletTracker(accountAddress string, network string) {

	//get watching list
	//

}
