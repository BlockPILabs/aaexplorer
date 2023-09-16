package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
)

type TransactionDecode struct {
	ent.Schema
}

func (TransactionDecode) Fields() []ent.Field {
	return []ent.Field{
		//time                     timestamp with time zone not null,
		field.Time("time").StructTag(`json:"time"`),
		//create_time              timestamp with time zone,
		field.Time("create_time").StructTag(`json:"createTime"`),
		//hash                     text,
		field.String("id").StorageKey("hash").StructTag(`json:"hash"`),
		//block_hash               text,
		field.String("block_hash").StructTag(`json:"blockHash"`),
		//block_number             numeric,
		field.Int64("block_number").StructTag(`json:"blockNumber"`),
		//nonce                    numeric,
		field.Int64("nonce").StructTag(`json:"nonce"`).GoType(decimal.Decimal{}),
		//transaction_index        numeric,
		field.Int64("transaction_index").StructTag(`json:"transactionIndex"`).GoType(decimal.Decimal{}),
		//from_addr                text,
		field.String("from_addr").StructTag(`json:"from_addr"`),
		//to_addr                  text,
		field.String("to_addr").StructTag(`json:"to_addr"`),
		//value                    numeric,
		field.Int64("value").StructTag(`json:"value"`).GoType(decimal.Decimal{}),
		//gas_price                numeric,
		field.Int64("gas_price").StructTag(`json:"gasPrice"`).GoType(decimal.Decimal{}),
		//gas                      numeric,
		field.Int64("gas").StructTag(`json:"gas"`).GoType(decimal.Decimal{}),
		//input                    text,
		field.String("input").StructTag(`json:"input"`),
		//r                        text,
		field.String("r").StructTag(`json:"r"`),
		//s                        text,
		field.String("s").StructTag(`json:"s"`),
		//v                        bigint,
		field.Int64("v").StructTag(`json:"v"`).GoType(decimal.Decimal{}),
		//chain_id                 bigint,
		field.Int64("chain_id").StructTag(`json:"chainId"`),
		//type                     text,
		field.String("type").StructTag(`json:"type"`),
		//max_fee_per_gas          numeric,
		field.Int64("max_fee_per_gas").StructTag(`json:"maxFeePerGas"`).GoType(decimal.Decimal{}).Nillable(),
		//max_priority_fee_per_gas numeric,
		field.Int64("max_priority_fee_per_gas").StructTag(`json:"maxPriorityFeePerGas"`).GoType(decimal.Decimal{}).Nillable(),
		//access_list              jsonb
		field.Other("access_list", &pgtype.JSONB{}).StructTag(`json:"accessList"`).SchemaType(map[string]string{
			dialect.Postgres: "jsonb",
		}),
		field.String("method"),
	}
}

func (TransactionDecode) Edges() []ent.Edge {
	return nil
}

func (TransactionDecode) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "transaction_decode"},
	}
}
