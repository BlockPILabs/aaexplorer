package vo

import "github.com/shopspring/decimal"

type HomeMevRequest struct {
	Network  string `json:"network"`
	BlockNum int64  `json:"blockNum"`
	PaginationRequest
}

type HomeMevResponse struct {
	MevInfos []MevInfo `json:"mevInfos"`
	Pagination
}

type MevInfo struct {
	Time         int64           `json:"time"`
	Type         string          `json:"type"`
	Victim       string          `json:"victim"`
	VictimType   string          `json:"victimType"`
	Attacker     string          `json:"attacker"`
	MevProfit    decimal.Decimal `json:"mevProfit"`
	MevProfitUsd decimal.Decimal `json:"mevProfitUsd"`
	BlockNum     int64           `json:"blockNum"`
}
