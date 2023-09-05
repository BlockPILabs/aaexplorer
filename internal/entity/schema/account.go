package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/jackc/pgtype"
)

type Account struct {
	ent.Schema
}

func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("address").StructTag(`json:"address"`),
		field.Bool("is_contract").StructTag(`json:"is_contract"`),
		field.Other("tag", &pgtype.TextArray{}).StructTag(`json:"tag"`).SchemaType(map[string]string{
			dialect.Postgres: "text[]",
		}),
		field.Other("label", &pgtype.TextArray{}).StructTag(`json:"label"`).SchemaType(map[string]string{
			dialect.Postgres: "text[]",
		}),
		field.String("abi").StructTag(`json:"abi"`),
	}
}

func (Account) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "account"},
	}
}
