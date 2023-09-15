package vo

type SearchAllRequest struct {
	Network           string `json:"network" params:"network" validate:"required,min=3"`
	Term              string `json:"term" params:"term" validate:"required,min=4"`
	SearchUserOpAndTx bool   `json:"searchUserOpAndTx" params:"searchUserOpAndTx"`
}

type SearchAllAccount struct {
	Address string `json:"address"`
	// AaType holds the value of the "aa_type" field.
	AaType string `json:"aaType"`
}

type SearchAllTransaction struct {
	Hash string `json:"txHash"`
}

type SearchAllBlock struct {
	BlockNumber int64  `json:"blockNumber"`
	BlockHash   string `json:"blockHash"`
}

type SearchAllResponseData struct {
	Type    string
	Records any
}

type SearchAllResponse struct {
	//WalletAccounts []*SearchAllAccount     `json:"walletAccount"`
	//Paymasters     []*SearchAllAccount     `json:"paymaster"`
	//Bundlers       []*SearchAllAccount     `json:"bundler"`
	//Transactions   []*SearchAllTransaction `json:"transactions"`
	//Blocks         []*SearchAllBlock       `json:"blocks"`
	Data []*SearchAllResponseData
}
