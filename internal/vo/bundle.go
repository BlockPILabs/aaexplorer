package vo

import "github.com/shopspring/decimal"

type BundlesVo struct {
	// TxHash holds the value of the "tx_hash" field.
	TxHash string `json:"txHash"`
	// BlockNumber holds the value of the "block_number" field.
	BlockNumber int64 `json:"blockNumber"`
	// Network holds the value of the "network" field.
	Network string `json:"network"`
	// Bundler holds the value of the "bundler" field.
	Bundler string `json:"bundler"`
	// EntryPoint holds the value of the "entry_point" field.
	EntryPoint string `json:"entryPoint"`
	// UserOpsNum holds the value of the "user_ops_num" field.
	UserOpsNum int64 `json:"userOpsNum"`
	// TxValue holds the value of the "tx_value" field.
	TxValue decimal.Decimal `json:"txValue"`
	// Fee holds the value of the "fee" field.
	Fee decimal.Decimal `json:"fee"`
	// Status holds the value of the "status" field.
	Status int `json:"status"`
	// TxTime holds the value of the "tx_time" field.
	TxTime int64 `json:"txTime"`
}
type GetBundlesRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
}

type GetBundlesResponse struct {
	Pagination
	Records []*BundlesVo `json:"records"`
}
