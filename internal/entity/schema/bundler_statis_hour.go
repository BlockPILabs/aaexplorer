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
	ID           int64           `json:"id"`
	Bundler      string          `json:"bundler"`
	Network      string          `json:"network"`
	UserOpsNum   int64           `json:"user_ops_num"`
	BundlesNum   int64           `json:"bundles_num"`
	GasCollected decimal.Decimal `json:"gas_collected"`
	StatisTime   time.Time       `json:"statis_time"`
	CreateTime   time.Time       `json:"create_time"`
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
