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
	UserAddress string `json:"userAddress"`
}

type ListWatchingAddressResponse struct {
	Pagination
	Monitors []*WatchingAddress
}

type WatchingAddress struct {
	ChainName       string  `json:"chainName"`
	AddressType     string  `json:"addressType"`
	MonitorAddress  string  `json:"monitorAddress"`
	Balance         string  `json:"balance"`
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
