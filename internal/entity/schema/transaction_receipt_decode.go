package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
)

type TransactionReceiptDecode struct {
	ent.Schema
}

func (TransactionReceiptDecode) Fields() []ent.Field {
	return []ent.Field{
		field.Time("time").
			StructTag(`json:"time"`),
		field.Time("create_time").
			StructTag(`json:"createTime"`),
		field.String("transaction_hash").
			StructTag(`json:"transactionHash"`),
		field.Int64("transaction_index").
			StructTag(`json:"transactionIndex"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 0)"}),
		field.String("block_hash").
			StructTag(`json:"blockHash"`),
		field.Int64("block_number").
			StructTag(`json:"blockNumber"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 0)"}),
		field.Int64("cumulative_gas_used").
			StructTag(`json:"cumulativeGasUsed"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("gas_used").
			StructTag(`json:"gasUsed"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.String("contract_address").
			StructTag(`json:"contractAddress"`),
		field.String("root").
			StructTag(`json:"root"`),
		field.String("status").
			StructTag(`json:"status"`),
		field.String("from_addr").
			StructTag(`json:"fromAddr"`),
		field.String("to_addr").
			StructTag(`json:"toAddr"`),
		field.Bytes("logs").
			StructTag(`json:"logs"`),
		field.String("logs_bloom").
			StructTag(`json:"logsBloom"`),
		field.String("revert_reason").
			StructTag(`json:"revertReason"`),
		field.String("type").
			StructTag(`json:"type"`),
		field.String("effective_gas_price").
			StructTag(`json:"effectiveGasPrice"`),
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
