package vo

import "time"

type BlocksVo struct {
	// Time holds the value of the "time" field.
	Time time.Time `json:"time"`
	// BlockNum holds the value of the "block_num" field.
	BlockNum int64 `json:"blockNum"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"createTime"`
	// Hash holds the value of the "hash" field.
	Hash string `json:"hash"`
	// Size holds the value of the "size" field.
	Size string `json:"size"`
	// Miner holds the value of the "miner" field.
	Miner string `json:"miner"`
	// Nonce holds the value of the "nonce" field.
	Nonce string `json:"nonce"`
	// Number holds the value of the "number" field.
	Number string `json:"number"`
	// Uncles holds the value of the "uncles" field.
	Uncles string `json:"uncles"`
	// GasUsed holds the value of the "gas_used" field.
	GasUsed string `json:"gasUsed"`
	// MixHash holds the value of the "mix_hash" field.
	MixHash string `json:"mixHash"`
	// GasLimit holds the value of the "gas_limit" field.
	GasLimit string `json:"gasLimit"`
	// ExtraData holds the value of the "extra_data" field.
	ExtraData string `json:"extraData"`
	// LogsBloom holds the value of the "logs_bloom" field.
	LogsBloom string `json:"logsBloom"`
	// StateRoot holds the value of the "state_root" field.
	StateRoot string `json:"stateRoot"`
	// Timestamp holds the value of the "timestamp" field.
	Timestamp string `json:"timestamp"`
	// Difficulty holds the value of the "difficulty" field.
	Difficulty string `json:"difficulty"`
	// ParentHash holds the value of the "parent_hash" field.
	ParentHash string `json:"parentHash"`
	// Sha3Uncles holds the value of the "sha3_uncles" field.
	Sha3Uncles string `json:"sha3Uncles"`
	// Withdrawals holds the value of the "withdrawals" field.
	Withdrawals string `json:"withdrawals"`
	// ReceiptsRoot holds the value of the "receipts_root" field.
	ReceiptsRoot string `json:"receiptsRoot"`
	// Transactions holds the value of the "transactions" field.
	Transactions string `json:"transactions"`
	// BaseFeePerGas holds the value of the "base_fee_per_gas" field.
	BaseFeePerGas string `json:"baseFeePerGas"`
	// TotalDifficulty holds the value of the "total_difficulty" field.
	TotalDifficulty string `json:"totalDifficulty"`
	// WithdrawalsRoot holds the value of the "withdrawals_root" field.
	WithdrawalsRoot string `json:"withdrawalsRoot"`
	// TransactionsRoot holds the value of the "transactions_root" field.
	TransactionsRoot string `json:"transactionsRoot"`
}
type GetBlocksRequest struct {
	PaginationRequest
	Network string `json:"network" params:"network" validate:"required,min=3"`
}

type GetBlocksResponse struct {
	Pagination
	Records []*BlocksVo `json:"records"`
}

type GetBlockRequest struct {
	Block string `json:"block" params:"block" validate:"required,min=3"`
}

type GetBlockResponse struct {
	Block *BlocksVo `json:"block"`
}
