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

type MevInfo struct {
	Id            int64
	TxHash        string
	BlockNumber   int64
	Network       string
	TxFromTag     string
	TxFrom        string
	TxTo          string
	GasFee        decimal.Decimal
	UserOpsGasFee decimal.Decimal
	Profit        decimal.Decimal
	TxTime        int64
	CreateTime    time.Time
	ent.Schema
}

func (MevInfo) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			StructTag(`json:"id"`),
		field.String("tx_hash").
			MaxLen(255).
			StructTag(`json:"txHash"`),
		field.Int64("block_number").
			StructTag(`json:"blockNumber"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.String("tx_from_tag").
			MaxLen(255).
			StructTag(`json:"txFromTag"`),
		field.String("tx_from").
			MaxLen(255).
			StructTag(`json:"txFrom"`),
		field.String("tx_to").
			MaxLen(255).
			StructTag(`json:"txTo"`),
		field.Int64("gas_fee").
			StructTag(`json:"gasFee"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("user_ops_gas_fee").
			StructTag(`json:"userOpsGasFee"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("profit").
			StructTag(`json:"profit"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("tx_time").
			StructTag(`json:"txTime"`),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
	}
}

func (MevInfo) Edges() []ent.Edge {
	return nil
}

func (MevInfo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "mev_info"},
	}
}
