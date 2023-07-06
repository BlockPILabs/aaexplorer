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
		field.Int8("id").StructTag(`json:""id")"`),
		field.String("bundler").StructTag(`json:""bundler")"`).MaxLen(255),
		field.String("network").StructTag(`json:""network")"`).MaxLen(255),
		field.Int8("user_ops_num").StructTag(`json:""userOpsNum")"`),
		field.Int8("bundles_num").StructTag(`json:""bundlesNum")"`),
		field.Int8("gas_collected").StructTag(`json:""gasCollected")"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(30, 20)"}),
		field.Int8("user_ops_num_d1").StructTag(`json:""userOpsNumD1")"`),
		field.Int8("bundles_num_d1").StructTag(`json:""bundlesNumD1")"`),
		field.Int8("gas_collected_d1").StructTag(`json:""gasCollectedD1")"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(30, 20)"}),
		field.Int8("user_ops_num_d7").StructTag(`json:""userOpsNumD7")"`),
		field.Int8("bundles_num_d7").StructTag(`json:""bundlesNumD7")"`),
		field.Int8("gas_collected_d7").StructTag(`json:""gasCollected_d7")"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(30, 20)"}),
		field.Int8("user_ops_num_d30").StructTag(`json:""userOpsNumD30")"`),
		field.Int8("bundles_num_d30").StructTag(`json:""bundlesNumD30")"`),
		field.Int8("gas_collected_d30").StructTag(`json:""gasCollectedD30")"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(30, 20)"}),
		field.Time("create_time").StructTag(`json:""createTime")"`).Default(time.Now).Immutable(),
		field.Time("update_time").StructTag(`json:""updateTime")"`).UpdateDefault(time.Now),
	}
}

func (BundlerInfo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "bundler_info"},
	}
}
