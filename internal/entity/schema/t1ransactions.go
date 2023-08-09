package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"time"
)

type T1ransactions struct {
	Time                 time.Time
	BlockNum             int64
	CreateTime           time.Time
	To                   string
	Gas                  string
	From                 string
	Hash                 string
	Type                 string
	Input                string
	Nonce                string
	Value                string
	ChainID              string
	GasPrice             string
	BlockHash            string
	AccessList           string
	BlockNumber          string
	MaxFeePerGas         string
	TransactionIndex     string
	MaxPriorityFeePerGas string
	ent.Schema
}

func (T1ransactions) Fields() []ent.Field {
	return []ent.Field{
		field.Time("time").
			StructTag(`json:"time"`),
		field.Int64("block_num").
			StructTag(`json:"blockNum"`),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
		field.String("to").
			MaxLen(255).
			StructTag(`json:"to"`),
		field.String("gas").
			MaxLen(255).
			StructTag(`json:"gas"`),
		field.String("from").
			MaxLen(255).
			StructTag(`json:"from"`),
		field.String("hash").
			MaxLen(255).
			StructTag(`json:"hash"`),
		field.String("type").
			MaxLen(255).
			StructTag(`json:"type"`),
		field.String("input").
			MaxLen(255).
			StructTag(`json:"input"`),
		field.String("nonce").
			MaxLen(255).
			StructTag(`json:"nonce"`),
		field.String("value").
			MaxLen(255).
			StructTag(`json:"value"`),
		field.String("chain_id").
			MaxLen(255).
			StructTag(`json:"chainId"`),
		field.String("gas_price").
			MaxLen(255).
			StructTag(`json:"gasPrice"`),
		field.String("block_hash").
			MaxLen(255).
			StructTag(`json:"blockHash"`),
		field.String("access_list").
			MaxLen(255).
			StructTag(`json:"accessList"`),
		field.String("block_number").
			MaxLen(255).
			StructTag(`json:"blockNumber"`),
		field.String("max_fee_per_gas").
			MaxLen(255).
			StructTag(`json:"maxFeePerGas"`),
		field.String("transaction_index").
			MaxLen(255).
			StructTag(`json:"transactionIndex"`),
		field.String("max_priority_fee_per_gas").
			MaxLen(255).
			StructTag(`json:"maxPriorityFeePerGas"`),
	}
}

func (T1ransactions) Edges() []ent.Edge {
	return nil
}

func (T1ransactions) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "t1ransactions"},
	}
}
