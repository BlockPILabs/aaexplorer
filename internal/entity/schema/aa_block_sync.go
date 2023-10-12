package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type AaBlockSync struct {
	ent.Schema
}

func (AaBlockSync) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").StorageKey("block_num"),
		field.Bool("block_scanned").Optional().Nillable(),
		field.Bool("tx_scanned").Optional().Nillable(),
		field.Bool("txr_scanned").Optional().Nillable(),
		field.Bool("scanned").Optional().Nillable(),
		field.Time("create_time"),
		field.Time("update_time"),
		field.Int64("scan_count"),
	}
}

func (AaBlockSync) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "aa_block_sync"},
	}
}
