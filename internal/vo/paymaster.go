package vo

import "github.com/shopspring/decimal"

type PaymastersVo struct {
	// Paymaster holds the value of the "paymaster" field.
	Paymaster      string `json:"paymaster"`
	PaymasterLabel string `json:"paymasterLabel"`
	// UserOpsNum holds the value of the "user_ops_num" field.
	UserOpsNum int64 `json:"userOpsNum"`
	// UserOpsNumD1 holds the value of the "user_ops_num_d1" field.
	UserOpsNumD1 int64 `json:"userOpsNumD1"`
	// Reserve holds the value of the "reserve" field.
	Reserve decimal.Decimal `json:"reserve"`
	// GasSponsored holds the value of the "gas_sponsored" field.
	GasSponsored decimal.Decimal `json:"gasSponsored"`
	// GasSponsoredUsd holds the value of the "gas_sponsored_usd" field.
	GasSponsoredUsd decimal.Decimal `json:"gasSponsoredUsd"`
}
type GetPaymastersRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
}

type GetPaymastersResponse struct {
	Pagination
	Records []*PaymastersVo `json:"records"`
}

type GetReserveRequest struct {
	PaginationRequest
	Paymaster string `json:"paymaster"`
	Network   string `json:"network"`
}

type GetReserveResponse struct {
	Pagination
	Records []*ReserveDetail
}

type ReserveDetail struct {
	Timestamp int64           `json:"timestamp"`
	Sender    string          `json:"sender"`
	Target    string          `json:"target"`
	Status    string          `json:"status"`
	Bundler   string          `json:"bundler"`
	Balance   decimal.Decimal `json:"balance"`
}

type GetPaymasterOverviewRequest struct {
	Paymaster string `json:"paymaster"`
	Network   string `json:"network"`
}

type GetPaymasterOverviewResponse struct {
	SponsorGasFeeUsd24h   decimal.Decimal `json:"sponsorGasFee24h"`
	UserOpsNum24h         int64           `json:"userOpsNum24h"`
	SponsorGasFeeUsdTotal decimal.Decimal `json:"sponsorGasFeeTotal"`
	UserOpsNumTotal       int64           `json:"userOpsNumTotal"`
	Dominance             decimal.Decimal `json:"dominance"`
	Rank                  int             `json:"rank"`
	AccountBalance        decimal.Decimal `json:"accountBalance"`
	EstimateTxnNum        int             `json:"estimateTxnNum"`
	Label                 []string        `json:"label"`
}
