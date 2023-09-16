package vo

import (
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
)

type TransactionRequestVo struct {
	Network string `json:"network" params:"network" validate:"required,min=3"`
	TxHash  string `json:"txHash" params:"txHash" validate:"required,min=3"`
}

type TransactionListRequestVo struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
	TxHash  string `json:"txHash" params:"txHash"`
}
type TransactionListResponse struct {
	Pagination
	Records []*TransactionRecord `json:"records"`
}

type TransactionRecord struct {
	Hash                 string           `json:"hash"`
	Time                 int64            `json:"time"`
	CreateTime           int64            `json:"createTime"`
	BlockHash            string           `json:"blockHash"`
	BlockNumber          int64            `json:"blockNumber"`
	Nonce                decimal.Decimal  `json:"nonce"`
	TransactionIndex     decimal.Decimal  `json:"transactionIndex"`
	FromAddr             string           `json:"fromAddr"`
	ToAddr               string           `json:"toAddr"`
	Value                decimal.Decimal  `json:"value"`
	GasPrice             decimal.Decimal  `json:"gasPrice"`
	Gas                  decimal.Decimal  `json:"gas"`
	Input                string           `json:"input"`
	R                    string           `json:"r"`
	S                    string           `json:"s"`
	V                    decimal.Decimal  `json:"v"`
	ChainID              int64            `json:"chainId"`
	Type                 string           `json:"type"`
	MaxFeePerGas         *decimal.Decimal `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *decimal.Decimal `json:"maxPriorityFeePerGas"`
	AccessList           *pgtype.JSONB    `json:"accessList"`
}
