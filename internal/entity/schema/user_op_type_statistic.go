package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

type UserOpTypeStatistic struct {
	Id            int64
	UserOpType    string
	UserOpSign    string
	Network       string
	StatisticType string
	OpNum         int64
	CreateTime    time.Time
	ent.Schema
}

func (UserOpTypeStatistic) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			StructTag(`json:"id"`),
		field.String("user_op_type").
			MaxLen(255).
			StructTag(`json:"userOpType"`),
		field.String("user_op_sign").
			MaxLen(255).
			StructTag(`json:"userOpSign"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.String("statistic_type").
			StructTag(`json:"statisticType"`),
		field.Int64("op_num").
			StructTag(`json:"opNum"`),
		field.Time("create_time").
			StructTag(`json:"createTime"`),
	}
}

func (UserOpTypeStatistic) Edges() []ent.Edge {
	return nil
}

func (UserOpTypeStatistic) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "user_op_type_statistic"},
	}
}
