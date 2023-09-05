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

type PaymasterStatisHour struct {
	ID              int64           `json:"id"`
	Paymaster       string          `json:"paymaster"`
	Network         string          `json:"network"`
	UserOpsNum      int64           `json:"user_ops_num"`
	Reserve         decimal.Decimal `json:"reserve"`
	GasSponsored    decimal.Decimal `json:"gas_sponsored"`
	GasSponsoredUsd decimal.Decimal `json:"gas_sponsored_usd"`
	StatisTime      time.Time       `json:"statis_time"`
	CreateTime      time.Time       `json:"create_time"`
	ent.Schema
}

func (PaymasterStatisHour) Fields() []ent.Field {
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
			StructTag(`json:"userOpsNum"`).Optional(),
		field.Int64("gas_sponsored").
			StructTag(`json:"gasSponsored"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("reserve").
			StructTag(`json:"reserve"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("reserve_usd").
			StructTag(`json:"reserveUsd"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Time("statis_time").
			StructTag(`json:"statisTime"`).
			Immutable(),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
	}
}

func (PaymasterStatisHour) Edges() []ent.Edge {
	return nil
}

func (PaymasterStatisHour) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "paymaster_statis_hour"},
	}
}
