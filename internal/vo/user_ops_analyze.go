package vo

type UserOpsTypeRequest struct {
	Network   string `json:"network"`
	TimeRange string `json:"timeRange"`
}

type UserOpsTypeResponse struct {
	UserOpTypes []*UserOpsType
}

type UserOpsType struct {
	UserOpType string `json:"userOpType"`
	Rate       string `json:"rate"`
}

type AAContractInteractRequest struct {
	Network   string `json:"network"`
	TimeRange string `json:"timeRange"`
}

type AAContractInteractResponse struct {
	AAContractInteract []*AAContractInteract
}

type AAContractInteract struct {
	ContractAddress string `json:"contractAddress"`
	Rate            string `json:"rate"`
}
