package schema

import (
	"entgo.io/ent"
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
		field.String("rpc").StructTag(`json:"rpc"`),
		field.Bool("is_testnet").StructTag(`json:"isTestnet"`),
		field.Time("create_time").Default(time.Now).StructTag(`json:"createTime"`).Immutable(),
		field.Time("update_time").UpdateDefault(time.Now).StructTag(`json:"updateTime"`),
	}
}

// Edges of the Network.
func (Network) Edges() []ent.Edge {
	return nil
}
