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
	ID                int64           `json:"id"`
	Bundler           string          `json:"bundler"`
	Network           string          `json:"network"`
	UserOpsNum        int64           `json:"user_ops_num"`
	BundlesNum        int64           `json:"bundles_num"`
	GasCollected      decimal.Decimal `json:"gas_collected"`
	FeeEarned         decimal.Decimal `json:"fee_earned"`
	TotalNum          int64           `json:"total_num"`
	StatisTime        time.Time       `json:"statis_time"`
	SuccessBundlesNum int64           `json:"success_bundles_num"`
	FailedBundlesNum  int64           `json:"failed_bundles_num"`
	CreateTime        time.Time       `json:"create_time"`
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
			StructTag(`json:"feeEarned"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional().Nillable(),
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
