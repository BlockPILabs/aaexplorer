package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Network holds the schema definition for the Network entity.
type Network struct {
	ent.Schema
}

// Fields of the Network.
func (Network) Fields() []ent.Field {
	return []ent.Field{
		field.Int("age"),
		field.String("name"),
		field.String("nickname").
			Unique(),
		field.JSON("test", []string{}),
	}
}

// Edges of the Network.
func (Network) Edges() []ent.Edge {
	return nil
}
