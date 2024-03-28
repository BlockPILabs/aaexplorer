package aa

type Log struct {
	Data                string   `json:"data"`
	Topics              []string `json:"topics"`
	Address             string   `json:"address"`
	Removed             bool     `json:"removed"`
	LogIndex            int64    `json:"logIndex"`
	BlockHash           string   `json:"blockHash"`
	BlockNumber         int64    `json:"blockNumber"`
	LogIndexRaw         string   `json:"logIndexRaw"`
	BlockNumberRaw      string   `json:"blockNumberRaw"`
	TransactionHash     string   `json:"transactionHash"`
	TransactionIndex    int64    `json:"transactionIndex"`
	TransactionIndexRaw string   `json:"transactionIndexRaw"`
}
