package schema

import (
	"entgo.io/ent"
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
		field.String("id").StorageKey("uuid").StructTag(`json:"uuid"`),
		field.Time("time").StructTag(`json:"time"`),
		field.String("user_ops_hash").StructTag(`json:"userOpsHash"`),
		field.String("tx_hash").StructTag(`json:"txHash"`),
		field.Int64("block_number").StructTag(`json:"blockNumber"`),
		field.String("network").StructTag(`json:"network"`),
		field.String("sender").StructTag(`json:"sender"`),
		field.String("target").StructTag(`json:"target"`),
		field.Int64("tx_value").StructTag(`json:"txValue"`).GoType(decimal.Decimal{}),
		field.String("source").StructTag(`json:"source"`),
		field.String("calldata").StructTag(`json:"calldata"`),
		field.Int64("tx_time").StructTag(`json:"txTime"`),
		field.Time("create_time").StructTag(`json:"createTime"`),
		field.Time("update_time").StructTag(`json:"updateTime"`),
		field.Int64("aa_index").StructTag(`json:"aaIndex"`),
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
