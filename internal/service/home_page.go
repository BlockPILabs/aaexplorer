package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/dailystatisticday"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/dailystatistichour"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"log"
	"time"
)

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
	} else if timeRange == config.RangeD7 {
		startTime := time.Now().Add(-7 * 24 * time.Hour)
		dailyStatisticDays, err := client.DailyStatisticDay.Query().Where(dailystatisticday.StatisticTimeGTE(startTime.UnixMilli()), dailystatisticday.NetworkEqualFold(network)).All(ctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		resp = getResponseDay(dailyStatisticDays)
	} else if timeRange == config.RangeD30 {
		startTime := time.Now().Add(-150 * 24 * time.Hour)
		dailyStatisticDays, err := client.DailyStatisticDay.Query().Where(dailystatisticday.StatisticTimeGTE(startTime.UnixMilli()), dailystatisticday.NetworkEqualFold(network)).All(ctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		resp = getResponseDay(dailyStatisticDays)
	}

	return resp, nil
}

func getResponseDay(days []*ent.DailyStatisticDay) *vo.DailyStatisticResponse {
	if len(days) == 0 {
		return nil
	}
	var resp vo.DailyStatisticResponse
	var details []*vo.DailyStatisticDetail
	for _, statisticDay := range days {
		resp.AccumulativeGasFeeUsd = resp.AccumulativeGasFeeUsd.Add(statisticDay.GasFeeUsd).Round(2)
		resp.AccumulativeGasFee = resp.AccumulativeGasFee.Add(statisticDay.GasFee).Round(2)
		resp.BundlerGasProfit = resp.BundlerGasProfit.Add(statisticDay.BundlerGasProfit).Round(2)
		resp.BundlerGasProfitUsd = resp.BundlerGasProfitUsd.Add(statisticDay.BundlerGasProfitUsd).Round(2)
		resp.PaymasterGasPaid = resp.PaymasterGasPaid.Add(statisticDay.PaymasterGasPaid).Round(2)
		resp.PaymasterGasPaidUsd = resp.PaymasterGasPaidUsd.Add(statisticDay.PaymasterGasPaidUsd).Round(2)
		resp.UserOpsNum = resp.UserOpsNum + statisticDay.UserOpsNum
		resp.ActiveAAWallet = resp.ActiveAAWallet + statisticDay.ActiveWallet

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
		details = append(details, detail)
	}

	resp.Details = details
	return &resp
}

func getResponseHour(hours []*ent.DailyStatisticHour) *vo.DailyStatisticResponse {
	if len(hours) == 0 {
		return nil
	}
	var resp *vo.DailyStatisticResponse
	var details []*vo.DailyStatisticDetail
	for _, statisticHour := range hours {
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
		details = append(details, detail)
	}

	resp.Details = details
	return resp
}
