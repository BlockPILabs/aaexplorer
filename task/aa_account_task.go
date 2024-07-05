package task

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaaccountdata"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaasset"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"time"
)

func InitAaAccountTask() {
	hourScheduler := chrono.NewDefaultTaskScheduler()
	_, err := hourScheduler.ScheduleWithCron(func(ctx context.Context) {
		AaAccountTask(ctx)
	}, "0 15 * * * *")

	if err == nil {
		logger.Info("AaAccountTask has been scheduled")
	}
}

func AaAccountTask(ctx context.Context) {
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
		accounts, err := client.AaAccountData.Query().Where(aaaccountdata.AaTypeIn("aa", "bundler", "paymaster")).All(ctx)
		if len(accounts) == 0 {
			continue
		}

		for _, account := range accounts {
			olds, err := client.AaAsset.Query().Where(aaasset.IDEqualFold(account.ID)).All(ctx)
			if err != nil {
				continue
			}
			if len(olds) > 0 {
				continue
			}
			aaAsset := client.AaAsset.Create().SetAssetValue(decimal.Zero).SetLastTime(0).SetCreateTime(time.Now()).SetUpdateTime(time.Now()).SetNetwork(network).SetBalance(decimal.Zero).SetID(account.ID)
			_, err = aaAsset.Save(ctx)
			if err != nil {
				logger.Error("AaAccountTask save aaAsset err ", "id", account.ID, "msg", err)
			}
		}

	}
}
