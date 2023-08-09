package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

type TransactionReceipt struct {
	ent.Schema
}

func (TransactionReceipt) Fields() []ent.Field {
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

func (TransactionReceipt) Edges() []ent.Edge {
	return nil
}

func (TransactionReceipt) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "transaction_receipt"},
	}
}
