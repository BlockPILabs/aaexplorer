package aa

type Log struct {
	Data                string   `json:"data"`
	Topics              []string `json:"topics"`
	Address             string   `json:"address"`
	Removed             bool     `json:"removed"`
	LogIndex            int      `json:"logIndex"`
	BlockHash           string   `json:"blockHash"`
	BlockNumber         int      `json:"blockNumber"`
	LogIndexRaw         string   `json:"logIndexRaw"`
	BlockNumberRaw      string   `json:"blockNumberRaw"`
	TransactionHash     string   `json:"transactionHash"`
	TransactionIndex    int      `json:"transactionIndex"`
	TransactionIndexRaw string   `json:"transactionIndexRaw"`
}
