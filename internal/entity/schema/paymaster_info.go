package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
	"time"
)

type PaymasterInfo struct {
	ent.Schema
}

// Fields of the FactoryInfo.
func (PaymasterInfo) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			MaxLen(255).
			StorageKey("paymaster").
			StructTag(`json:"paymaster"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.Int64("user_ops_num").
			StructTag(`json:"userOpsNum"`).Optional(),
		field.Int64("gas_sponsored").
			StructTag(`json:"gasSponsored"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("reserve").
			StructTag(`json:"reserve"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("reserve_usd").
			StructTag(`json:"reserveUsd"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("gas_sponsored_usd").
			StructTag(`json:"gasSponsoredUsd"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("user_ops_num_d1").
			StructTag(`json:"userOpsNumD1"`).Optional(),
		field.Int64("gas_sponsored_d1").
			StructTag(`json:"gasSponsoredD1"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("gas_sponsored_usd_d1").
			StructTag(`json:"gasSponsoredUsdD1"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("user_ops_num_d7").
			StructTag(`json:"userOpsNumD7"`).Optional(),
		field.Int64("gas_sponsored_d7").
			StructTag(`json:"gasSponsoredD7"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("gas_sponsored_usd_d7").
			StructTag(`json:"gasSponsoredUsdD7"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("user_ops_num_d30").
			StructTag(`json:"userOpsNumD30"`).Optional(),
		field.Int64("gas_sponsored_d30").
			StructTag(`json:"gasSponsoredD30"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("gas_sponsored_usd_d30").
			StructTag(`json:"gasSponsoredUsdD30"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
		field.Time("update_time").
			Optional().
			StructTag(`json:"updateTime"`),
	}
}

func (PaymasterInfo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "paymaster_info"},
	}
}

func (PaymasterInfo) Edges() []ent.Edge {
	return []ent.Edge{

		edge.To("account", Account.Type).
			StorageKey(
				edge.Column("address"),
			).
			Unique(),
	}
}
