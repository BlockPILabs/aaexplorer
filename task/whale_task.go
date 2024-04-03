package task

import (
	"context"
	constConfig "github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaasset"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"time"
)

func InitWhaleStatistic(ctx context.Context) {
	hourScheduler := chrono.NewDefaultTaskScheduler()
	_, err := hourScheduler.ScheduleWithCron(func(ctx context.Context) {
		doWhaleHourStatistic(ctx)
	}, "0 20 * * * *")

	if err == nil {
		logger.Info("whaleHourStatistic has been scheduled")
	}

	dayScheduler := chrono.NewDefaultTaskScheduler()
	_, err = dayScheduler.ScheduleWithCron(func(ctx context.Context) {
		doWhaleDayStatistic(ctx)
	}, "0 35 0 * * *")

	if err == nil {
		logger.Info("whaleDayStatistic has been scheduled")
	}
}

func doWhaleHourStatistic(ctx context.Context) {
	cli, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	records, err := cli.Network.Query().All(context.Background())
	if len(records) == 0 {
		return
	}
	for _, record := range records {
		network := record.ID
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		assets, err := client.AaAsset.Query().Order(ent.Desc(aaasset.FieldAssetValue)).Limit(constConfig.WhaleNum).All(ctx)
		if err != nil {
			logger.Error("doWhaleHourStatistic err, ", "network", network, "msg", err)
			continue
		}
		if len(assets) == 0 {
			continue
		}
		totalValue := decimal.Zero
		for _, asset := range assets {
			totalValue = totalValue.Add(asset.AssetValue)
		}
		now := time.Now()
		hourStart := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
		whaleHour := client.WhaleStatisticHour.Create().SetWhaleNum(constConfig.WhaleNum).SetNetwork(network).SetStatisticTime(hourStart.UnixMilli()).SetTotalUsd(totalValue).SetCreateTime(time.Now())
		_, err = whaleHour.Save(ctx)
		if err != nil {
			logger.Error("doWhaleHourStatistic save err, ", "hour", hourStart, "msg", err)
		} else {
			logger.Info("doWhaleHourStatistic save success, ", "hour", hourStart)
		}

	}
}

func doWhaleDayStatistic(ctx context.Context) {
	cli, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	records, err := cli.Network.Query().All(context.Background())
	if len(records) == 0 {
		return
	}
	for _, record := range records {
		network := record.ID
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		assets, err := client.AaAsset.Query().Order(ent.Desc(aaasset.FieldAssetValue)).Limit(constConfig.WhaleNum).All(ctx)
		if err != nil {
			logger.Error("doWhaleDayStatistic err, ", "network", network, "msg", err)
			continue
		}
		if len(assets) == 0 {
			continue
		}
		totalValue := decimal.Zero
		for _, asset := range assets {
			totalValue = totalValue.Add(asset.AssetValue)
		}
		now := time.Now()
		dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		whaleDay := client.WhaleStatisticDay.Create().SetWhaleNum(constConfig.WhaleNum).SetNetwork(network).SetStatisticTime(dayStart.UnixMilli()).SetTotalUsd(totalValue).SetCreateTime(time.Now())
		_, err = whaleDay.Save(ctx)
		if err != nil {
			logger.Error("doWhaleDayStatistic save err, ", "day", dayStart, "msg", err)
		} else {
			logger.Info("doWhaleDayStatistic save success, ", "day", dayStart)
		}

	}
}
