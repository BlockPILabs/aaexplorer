package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

type PaymasterStatisDay struct {
	ID           int64     `json:"id"`
	Paymaster    string    `json:"paymaster"`
	Network      string    `json:"network"`
	UserOpsNum   int64     `json:"user_ops_num"`
	GasSponsored float64   `json:"gas_sponsored"`
	StatisTime   time.Time `json:"statis_time"`
	CreateTime   time.Time `json:"create_time"`
	ent.Schema
}

func (PaymasterStatisDay) Fields() []ent.Field {
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
		field.Time("statis_time").
			StructTag(`json:"statisTime"`).
			Immutable(),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
	}
}

func (PaymasterStatisDay) Edges() []ent.Edge {
	return nil
}

func (PaymasterStatisDay) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "paymaster_statis_day"},
	}
}
