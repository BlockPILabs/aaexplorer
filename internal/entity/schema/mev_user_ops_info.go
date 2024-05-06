package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
)

type MevUserOpsInfo struct {
	ent.Schema
}

func (MevUserOpsInfo) Fields() []ent.Field {
	return []ent.Field{
		//time                     timestamp with time zone not null,
		field.Time("time").StructTag(`json:"time"`),
		//create_time              timestamp with time zone,
		field.Time("create_time").StructTag(`json:"createTime"`),
		//hash                     text,
		field.String("id").StorageKey("user_ops_hash").StructTag(`json:"userOpsHash"`),
		field.String("tx_hash").StructTag(`json:"txHash"`),
		//block_hash               text,
		field.String("block_hash").StructTag(`json:"blockHash"`).Optional(),
		//block_number             numeric,
		field.Int64("block_number").StructTag(`json:"blockNumber"`),
		//transaction_index        numeric,
		field.Int64("transaction_index").StructTag(`json:"transactionIndex"`).GoType(decimal.Decimal{}),
		//from_addr                text,
		field.String("from_addr").StructTag(`json:"fromAddr"`),
		//to_addr                  text,
		field.String("to_addr").StructTag(`json:"toAddr"`),
		//value                    numeric,
		field.Int64("value").StructTag(`json:"value"`).GoType(decimal.Decimal{}),
		//gas_price                numeric,
		field.Int64("gas_price").StructTag(`json:"gasPrice"`).GoType(decimal.Decimal{}),
		//gas                      numeric,
		field.Int64("gas").StructTag(`json:"gas"`).GoType(decimal.Decimal{}),
		//chain_id                 bigint,
		field.String("sender").StructTag(`json:"sender"`),

		field.String("target").StructTag(`json:"target"`),
		field.String("victim_sender").StructTag(`json:"victimSender"`),
		//type                     text,
		field.Int64("nonce").StructTag(`json:"nonce"`),
		field.String("bundler").StructTag(`json:"bundler"`),
		field.String("paymaster").StructTag(`json:"paymaster"`),
	}
}

func (MevUserOpsInfo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "mev_user_ops_info_transaction"},
	}
}
