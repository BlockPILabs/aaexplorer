package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
)

type MevTransaction struct {
	ent.Schema
}

func (MevTransaction) Fields() []ent.Field {
	return []ent.Field{
		//time                     timestamp with time zone not null,
		field.Time("time").StructTag(`json:"time"`),
		//create_time              timestamp with time zone,
		field.Time("create_time").StructTag(`json:"createTime"`),
		//hash                     text,
		field.String("id").StorageKey("tx_hash").StructTag(`json:"txHash"`),
		//block_hash               text,
		field.String("block_hash").StructTag(`json:"blockHash"`).Optional(),
		//block_number             numeric,
		field.Int64("block_number").StructTag(`json:"blockNumber"`),
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
		field.String("mev_type").StructTag(`json:"mevType"`),
		//chain_id                 bigint,
		field.String("victim").StructTag(`json:"victim"`),
		field.String("victim_type").StructTag(`json:"victimType"`),
		field.String("victim_tx_hash").StructTag(`json:"victimTxHash"`),
		field.Int64("victim_block_number").StructTag(`json:"victimBlockNumber"`),
		field.String("victim_from_addr").StructTag(`json:"victimFromAddr"`),
		field.String("victim_to_addr").StructTag(`json:"victimToAddr"`),
		//type                     text,
		field.String("attacker").StructTag(`json:"attacker"`),
		field.Int64("bundler_loss").StructTag(`json:"bundlerLoss"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("bundler_loss_usd").StructTag(`json:"bundlerLossUsd"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("mev_profit").StructTag(`json:"mevProfit"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("mev_profit_usd").StructTag(`json:"mevProfitUsd"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
	}
}

func (MevTransaction) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "mev_transaction"},
	}
}
