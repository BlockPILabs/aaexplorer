package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aauseropsinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/dailystatisticday"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/dailystatistichour"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/shopspring/decimal"
	"log"
	"math"
	"sort"
	"time"
)

const DaySecond = 24 * 3600

func GetDailyStatistic(ctx context.Context, req vo.DailyStatisticRequest) (*vo.DailyStatisticResponse, error) {

	network := req.Network
	client, err := entity.Client(ctx, network)
	if err != nil {
		return nil, err
	}
	timeRange := req.TimeRange
	var resp *vo.DailyStatisticResponse

	if timeRange == config.RangeH24 {
		startTime := time.Now().Add(-24 * time.Hour)
		dailyStatisticHours, err := client.DailyStatisticHour.Query().Where(dailystatistichour.StatisticTimeGTE(startTime.UnixMilli()), dailystatistichour.NetworkEqualFold(network)).All(ctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		resp = getResponseHour(dailyStatisticHours)
		resp.Ups = decimal.NewFromInt(resp.UserOpsNum).DivRound(decimal.NewFromInt(DaySecond), 6)
	} else if timeRange == config.RangeD7 {
		startTime := time.Now().Add(-7 * 24 * time.Hour)
		dailyStatisticDays, err := client.DailyStatisticDay.Query().Where(dailystatisticday.StatisticTimeGTE(startTime.UnixMilli()), dailystatisticday.NetworkEqualFold(network)).All(ctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		resp = getResponseDay(dailyStatisticDays, 7)
		resp.Ups = decimal.NewFromInt(resp.UserOpsNum).DivRound(decimal.NewFromInt(7*DaySecond), 6)
	} else if timeRange == config.RangeD30 {
		startTime := time.Now().Add(-70 * 24 * time.Hour)
		dailyStatisticDays, err := client.DailyStatisticDay.Query().Where(dailystatisticday.StatisticTimeGTE(startTime.UnixMilli()), dailystatisticday.NetworkEqualFold(network)).All(ctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		resp = getResponseDay(dailyStatisticDays, 30)
		resp.Ups = decimal.NewFromInt(resp.UserOpsNum).DivRound(decimal.NewFromInt(30*DaySecond), 6)
	}

	return resp, nil
}

func getResponseDay(days []*ent.DailyStatisticDay, day int) *vo.DailyStatisticResponse {
	if len(days) == 0 {
		return nil
	}
	var resp vo.DailyStatisticResponse
	var details []*vo.DailyStatisticDetail
	var statisticTimeMap = make(map[int64]bool)
	for _, statisticDay := range days {
		if statisticTimeMap[statisticDay.StatisticTime] {
			continue
		}
		resp.AccumulativeGasFeeUsd = resp.AccumulativeGasFeeUsd.Add(statisticDay.GasFeeUsd).Round(2)
		resp.AccumulativeGasFee = resp.AccumulativeGasFee.Add(statisticDay.GasFee).Round(2)
		resp.BundlerGasProfit = resp.BundlerGasProfit.Add(statisticDay.BundlerGasProfit).Round(2)
		resp.BundlerGasProfitUsd = resp.BundlerGasProfitUsd.Add(statisticDay.BundlerGasProfitUsd).Round(2)
		resp.PaymasterGasPaid = resp.PaymasterGasPaid.Add(statisticDay.PaymasterGasPaid).Round(2)
		resp.PaymasterGasPaidUsd = resp.PaymasterGasPaidUsd.Add(statisticDay.PaymasterGasPaidUsd).Round(2)
		resp.UserOpsNum = resp.UserOpsNum + statisticDay.UserOpsNum
		resp.ActiveAAWallet = resp.ActiveAAWallet + statisticDay.ActiveWallet
		resp.LastStatisticTime = statisticDay.CreateTime.UnixMilli()

		detail := &vo.DailyStatisticDetail{
			Time:                  statisticDay.StatisticTime,
			AccumulativeGasFeeUsd: statisticDay.GasFeeUsd.Round(2),
			AccumulativeGasFee:    statisticDay.GasFee.Round(2),
			BundlerGasProfit:      statisticDay.BundlerGasProfit.Round(2),
			BundlerGasProfitUsd:   statisticDay.BundlerGasProfitUsd.Round(2),
			PaymasterGasPaid:      statisticDay.PaymasterGasPaid.Round(2),
			PaymasterGasPaidUsd:   statisticDay.PaymasterGasPaidUsd.Round(2),
			UserOpsNum:            statisticDay.UserOpsNum,
			ActiveAAWallet:        statisticDay.ActiveWallet,
		}
		statisticTimeMap[statisticDay.StatisticTime] = true
		details = append(details, detail)
	}
	sort.Sort(vo.ByDailyStatisticTime(details))
	resp.Details = details
	return &resp
}

func getResponseHour(hours []*ent.DailyStatisticHour) *vo.DailyStatisticResponse {
	if len(hours) == 0 {
		return nil
	}
	var resp *vo.DailyStatisticResponse
	var details []*vo.DailyStatisticDetail
	var statisticTimeMap = make(map[int64]bool)
	for _, statisticHour := range hours {
		if statisticTimeMap[statisticHour.StatisticTime] {
			continue
		}
		resp.AccumulativeGasFeeUsd = resp.AccumulativeGasFeeUsd.Add(statisticHour.GasFeeUsd).Round(2)
		resp.AccumulativeGasFee = resp.AccumulativeGasFee.Add(statisticHour.GasFee).Round(2)
		resp.BundlerGasProfit = resp.BundlerGasProfit.Add(statisticHour.BundlerGasProfit).Round(2)
		resp.BundlerGasProfitUsd = resp.BundlerGasProfitUsd.Add(statisticHour.BundlerGasProfitUsd).Round(2)
		resp.PaymasterGasPaid = resp.PaymasterGasPaid.Add(statisticHour.PaymasterGasPaid).Round(2)
		resp.PaymasterGasPaidUsd = resp.PaymasterGasPaidUsd.Add(statisticHour.PaymasterGasPaidUsd).Round(2)
		resp.UserOpsNum = resp.UserOpsNum + statisticHour.UserOpsNum
		resp.ActiveAAWallet = resp.ActiveAAWallet + statisticHour.ActiveWallet

		detail := &vo.DailyStatisticDetail{
			Time:                  statisticHour.StatisticTime,
			AccumulativeGasFeeUsd: statisticHour.GasFeeUsd.Round(2),
			AccumulativeGasFee:    statisticHour.GasFee.Round(2),
			BundlerGasProfit:      statisticHour.BundlerGasProfit.Round(2),
			BundlerGasProfitUsd:   statisticHour.BundlerGasProfitUsd.Round(2),
			PaymasterGasPaid:      statisticHour.PaymasterGasPaid.Round(2),
			PaymasterGasPaidUsd:   statisticHour.PaymasterGasPaidUsd.Round(2),
			UserOpsNum:            statisticHour.UserOpsNum,
			ActiveAAWallet:        statisticHour.ActiveWallet,
		}
		statisticTimeMap[statisticHour.StatisticTime] = true
		details = append(details, detail)
	}
	sort.Sort(vo.ByDailyStatisticTime(details))
	resp.Details = details
	return resp
}

func GetAATxnDominance(ctx context.Context, req vo.AATxnDominanceRequest) (*vo.AATxnDominanceResponse, error) {

	network := req.Network
	client, err := entity.Client(ctx, network)
	if err != nil {
		return nil, err
	}
	timeRange := req.TimeRange
	var resp *vo.AATxnDominanceResponse

	if timeRange == config.RangeH24 {
		startTime := time.Now().Add(-24 * time.Hour)
		dailyStatisticHours, err := client.DailyStatisticHour.Query().Where(dailystatistichour.StatisticTimeGTE(startTime.UnixMilli()), dailystatistichour.NetworkEqualFold(network)).All(ctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		resp = getDominanceResponseHour(dailyStatisticHours)
	} else if timeRange == config.RangeD7 {
		startTime := time.Now().Add(-7 * 24 * time.Hour)
		dailyStatisticDays, err := client.DailyStatisticDay.Query().Where(dailystatisticday.StatisticTimeGTE(startTime.UnixMilli()), dailystatisticday.NetworkEqualFold(network)).All(ctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		resp = getDominanceResponseDay(dailyStatisticDays)
	} else if timeRange == config.RangeD30 {
		startTime := time.Now().Add(-150 * 24 * time.Hour)
		dailyStatisticDays, err := client.DailyStatisticDay.Query().Where(dailystatisticday.StatisticTimeGTE(startTime.UnixMilli()), dailystatisticday.NetworkEqualFold(network)).All(ctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		resp = getDominanceResponseDay(dailyStatisticDays)
	}

	return resp, nil
}

func getDominanceResponseDay(days []*ent.DailyStatisticDay) *vo.AATxnDominanceResponse {
	if len(days) == 0 {
		return nil
	}
	var resp = vo.AATxnDominanceResponse{}
	var details []*vo.AATxnDominanceDetail
	var statisticMap = make(map[int64]bool)
	for _, statisticDay := range days {
		if statisticMap[statisticDay.StatisticTime] {
			continue
		}
		rate := getRate(statisticDay.TxNum, statisticDay.AaTxNum)
		detail := &vo.AATxnDominanceDetail{
			Time:      statisticDay.StatisticTime,
			Dominance: rate,
		}
		statisticMap[statisticDay.StatisticTime] = true
		details = append(details, detail)
	}
	sort.Sort(vo.ByDominanceTime(details))
	resp.DominanceDetails = details
	return &resp
}

func getDominanceResponseHour(hours []*ent.DailyStatisticHour) *vo.AATxnDominanceResponse {
	if len(hours) == 0 {
		return nil
	}
	var resp = vo.AATxnDominanceResponse{}
	var details []*vo.AATxnDominanceDetail
	var statisticMap = make(map[int64]bool)
	for _, statisticHour := range hours {
		if statisticMap[statisticHour.StatisticTime] {
			continue
		}
		rate := getRate(statisticHour.TxNum, statisticHour.AaTxNum)
		detail := &vo.AATxnDominanceDetail{
			Time:      statisticHour.StatisticTime,
			Dominance: rate,
		}
		statisticMap[statisticHour.StatisticTime] = true
		details = append(details, detail)
	}
	sort.Sort(vo.ByDominanceTime(details))
	resp.DominanceDetails = details
	return &resp
}

func getRate(txNum int64, aaTxNum int64) decimal.Decimal {
	if txNum == 0 {
		return decimal.Zero
	}

	return decimal.NewFromInt(aaTxNum).DivRound(decimal.NewFromInt(txNum), 4)
}

func GetLatestUserOps(ctx context.Context, req vo.LatestUserOpsRequest) (*vo.LatestUserOpsResponse, error) {
	network := req.Network
	client, err := entity.Client(ctx, network)
	if err != nil {
		return nil, err
	}
	var resp = &vo.LatestUserOpsResponse{}

	var ago24h = time.Now().Second() - 24*3600

	count, err := client.AAUserOpsInfo.Query().Where(aauseropsinfo.TxTimeGTE(int64(ago24h))).Count(context.Background())
	allGas, err := client.AAUserOpsInfo.Query().Where(aauseropsinfo.TxTimeGTE(int64(ago24h))).Aggregate(ent.Sum(aauseropsinfo.FieldActualGasCost)).Strings(context.Background())
	if err != nil {
		resp.AverageGasCost24h = decimal.Zero
		return resp, nil
	}
	if len(allGas) == 0 || count == 0 {
		resp.AverageGasCost24h = decimal.Zero
		return resp, nil
	}
	start, err := decimal.NewFromString(allGas[0])
	if err != nil {
		resp.AverageGasCost24h = decimal.Zero
		return resp, nil
	}

	averageGas := rayDiv(start).DivRound(decimal.NewFromInt(int64(count)), 4)
	resp.AverageGasCost24h = averageGas
	return resp, nil
}

func rayDiv(gas decimal.Decimal) decimal.Decimal {
	return gas.DivRound(decimal.NewFromFloat(math.Pow10(18)), 18)
}
