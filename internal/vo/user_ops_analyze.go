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

type ByContractNum []*AAContractInteract

func (b ByContractNum) Len() int      { return len(b) }
func (b ByContractNum) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByContractNum) Less(i, j int) bool {
	return b[i].Rate.Cmp(b[j].Rate) > 0
}

type ByUserOpsTypeNum []*UserOpsType

func (b ByUserOpsTypeNum) Len() int      { return len(b) }
func (b ByUserOpsTypeNum) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByUserOpsTypeNum) Less(i, j int) bool {
	return b[i].Rate.Cmp(b[j].Rate) > 0
}
