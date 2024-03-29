package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type BlockSync struct {
	ent.Schema
}

func (BlockSync) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").StorageKey("block_num"),
		field.Bool("scanned").Optional().Nillable(),
		field.Time("create_time"),
		field.Time("update_time"),
	}
}

func (BlockSync) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "block_sync"},
	}
}

type TransactionReceiptBlockSync struct {
	ent.Schema
}

func (TransactionReceiptBlockSync) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").StorageKey("block_num"),
		field.Bool("scanned").Optional().Nillable(),
		field.Time("create_time"),
		field.Time("update_time"),
	}
}

func (TransactionReceiptBlockSync) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "transaction_receipt_block_sync"},
	}
}

type TransactionSync struct {
	ent.Schema
}

func (TransactionSync) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").StorageKey("block_num"),
		field.Bool("scanned").Optional().Nillable(),
		field.Time("create_time"),
		field.Time("update_time"),
	}
}

func (TransactionSync) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "transaction_sync"},
	}
}
