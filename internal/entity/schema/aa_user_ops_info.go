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

type AAUserOpsInfo struct {
	ent.Schema
}

func (AAUserOpsInfo) Fields() []ent.Field {
	return []ent.Field{

		field.Time("time").StructTag(`json:"time"`),
		field.String("id").StorageKey("user_operation_hash").StructTag(`json:"userOperationHash"`),
		field.String("tx_hash").StructTag(`json:"txHash"`),
		field.Int64("block_number").StructTag(`json:"blockNumber"`),
		field.String("network").StructTag(`json:"network"`),
		field.String("sender").StructTag(`json:"sender"`),
		field.String("target").StructTag(`json:"target"`),
		field.Other("targets", &pgtype.TextArray{}).StructTag(`json:"targets"`).SchemaType(map[string]string{
			dialect.Postgres: "varchar(127)[]",
		}),
		field.Int("targets_count").StructTag(`json:"targetsCount"`),
		field.Int64("tx_value").StructTag(`json:"txValue"`).GoType(decimal.Decimal{}),
		field.Int64("fee").StructTag(`json:"fee"`).GoType(decimal.Decimal{}),
		field.String("bundler").StructTag(`json:"bundler"`),
		field.String("entry_point").StructTag(`json:"entryPoint"`),
		field.String("factory").StructTag(`json:"factory"`),
		field.String("paymaster").StructTag(`json:"paymaster"`),
		field.String("paymaster_and_data").StructTag(`json:"paymasterAndData"`),
		field.String("signature").StructTag(`json:"signature"`),
		field.String("calldata").StructTag(`json:"calldata"`),
		field.String("calldata_contract").StructTag(`json:"calldataContract"`),
		field.Int64("nonce").StructTag(`json:"nonce"`),
		field.Int64("call_gas_limit").StructTag(`json:"callGasLimit"`),
		field.Int64("pre_verification_gas").StructTag(`json:"preVerificationGas"`),
		field.Int64("verification_gas_limit").StructTag(`json:"verificationGasLimit"`),
		field.Int64("max_fee_per_gas").StructTag(`json:"max_fee_per_gas"`),
		field.Int64("max_priority_fee_per_gas").StructTag(`json:"maxPriorityFeePerGas"`),
		field.Int64("tx_time").StructTag(`json:"txTime"`),
		field.String("init_code").StructTag(`json:"initCode"`),
		field.Int32("status").StructTag(`json:"status"`),
		field.String("source").StructTag(`json:"source"`),
		field.Int64("actual_gas_cost").StructTag(`json:"actualGasCost"`),
		field.Int64("actual_gas_used").StructTag(`json:"actualGasUsed"`),
		field.Time("create_time").StructTag(`json:"createTime"`),
		field.Time("update_time").StructTag(`json:"updateTime"`),
		field.Int64("usd_amount").StructTag(`json:"usdAmount"`).GoType(decimal.Decimal{}).Optional().Nillable(),
		field.Int("aa_index").StructTag(`json:"aaIndex"`),
		field.Int64("fee_usd").StructTag(`json:"feeUsd"`).GoType(decimal.Decimal{}).Optional(),
		field.Int64("tx_value_usd").StructTag(`json:"txValueUsd"`).GoType(decimal.Decimal{}).Optional(),
	}
}

func (AAUserOpsInfo) Edges() []ent.Edge {
	return nil
}

func (AAUserOpsInfo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "aa_user_ops_info"},
	}
}
