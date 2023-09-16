package vo

import (
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
)

type AaTransactionRequestVo struct {
	Network string `json:"network" params:"network" validate:"required,min=3"`
	TxHash  string `json:"txHash" params:"txHash" validate:"required,min=3"`
}

type AaTransactionListRequestVo struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
	TxHash  string `json:"txHash" params:"txHash"`
	Address string `json:"address" params:"address"`
}
type AaTransactionListResponse struct {
	Pagination
	Records []*AaTransactionRecord `json:"records"`
}

type AaTransactionRecord struct {
	HASH                     *string          `json:"hash"`
	TIME                     int64            `json:"time"`
	BLOCK_HASH               *string          `json:"blockHash"`
	BLOCK_NUMBER             *int64           `json:"blockNumber"`
	USEROP_COUNT             *int64           `json:"useropCount"`
	IS_MEV                   *bool            `json:"isMev"`
	BUNDLER_PROFIT           *decimal.Decimal `json:"bundlerProfit"`
	NONCE                    *decimal.Decimal `json:"nonce"`
	TRANSACTION_INDEX        *decimal.Decimal `json:"transactionIndex"`
	FROM_ADDR                *string          `json:"fromAddr"`
	TO_ADDR                  *string          `json:"toAddr"`
	VALUE                    *decimal.Decimal `json:"value"`
	GAS_PRICE                *decimal.Decimal `json:"gasPrice"`
	GAS                      *decimal.Decimal `json:"gas"`
	INPUT                    *string          `json:"input"`
	R                        *string          `json:"r"`
	S                        *string          `json:"s"`
	V                        *decimal.Decimal `json:"v"`
	CHAIN_ID                 *int64           `json:"chainId"`
	TYPE                     *string          `json:"type"`
	MAX_FEE_PER_GAS          *decimal.Decimal `json:"maxFeePerGas"`
	MAX_PRIORITY_FEE_PER_GAS *decimal.Decimal `json:"maxPriorityFeePerGas"`
	ACCESS_LIST              *pgtype.JSONB    `json:"accessList"`
}
