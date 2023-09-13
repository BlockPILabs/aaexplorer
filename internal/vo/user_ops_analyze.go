package vo

import "github.com/shopspring/decimal"

type UserOpsTypeRequest struct {
	Network   string `json:"network"`
	TimeRange string `json:"timeRange"`
}

type UserOpsTypeResponse struct {
	UserOpTypes []*UserOpsType
}

type UserOpsType struct {
	UserOpType string          `json:"userOpType"`
	Rate       decimal.Decimal `json:"rate"`
}

type AAContractInteractRequest struct {
	Network   string `json:"network"`
	TimeRange string `json:"timeRange"`
}

type AAContractInteractResponse struct {
	TotalNum           int64 `json:"totalNum"`
	AAContractInteract []*AAContractInteract
}

type AAContractInteract struct {
	ContractAddress string          `json:"contractAddress"`
	Rate            decimal.Decimal `json:"rate"`
	SingleNum       int64           `json:"singleNum"`
}
