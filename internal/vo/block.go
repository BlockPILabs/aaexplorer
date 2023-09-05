package vo

import (
	"github.com/shopspring/decimal"
	"time"
)

type BlocksVo struct {
	// ID of the ent.
	ID int64 `json:"number"`
	//// Time holds the value of the "time" field.
	//Time time.Time `json:"time"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"createTime"`
	//// Hash holds the value of the "hash" field.
	//Hash string `json:"hash"`
	//// ParentHash holds the value of the "parent_hash" field.
	//ParentHash string `json:"parentHash"`
	//// Nonce holds the value of the "nonce" field.
	//Nonce decimal.Decimal `json:"nonce"`
	//// Sha3Uncles holds the value of the "sha3_uncles" field.
	//Sha3Uncles string `json:"sha3Uncles"`
	//// LogsBloom holds the value of the "logs_bloom" field.
	//LogsBloom string `json:"logsBloom"`
	//// TransactionsRoot holds the value of the "transactions_root" field.
	//TransactionsRoot string `json:"transactionsRoot"`
	//// StateRoot holds the value of the "state_root" field.
	//StateRoot string `json:"stateRoot"`
	//// ReceiptsRoot holds the value of the "receipts_root" field.
	//ReceiptsRoot string `json:"receiptsRoot"`
	//// Miner holds the value of the "miner" field.
	//Miner string `json:"miner"`
	//// MixHash holds the value of the "mix_hash" field.
	//MixHash string `json:"mixHash"`
	//// Difficulty holds the value of the "difficulty" field.
	//Difficulty decimal.Decimal `json:"difficulty"`
	//// TotalDifficulty holds the value of the "total_difficulty" field.
	//TotalDifficulty decimal.Decimal `json:"totalDifficulty"`
	//// ExtraData holds the value of the "extra_data" field.
	//ExtraData string `json:"extraData"`
	//// Size holds the value of the "size" field.
	//Size decimal.Decimal `json:"size"`
	//// GasLimit holds the value of the "gas_limit" field.
	//GasLimit decimal.Decimal `json:"gasLimit"`
	//// GasUsed holds the value of the "gas_used" field.
	//GasUsed decimal.Decimal `json:"gasUsed"`
	//// Timestamp holds the value of the "timestamp" field.
	//Timestamp decimal.Decimal `json:"timestamp"`
	//// TransactionCount holds the value of the "transaction_count" field.
	//TransactionCount decimal.Decimal `json:"transactionCount"`
	//// Uncles holds the value of the "uncles" field.
	//Uncles []string `json:"uncles"`
	//// BaseFeePerGas holds the value of the "base_fee_per_gas" field.
	//BaseFeePerGas decimal.Decimal `json:"baseFeePerGas"`
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
	Block   string `json:"block" params:"block" validate:"required,min=3"`
	Network string `json:"network" params:"network" validate:"required,min=3"`
}

type BlockVo struct {
	// ID of the ent.
	ID int64 `json:"number"`
	// Time holds the value of the "time" field.
	Time time.Time `json:"time"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"createTime"`
	// Hash holds the value of the "hash" field.
	Hash string `json:"hash"`
	// ParentHash holds the value of the "parent_hash" field.
	ParentHash string `json:"parentHash"`
	// Nonce holds the value of the "nonce" field.
	Nonce decimal.Decimal `json:"nonce"`
	// Sha3Uncles holds the value of the "sha3_uncles" field.
	Sha3Uncles string `json:"sha3Uncles"`
	// LogsBloom holds the value of the "logs_bloom" field.
	LogsBloom string `json:"logsBloom"`
	// TransactionsRoot holds the value of the "transactions_root" field.
	TransactionsRoot string `json:"transactionsRoot"`
	// StateRoot holds the value of the "state_root" field.
	StateRoot string `json:"stateRoot"`
	// ReceiptsRoot holds the value of the "receipts_root" field.
	ReceiptsRoot string `json:"receiptsRoot"`
	// Miner holds the value of the "miner" field.
	Miner string `json:"miner"`
	// MixHash holds the value of the "mix_hash" field.
	MixHash string `json:"mixHash"`
	// Difficulty holds the value of the "difficulty" field.
	Difficulty decimal.Decimal `json:"difficulty"`
	// TotalDifficulty holds the value of the "total_difficulty" field.
	TotalDifficulty decimal.Decimal `json:"totalDifficulty"`
	// ExtraData holds the value of the "extra_data" field.
	ExtraData string `json:"extraData"`
	// Size holds the value of the "size" field.
	Size decimal.Decimal `json:"size"`
	// GasLimit holds the value of the "gas_limit" field.
	GasLimit decimal.Decimal `json:"gasLimit"`
	// GasUsed holds the value of the "gas_used" field.
	GasUsed decimal.Decimal `json:"gasUsed"`
	// Timestamp holds the value of the "timestamp" field.
	Timestamp decimal.Decimal `json:"timestamp"`
	// TransactionCount holds the value of the "transaction_count" field.
	TransactionCount decimal.Decimal `json:"transactionCount"`
	// Uncles holds the value of the "uncles" field.
	Uncles []string `json:"uncles"`
	// BaseFeePerGas holds the value of the "base_fee_per_gas" field.
	BaseFeePerGas decimal.Decimal `json:"baseFeePerGas"`
}
type GetBlockResponse struct {
	Block *BlockVo `json:"block"`
}
