package vo

import (
	"github.com/shopspring/decimal"
)

type UserOpVo struct {
	// Time holds the value of the "time" field.
	Time int64 `json:"time"`
	// UserOperationHash holds the value of the "user_operation_hash" field.
	UserOperationHash string `json:"userOperationHash"`
	// TxHash holds the value of the "tx_hash" field.
	TxHash string `json:"txHash"`
	// BlockNumber holds the value of the "block_number" field.
	BlockNumber int64 `json:"blockNumber"`
	// Network holds the value of the "network" field.
	Network string `json:"network"`
	// Sender holds the value of the "sender" field.
	Sender      string `json:"sender"`
	SenderLabel string `json:"senderLabel"`
	// Target holds the value of the "target" field.
	Target      string `json:"target"`
	TargetLabel string `json:"targetLabel"`
	// TxValue holds the value of the "tx_value" field.
	TxValue decimal.Decimal `json:"txValue"`
	// Fee holds the value of the "fee" field.
	Fee decimal.Decimal `json:"fee"`
	// Status holds the value of the "status" field.
	Status int32 `json:"status"`
	// Source holds the value of the "source" field.
	Source string `json:"source"`
	// Targets holds the value of the "targets" field.
	Targets []string `json:"targets"`
	// TargetsCount holds the value of the "targets_count" field.
	TargetsCount int `json:"targetsCount"`
	// Bundler holds the value of the "bundler" field.
	Bundler      string `json:"bundler"`
	BundlerLabel string `json:"bundlerLabel"`
	// Paymaster holds the value of the "paymaster" field.
	Paymaster      string `json:"paymaster"`
	PaymasterLabel string `json:"paymasterLabel"`
}
type GetUserOpsRequest struct {
	PaginationRequest
	Network           string `json:"network" params:"network" validate:"required,min=3"`
	LatestBlockNumber int64  `json:"latestBlockNumber" params:"latestBlockNumber" validate:"min=0"`
	BlockNumber       int64  `json:"blockNumber" params:"blockNumber" validate:"min=0"`
	StartTime         int64  `json:"startTime"`
	EndTime           int64  `json:"endTime"`
	TxHash            string `json:"txHash" params:"txHash" validate:"omitempty,txHash"`
	Bundler           string `json:"bundler" params:"bundler" validate:"omitempty,hexAddress"`
	Paymaster         string `json:"paymaster" params:"paymaster" validate:"omitempty,hexAddress"`
	Factory           string `json:"factory" params:"factory" validate:"omitempty,hexAddress"`
	Account           string `json:"account" params:"account" validate:"omitempty,hexAddress"`
	HashTerm          string `json:"hashTerm" params:"hashTerm"`
}

type GetUserOpsResponse struct {
	Pagination
	Records []*UserOpVo `json:"records"`
}

type UserOpsAnalysisRequestVo struct {
	Network           string `json:"network" params:"network" validate:"required,min=3"`
	UserOperationHash string `json:"userOperationHash" params:"userOperationHash" validate:"required,txHash"`
	TxHash            string `json:"txHash" params:"txHash"`
}

type UserOpsAnalysisListRequestVo struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
	TxHash  string `json:"txHash" params:"txHash"`
}
type UserOpsAnalysisListResponse struct {
	Pagination
	Records []*UserOpsAnalysisRecord `json:"records"`
}
type UserOpsAnalysisRecord struct {
	UserOperationHash    string           `json:"userOperationHash"`
	Time                 int64            `json:"time"`
	TxHash               string           `json:"txHash"`
	BlockNumber          int64            `json:"blockNumber"`
	Network              string           `json:"network"`
	Sender               string           `json:"sender"`
	Target               string           `json:"target"`
	Targets              []string         `json:"targets"`
	TargetsCount         int              `json:"targetsCount"`
	TxValue              decimal.Decimal  `json:"txValue"`
	Fee                  decimal.Decimal  `json:"fee"`
	Bundler              string           `json:"bundler"`
	BundlerLabel         []string         `json:"bundlerLabel"`
	EntryPoint           string           `json:"entryPoint"`
	Factory              string           `json:"factory"`
	Paymaster            string           `json:"paymaster"`
	PaymasterLabel       []string         `json:"paymasterLabel"`
	PaymasterAndData     string           `json:"paymasterAndData"`
	Signature            string           `json:"signature"`
	Calldata             string           `json:"calldata"`
	CalldataContract     string           `json:"calldataContract"`
	Nonce                int64            `json:"nonce"`
	CallGasLimit         int64            `json:"callGasLimit"`
	PreVerificationGas   int64            `json:"preVerificationGas"`
	VerificationGasLimit int64            `json:"verificationGasLimit"`
	MaxFeePerGas         int64            `json:"maxFeePerGas"`
	MaxPriorityFeePerGas int64            `json:"maxPriorityFeePerGas"`
	TxTime               int64            `json:"txTime"`
	InitCode             string           `json:"initCode"`
	Status               int32            `json:"status"`
	Source               string           `json:"source"`
	ActualGasCost        int64            `json:"actualGasCost"`
	ActualGasUsed        int64            `json:"actualGasUsed"`
	CreateTime           int64            `json:"createTime"`
	UpdateTime           int64            `json:"updateTime"`
	UsdAmount            *decimal.Decimal `json:"usdAmount"`
	ConfirmBlock         int64            `json:"confirmBlock"`
	AaIndex              int              `json:"aaIndex"`
	FeeUsd               decimal.Decimal  `json:"feeUsd"`
	TxValueUsd           decimal.Decimal  `json:"txValueUsd"`
	BundlerProfit        decimal.Decimal  `json:"bundlerProfit"`
	BundlerProfitUsd     decimal.Decimal  `json:"bundlerProfitUsd"`
	CallData             []CallDataInfo   `json:"callData"`
}

type CallDataInfo struct {
	Time        int64           `json:"time"`
	UserOpsHash string          `json:"userOpsHash"`
	TxHash      string          `json:"txHash"`
	BlockNumber int64           `json:"blockNumber"`
	Network     string          `json:"network"`
	Sender      string          `json:"sender"`
	Target      string          `json:"target"`
	TxValue     decimal.Decimal `json:"txValue"`
	Source      string          `json:"source"`
	Calldata    string          `json:"calldata"`
	TxTime      int64           `json:"txTime"`
	CreateTime  int64           `json:"createTime"`
	UpdateTime  int64           `json:"updateTime"`
	AaIndex     int             `json:"aaIndex"`
}
