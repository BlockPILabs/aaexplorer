package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

type T1ransactionReceipt struct {
	Time              time.Time
	BlockNum          int64
	CreateTime        time.Time
	BlockHash         string
	BlockNumber       string
	ContractAddress   string
	CumulativeGasUsed string
	EffectiveGasPrice string
	From              string
	GasUsed           string
	Logs              string
	LogsBloom         string
	Status            string
	To                string
	TransactionHash   string
	TransactionIndex  string
	Type              string
	ent.Schema
}

func (T1ransactionReceipt) Fields() []ent.Field {
	return []ent.Field{
		field.Time("time").
			StructTag(`json:"time"`),
		field.Int64("block_num").
			StructTag(`json:"blockNum"`),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
		field.String("block_hash").
			MaxLen(255).
			StructTag(`json:"blockHash"`),
		field.String("block_number").
			MaxLen(255).
			StructTag(`json:"blockNumber"`),
		field.String("contract_address").
			MaxLen(255).
			StructTag(`json:"contractAddress"`),
		field.String("cumulative_gas_used").
			MaxLen(255).
			StructTag(`json:"cumulativeGasUsed"`),
		field.String("effective_gas_price").
			MaxLen(255).
			StructTag(`json:"effective_gas_price"`),
		field.String("from").
			MaxLen(255).
			StructTag(`json:"from"`),
		field.String("gas_used").
			MaxLen(255).
			StructTag(`json:"gasUsed"`),
		field.String("logs").
			MaxLen(255).
			StructTag(`json:"logs"`),
		field.String("logs_bloom").
			MaxLen(255).
			StructTag(`json:"logsBloom"`),
		field.String("status").
			MaxLen(255).
			StructTag(`json:"status"`),
		field.String("to").
			MaxLen(255).
			StructTag(`json:"to"`),
		field.String("transaction_hash").
			MaxLen(255).
			StructTag(`json:"transactionHash"`),
		field.String("transaction_index").
			MaxLen(255).
			StructTag(`json:"transactionIndex"`),
		field.String("type").
			MaxLen(255).
			StructTag(`json:"type"`),
	}
}

func (T1ransactionReceipt) Edges() []ent.Edge {
	return nil
}

func (T1ransactionReceipt) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "t1ransaction_receipt"},
	}
}
