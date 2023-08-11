package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

type BlockData struct {
	Time             time.Time `db:"time"`
	BlockNum         int64     `db:"block_num"`
	CreateTime       time.Time `db:"create_time"`
	Hash             string    `db:"hash"`
	Size             string    `db:"size"`
	Miner            string    `db:"miner"`
	Nonce            string    `db:"nonce"`
	Number           string    `db:"number"`
	Uncles           string    `db:"uncles" json:"uncles"`
	GasUsed          string    `db:"gas_used"`
	MixHash          string    `db:"mix_hash"`
	GasLimit         string    `db:"gas_limit"`
	ExtraData        string    `db:"extra_data"`
	LogsBloom        string    `db:"logs_bloom"`
	StateRoot        string    `db:"state_root"`
	Timestamp        string    `db:"timestamp"`
	Difficulty       string    `db:"difficulty"`
	ParentHash       string    `db:"parent_hash"`
	Sha3Uncles       string    `db:"sha3_uncles"`
	Withdrawals      string    `db:"withdrawals" json:"withdrawals"`
	ReceiptsRoot     string    `db:"receipts_root"`
	BaseFeePerGas    string    `db:"base_fee_per_gas"`
	TotalDifficulty  string    `db:"total_difficulty"`
	WithdrawalsRoot  string    `db:"withdrawals_root"`
	TransactionsRoot string    `db:"transactions_root"`
	ent.Schema
}

func (BlockData) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			StructTag(`json:"blockNum"`).Unique().StorageKey("block_num"),
		field.Time("time").
			StructTag(`json:"time"`),
		field.Time("create_time").
			StructTag(`json:"createTime"`),
		field.String("hash").
			MaxLen(255).
			StructTag(`json:"hash"`),
		field.String("size").
			MaxLen(255).
			StructTag(`json:"size"`),
		field.String("miner").
			MaxLen(255).
			StructTag(`json:"miner"`),
		field.String("nonce").
			MaxLen(255).
			StructTag(`json:"nonce"`),
		field.String("number").
			MaxLen(255).
			StructTag(`json:"number"`),
		field.String("uncles").
			StructTag(`json:"uncles"`),
		field.String("gas_used").
			MaxLen(255).
			StructTag(`json:"gasUsed"`),
		field.String("mix_hash").
			MaxLen(255).
			StructTag(`json:"mixHash"`),
		field.String("gas_limit").
			MaxLen(255).
			StructTag(`json:"gasLimit"`),
		field.String("extra_data").
			MaxLen(255).
			StructTag(`json:"extraData"`),
		field.String("logs_bloom").
			MaxLen(255).
			StructTag(`json:"logsBloom"`),
		field.String("state_root").
			MaxLen(255).
			StructTag(`json:"stateRoot"`),
		field.String("timestamp").
			MaxLen(255).
			StructTag(`json:"timestamp"`),
		field.String("difficulty").
			MaxLen(255).
			StructTag(`json:"difficulty"`),
		field.String("parent_hash").
			MaxLen(255).
			StructTag(`json:"parentHash"`),
		field.String("sha3_uncles").
			MaxLen(255).
			StructTag(`json:"sha3Uncles"`),
		field.String("withdrawals").
			StructTag(`json:"withdrawals"`),
		field.String("receipts_root").
			MaxLen(255).
			StructTag(`json:"receiptsRoot"`),
		field.String("base_fee_per_gas").
			MaxLen(255).
			StructTag(`json:"baseFeePerGas"`),
		field.String("total_difficulty").
			MaxLen(255).
			StructTag(`json:"totalDifficulty"`),
		field.String("withdrawals_root").
			MaxLen(255).
			StructTag(`json:"withdrawalsRoot"`),
		field.String("transactions_root").
			MaxLen(255).
			StructTag(`json:"transactionsRoot"`),
	}
}

func (BlockData) Edges() []ent.Edge {
	return nil
}

func (BlockData) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "block_data"},
	}
}
