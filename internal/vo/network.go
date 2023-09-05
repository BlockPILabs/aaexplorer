package vo

type NetworkVo struct {
	// Name holds the value of the "name" field.
	Name string `json:"name"`
	// Network holds the value of the "network" field.
	Network string `json:"network"`
	ChainID int64  `json:"chainId"`
	// IsTestnet holds the value of the "is_testnet" field.
	IsTestnet bool `json:"isTestnet"`
}
type GetNetworksResponse struct {
	Pagination
	Records []*NetworkVo `json:"records"`
}
