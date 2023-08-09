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

type UserAssetInfo struct {
	ent.Schema
}

func (UserAssetInfo) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			StructTag(`json:"id"`),
		field.String("account_address").
			MaxLen(255).
			StructTag(`json:"accountAddress"`),
		field.String("contract_address").
			MaxLen(255).
			StructTag(`json:"contractAddress"`),
		field.String("symbol").
			MaxLen(255).
			StructTag(`json:"symbol"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.Int64("amount").StructTag(`json:"amount"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("last_time").
			StructTag(`json:"lastTime"`),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
	}
}

func (UserAssetInfo) Edges() []ent.Edge {
	return nil
}

func (UserAssetInfo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "user_asset_info"},
	}
}
