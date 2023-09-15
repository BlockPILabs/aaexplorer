package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

type AssetChangeTrace struct {
	Id             int64
	TxHash         string
	BlockNumber    int64
	Network        string
	Address        string
	AddressType    int
	LastChangeTime time.Time
	SyncFlag       int
	CreateTime     time.Time
	UpdateTime     time.Time
	ent.Schema
}

func (AssetChangeTrace) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			StructTag(`json:"id"`),
		field.String("tx_hash").
			MaxLen(255).
			StructTag(`json:"txHash"`),
		field.Int64("block_number").
			StructTag(`json:"blockNumber"`),
		field.String("address").
			MaxLen(255).
			StructTag(`json:"address"`),
		field.Int("address_type").
			StructTag(`json:"addressType"`),
		field.Time("last_change_time").
			StructTag(`json:"lastChangeTime"`),
		field.Int("sync_flag").
			StructTag(`json:"syncFlag"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
		field.Time("update_time").
			Default(time.Now).
			StructTag(`json:"updateTime"`),
	}
}

func (AssetChangeTrace) Edges() []ent.Edge {
	return nil
}

func (AssetChangeTrace) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset_change_trace"},
	}
}
