package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
)

type AaTransactionInfo struct {
	ent.Schema
}

func (AaTransactionInfo) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("hash").StructTag(`json:"hash"`),
		field.Time("time").StructTag(`json:"time"`),
		field.String("block_hash").StructTag(`json:"blockHash"`),
		field.Int64("block_number").StructTag(`json:"blockNumber"`),
		field.Int64("userop_count").StructTag(`json:"useropCount"`),
		field.Bool("is_mev").StructTag(`json:"isMev"`),
		field.Int64("bundler_profit").StructTag(`json:"bundlerProfit"`).GoType(decimal.Decimal{}),
		field.Time("create_time").StructTag(`json:"createTime"`),
	}
}

func (AaTransactionInfo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "aa_transaction_info"},
	}
}
