package vo

import (
	"github.com/shopspring/decimal"
	"time"
)

type DailyStatisticRequest struct {
	Network   string `json:"network"`
	TimeRange string `json:"timeRange"`
}

type DailyStatisticResponse struct {
	UserOpsNum            int64           `json:"userOpsNum"`
	ActiveAAWallet        int64           `json:"activeAAWallet"`
	AccumulativeGasFee    decimal.Decimal `json:"accumulativeGasFee"`
	AccumulativeGasFeeUsd decimal.Decimal `json:"accumulativeGasFeeUsd"`
	PaymasterGasPaid      decimal.Decimal `json:"paymasterGasPaid"`
	PaymasterGasPaidUsd   decimal.Decimal `json:"paymasterGasPaidUsd"`
	BundlerGasProfit      decimal.Decimal `json:"bundlerGasProfit"`
	BundlerGasProfitUsd   decimal.Decimal `json:"bundlerGasProfitUsd"`
	Details               []*DailyStatisticDetail
}

type DailyStatisticDetail struct {
	Time                  time.Time       `json:"time"`
	UserOpsNum            int64           `json:"userOpsNum"`
	ActiveAAWallet        int64           `json:"activeAAWallet"`
	AccumulativeGasFee    decimal.Decimal `json:"accumulativeGasFee"`
	AccumulativeGasFeeUsd decimal.Decimal `json:"accumulativeGasFeeUsd"`
	PaymasterGasPaid      decimal.Decimal `json:"paymasterGasPaid"`
	PaymasterGasPaidUsd   decimal.Decimal `json:"paymasterGasPaidUsd"`
	BundlerGasProfit      decimal.Decimal `json:"bundlerGasProfit"`
	BundlerGasProfitUsd   decimal.Decimal `json:"bundlerGasProfitUsd"`
}
