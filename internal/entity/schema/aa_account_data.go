package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
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
		field.Int64("user_ops_num").StructTag(`json:"userOpsNum"`).Optional(),
		field.Int64("total_balance_usd").StructTag(`json:"totalBalanceUsd"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("last_time").StructTag(`json:"lastTime"`).Optional(),
		field.Time("update_time").StructTag(`json:"updateTime"`).Optional(),
	}
}

func (AaAccountData) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "aa_account_data"},
	}
}

func (AaAccountData) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("account", Account.Type).
			StorageKey(
				edge.Column("address"),
			).
			Unique(),
	}
}
