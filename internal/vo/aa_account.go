package vo

import (
	"github.com/shopspring/decimal"
)

type AaAccountRequestVo struct {
	Network string `json:"network" params:"network" validate:"required,min=3"`
	Address string `json:"address" params:"address" validate:"required,min=3"`
}

type AaAccountRecord struct {
	Address     *string          `json:"address"`
	AaType      *string          `json:"aaType"`
	Factory     *string          `json:"factory"`
	FactoryTime int64            `json:"factoryTime"`
	TotalAmount *decimal.Decimal `json:"totalAmount"`
}

type AaAccountNetworkRequestVo struct {
	Address string `json:"address" params:"address" validate:"required,min=3"`
}

type AaAccountNetworkResponseVo struct {
	Chains []string `json:"chains"`
}
