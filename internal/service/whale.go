package service

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaasset"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aatransactioninfo"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/token"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/whalestatisticday"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/whalestatistichour"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/shopspring/decimal"
	"log"
	"sort"
	"time"
)

func GetWhaleOverview(ctx context.Context, req vo.WhaleOverviewRequest) (*vo.WhaleOverviewResponse, error) {
	network := req.Network
	var resp = &vo.WhaleOverviewResponse{}

	client, err := entity.Client(ctx, network)
	if err != nil {
		return nil, err
	}
	aaAssets, err := client.AaAsset.Query().Order(ent.Desc(aaasset.FieldAssetValue)).Limit(config.WhaleNum).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(aaAssets) == 0 {
		return nil, nil
	}
	minAssetValue := aaAssets[len(aaAssets)-1].AssetValue
	noWhaleCount, err := client.AaAsset.Query().Where(aaasset.AssetValueLTE(minAssetValue)).Count(ctx)
	if err != nil {
		return nil, nil
	}
	totalAssetValue := decimal.Zero
	totalEth := decimal.Zero
	var whaleAddresss []string
	for _, asset := range aaAssets {
		totalAssetValue = totalAssetValue.Add(asset.AssetValue)
		totalEth = totalEth.Add(asset.Balance)
		whaleAddresss = append(whaleAddresss, asset.ID)
	}

	txStartTimeMs := time.Now().UnixMilli() - config.WhaleTxDay*DaySecond*1000
	txStartTime := time.UnixMilli(txStartTimeMs)
	allCount, err := client.AaTransactionInfo.Query().Where(aatransactioninfo.TimeGTE(txStartTime)).Count(ctx)
	if err != nil {
		allCount = 0
	}
	whaleCount, err := client.AaTransactionInfo.Query().Where(aatransactioninfo.TimeGTE(txStartTime), aatransactioninfo.IDIn(whaleAddresss[:]...)).Count(ctx)
	if err != nil {
		whaleCount = 0
	}

	txDominance := decimal.Zero
	if allCount != 0 {
		txDominance = decimal.NewFromInt(int64(whaleCount)).DivRound(decimal.NewFromInt(int64(allCount)), 6)
	}

	baseTokens, err := client.Token.Query().Where(token.TypeEQ("base")).All(ctx)
	if err != nil {
		return nil, nil
	}
	ethPrice := decimal.Zero
	if len(baseTokens) > 0 {
		ethPrice = baseTokens[0].TokenPrice
	}
	ratio := decimal.NewFromInt(int64(len(aaAssets))).Div(decimal.NewFromInt(int64(noWhaleCount)).Add(decimal.NewFromInt(int64(len(aaAssets))))).RoundDown(6)
	totalEthValue := totalEth.Mul(ethPrice).RoundDown(6)

	resp.TotalAssetUsd = totalAssetValue
	resp.TotalEthUsd = totalEthValue
	resp.Ratio = ratio
	resp.TxDominance = txDominance
	return resp, nil
}

func GetWhaleChart(ctx context.Context, req vo.WhaleChartRequest) (*vo.WhaleChartResponse, error) {
	network := req.Network
	var resp = &vo.WhaleChartResponse{}
	client, err := entity.Client(ctx, network)
	if err != nil {
		return nil, err
	}

	timeRange := req.TimeRange

	if timeRange == config.RangeH24 {
		startTime := time.Now().Add(-24 * time.Hour)
		whaleStatisticHours, err := client.WhaleStatisticHour.Query().Where(whalestatistichour.StatisticTimeGTE(startTime.UnixMilli()), whalestatistichour.NetworkEqualFold(network)).All(ctx)
		if err != nil {
			return nil, err
		}
		resp = getWhaleResponseHour(whaleStatisticHours)
		if resp == nil {
			return nil, nil
		}
	} else if timeRange == config.RangeD7 {
		startTime := time.Now().Add(-8 * 24 * time.Hour)
		whaleStatisticDays, err := client.WhaleStatisticDay.Query().Where(whalestatisticday.StatisticTimeGTE(startTime.UnixMilli()), whalestatisticday.NetworkEqualFold(network)).All(ctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		resp = getWhaleResponseDay(whaleStatisticDays)
		if resp == nil {
			return nil, nil
		}
	} else if timeRange == config.RangeD30 {
		startTime := time.Now().Add(-31 * 24 * time.Hour)
		whaleStatisticDays, err := client.WhaleStatisticDay.Query().Where(whalestatisticday.StatisticTimeGTE(startTime.UnixMilli()), whalestatisticday.NetworkEqualFold(network)).All(ctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		resp = getWhaleResponseDay(whaleStatisticDays)
		if resp == nil {
			return nil, nil
		}
	}
	return resp, nil
}

func getWhaleResponseDay(days []*ent.WhaleStatisticDay) *vo.WhaleChartResponse {

	if len(days) == 0 {
		return nil
	}
	var resp vo.WhaleChartResponse
	var details []*vo.WhaleChartDetail

	for _, statisticDay := range days {

		detail := &vo.WhaleChartDetail{
			Time:  statisticDay.StatisticTime,
			Value: statisticDay.TotalUsd,
		}
		details = append(details, detail)
	}
	sort.Sort(vo.ByWhaleDailyStatisticTime(details))
	resp.Details = details
	return &resp
}

func getWhaleResponseHour(hours []*ent.WhaleStatisticHour) *vo.WhaleChartResponse {
	if len(hours) == 0 {
		return nil
	}
	var resp vo.WhaleChartResponse
	var details []*vo.WhaleChartDetail

	for _, statisticHour := range hours {

		detail := &vo.WhaleChartDetail{
			Time:  statisticHour.StatisticTime,
			Value: statisticHour.TotalUsd,
		}
		details = append(details, detail)
	}
	sort.Sort(vo.ByWhaleDailyStatisticTime(details))
	resp.Details = details
	return &resp
}
