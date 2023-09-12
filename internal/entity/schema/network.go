package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

// Network holds the schema definition for the Network entity.
type Network struct {
	ent.Schema
}

// Fields of the Network.
func (Network) Fields() []ent.Field {
	return []ent.Field{
		//field.Int64("id").StructTag(`json:"id"`),
		field.String("id").StorageKey("network").Unique().MaxLen(127).StructTag(`json:"network"`),
		field.Int64("chain_id").StructTag(`json:"chainId"`),
		field.String("name").MaxLen(255).StructTag(`json:"name"`),
		field.String("chain_name").MaxLen(255).StructTag(`json:"chainName"`),
		field.String("http_rpc").StructTag(`json:"httpRpc"`),
		field.Bool("is_testnet").StructTag(`json:"isTestnet"`),
		field.Time("create_time").Default(time.Now).StructTag(`json:"createTime"`).Immutable(),
		field.Time("update_time").Optional().Default(time.Now).UpdateDefault(time.Now).StructTag(`json:"updateTime"`),
		field.Time("delete_time").Optional().StructTag(`json:"deleteTime"`),
		field.String("scan").StructTag(`json:"scan"`),
		field.String("scan_tx").StructTag(`json:"scanTx"`),
		field.String("scan_block").StructTag(`json:"scanBlock"`),
		field.String("scan_address").StructTag(`json:"scanAddress"`),
	}
}

// Edges of the Network.
func (Network) Edges() []ent.Edge {
	return nil
}

func (Network) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "network"},
	}
}
