package vo

type PaymastersVo struct {
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
	// Status holds the value of the "status" field.
	Status int `json:"status"`
	// TxTime holds the value of the "tx_time" field.
	TxTime int64 `json:"txTime"`
}
type GetPaymastersRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
}

type GetPaymastersResponse struct {
	Pagination
	Records []*PaymastersVo `json:"records"`
}