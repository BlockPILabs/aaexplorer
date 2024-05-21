package vo

import "github.com/shopspring/decimal"

type AddMonitorRequest struct {
	Network        string `json:"network"`
	MonitorAddress string `json:"monitorAddress"`
	Sign           string `json:"sign"`
	Message        string `json:"message"`
}

type AddMonitorResponse struct {
}

type RemoveMonitorRequest struct {
	Network        string `json:"network"`
	MonitorAddress string `json:"monitorAddress"`
	Sign           string `json:"sign"`
	Message        string `json:"message"`
}

type RemoveMonitorResponse struct {
}

type ListWatchingAddressRequest struct {
	PaginationRequest
	Network     string `json:"network" params:"network" validate:"required,min=3"`
	UserAddress string `json:"userAddress" param:"userAddress" validate:"required"`
}

type ListWatchingAddressResponse struct {
	Pagination
	Monitors []*WatchingAddress `json:"monitors"`
}

type WatchingAddress struct {
	Network         string  `json:"network"`
	AddressType     string  `json:"addressType"`
	MonitorAddress  string  `json:"monitorAddress"`
	Balance         float64 `json:"balance"`
	Profits24H      float64 `json:"profits24H"`
	SponsoredGas24H float64 `json:"sponsoredGas24H"`
	TotalUserOps    int64   `json:"totalUserOps"`
}
type AssetDetailRequest struct {
	Network     string `json:"network"`
	UserAddress string `json:"userAddress"`
}

type AssetDetailResponse struct {
	TotalAssetUsd decimal.Decimal `json:"totalAssetUsd"`
	AssetDetails  []AssetDetail   `json:"assetDetails"`
}

type AssetDetail struct {
	Percent   decimal.Decimal `json:"percent"`
	Symbol    string          `json:"symbol"`
	Network   string          `json:"network"`
	Amount    decimal.Decimal `json:"amount"`
	AmountUsd decimal.Decimal `json:"amountUsd"`
}

type MevInfoRequest struct {
	Network string `json:"network"`
}

type MevInfoResponse struct {
	AttackerTotalProfit float64 `json:"attackerTotalProfit"`
	MevUserOpsRatio     float64 `json:"mevUserOpsRatio"`
	AttackerAccounts    int     `json:"attackerAccounts"`
}

type AccountTypeRequest struct {
	AccountAddress string `json:"accountAddress"`
	Network        string `json:"network"`
}

type AccountTypeResponse struct {
	AccountAddress string `json:"accountAddress"`
	Type           string `json:"type"`
}
