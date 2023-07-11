package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

type PaymasterInfo struct {
	ent.Schema
}

// Fields of the FactoryInfo.
func (PaymasterInfo) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			StructTag(`json:"id"`),
		field.String("paymaster").
			MaxLen(255).
			StructTag(`json:"paymaster"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.Int64("user_ops_num").
			StructTag(`json:"userOpsNum"`),
		field.Float32("gas_sponsored").
			StructTag(`json:"gasSponsored"`),
		field.Int64("user_ops_num_d1").
			StructTag(`json:"userOpsNumD1"`),
		field.Float32("gas_sponsored_d1").
			StructTag(`json:"gasSponsoredD1"`),
		field.Int64("user_ops_num_d7").
			StructTag(`json:"userOpsNumD7"`),
		field.Float32("gas_sponsored_d7").
			StructTag(`json:"gasSponsoredD7"`),
		field.Int64("user_ops_num_d30").
			StructTag(`json:"userOpsNumD30"`),
		field.Float("gas_sponsored_d30").
			StructTag(`json:"gasSponsoredD30"`),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
		field.Time("update_time").
			Optional().
			StructTag(`json:"updateTime"`),
	}
}

func (PaymasterInfo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "paymaster_info"},
	}
}

func (PaymasterInfo) Edges() []ent.Edge {
	return nil
}
