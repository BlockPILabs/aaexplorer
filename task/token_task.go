package task

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/token"
	"github.com/BlockPILabs/aaexplorer/third/cmc"
	"github.com/shopspring/decimal"
	"time"
)
import constConfig "github.com/BlockPILabs/aaexplorer/config"

func TokenTask(ctx context.Context) {
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
				logger.Error("TokenTask get token err ", "symbol", "msg", symbol, err)
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
				SetNetwork(constConfig.EthNetwork).SetUpdateTime(now).SetCreateTime(now).SetLastTime(now.UnixMilli()).SetSymbol(symbol).
				SetContractAddress(contractAddress).SetMarketRank(tokenInfo.Rank).SetFullName(tokenInfo.Name).SetType("").Save(ctx)

			if err != nil {
				logger.Error("TokenTask save token err ", "symbol", "msg", symbol, err)
			} else {
				logger.Info("TokenTask save token success", "symbol", symbol)
			}

		}
	}

}
