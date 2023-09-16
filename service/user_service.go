package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aaaccountdata"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/tokenpriceinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/userassetinfo"
	"github.com/BlockPILabs/aa-scan/third/moralis"
	"github.com/shopspring/decimal"
	"log"
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
	client, err := entity.Client(context.Background(), network)
	if err != nil {
		return nil
	}
	userAssetInfos, err := client.UserAssetInfo.Query().Where(userassetinfo.AccountAddressEqualFold(accountAddress), userassetinfo.NetworkEqualFold(network)).All(context.Background())
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
	var hasNative = false

	for _, asset := range userAssetInfos {
		if asset.ContractAddress == config.ZeroAddress {
			hasNative = true
		}
		if asset.Amount.Cmp(decimal.Zero) == 0 {
			continue
		}
		tokenPrice, err := client.TokenPriceInfo.Query().Where(tokenpriceinfo.ContractAddressEqualFold(asset.ContractAddress)).Limit(1).All(context.Background())
		if err != nil {
			continue
		}
		var detail = &WalletBalance{
			Symbol:          asset.Symbol,
			ContractAddress: asset.ContractAddress,
			Amount:          asset.Amount,
		}
		if tokenPrice == nil || len(tokenPrice) == 0 {
			onePrice := moralis.GetTokenPrice(asset.ContractAddress, network)
			var usdPrice = decimal.Zero
			if onePrice != nil {
				usdPrice = onePrice.UsdPrice
			}
			saveTokenPrice(asset.ContractAddress, asset.Symbol, usdPrice, client, asset.Network)
			detail.ValueUsd = usdPrice.Mul(asset.Amount)

		} else {
			onePrice := tokenPrice[0]
			detail.ValueUsd = onePrice.TokenPrice.Mul(asset.Amount)
		}
		if detail.ValueUsd.Equal(decimal.Zero) {
			continue
		}
		details = append(details, detail)
		totalBalance = totalBalance.Add(detail.ValueUsd)
	}
	if !hasNative {
		nativeBalance := moralis.GetNativeTokenBalance(accountAddress, network)
		nativePrice := GetNativePrice(network)
		nativeUsd := nativePrice.Mul(nativeBalance)
		totalBalance = totalBalance.Add(nativeUsd)
		var detail = &WalletBalance{
			Symbol:          moralis.GetNativeName(network),
			ContractAddress: config.ZeroAddress,
			Amount:          nativeBalance,
			ValueUsd:        nativeUsd,
			Percent:         nativeUsd.DivRound(totalBalance, 4),
		}
		details = append(details, detail)
		client.UserAssetInfo.Create().SetSymbol(detail.Symbol).SetNetwork(network).SetContractAddress(detail.ContractAddress).SetAccountAddress(accountAddress).SetLastTime(time.Now().UnixMilli()).Save(context.Background())
	}

	for _, one := range details {
		one.Percent = one.ValueUsd.DivRound(totalBalance, 4)
	}

	return details
}

func saveTokenPrice(contractAddress string, symbol string, usdPrice decimal.Decimal, client *ent.Client, network string) {
	_, err := client.TokenPriceInfo.Create().
		SetTokenPrice(usdPrice).
		SetNetwork(network).
		SetContractAddress(contractAddress).
		SetSymbol(symbol).
		SetLastTime(time.Now().UnixMilli()).Save(context.Background())
	if err != nil {
		log.Println(err)
	}
}

func GetNativePrice(network string) decimal.Decimal {
	client, err := entity.Client(context.Background(), network)
	if err != nil {
		return decimal.Zero
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
		return decimal.Zero
	}
	if len(prices) == 0 {
		token := moralis.GetTokenPrice(contract, network)
		if token == nil {
			return decimal.Zero
		}
		client.TokenPriceInfo.Create().
			SetTokenPrice(token.UsdPrice).
			SetNetwork(network).
			SetContractAddress(contract).
			SetLastTime(time.Now().UnixMilli()).
			SetSymbol(token.TokenSymbol).Save(context.Background())

		return token.UsdPrice
	}

	return prices[0].TokenPrice
}

func GetTokenPrice(tokenAddress string, network string) decimal.Decimal {
	client, err := entity.Client(context.Background(), network)
	if err != nil {
		return decimal.Zero
	}
	prices, err := client.TokenPriceInfo.Query().Where(tokenpriceinfo.ContractAddressEqualFold(tokenAddress), tokenpriceinfo.NetworkEqualFold(network)).All(context.Background())
	if err != nil {
		return decimal.Zero
	}
	if len(prices) == 0 {
		token := moralis.GetTokenPrice(tokenAddress, network)
		if token == nil {
			saveTokenPrice(tokenAddress, "", decimal.Zero, client, network)
			return decimal.Zero
		}
		saveTokenPrice(tokenAddress, token.TokenSymbol, token.UsdPrice, client, network)

		return token.UsdPrice
	}

	return prices[0].TokenPrice
}

func GetTotalBalance(address string, network string) decimal.Decimal {
	client, err := entity.Client(context.Background(), network)
	if err != nil {
		return decimal.Zero
	}
	existAccount, err := client.AaAccountData.Query().Where(aaaccountdata.IDEQ(address)).Limit(1).All(context.Background())
	if len(existAccount) > 0 {
		return existAccount[0].TotalBalanceUsd
	}
	details := GetWalletBalanceDetail(address, network)

	if len(details) == 0 {
		updateBalance(address, client, decimal.Zero)
		return decimal.Zero
	}
	var totalBalance = decimal.Zero
	for _, detail := range details {
		totalBalance = totalBalance.Add(detail.ValueUsd)
	}
	updateBalance(address, client, totalBalance)
	return totalBalance
}

func updateBalance(address string, client *ent.Client, balance decimal.Decimal) {
	accountData, _ := client.AaAccountData.Query().Where(aaaccountdata.IDEQ(address)).All(context.Background())
	if accountData != nil && len(accountData) > 0 {
		client.AaAccountData.Update().SetTotalBalanceUsd(balance).SetLastTime(time.Now().UnixMilli()).Where(aaaccountdata.IDEQ(accountData[0].ID)).Exec(context.Background())
	}
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
