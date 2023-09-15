package vo

import (
	"github.com/shopspring/decimal"
)

type DailyStatisticRequest struct {
	Network   string `json:"network"`
	TimeRange string `json:"timeRange"`
}

type DailyStatisticResponse struct {
	UserOpsNum            int64           `json:"userOpsNum"`
	Ups                   decimal.Decimal `json:"ups"`
	ActiveAAWallet        int64           `json:"activeAAWallet"`
	AccumulativeGasFee    decimal.Decimal `json:"accumulativeGasFee"`
	AccumulativeGasFeeUsd decimal.Decimal `json:"accumulativeGasFeeUsd"`
	PaymasterGasPaid      decimal.Decimal `json:"paymasterGasPaid"`
	PaymasterGasPaidUsd   decimal.Decimal `json:"paymasterGasPaidUsd"`
	BundlerGasProfit      decimal.Decimal `json:"bundlerGasProfit"`
	BundlerGasProfitUsd   decimal.Decimal `json:"bundlerGasProfitUsd"`
	LastStatisticTime     int64           `json:"lastStatisticTime"`
	Details               []*DailyStatisticDetail
}

type DailyStatisticDetail struct {
	Time                  int64           `json:"time"`
	UserOpsNum            int64           `json:"userOpsNum"`
	ActiveAAWallet        int64           `json:"activeAAWallet"`
	AccumulativeGasFee    decimal.Decimal `json:"accumulativeGasFee"`
	AccumulativeGasFeeUsd decimal.Decimal `json:"accumulativeGasFeeUsd"`
	PaymasterGasPaid      decimal.Decimal `json:"paymasterGasPaid"`
	PaymasterGasPaidUsd   decimal.Decimal `json:"paymasterGasPaidUsd"`
	BundlerGasProfit      decimal.Decimal `json:"bundlerGasProfit"`
	BundlerGasProfitUsd   decimal.Decimal `json:"bundlerGasProfitUsd"`
}

type ByDailyStatisticTime []*DailyStatisticDetail

func (b ByDailyStatisticTime) Len() int           { return len(b) }
func (b ByDailyStatisticTime) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByDailyStatisticTime) Less(i, j int) bool { return b[i].Time < b[j].Time }

type AATxnDominanceRequest struct {
	Network   string `json:"network"`
	TimeRange string `json:"timeRange"`
}

type AATxnDominanceResponse struct {
	DominanceDetails []*AATxnDominanceDetail
}

type AATxnDominanceDetail struct {
	Time      int64           `json:"time"`
	Dominance decimal.Decimal `json:"dominance"`
}

type LatestUserOpsRequest struct {
	Network string `json:"network"`
}

type LatestUserOpsResponse struct {
	AverageProcessTime24h decimal.Decimal `json:"averageProcessTime24h"`
	AverageGasCost24h     decimal.Decimal `json:"averageGasCost24h"`
	PendingTransactionNum int64           `json:"pendingTransactionNum"`
}

type ByDominanceTime []*AATxnDominanceDetail

func (b ByDominanceTime) Len() int           { return len(b) }
func (b ByDominanceTime) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByDominanceTime) Less(i, j int) bool { return b[i].Time < b[j].Time }
