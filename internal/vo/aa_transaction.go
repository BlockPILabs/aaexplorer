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
	Hash                 string           `json:"hash"`
	Time                 int64            `json:"time"`
	BlockHash            string           `json:"blockHash"`
	BlockNumber          int64            `json:"blockNumber"`
	UseropCount          int64            `json:"useropCount"`
	IsMev                bool             `json:"isMev"`
	BundlerProfit        decimal.Decimal  `json:"bundlerProfit"`
	BundlerProfitUsd     decimal.Decimal  `json:"bundlerProfitUsd"`
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
	ChainID              int64            `json:"chainID"`
	Type                 string           `json:"type"`
	MaxFeePerGas         *decimal.Decimal `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *decimal.Decimal `json:"maxPriorityFeePerGas"`
	AccessList           *pgtype.JSONB    `json:"accessList"`
	Method               string           `json:"method"`

	ContractAddress   string          `json:"contractAddress"`
	CumulativeGasUsed int64           `json:"cumulativeGasUsed"`
	EffectiveGasPrice string          `json:"effective_gas_price"`
	GasUsed           decimal.Decimal `json:"gasUsed"`
	//Logs              string          `json:"logs"`
	//LogsBloom         string          `json:"logsBloom"`
	Status string `json:"status"`

	TokenPriceUsd decimal.Decimal `json:"tokenPriceUsd"`
	GasPriceUsd   decimal.Decimal `json:"gasPriceUsd"`
	ValueUsd      decimal.Decimal `json:"valueUsd"`
}
