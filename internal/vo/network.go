package vo

type NetworkVo struct {
	// Name holds the value of the "name" field.
	Name string `json:"name"`
	// ChainName holds the value of the "chain_name" field.
	ChainName string `json:"chainName"`
	// Network holds the value of the "network" field.
	Network string `json:"network"`
	ChainID int64  `json:"chainId"`
	// IsTestnet holds the value of the "is_testnet" field.
	IsTestnet bool `json:"isTestnet"`
	// Scan holds the value of the "scan" field.
	Scan string `json:"scan"`
	// ScanTx holds the value of the "scan_tx" field.
	ScanTx string `json:"scanTx"`
	// ScanBlock holds the value of the "scan_block" field.
	ScanBlock string `json:"scanBlock"`
	// ScanAddress holds the value of the "scan_address" field.
	ScanAddress string `json:"scanAddress"`
	// ScanName holds the value of the "scan_name" field.
	ScanName string `json:"scanName"`
}
type GetNetworksResponse struct {
	Pagination
	Records []*NetworkVo `json:"records"`
}
