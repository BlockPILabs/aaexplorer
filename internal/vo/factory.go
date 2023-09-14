package vo

type FactoryVo struct {
	// ID of the ent.
	ID string `json:"factory"`
	// AccountNum holds the value of the "account_num" field.
	AccountNum int `json:"accountNum"`
	// AccountDeployNum holds the value of the "account_deploy_num" field.
	AccountDeployNum int `json:"accountDeployNum"`
}
type GetFactoriesRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
}

type GetFactoriesResponse struct {
	Pagination
	Records []*FactoryVo `json:"records"`
}
