package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

// Network holds the schema definition for the Network entity.
type Network struct {
	ent.Schema
}

// Fields of the Network.
func (Network) Fields() []ent.Field {
	return []ent.Field{
		field.Int8("id").StructTag(`json:"id"`),
		field.String("name").MaxLen(255).StructTag(`json:"name"`),
		field.String("network").Unique().MaxLen(127).StructTag(`json:"network"`),
		field.String("logo").MaxLen(255).StructTag(`json:"logo"`),
		field.String("http_rpc").StructTag(`json:"http_rpc"`),
		field.Bool("is_testnet").StructTag(`json:"isTestnet"`),
		field.Time("create_time").Default(time.Now).StructTag(`json:"createTime"`).Immutable(),
		field.Time("update_time").Optional().UpdateDefault(time.Now).StructTag(`json:"updateTime"`),
	}
}

// Edges of the Network.
func (Network) Edges() []ent.Edge {
	return nil
}

func (Network) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "network"},
	}
}
