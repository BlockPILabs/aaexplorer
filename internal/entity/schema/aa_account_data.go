package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type AaAccountData struct {
	ent.Schema
}

func (AaAccountData) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("address").StructTag(`json:"address"`).Unique(),
		field.String("aa_type").StructTag(`json:"aaType"`),
		field.String("factory").StructTag(`json:"factory"`),
		field.Time("factory_time").StructTag(`json:"factoryTime"`),
		field.Time("update_time").StructTag(`json:"updateTime"`).Optional(),
	}
}

func (Account) AaAccountData() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "aa_account_data"},
	}
}
