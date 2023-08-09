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

type AAUserOpsCalldata struct {
	Time        time.Time
	Uuid        string
	UserOpsHash string
	TxHash      string
	BlockNumber int64
	Network     string
	Sender      string
	Target      string
	TxValue     decimal.Decimal
	Source      string
	Calldata    string
	TxTime      int64
	ent.Schema
}

func (AAUserOpsCalldata) Fields() []ent.Field {
	return []ent.Field{
		field.Time("time").
			StructTag(`json:"time"`),
		field.String("uuid").
			MaxLen(255).
			StructTag(`json:"uuid"`),
		field.String("user_ops_hash").
			MaxLen(255).
			StructTag(`json:"userOpsHash"`),
		field.String("tx_hash").
			MaxLen(255).
			StructTag(`json:"txHash"`),
		field.Int64("block_number").
			StructTag(`json:"blockNumber"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.String("sender").
			MaxLen(255).
			StructTag(`json:"sender"`),
		field.String("target").
			MaxLen(255).
			StructTag(`json:"target"`),
		field.Int64("tx_value").
			StructTag(`json:"txValue"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.String("source").
			MaxLen(255).
			StructTag(`json:"source"`),
		field.String("calldata").
			StructTag(`json:"calldata"`),
		field.Int64("tx_time").
			StructTag(`json:"txTime"`),
	}
}

func (AAUserOpsCalldata) Edges() []ent.Edge {
	return nil
}

func (AAUserOpsCalldata) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "aa_user_ops_calldata"},
	}
}
