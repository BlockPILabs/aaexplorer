package vo

import (
	"time"
)

type BundleVo struct {
	// ID of the ent.
	ID int64 `json:"id"`
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
	TxValue float32 `json:"txValue"`
	// Fee holds the value of the "fee" field.
	Fee float32 `json:"fee"`
	// GasPrice holds the value of the "gas_price" field.
	GasPrice string `json:"gasPrice"`
	// GasLimit holds the value of the "gas_limit" field.
	GasLimit int64 `json:"gasLimit"`
	// Status holds the value of the "status" field.
	Status int `json:"status"`
	// TxTime holds the value of the "tx_time" field.
	TxTime int64 `json:"txTime"`
	// TxTimeFormat holds the value of the "tx_time_format" field.
	TxTimeFormat string `json:"txTimeFormat"`
	// Beneficiary holds the value of the "beneficiary" field.
	Beneficiary string `json:"beneficiary"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"createTime"`
}
type GetBundlesRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
}

type GetBundlesResponse struct {
	Pagination
	Records []*BundleVo `json:"records"`
}
