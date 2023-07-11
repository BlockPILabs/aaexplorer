package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

type FactoryInfo struct {
	ent.Schema
}

func (FactoryInfo) Fields() []ent.Field {
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
		field.Int("account_num").
			StructTag(`json:"accountNum"`),
		field.Int("account_deploy_num").
			StructTag(`json:"accountDeployNum"`),
		field.Int("account_num_d1").
			StructTag(`json:"accountNumD1"`),
		field.Int("account_deploy_num_d1").
			StructTag(`json:"accountDeployNumD1"`),
		field.Int("account_num_d7").
			StructTag(`json:"accountNumD7"`),
		field.Int("account_deploy_num_d7").
			StructTag(`json:"accountDeployNumD7"`),
		field.Int("account_num_d30").
			StructTag(`json:"accountNumD30"`),
		field.Int("account_deploy_num_d30").
			StructTag(`json:"accountDeployNumD30"`),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
		field.Time("update_time").
			Optional().
			StructTag(`json:"updateTime"`),
	}
}

func (FactoryInfo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "factory_info"},
	}
}

func (FactoryInfo) Edges() []ent.Edge {
	return nil
}
