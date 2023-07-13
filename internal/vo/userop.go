package vo

import "github.com/shopspring/decimal"

type UserOpVo struct {
	// UserOperationHash holds the value of the "user_operation_hash" field.
	UserOperationHash string `json:"userOperationHash"`
	// TxHash holds the value of the "tx_hash" field.
	TxHash string `json:"txHash"`
	// BlockNumber holds the value of the "block_number" field.
	BlockNumber int64 `json:"blockNumber"`
	// Network holds the value of the "network" field.
	Network string `json:"network"`
	// Sender holds the value of the "sender" field.
	Sender string `json:"sender"`
	// Target holds the value of the "target" field.
	Target string `json:"target"`
	// TxValue holds the value of the "tx_value" field.
	TxValue decimal.Decimal `json:"txValue"`
	// Fee holds the value of the "fee" field.
	Fee decimal.Decimal `json:"fee"`
	// TxTime holds the value of the "tx_time" field.
	TxTime int64 `json:"txTime"`
	// TxTimeFormat holds the value of the "tx_time_format" field.
	TxTimeFormat string `json:"txTimeFormat"`
	// InitCode holds the value of the "init_code" field.
	InitCode string `json:"initCode"`
	// Status holds the value of the "status" field.
	Status int `json:"status"`
}
type GetUserOpsRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
}

type GetUserOpsResponse struct {
	Pagination
	Records []*UserOpVo `json:"records"`
}
