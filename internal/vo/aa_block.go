package vo

import (
	"github.com/shopspring/decimal"
)

type GetAaBlocksRequest struct {
	PaginationRequest
	Network           string `json:"network" params:"network" validate:"required,min=3"`
	LatestBlockNumber int64  `json:"latestBlockNumber" params:"latestBlockNumber" validate:"min=0"`
}

type GetAaBlocksResponse struct {
	Pagination
	Records []*AaBlocksVo `json:"records"`
}

type AaBlocksVo struct {
	// ID of the ent.
	Number int64 `json:"number,omitempty"`
	// Time holds the value of the "time" field.
	Time int64 `json:"time,omitempty"`
	// Hash holds the value of the "hash" field.
	Hash string `json:"hash,omitempty"`
	// UseropCount holds the value of the "userop_count" field.
	UseropCount int `json:"useropCount,omitempty"`
	// UseropMevCount holds the value of the "userop_mev_count" field.
	UseropMevCount int `json:"useropMevCount,omitempty"`
	// BundlerProfit holds the value of the "bundler_profit" field.
	BundlerProfit    decimal.Decimal `json:"bundlerProfit,omitempty"`
	BundlerProfitUsd decimal.Decimal `json:"bundlerProfitUsd,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime int64 `json:"createTime,omitempty"`
}
