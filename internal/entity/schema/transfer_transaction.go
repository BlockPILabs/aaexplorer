package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
)

type TransferTransaction struct {
	ent.Schema
}

func (TransferTransaction) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").StorageKey("id").StructTag(`json:"id"`),
		//time                     timestamp with time zone not null,
		field.Time("time").StructTag(`json:"time"`),
		//create_time              timestamp with time zone,
		field.Time("create_time").StructTag(`json:"createTime"`),
		//hash                     text,
		field.String("tx_hash").StructTag(`json:"txHash"`),
		//block_hash               text,
		field.String("block_hash").StructTag(`json:"blockHash"`).Optional(),
		//block_number             numeric,
		field.Int64("block_number").StructTag(`json:"blockNumber"`),
		//transaction_index        numeric,
		field.Int64("transaction_index").StructTag(`json:"transactionIndex"`),
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
		field.Int64("transfer_value").StructTag(`json:"transferValue"`).GoType(decimal.Decimal{}),
		field.String("token_symbol").StructTag(`json:"tokenSymbol"`),
		field.String("token_address").StructTag(`json:"tokenAddress"`),
		field.String("token_url").StructTag(`json:"tokenUrl"`),
	}
}

func (TransferTransaction) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "transfer_transaction"},
	}
}
