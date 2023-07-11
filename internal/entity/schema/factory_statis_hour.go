package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

type FactoryStatisHour struct {
	ID               int64     `json:"id"`
	Factory          string    `json:"factory"`
	Network          string    `json:"network"`
	AccountNum       int64     `json:"account_num"`
	AccountDeployNum int64     `json:"account_deploy_num"`
	StatisTime       time.Time `json:"statis_time"`
	CreateTime       time.Time `json:"create_time"`
	ent.Schema
}

func (FactoryStatisHour) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			StructTag(`json:"id"`),
		field.String("factory").
			MaxLen(255).
			StructTag(`json:"factory"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.Int64("account_num").
			StructTag(`json:"accountNum"`),
		field.Int64("account_deploy_num").
			StructTag(`json:"accountDeployNum"`),
		field.Time("statis_time").
			StructTag(`json:"statisTime"`).
			Immutable(),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
	}
}

// Edges of the FactoryStatisHour.
func (FactoryStatisHour) Edges() []ent.Edge {
	return nil
}

// Annotations of the FactoryStatisHour.
func (FactoryStatisHour) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "factory_statis_hour"},
	}
}
