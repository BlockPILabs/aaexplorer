package vo

import (
	"github.com/shopspring/decimal"
	"time"
)

type AaAccountRequestVo struct {
	Network string `json:"network" params:"network" validate:"required,min=3"`
	Address string `json:"address" params:"address" validate:"required,min=3"`
}

type AaAccountRecord struct {
	Address     *string          `json:"address"`
	AaType      *string          `json:"aaType"`
	Factory     *string          `json:"factory"`
	FactoryTime *time.Time       `json:"factoryTime"`
	TotalAmount *decimal.Decimal `json:"totalAmount"`
}

// AaAccountData is the model entity for the AaAccountData schema.
type AaAccountDataVo struct {
	// ID of the ent.
	ID           string `json:"address"`
	AddressLable string `json:"addressLable"`
	// AaType holds the value of the "aa_type" field.
	AaType string `json:"aaType"`
	// Factory holds the value of the "factory" field.
	Factory string `json:"factory"`
	// FactoryTime holds the value of the "factory_time" field.
	FactoryTime int64 `json:"factoryTime"`
	// UserOpsNum holds the value of the "user_ops_num" field.
	UserOpsNum int64 `json:"userOpsNum"`
	// TotalBalanceUsd holds the value of the "total_balance_usd" field.
	TotalBalanceUsd decimal.Decimal `json:"totalBalanceUsd"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime int64 `json:"updateTime"`
}

type GetAccountsRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
	Address string `json:"address" params:"address"`
	Factory string `json:"factory" params:"factory"`
}

type GetAccountsResponse struct {
	Pagination
	Records []*AaAccountDataVo `json:"records"`
}

type AaAccountNetworkRequestVo struct {
	Address string `json:"address" params:"address" validate:"required,min=3"`
}

type AaAccountNetworkResponseVo struct {
	Chains []string `json:"chains"`
}
