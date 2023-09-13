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

type BundlerInfo struct {
	ent.Schema
}

func (BundlerInfo) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StructTag(`json:"bundler"`).StorageKey("bundler").MaxLen(255),
		field.String("network").StructTag(`json:"network"`).MaxLen(255),
		field.Int64("user_ops_num").StructTag(`json:"userOpsNum"`).Optional(),
		field.Int64("bundles_num").StructTag(`json:"bundlesNum"`).Optional(),
		field.Int64("gas_collected").StructTag(`json:"gasCollected"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("fee_earned").StructTag(`json:"feeEarned"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("fee_earned_usd").StructTag(`json:"feeEarnedUsd"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("success_rate").StructTag(`json:"successRate"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 4)"}).Optional(),
		field.Int64("bundle_rate").StructTag(`json:"bundleRate"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 4)"}).Optional(),
		field.Int64("success_bundles_num").
			StructTag(`json:"successBundlesNum"`).Optional(),
		field.Int64("failed_bundles_num").
			StructTag(`json:"failedBundlesNum"`).Optional(),

		field.Int64("user_ops_num_d1").StructTag(`json:"userOpsNumD1"`).Optional(),
		field.Int64("bundles_num_d1").StructTag(`json:"bundlesNumD1"`).Optional(),
		field.Int64("gas_collected_d1").StructTag(`json:"gasCollectedD1"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("fee_earned_d1").StructTag(`json:"feeEarnedD1"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("fee_earned_usd_d1").StructTag(`json:"feeEarnedUsdD1"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("success_rate_d1").StructTag(`json:"successRateD1"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 4)"}).Optional(),
		field.Int64("bundle_rate_d1").StructTag(`json:"bundleRateD1"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 4)"}).Optional(),

		field.Int64("user_ops_num_d7").StructTag(`json:"userOpsNumD7"`).Optional(),
		field.Int64("bundles_num_d7").StructTag(`json:"bundlesNumD7"`).Optional(),
		field.Int64("gas_collected_d7").StructTag(`json:"gasCollected_d7"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("fee_earned_d7").StructTag(`json:"feeEarnedD7"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("fee_earned_usd_d7").StructTag(`json:"feeEarnedUsdD7"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("success_rate_d7").StructTag(`json:"successRateD7"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 4)"}).Optional(),
		field.Int64("bundle_rate_d7").StructTag(`json:"bundleRateD7"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 4)"}).Optional(),

		field.Int64("user_ops_num_d30").StructTag(`json:"userOpsNumD30"`).Optional(),
		field.Int64("bundles_num_d30").StructTag(`json:"bundlesNumD30"`).Optional(),
		field.Int64("gas_collected_d30").StructTag(`json:"gasCollectedD30"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("fee_earned_d30").StructTag(`json:"feeEarnedD30"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("fee_earned_usd_d30").StructTag(`json:"feeEarnedUsdD30"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}).Optional(),
		field.Int64("success_rate_d30").StructTag(`json:"successRateD30"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 4)"}).Optional(),
		field.Int64("bundle_rate_d30").StructTag(`json:"bundleRateD30"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 4)"}).Optional(),
		field.Time("create_time").StructTag(`json:"createTime"`).Default(time.Now).Immutable(),
		field.Time("update_time").StructTag(`json:"updateTime"`).UpdateDefault(time.Now).Optional(),
	}
}

func (BundlerInfo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "bundler_info"},
	}
}

func (BundlerInfo) Edges() []ent.Edge {
	return nil
}
