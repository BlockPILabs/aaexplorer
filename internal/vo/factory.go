package vo

import (
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
)

type FactoryVo struct {
	// ID of the ent.
	ID           string `json:"factory"`
	FactoryLabel string `json:"factoryLabel"`
	// AccountNum holds the value of the "account_num" field.
	AccountNum int `json:"accountNum"`
	// AccountNumD1 holds the value of the "account_num_d1" field.
	AccountNumD1 int `json:"accountNumD1"`
	// Dominance holds the value of the "dominance" field.
	Dominance decimal.Decimal `json:"dominance"`
}
type FactoryDbVo struct {
	// ID of the ent.
	ID    string            `json:"factory"`
	Label *pgtype.TextArray `json:"label"`
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

type GetFactoryAccountsRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
	Factory string `json:"factory" params:"factory" validate:"required,len=42"`
}

type GetFactoryRequest struct {
	Network string `json:"network" params:"network" validate:"required,min=3"`
	Factory string `json:"factory" params:"factory" validate:"required,len=42"`
}
type GetFactoryResponse struct {
	TotalAccountDeployNum int   `json:"accountDeployNum"`
	AccountDeployNumD1    int   `json:"accountDeployNumD1"`
	UserOpsNum            int64 `json:"userOpsNum"`
	Rank                  int64 `json:"rank"`
}
