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

type WhaleStatisticDay struct {
	Id            int64
	Network       string
	WhaleNum      int64
	TotalUsd      decimal.Decimal
	StatisticTime time.Time
	CreateTime    time.Time
	ent.Schema
}

func (WhaleStatisticDay) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			StructTag(`json:"id"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.Int64("whale_num").
			StructTag(`json:"whaleNum"`),
		field.Int64("total_usd").
			StructTag(`json:"totalUsd"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 6)"}),
		field.Time("statistic_time").
			StructTag(`json:"statisticTime"`),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
	}
}

func (WhaleStatisticDay) Edges() []ent.Edge {
	return nil
}

func (WhaleStatisticDay) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "whale_statistic_day"},
	}
}
