package vo

type SearchAllRequest struct {
	Network string `json:"network" params:"network" validate:"required,min=3"`
	Term    string `json:"term" params:"term" validate:"required,min=1"`
}

type SearchAllAccount struct {
	Address string `json:"address"`
}

type SearchAllTransaction struct {
	TxHash string `json:"txHash"`
}

type SearchAllBlock struct {
	BlockNumber int64 `json:"blockHash"`
}

type SearchAllResponse struct {
	WalletAccounts []*SearchAllAccount     `json:"walletAccount"`
	Paymasters     []*SearchAllAccount     `json:"paymaster"`
	Bundlers       []*SearchAllAccount     `json:"bundler"`
	Transactions   []*SearchAllTransaction `json:"transactions"`
	Blocks         []*SearchAllBlock       `json:"blocks"`
}
