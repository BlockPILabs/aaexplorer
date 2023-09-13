package vo

import "github.com/shopspring/decimal"

type TopRequest struct {
	Type string
}

type TopBundlerRequest struct {
	PaginationRequest
	Network string `json:"network"`
}

type TopPaymasterRequest struct {
	PaginationRequest
	Network string `json:"network"`
}

type TopFactoryRequest struct {
	PaginationRequest
	Network string `json:"network"`
}

type TopBundlerResponse struct {
	Pagination
	BundlerDetails []*BundlerDetail
}

type BundlerDetail struct {
	Address         string          `json:"address"`
	Bundles         int64           `json:"bundles"`
	Success24H      decimal.Decimal `json:"success24H"`
	FeeEarned24H    decimal.Decimal `json:"feeEarned24H"`
	FeeEarnedUsd24H decimal.Decimal `json:"feeEarnedUsd24H"`
}

type TopPaymasterResponse struct {
	Pagination
	PaymasterDetails []*PaymasterDetail
}

type PaymasterDetail struct {
	Address         string          `json:"address"`
	ReserveUsd      decimal.Decimal `json:"reserveUsd"`
	GasSponsored    decimal.Decimal `json:"gasSponsored"`
	GasSponsoredUsd decimal.Decimal `json:"gasSponsoredUsd"`
}

type TopFactoryResponse struct {
	Pagination
	FactoryDetails []*FactoryDetail
}

type FactoryDetail struct {
	Address       string `json:"address"`
	ActiveAccount int64  `json:"activeAccount"`
	TotalAccount  int64  `json:"totalAccount"`
}
