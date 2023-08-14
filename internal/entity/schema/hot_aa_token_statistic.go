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

type HotAATokenStatistic struct {
	Id            int64
	TokenSymbol   string
	Network       string
	StatisticType string
	Volume        decimal.Decimal
	CreateTime    time.Time
	ent.Schema
}

func (HotAATokenStatistic) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			StructTag(`json:"id"`),
		field.String("tokenSymbol").
			MaxLen(255).
			StructTag(`json:"tokenSymbol"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.String("statisticType").
			MaxLen(255).
			StructTag(`json:"statisticType"`),
		field.Int64("volume").
			StructTag(`json:"volume"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 8)"}),
		field.Time("createTime").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
	}
}

func (HotAATokenStatistic) Edges() []ent.Edge {
	return nil
}

func (HotAATokenStatistic) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "hot_aa_token_statistic"},
	}
}
