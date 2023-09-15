package vo

import (
	"github.com/shopspring/decimal"
	"time"
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
	Time time.Time `json:"time,omitempty"`
	// Hash holds the value of the "hash" field.
	Hash string `json:"hash,omitempty"`
	// UseropCount holds the value of the "userop_count" field.
	UseropCount int `json:"userop_count,omitempty"`
	// UseropMevCount holds the value of the "userop_mev_count" field.
	UseropMevCount int `json:"userop_mev_count,omitempty"`
	// BundlerProfit holds the value of the "bundler_profit" field.
	BundlerProfit decimal.Decimal `json:"bundler_profit,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
}
