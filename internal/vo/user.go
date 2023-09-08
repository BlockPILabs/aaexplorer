package vo

import "github.com/shopspring/decimal"

type UserBalanceRequest struct {
	Network        string `json:"network"`
	AccountAddress string `json:"accountAddress"`
}

type UserBalanceResponse struct {
	TotalUsd       decimal.Decimal `json:"totalUsd"`
	BalanceDetails []*BalanceDetail
}

type BalanceDetail struct {
	TokenAmount   decimal.Decimal `json:"tokenAmount"`
	Percentage    string          `json:"percentage"`
	TokenValueUsd decimal.Decimal `json:"tokenValueUsd"`
}
