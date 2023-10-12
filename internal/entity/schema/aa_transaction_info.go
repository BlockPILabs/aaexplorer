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

type AaTransactionInfo struct {
	ent.Schema
}

func (AaTransactionInfo) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("hash").StructTag(`json:"hash"`),
		field.Time("time").StructTag(`json:"time"`),
		field.String("block_hash").StructTag(`json:"blockHash"`),
		field.Int64("block_number").StructTag(`json:"blockNumber"`),
		field.Int64("userop_count").StructTag(`json:"useropCount"`),
		field.Bool("is_mev").StructTag(`json:"isMev"`),
		field.Int64("bundler_profit").StructTag(`json:"bundlerProfit"`).GoType(decimal.Decimal{}),
		field.Time("create_time").StructTag(`json:"createTime"`),
		field.Int64("bundler_profit_usd").GoType(decimal.Decimal{}).Optional(),

		field.Int64("nonce").StructTag(`json:"nonce"`).GoType(decimal.Decimal{}),
		field.Int64("transaction_index").StructTag(`json:"transactionIndex"`).GoType(decimal.Decimal{}),
		field.String("from_addr").StructTag(`json:"from_addr"`),
		field.String("to_addr").StructTag(`json:"to_addr"`),
		field.Int64("value").StructTag(`json:"value"`).GoType(decimal.Decimal{}),
		field.Int64("gas_price").StructTag(`json:"gasPrice"`).GoType(decimal.Decimal{}),
		field.Int64("gas").StructTag(`json:"gas"`).GoType(decimal.Decimal{}),
		field.String("input").StructTag(`json:"input"`),
		field.String("r").StructTag(`json:"r"`),
		field.String("s").StructTag(`json:"s"`),
		field.Int64("v").StructTag(`json:"v"`).GoType(decimal.Decimal{}),
		field.Int64("chain_id").StructTag(`json:"chainId"`),
		field.String("type").StructTag(`json:"type"`),
		field.Int64("max_fee_per_gas").StructTag(`json:"maxFeePerGas"`).GoType(decimal.Decimal{}).Nillable(),
		field.Int64("max_priority_fee_per_gas").StructTag(`json:"maxPriorityFeePerGas"`).GoType(decimal.Decimal{}).Nillable(),
		field.Other("access_list", &pgtype.JSONB{}).StructTag(`json:"accessList"`).SchemaType(map[string]string{
			dialect.Postgres: "jsonb",
		}),
		field.String("method"),

		field.String("contract_address").MaxLen(255).StructTag(`json:"contractAddress"`),
		field.Int64("cumulative_gas_used").StructTag(`json:"cumulativeGasUsed"`),
		field.String("effective_gas_price").StructTag(`json:"effective_gas_price"`),
		field.Int64("gas_used").StructTag(`json:"gasUsed"`).GoType(decimal.Decimal{}),
		field.String("logs").MaxLen(255).StructTag(`json:"logs"`),
		field.String("logs_bloom").MaxLen(255).StructTag(`json:"logsBloom"`),
		field.String("status").MaxLen(255).StructTag(`json:"status"`),
	}
}

func (AaTransactionInfo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "aa_transaction_info"},
	}
}

/*func (AaTransactionInfo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("aatx", TransactionDecode.Type).StorageKey().Field("id").Unique(),
	}
}
*/
