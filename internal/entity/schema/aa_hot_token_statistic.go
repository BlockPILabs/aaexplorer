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

type AAHotTokenStatistic struct {
	Id              int64
	Symbol          string
	ContractAddress string
	Network         string
	StatisticType   string
	UsdAmount       decimal.Decimal
	CreateTime      time.Time
	ent.Schema
}

func (AAHotTokenStatistic) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			StructTag(`json:"id"`),
		field.String("symbol").
			MaxLen(255).
			StructTag(`json:"symbol"`),
		field.String("contract_address").
			MaxLen(255).
			StructTag(`json:"contractAddress"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.String("statistic_type").
			MaxLen(255).
			StructTag(`json:"statisticType"`),
		field.Int64("usd_amount").
			StructTag(`json:"usdAmount"`).
			GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
	}
}

func (AAHotTokenStatistic) Edges() []ent.Edge {
	return nil
}

func (AAHotTokenStatistic) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "aa_hot_token_statistic"},
	}
}
