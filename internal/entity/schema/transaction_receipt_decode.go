package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
	"time"
)

type TransactionReceiptDecode struct {
	ent.Schema
}

func (TransactionReceiptDecode) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("transaction_hash").
			MaxLen(255).
			StructTag(`json:"transactionHash"`),
		field.Time("time").
			StructTag(`json:"time"`),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
		field.String("block_hash").
			MaxLen(255).
			StructTag(`json:"blockHash"`),
		field.Int64("block_number").
			StructTag(`json:"blockNumber"`),
		field.String("contract_address").
			MaxLen(255).
			StructTag(`json:"contractAddress"`),
		field.Int64("cumulative_gas_used").
			StructTag(`json:"cumulativeGasUsed"`),
		field.String("effective_gas_price").
			StructTag(`json:"effective_gas_price"`),
		field.String("from_addr").
			MaxLen(255).
			StructTag(`json:"from"`),
		field.Int64("gas_used").
			StructTag(`json:"gasUsed"`).GoType(decimal.Decimal{}),
		field.Text("logs").StructTag(`json:"logs"`),
		field.Text("logs_bloom").
			StructTag(`json:"logsBloom"`),
		field.String("status").
			MaxLen(255).
			StructTag(`json:"status"`),
		field.String("to_addr").
			MaxLen(255).
			StructTag(`json:"to"`),
		field.Int64("transaction_index").
			StructTag(`json:"transactionIndex"`),
		field.String("type").
			MaxLen(255).
			StructTag(`json:"type"`),
	}
}

func (TransactionReceiptDecode) Edges() []ent.Edge {
	return nil
}

func (TransactionReceiptDecode) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "transaction_receipt_decode"},
	}
}
