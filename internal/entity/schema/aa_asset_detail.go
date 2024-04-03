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

type AaAssetDetail struct {
	ent.Schema
}

func (AaAssetDetail) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			StructTag(`json:"id"`),
		field.String("user_address").StructTag(`json:"user_address"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.String("contract_address").
			MaxLen(255).
			StructTag(`json:"contractAddress"`),
		field.String("symbol").
			MaxLen(255).
			StructTag(`json:"symbol"`),
		field.String("is_native").
			MaxLen(255).
			StructTag(`json:"isNative"`),
		field.Int64("asset_amount").StructTag(`json:"assetAmount"`).GoType(decimal.Zero).
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

func (AaAssetDetail) Edges() []ent.Edge {
	return nil
}

func (AaAssetDetail) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "aa_asset_detail"},
	}
}
