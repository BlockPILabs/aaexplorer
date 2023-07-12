package vo

import "time"

type UserOpVo struct {
	// ID of the ent.
	ID int64 `json:"id"`
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
	TxValue float32 `json:"txValue"`
	// Fee holds the value of the "fee" field.
	Fee float32 `json:"fee"`
	// Bundler holds the value of the "bundler" field.
	Bundler string `json:"bundler"`
	// EntryPoint holds the value of the "entry_point" field.
	EntryPoint string `json:"entryPoint"`
	// Factory holds the value of the "factory" field.
	Factory string `json:"factory"`
	// Paymaster holds the value of the "paymaster" field.
	Paymaster string `json:"paymaster"`
	// PaymasterAndData holds the value of the "paymaster_and_data" field.
	PaymasterAndData string `json:"paymasterAndData"`
	// Signature holds the value of the "signature" field.
	Signature string `json:"signature"`
	// Calldata holds the value of the "calldata" field.
	Calldata string `json:"calldata"`
	// Nonce holds the value of the "nonce" field.
	Nonce int64 `json:"nonce"`
	// CallGasLimit holds the value of the "call_gas_limit" field.
	CallGasLimit int64 `json:"callGasLimit"`
	// PreVerificationGas holds the value of the "pre_verification_gas" field.
	PreVerificationGas int64 `json:"preVerificationGas"`
	// VerificationGasLimit holds the value of the "verification_gas_limit" field.
	VerificationGasLimit int64 `json:"verificationGasLimit"`
	// MaxFeePerGas holds the value of the "max_fee_per_gas" field.
	MaxFeePerGas int64 `json:"maxFeePerGas"`
	// MaxPriorityFeePerGas holds the value of the "max_priority_fee_per_gas" field.
	MaxPriorityFeePerGas int64 `json:"maxPriorityFeePerGas"`
	// TxTime holds the value of the "tx_time" field.
	TxTime int64 `json:"txTime"`
	// TxTimeFormat holds the value of the "tx_time_format" field.
	TxTimeFormat string `json:"txTimeFormat"`
	// InitCode holds the value of the "init_code" field.
	InitCode string `json:"initCode"`
	// Status holds the value of the "status" field.
	Status int `json:"status"`
	// Source holds the value of the "source" field.
	Source string `json:"source"`
	// ActualGasCost holds the value of the "actual_gas_cost" field.
	ActualGasCost int64 `json:"actualGasCost"`
	// ActualGasUsed holds the value of the "actual_gas_used" field.
	ActualGasUsed int64 `json:"actualGasUsed"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"createTime"`
}
type GetUserOpsRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
}

type GetUserOpsResponse struct {
	Pagination
	Records []*UserOpVo `json:"records"`
}
