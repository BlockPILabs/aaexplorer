package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
	"time"
)

type BundlerStatisHour struct {
	ent.Schema
}

func (BundlerStatisHour) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			StructTag(`json:"id"`),
		field.String("bundler").
			MaxLen(255).
			StructTag(`json:"bundler"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.Int64("user_ops_num").
			StructTag(`json:"userOpsNum"`),
		field.Int64("bundles_num").
			StructTag(`json:"bundlesNum"`),
		field.Int64("gas_collected").
			StructTag(`json:"gasCollected"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("fee_earned").
			StructTag(`json:"feeEarned"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("total_num").
			StructTag(`json:"totalNum"`).Optional(),
		field.Int64("success_bundles_num").
			StructTag(`json:"successBundlesNum"`),
		field.Int64("failed_bundles_num").
			StructTag(`json:"failedBundlesNum"`),
		field.Time("statis_time").
			StructTag(`json:"statisTime"`).
			Immutable(),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
	}
}

func (BundlerStatisHour) Edges() []ent.Edge {
	return nil
}

func (BundlerStatisHour) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "bundler_statis_hour"},
	}
}
