package vo

import "github.com/shopspring/decimal"

type BundlerVo struct {
	// ID of the ent.
	ID int64 `json:"id"`
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
	// UserOpsNumD1 holds the value of the "user_ops_num_d1" field.
	UserOpsNumD1 int64 `json:"userOpsNumD1"`
	// BundlesNumD1 holds the value of the "bundles_num_d1" field.
	BundlesNumD1 int64 `json:"bundlesNumD1"`
	// GasCollectedD1 holds the value of the "gas_collected_d1" field.
	GasCollectedD1 decimal.Decimal `json:"gasCollectedD1"`
	// UserOpsNumD7 holds the value of the "user_ops_num_d7" field.
	UserOpsNumD7 int64 `json:"userOpsNumD7"`
	// BundlesNumD7 holds the value of the "bundles_num_d7" field.
	BundlesNumD7 int64 `json:"bundlesNumD7"`
	// GasCollectedD7 holds the value of the "gas_collected_d7" field.
	GasCollectedD7 decimal.Decimal `json:"gasCollected_d7"`
	// UserOpsNumD30 holds the value of the "user_ops_num_d30" field.
	UserOpsNumD30 int64 `json:"userOpsNumD30"`
	// BundlesNumD30 holds the value of the "bundles_num_d30" field.
	BundlesNumD30 int64 `json:"bundlesNumD30"`
	// GasCollectedD30 holds the value of the "gas_collected_d30" field.
	GasCollectedD30 decimal.Decimal `json:"gasCollectedD30"`
}
type GetBundlersRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
}

type GetBundlersResponse struct {
	Pagination
	Records []*BundlerVo `json:"records"`
}
