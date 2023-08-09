package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

type AAContractInteract struct {
	Id              int64
	ContractAddress string
	Network         string
	StatisticType   string
	InteractNum     int64
	CreateTime      time.Time
	ent.Schema
}

func (AAContractInteract) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			StructTag(`json:"id"`),
		field.String("contract_address").
			MaxLen(255).
			StructTag(`json:"contractAddress"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.String("statistic_type").
			MaxLen(255).
			StructTag(`json:"statisticType"`),
		field.Int64("interact_num").
			StructTag(`json:"interactNum"`),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
	}
}

func (AAContractInteract) Edges() []ent.Edge {
	return nil
}

func (AAContractInteract) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "aa_contract_interact"},
	}
}
