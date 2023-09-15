package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
)

type AaBlockInfo struct {
	ent.Schema
}

func (AaBlockInfo) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").StorageKey("number"),
		field.Time("time"),
		field.String("hash"),
		field.Int("userop_count"),
		field.Int("userop_mev_count"),
		field.Int64("bundler_profit").GoType(decimal.Decimal{}),
		field.Int64("bundler_profit_usd").GoType(decimal.Decimal{}),
		field.Time("create_time"),
	}
}

func (AaBlockInfo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "aa_block_info"},
	}
}
