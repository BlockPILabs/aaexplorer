package vo

type FactoryVo struct {
	// ID of the ent.
	ID string `json:"factory"`
	// AccountNum holds the value of the "account_num" field.
	AccountNum int `json:"accountNum"`
	// AccountNumD1 holds the value of the "account_num_d1" field.
	AccountNumD1 int `json:"accountNumD1"`
}
type GetFactoriesRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
}

type GetFactoriesResponse struct {
	Pagination
	Records []*FactoryVo `json:"records"`
}
