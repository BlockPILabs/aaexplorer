package vo

import (
	"github.com/shopspring/decimal"
)

type BundlersVo struct {
	// Bundler holds the value of the "bundler" field.
	Bundler      string `json:"bundler"`
	BundlerLabel string `json:"bundlerLabel"`
	// BundlesNum holds the value of the "bundles_num" field.
	BundlesNum int64 `json:"bundlesNum"`

	// BundleRate holds the value of the "bundle_rate" field.
	BundleRate decimal.Decimal `json:"bundleRate"`

	// UserOpsNum holds the value of the "user_ops_num" field.
	UserOpsNum int64 `json:"userOpsNum"`
	// SuccessRate holds the value of the "success_rate" field.
	SuccessRate decimal.Decimal `json:"successRate"`
	// SuccessRateD1 holds the value of the "success_rate_d1" field.
	SuccessRateD1 decimal.Decimal `json:"successRateD1"`
	// BundlesNumD1 holds the value of the "bundles_num_d1" field.
	BundlesNumD1 int64 `json:"bundlesNumD1"`
	// FeeEarnedD1 holds the value of the "fee_earned_d1" field.
	FeeEarnedD1 decimal.Decimal `json:"feeEarnedD1"`
	// FeeEarnedUsdD1 holds the value of the "fee_earned_usd_d1" field.
	FeeEarnedUsdD1 decimal.Decimal `json:"feeEarnedUsdD1"`
}
type GetBundlersRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
}

type GetBundlersResponse struct {
	Pagination
	Records []*BundlersVo `json:"records"`
}

type GetBundlerRequest struct {
	Network string `json:"network" params:"network" validate:"required,min=3"`
	Bundler string `json:"bundler" params:"bundler" validate:"required,hexAddress"`
}
type GetBundlerResponse struct {
	// FeeEarnedUsdD1 holds the value of the "fee_earned_usd_d1" field.
	FeeEarnedUsdD1 decimal.Decimal `json:"feeEarnedUsdD1"`
	// FeeEarnedUsd holds the value of the "fee_earned_usd" field.
	FeeEarnedUsd decimal.Decimal `json:"feeEarnedUsd"`
	// SuccessRateD1 holds the value of the "success_rate_d1" field.
	SuccessRateD1 decimal.Decimal `json:"successRateD1"`
	// SuccessRate holds the value of the "success_rate" field.
	SuccessRate decimal.Decimal `json:"successRate"`
	// BundleRate holds the value of the "bundle_rate" field.
	BundleRate    decimal.Decimal `json:"bundleRate"`
	Rank          int64           `json:"rank"`
	TotalBundlers int64           `json:"totalBundlers"`
	Label         []string        `json:"label"`
}

type ListMEVBundlersRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
}

type ListBlockMEVBundlersRequest struct {
	PaginationRequest
	Network     string `json:"network" params:"network" validate:"required,min=3"`
	BlockNumber int64  `json:"blockNumber" params:"blockNumber" validate:"required"`
}

type ListBundlerMEVBundlersRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
	Bundler string `json:"bundler" params:"bundler" validate:"required"`
}

type ListMEVBundlersResponse struct {
	Pagination
	Records []*MEVBundlerAssets `json:"records"`
}

type MEVBundlerAssets struct {
	Timestamp        int64           `json:"timestamp"`
	UserOpHash       string          `json:"userOpHash"`
	TxHash           string          `json:"txHash"`
	MevType          string          `json:"mevType"`
	Victim           string          `json:"victim"`
	VictimType       string          `json:"victimType"`
	Attacker         string          `json:"attacker"`
	BundlerLoss      decimal.Decimal `json:"bundlerLoss"`
	BundlerLossInUsd decimal.Decimal `json:"bundlerLossInUsd"`
	MevProfits       decimal.Decimal `json:"mevProfits"`
	MevProfitsInUsd  decimal.Decimal `json:"mevProfitsInUsd"`
}

type BundlerTransferRequest struct {
	Network string `json:"network"`
	Address string `json:"address"`
	PaginationRequest
}

type BundlerTransferResponse struct {
	Pagination
	TransferList []TransferInfo `json:"transferList"`
}

type TransferInfo struct {
	Id          int64           `json:"id"`
	TxnHash     string          `json:"txnHash"`
	Source      string          `json:"source"`
	Timestamp   int64           `json:"timestamp"`
	From        string          `json:"from"`
	To          string          `json:"to"`
	Value       decimal.Decimal `json:"value"`
	TokenSymbol string          `json:"tokenSymbol"`
	TokenImage  string          `json:"tokenImage"`
}

type BundlerBalanceRequest struct {
	Network string `json:"network"`
	Address string `json:"address"`
}

type BundlerBalanceResponse struct {
	TotalUsd     decimal.Decimal `json:"totalUsd"`
	AssetDetails []AssetInfo     `json:"assetDetails"`
}

type AssetInfo struct {
	Percent   decimal.Decimal `json:"percent"`
	Symbol    string          `json:"symbol"`
	TokenUrl  string          `json:"tokenImage"`
	Amount    decimal.Decimal `json:"amount"`
	AmountUsd decimal.Decimal `json:"amountUsd"`
}
