package vo

import "github.com/shopspring/decimal"

type TopRequest struct {
	Type string
}

type TopBundlerResponse struct {
	Address         string
	Bundles         int64
	Success24H      decimal.Decimal
	FeeEarned24H    decimal.Decimal
	FeeEarnedUsd24H decimal.Decimal
}

type TopPaymasterResponse struct {
	Address         string
	Reserve         decimal.Decimal
	GasSponsored    decimal.Decimal
	GasSponsoredUsd decimal.Decimal
}

type TopFactoryResponse struct {
	Address       string
	ActiveAccount int64
	TotalAccount  int64
}
