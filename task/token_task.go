package task

import (
	"context"
	constConfig "github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/token"
	"github.com/BlockPILabs/aaexplorer/third/cmc"
	"github.com/BlockPILabs/aaexplorer/third/schedule"
	"github.com/shopspring/decimal"
	"time"
)

func init() {
	schedule.Add("token_task", func(ctx context.Context) {
		TokenTask(ctx)
	}).ScheduleWithCron("0 0 0 * * *")
}

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
