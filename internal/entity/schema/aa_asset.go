package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
	"time"
)

type AaAsset struct {
	ent.Schema
}

func (AaAsset) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("user_address").StructTag(`json:"user_address"`).Unique(),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.Int64("balance").StructTag(`json:"balance"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("asset_value").StructTag(`json:"assetValue"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("last_time").
			StructTag(`json:"lastTime"`),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
		field.Time("update_time").
			Default(time.Now).
			UpdateDefault(time.Now).
			StructTag(`json:"updateTime"`).
			Immutable(),
	}
}

func (AaAsset) Edges() []ent.Edge {
	return nil
}

func (AaAsset) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "aa_asset"},
	}
}
