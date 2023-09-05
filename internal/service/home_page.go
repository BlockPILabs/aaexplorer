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
		dailyStatisticHours, err := client.DailyStatisticHour.Query().Where(dailystatistichour.StatisticTimeGTE(startTime), dailystatistichour.NetworkEqualFold(network)).All(ctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		resp = getResponseHour(dailyStatisticHours)
	} else if timeRange == config.RangeD7 {
		startTime := time.Now().Add(-7 * 24 * time.Hour)
		dailyStatisticDays, err := client.DailyStatisticDay.Query().Where(dailystatisticday.StatisticTimeGTE(startTime), dailystatisticday.NetworkEqualFold(network)).All(ctx)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		resp = getResponseDay(dailyStatisticDays)
	} else if timeRange == config.RangeD30 {
		startTime := time.Now().Add(-30 * 24 * time.Hour)
		dailyStatisticDays, err := client.DailyStatisticDay.Query().Where(dailystatisticday.StatisticTimeGTE(startTime), dailystatisticday.NetworkEqualFold(network)).All(ctx)
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
	var resp *vo.DailyStatisticResponse
	var details []*vo.DailyStatisticDetail
	for _, statisticDay := range days {
		resp.AccumulativeGasFeeUsd = resp.AccumulativeGasFeeUsd.Add(statisticDay.GasFee)
		resp.AccumulativeGasFee = resp.AccumulativeGasFee.Add(statisticDay.GasFee)
		resp.BundlerGasProfit = resp.BundlerGasProfit.Add(statisticDay.BundlerGasProfit)
		resp.BundlerGasProfitUsd = resp.BundlerGasProfitUsd.Add(statisticDay.BundlerGasProfitUsd)
		resp.PaymasterGasPaid = resp.PaymasterGasPaid.Add(statisticDay.PaymasterGasPaid)
		resp.PaymasterGasPaidUsd = resp.PaymasterGasPaidUsd.Add(statisticDay.PaymasterGasPaidUsd)
		resp.UserOpsNum = resp.UserOpsNum + statisticDay.UserOpsNum
		resp.ActiveAAWallet = resp.ActiveAAWallet + statisticDay.ActiveWallet

		detail := &vo.DailyStatisticDetail{
			Time:                  statisticDay.StatisticTime,
			AccumulativeGasFeeUsd: statisticDay.GasFee,
			AccumulativeGasFee:    statisticDay.GasFee,
			BundlerGasProfit:      statisticDay.BundlerGasProfit,
			BundlerGasProfitUsd:   statisticDay.BundlerGasProfitUsd,
			PaymasterGasPaid:      statisticDay.PaymasterGasPaid,
			PaymasterGasPaidUsd:   statisticDay.PaymasterGasPaidUsd,
			UserOpsNum:            statisticDay.UserOpsNum,
			ActiveAAWallet:        statisticDay.ActiveWallet,
		}
		details = append(details, detail)
	}

	resp.Details = details
	return resp
}

func getResponseHour(hours []*ent.DailyStatisticHour) *vo.DailyStatisticResponse {
	if len(hours) == 0 {
		return nil
	}
	var resp *vo.DailyStatisticResponse
	var details []*vo.DailyStatisticDetail
	for _, statisticHour := range hours {
		resp.AccumulativeGasFeeUsd = resp.AccumulativeGasFeeUsd.Add(statisticHour.GasFee)
		resp.AccumulativeGasFee = resp.AccumulativeGasFee.Add(statisticHour.GasFee)
		resp.BundlerGasProfit = resp.BundlerGasProfit.Add(statisticHour.BundlerGasProfit)
		resp.BundlerGasProfitUsd = resp.BundlerGasProfitUsd.Add(statisticHour.BundlerGasProfitUsd)
		resp.PaymasterGasPaid = resp.PaymasterGasPaid.Add(statisticHour.PaymasterGasPaid)
		resp.PaymasterGasPaidUsd = resp.PaymasterGasPaidUsd.Add(statisticHour.PaymasterGasPaidUsd)
		resp.UserOpsNum = resp.UserOpsNum + statisticHour.UserOpsNum
		resp.ActiveAAWallet = resp.ActiveAAWallet + statisticHour.ActiveWallet

		detail := &vo.DailyStatisticDetail{
			Time:                  statisticHour.StatisticTime,
			AccumulativeGasFeeUsd: statisticHour.GasFee,
			AccumulativeGasFee:    statisticHour.GasFee,
			BundlerGasProfit:      statisticHour.BundlerGasProfit,
			BundlerGasProfitUsd:   statisticHour.BundlerGasProfitUsd,
			PaymasterGasPaid:      statisticHour.PaymasterGasPaid,
			PaymasterGasPaidUsd:   statisticHour.PaymasterGasPaidUsd,
			UserOpsNum:            statisticHour.UserOpsNum,
			ActiveAAWallet:        statisticHour.ActiveWallet,
		}
		details = append(details, detail)
	}

	resp.Details = details
	return resp
}
