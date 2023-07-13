package vo

import "github.com/shopspring/decimal"

type BundlersVo struct {
	// Bundler holds the value of the "bundler" field.
	Bundler string `json:"bundler"`
	// Network holds the value of the "network" field.
	Network string `json:"network"`
	// UserOpsNum holds the value of the "user_ops_num" field.
	UserOpsNum int64 `json:"userOpsNum"`
	// BundlesNum holds the value of the "bundles_num" field.
	BundlesNum int64 `json:"bundlesNum"`
	// GasCollected holds the value of the "gas_collected" field.
	GasCollected decimal.Decimal `json:"gasCollected"`
}
type GetBundlersRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
}

type GetBundlersResponse struct {
	Pagination
	Records []*BundlersVo `json:"records"`
}
