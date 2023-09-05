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

type BlockDataDecode struct {
	ent.Schema
}

func (BlockDataDecode) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").StorageKey("number").StructTag(`json:"number"`).GoType(decimal.Decimal{}),
		field.Time("time").StructTag(`json:"time"`),
		field.Time("create_time").StructTag(`json:"createTime"`),
		field.String("hash").NotEmpty().MaxLen(255).StructTag(`json:"hash"`),
		field.String("parent_hash").StructTag(`json:"parentHash"`),
		field.Int64("nonce").StructTag(`json:"nonce"`).GoType(decimal.Decimal{}),
		field.String("sha3_uncles").StructTag(`json:"sha3Uncles"`),
		field.String("logs_bloom").StructTag(`json:"logsBloom"`),
		field.String("transactions_root").StructTag(`json:"transactionsRoot"`),
		field.String("state_root").StructTag(`json:"stateRoot"`),
		field.String("receipts_root").StructTag(`json:"receiptsRoot"`),
		field.String("miner").StructTag(`json:"miner"`),
		field.String("mix_hash").StructTag(`json:"mixHash"`),
		field.Int64("difficulty").StructTag(`json:"difficulty"`).GoType(decimal.Decimal{}),
		field.Int64("total_difficulty").StructTag(`json:"totalDifficulty"`).GoType(decimal.Decimal{}),
		field.String("extra_data").StructTag(`json:"extraData"`),
		field.Int64("size").StructTag(`json:"size"`).GoType(decimal.Decimal{}),
		field.Int64("gas_limit").StructTag(`json:"gasLimit"`).GoType(decimal.Decimal{}),
		field.Int64("gas_used").StructTag(`json:"gasUsed"`).GoType(decimal.Decimal{}),
		field.Int64("timestamp").StructTag(`json:"timestamp"`).GoType(decimal.Decimal{}),
		field.Int64("transaction_count").StructTag(`json:"transactionCount"`).GoType(decimal.Decimal{}),
		field.Other("uncles", &pgtype.TextArray{}).StructTag(`json:"uncles"`).SchemaType(map[string]string{
			dialect.Postgres: "text[]",
		}),
		field.Int64("base_fee_per_gas").StructTag(`json:"baseFeePerGas"`).GoType(decimal.Decimal{}),
	}
}

func (BlockDataDecode) Edges() []ent.Edge {
	return nil
}

func (BlockDataDecode) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "block_data_decode"},
	}
}
