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

type UserOpsCalldata struct {
	ent.Schema
}

func (UserOpsCalldata) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Comment("Unique identifier").
			StructTag(`json:"id"`),
		field.String("user_ops_hash").
			MaxLen(255).
			Optional().
			Comment("User ops hash").
			StructTag(`json:"userOpsHash"`),
		field.String("tx_hash").
			MaxLen(255).
			Optional().
			Comment("Transaction hash").
			StructTag(`json:"txHash"`),
		field.Int64("block_number").
			Optional().
			Comment("Block number").
			StructTag(`json:"blockNumber"`),
		field.String("network").
			MaxLen(255).
			Optional().
			Comment("Network").
			StructTag(`json:"network"`),
		field.String("sender").
			MaxLen(255).
			Optional().
			Comment("Sender address").
			StructTag(`json:"sender"`),
		field.String("target").
			MaxLen(255).
			Optional().
			Comment("Target address").
			StructTag(`json:"target"`),
		field.Float("tx_value").
			Comment("Transaction value").
			StructTag(`json:"txValue"`).
			GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.String("source").
			MaxLen(255).
			Optional().
			Comment("Source").
			StructTag(`json:"source"`),
		field.String("calldata").
			MaxLen(255).
			Optional().
			Comment("Calldata").
			StructTag(`json:"calldata"`),
		field.Int64("tx_time").
			Optional().
			Comment("Transaction time").
			StructTag(`json:"txTime"`),
		field.Time("create_time").
			Default(time.Now).
			Immutable().
			Comment("Creation time").
			StructTag(`json:"createTime"`),
	}
}

func (UserOpsCalldata) Edges() []ent.Edge {
	return nil
}

func (UserOpsCalldata) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "user_ops_calldata"},
	}
}
