package vo

import "github.com/shopspring/decimal"

type BundlersVo struct {
	// Bundler holds the value of the "bundler" field.
	Bundler      string `json:"bundler"`
	BundlerLabel string `json:"bundlerLabel"`
	// BundlesNum holds the value of the "bundles_num" field.
	BundlesNum int64 `json:"bundlesNum"`
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
