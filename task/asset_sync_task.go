package task

import (
	"context"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/assetchangetrace"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/tokenpriceinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/userassetinfo"
	"github.com/BlockPILabs/aa-scan/third/moralis"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"log"
	"math"
	"time"
)

func AssetTask() {
	dayScheduler := chrono.NewDefaultTaskScheduler()
	_, err := dayScheduler.ScheduleWithCron(func(ctx context.Context) {
		AssetSync()
	}, "0 15 0 * * ?")

	if err == nil {
		log.Print("AssetSyncTask has been scheduled")
	}
}

func AssetSync() {
	cli, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	networks, err := cli.Network.Query().All(context.Background())
	if len(networks) == 0 {
		return
	}
	hour6Ago := time.Now().UnixMilli() - 6*3600*1000
	for _, record := range networks {
		network := record.Name
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			return
		}
		for {
			changeTraces, err := client.AssetChangeTrace.Query().Where(assetchangetrace.SyncFlagEQ(0), assetchangetrace.LastChangeTimeLT(hour6Ago)).Limit(100).All(context.Background())
			if err != nil {
				return
			}
			if len(changeTraces) == 0 {
				return
			}
			var accounts []*ent.AssetChangeTrace
			var tokens []*ent.AssetChangeTrace
			var accountAddrs []string
			var tokenAddrs []string
			var changes []int64

			for _, changeTrace := range changeTraces {
				if changeTrace.AddressType == config.AddressTypeAccount {
					accounts = append(accounts, changeTrace)
					accountAddrs = append(accountAddrs, changeTrace.Address)
				} else {
					tokens = append(tokens, changeTrace)
					tokenAddrs = append(tokenAddrs, changeTrace.Address)
				}
				changes = append(changes, changeTrace.ID)
			}
			syncAccountBalance(client, accounts)
			syncTokenPrice(client, tokens)
			syncWTokenPrice(client)
			client.AssetChangeTrace.Update().Where(assetchangetrace.IDIn(changes[:]...)).SetSyncFlag(1).SetLastChangeTime(time.Now().UnixMilli()).Exec(context.Background())
		}
	}

	//syncAssetValue(client, accountAddrs, tokenAddrs)
}

func syncWTokenPrice(client *ent.Client) {
	var wtokens = make(map[string]string)
	wtokens[config.WBNB] = config.BSC
	wtokens[config.WETH] = config.Eth
	wtokens[config.WMATIC] = config.Polygon
	for token, value := range wtokens {
		tokenPrice := moralis.GetTokenPrice(token, value)
		curMillis := time.Now().UnixMilli()
		exist, err := client.TokenPriceInfo.Query().Where(tokenpriceinfo.ContractAddressEqualFold(token), tokenpriceinfo.NetworkEqualFold(value)).Exist(context.Background())
		if err != nil {
			continue
		}
		if exist {
			client.TokenPriceInfo.Update().
				Where(tokenpriceinfo.ContractAddressEQ(token)).
				SetTokenPrice(tokenPrice.UsdPrice).
				SetLastTime(curMillis).
				Exec(context.Background())
		} else {
			client.TokenPriceInfo.Create().
				SetTokenPrice(tokenPrice.UsdPrice).
				SetSymbol(tokenPrice.TokenSymbol).
				SetContractAddress(token).
				SetNetwork(value).
				SetLastTime(curMillis).Save(context.Background())
		}

	}

}

func syncAssetValue(client *ent.Client, accounts []string, tokens []string) {
	assetInfos, err := client.UserAssetInfo.Query().Where(userassetinfo.AccountAddressIn(accounts[:]...)).All(context.Background())
	if err != nil {
		return
	}
	if len(assetInfos) == 0 {
		return
	}
	var contractMap = make(map[string]decimal.Decimal)
	for _, assetInfo := range assetInfos {
		contractAddress := assetInfo.ContractAddress
		_, contractOk := contractMap[contractAddress]
		if !contractOk {
			client.TokenPriceInfo.Query().Where(tokenpriceinfo.ContractAddressEqualFold(contractAddress)).All(context.Background())
		}
	}
}

func syncTokenPrice(client *ent.Client, tokens []*ent.AssetChangeTrace) map[string]decimal.Decimal {
	if len(tokens) == 0 {
		return nil
	}
	var priceMap = make(map[string]decimal.Decimal)
	for _, token := range tokens {
		tokenPrice := moralis.GetTokenPrice(token.Address, token.Network)
		curMillis := time.Now().UnixMilli()
		exist, err := client.TokenPriceInfo.Query().Where(tokenpriceinfo.ContractAddressEqualFold(token.Address), tokenpriceinfo.NetworkEqualFold(token.Network)).Exist(context.Background())
		if err != nil {
			continue
		}
		if exist {
			client.TokenPriceInfo.Update().
				Where(tokenpriceinfo.ContractAddressEQ(token.Address)).
				SetTokenPrice(tokenPrice.UsdPrice).
				SetLastTime(curMillis).
				Exec(context.Background())
		} else {
			client.TokenPriceInfo.Create().
				SetTokenPrice(tokenPrice.UsdPrice).
				SetSymbol(tokenPrice.TokenSymbol).
				SetContractAddress(token.Address).
				SetNetwork(token.Network).
				SetLastTime(curMillis).Save(context.Background())
		}

		priceMap[token.Address] = tokenPrice.UsdPrice
	}

	return priceMap

}

func syncAccountBalance(client *ent.Client, accounts []*ent.AssetChangeTrace) map[string][]*moralis.TokenBalance {
	if len(accounts) == 0 {
		return nil
	}
	var accountMap = make(map[string][]*moralis.TokenBalance)
	for _, account := range accounts {
		tokenBalances := moralis.GetTokenBalance(account.Address, account.Network)
		nativeTokenBalance := moralis.GetNativeTokenBalance(account.Address, account.Network)
		tokenBalances = addNativeToken(tokenBalances, nativeTokenBalance, account.Address, account.Network, client)
		if len(tokenBalances) == 0 {
			continue
		}
		curMillis := time.Now().UnixMilli()
		var userAssetInfoCreates []*ent.UserAssetInfoCreate
		for _, tokenBalance := range tokenBalances {
			userAssetCreate := client.UserAssetInfo.Create().
				SetAccountAddress(account.Address).
				SetContractAddress(tokenBalance.TokenAddress).
				SetSymbol(tokenBalance.Symbol).
				SetNetwork(account.Network).
				SetAmount(tokenBalance.Balance.DivRound(decimal.New(int64(math.Pow10(int(tokenBalance.Decimals))), 0), tokenBalance.Decimals)).
				SetLastTime(curMillis)
			userAssetInfoCreates = append(userAssetInfoCreates, userAssetCreate)
		}

		client.UserAssetInfo.Delete().Where(userassetinfo.AccountAddressEqualFold(account.Address), userassetinfo.NetworkEQ(account.Network)).Exec(context.Background())

		err := client.UserAssetInfo.CreateBulk(userAssetInfoCreates...).Exec(context.Background())
		if err != nil {
			log.Println(err)
		}
		accountMap[account.Address] = tokenBalances
	}

	return accountMap
}

func addNativeToken(balances []*moralis.TokenBalance, native decimal.Decimal, address string, network string, client *ent.Client) []*moralis.TokenBalance {
	if native.Equal(decimal.Zero) {
		return balances
	}
	if balances == nil || len(balances) == 0 {
		balances = []*moralis.TokenBalance{}
	}
	userAssetCreate := &moralis.TokenBalance{
		Balance:      native,
		TokenAddress: config.ZeroAddress,
		Name:         moralis.GetNativeName(network),
		Decimals:     config.EvmDecimal,
	}
	balances = append(balances, userAssetCreate)
	return balances
}

func GetWToken(network string) string {
	if network == config.EthNetwork {
		return config.WETH
	} else if network == config.BscNetwork {
		return config.WBNB
	} else if network == config.PolygonNetwork {
		return config.WMATIC
	}

	return ""
}
