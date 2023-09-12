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

type AATxnDominanceRequest struct {
	Network   string `json:"network"`
	TimeRange string `json:"timeRange"`
}

type AATxnDominanceResponse struct {
	DominanceDetails []*AATxnDominanceDetail
}

type AATxnDominanceDetail struct {
	Time      int64  `json:"time"`
	Dominance string `json:"dominance"`
}

type LatestUserOpsRequest struct {
	Network string `json:"network"`
}

type LatestUserOpsResponse struct {
}
