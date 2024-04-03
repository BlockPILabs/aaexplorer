package task

import (
	"context"
	constConfig "github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/token"
	"github.com/BlockPILabs/aaexplorer/third/cmc"
	"github.com/BlockPILabs/aaexplorer/third/schedule"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"log"
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
	}, "0 1 2 * * *")

	if err == nil {
		log.Print("RefreshPrice has been scheduled")
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
				continue
			}
			price := cmc.GetTokenPrice(one.Symbol)
			err := client.Token.Update().SetTokenPrice(price).SetLastTime(now).Where(token.IDEQ(one.ID)).Exec(ctx)
			if err != nil {
				logger.Error("RefreshTokenPrice err ", "symbol", one.Symbol, "msg", err)
			} else {
				logger.Error("RefreshTokenPrice success ", "symbol", one.Symbol, "msg", err)

			}
			time.Sleep(time.Second)
		}
	}

}
