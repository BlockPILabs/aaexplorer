package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

type Monitor struct {
	ent.Schema
}

func (Monitor) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").StorageKey("id").StructTag(`json:"id"`),
		field.String("user_address").StructTag(`json:"userAddress"`),
		field.String("monitor_address").StructTag(`json:"monitorAddress"`),
		field.String("monitor_address_type").StructTag(`json:"monitorAddressType"`),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
	}
}

func (Monitor) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "monitor"},
	}
}
