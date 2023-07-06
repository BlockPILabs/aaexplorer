package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
	"time"
)

type BundlerInfo struct {
	ent.Schema
}

func (BundlerInfo) Fields() []ent.Field {
	return []ent.Field{
		field.Int8("id"),
		field.String("bundler").MaxLen(255),
		field.String("network").MaxLen(255),
		field.Int8("user_ops_num"),
		field.Int8("bundles_num"),
		field.Int8("gas_collected").GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(30, 20)"}),
		field.Int8("user_ops_num_d1"),
		field.Int8("bundles_num_d1"),
		field.Int8("gas_collected_d1").GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(30, 20)"}),
		field.Int8("user_ops_num_d7"),
		field.Int8("bundles_num_d7"),
		field.Int8("gas_collected_d7").GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(30, 20)"}),
		field.Int8("user_ops_num_d30"),
		field.Int8("bundles_nunum_d30"),
		field.Int8("gas_collected_d30").GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(30, 20)"}),
		field.Time("create_time").Default(time.Now).Immutable(),
		field.Time("update_time").UpdateDefault(time.Now),
	}
}
