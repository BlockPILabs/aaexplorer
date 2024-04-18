package task

import (
	"context"
	constConfig "github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaasset"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaassetdetail"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/token"
	"github.com/BlockPILabs/aaexplorer/third/cmc"
	"github.com/BlockPILabs/aaexplorer/third/schedule"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"log"
	"strings"
	"time"
)

func init() {
	schedule.Add("token_task", func(ctx context.Context) {
		TokenTask(ctx)
	}).ScheduleWithCron("0 0 0 * * *")
}

const PriceExpire = 20 * 3600 * 1000

func TokenTask(ctx context.Context) {
	tokenList := cmc.GetTopToken(constConfig.TopNum)
	if len(tokenList) == 0 {
		return
	}
	cli, err := entity.Client(ctx)
	if err != nil {
		return
	}

	networks, err := cli.Network.Query().All(ctx)

	for _, net := range networks {
		network := net.ID
		client, err := entity.Client(ctx, network)
		if err != nil {
			continue
		}

		for _, tokenItem := range tokenList {
			symbol := tokenItem.Symbol
			tokens, err := client.Token.Query().Where(token.SymbolEqualFold(symbol), token.NetworkEqualFold(network)).All(ctx)
			if err != nil {
				logger.Error("TokenTask get token err ", "symbol", "msg", symbol, err)
				continue
			}
			if len(tokens) > 0 {
				for _, oldToken := range tokens {
					if oldToken.MarketRank != tokenItem.Rank {
						_, err = client.Token.Update().Where(token.IDEQ(oldToken.ID)).SetMarketRank(tokenItem.Rank).Save(ctx)
						if err != nil {
							logger.Error("TokenTask update token ranking err ", "symbol", "msg", symbol, err)
						}
					}
				}
				continue
			}
			var contractAddress string
			platfrom := tokenItem.Platform
			if platfrom != nil && platfrom.Slug == network {
				contractAddress = platfrom.TokenAddress
			}

			now := time.Now()
			_, err = client.Token.Create().SetTokenPrice(decimal.Zero).
				SetNetwork(network).SetUpdateTime(now).SetCreateTime(now).SetLastTime(now.UnixMilli()).
				SetSymbol(symbol).SetContractAddress(contractAddress).SetMarketRank(tokenItem.Rank).SetFullName(tokenItem.Name).
				SetType("").Save(ctx)
			if err != nil {
				logger.Error("TokenTask save token err ", "symbol", "msg", symbol, err)
			} else {
				logger.Info("TokenTask save token success", "symbol", symbol)
			}
		}
	}
}

func InitRefreshToken(ctx context.Context) {
	hourScheduler := chrono.NewDefaultTaskScheduler()
	_, err := hourScheduler.ScheduleWithCron(func(ctx context.Context) {
		RefreshToken(ctx)
	}, "0 1 1 * * *")

	if err == nil {
		log.Print("RefreshToken has been scheduled")
	}
}

func InitRefreshPrice(ctx context.Context) {
	hourScheduler := chrono.NewDefaultTaskScheduler()
	_, err := hourScheduler.ScheduleWithCron(func(ctx context.Context) {
		RefreshPrice(ctx)
		RefreshOldAsset(ctx)
	}, "0 1 2 * * *")

	if err == nil {
		log.Print("RefreshPrice has been scheduled")
	}
}

func RefreshOldAsset(ctx context.Context) {
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
		aas, err := client.AaAsset.Query().All(ctx)
		if err != nil {
			logger.Error("AssetRefreshTask query asset err ", "msg", err)
			continue
		}
		if len(aas) == 0 {
			continue
		}
		tokens, err := client.Token.Query().All(ctx)
		if len(tokens) == 0 {
			continue
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

		for _, aa := range aas {
			address := aa.ID
			balance := aa.Balance
			details, err := client.AaAssetDetail.Query().Where(aaassetdetail.UserAddressEqualFold(address)).All(ctx)
			if err != nil {
				continue
			}
			if len(details) == 0 && balance.Cmp(decimal.Zero) == 0 {
				continue
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
				logger.Info("RefreshOldAsset update asset err, ", "userAddress", address, "network", network, "msg", err)
			} else {
				logger.Info("RefreshOldAsset update asset success, ", "userAddress", address, "network", network)
			}
		}
	}
}

func RefreshToken(ctx context.Context) {
	tokenInfos := cmc.GetTopToken(constConfig.TopNum)
	if len(tokenInfos) == 0 {
		return
	}
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
		client, err := entity.Client(ctx, network)
		if err != nil {
			continue
		}
		for _, tokenInfo := range tokenInfos {
			symbol := tokenInfo.Symbol
			tokens, err := client.Token.Query().Where(token.SymbolEqualFold(symbol), token.NetworkEqualFold(network)).All(ctx)
			if err != nil {
				logger.Error("RefreshToken get token err ", "symbol", "msg", symbol, err)
				continue
			}
			if len(tokens) > 0 {
				continue
			}

			contractAddress := ""
			platform := tokenInfo.Platform
			if platform != nil && platform.Slug == network {
				contractAddress = platform.TokenAddress
			}
			now := time.Now()
			_, err = client.Token.Create().SetTokenPrice(decimal.Zero).
				SetNetwork(network).SetUpdateTime(now).SetCreateTime(now).SetLastTime(now.UnixMilli()).SetSymbol(symbol).
				SetContractAddress(contractAddress).SetMarketRank(tokenInfo.Rank).SetFullName(tokenInfo.Name).SetType("").Save(ctx)

			if err != nil {
				logger.Error("RefreshToken save token err ", "symbol", "msg", symbol, err)
			} else {
				logger.Info("RefreshToken save token success", "symbol", symbol)
			}

		}
	}
}

func RefreshPrice(ctx context.Context) {
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
		client, err := entity.Client(ctx, network)
		if err != nil {
			continue
		}
		tokens, err := client.Token.Query().Where(token.NetworkEqualFold(network)).All(ctx)
		if len(tokens) == 0 {
			continue
		}
		now := time.Now().UnixMilli()
		for _, one := range tokens {
			if now-one.LastTime < PriceExpire {
				//continue
			}
			price := cmc.GetTokenPrice(one.Symbol)
			logger.Info("RefreshPrice get price success, ", "symbol", one.Symbol, "price", price)
			if price.Cmp(decimal.Zero) == 0 {
				continue
			}
			err := client.Token.Update().SetTokenPrice(price).SetLastTime(now).Where(token.IDEQ(one.ID)).Exec(ctx)
			if err != nil {
				logger.Error("RefreshTokenPrice err ", "symbol", one.Symbol, "msg", err)
			} else {
				logger.Info("RefreshTokenPrice success ", "symbol", one.Symbol)

			}
			time.Sleep(3 * time.Second)
		}
	}

}
