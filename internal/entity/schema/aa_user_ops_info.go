package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
	"time"
)

type AAUserOpsInfo struct {
	Time                 time.Time `db:"time"`
	UserOperationHash    string    `db:"user_operation_hash"`
	TxHash               string    `db:"tx_hash"`
	BlockNumber          int64     `db:"block_number"`
	Network              string    `db:"network"`
	Sender               string    `db:"sender"`
	Target               string    `db:"target"`
	TxValue              float64   `db:"tx_value"`
	Fee                  float64   `db:"fee"`
	Bundler              string    `db:"bundler"`
	EntryPoint           string    `db:"entry_point"`
	Factory              string    `db:"factory"`
	Paymaster            string    `db:"paymaster"`
	PaymasterAndData     string    `db:"paymaster_and_data"`
	Signature            string    `db:"signature"`
	Calldata             string    `db:"calldata"`
	CalldataContract     string    `db:"calldata_contract"`
	Nonce                int64     `db:"nonce"`
	CallGasLimit         int64     `db:"call_gas_limit"`
	PreVerificationGas   int64     `db:"pre_verification_gas"`
	VerificationGasLimit int64     `db:"verification_gas_limit"`
	MaxFeePerGas         int64     `db:"max_fee_per_gas"`
	MaxPriorityFeePerGas int64     `db:"max_priority_fee_per_gas"`
	TxTime               int64     `db:"tx_time"`
	TxTimeFormat         string    `db:"tx_time_format"`
	InitCode             string    `db:"init_code"`
	Status               int       `db:"status"`
	Source               string    `db:"source"`
	ActualGasCost        int64     `db:"actual_gas_cost"`
	ActualGasUsed        int64     `db:"actual_gas_used"`
	ent.Schema
}

func (AAUserOpsInfo) Fields() []ent.Field {
	return []ent.Field{
		field.Time("time").
			StructTag(`json:"time"`),
		field.String("user_operation_hash").
			MaxLen(255).
			StructTag(`json:"userOperationHash"`),
		field.String("tx_hash").
			MaxLen(255).
			StructTag(`json:"txHash"`),
		field.Int64("block_number").
			StructTag(`json:"blockNumber"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.String("sender").
			MaxLen(255).
			StructTag(`json:"sender"`),
		field.String("target").
			MaxLen(255).
			StructTag(`json:"target"`),
		field.Int64("tx_value").
			StructTag(`json:"txValue"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("fee").
			StructTag(`json:"fee"`).GoType(decimal.Zero).
			SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.String("bundler").
			MaxLen(255).
			StructTag(`json:"bundler"`),
		field.String("entry_point").
			MaxLen(255).
			StructTag(`json:"entryPoint"`),
		field.String("factory").
			MaxLen(255).
			StructTag(`json:"factory"`),
		field.String("paymaster").
			MaxLen(255).
			StructTag(`json:"paymaster"`),
		field.String("paymaster_and_data").
			StructTag(`json:"paymasterAndData"`),
		field.String("signature").
			StructTag(`json:"signature"`),
		field.String("calldata").
			StructTag(`json:"calldata"`),
		field.String("calldata_contract").
			StructTag(`json:"calldataContract"`),
		field.Int64("nonce").
			StructTag(`json:"nonce"`),
		field.Int64("call_gas_limit").
			StructTag(`json:"callGasLimit"`),
		field.Int64("pre_verification_gas").
			StructTag(`json:"preVerificationGas"`),
		field.Int64("verification_gas_limit").
			StructTag(`json:"verificationGasLimit"`),
		field.Int64("max_fee_per_gas").
			StructTag(`json:"maxFeePerGas"`),
		field.Int64("max_priority_fee_per_gas").
			StructTag(`json:"maxPriorityFeePerGas"`),
		field.Int64("tx_time").
			StructTag(`json:"txTime"`),
		field.String("tx_time_format").
			StructTag(`json:"txTimeFormat"`),
		field.String("init_code").
			StructTag(`json:"initCode"`),
		field.Int("status").
			StructTag(`json:"status"`),
		field.String("source").
			MaxLen(255).
			StructTag(`json:"source"`),
		field.Int64("actual_gas_cost").
			StructTag(`json:"actualGasCost"`),
		field.Int64("actual_gas_used").
			StructTag(`json:"actualGasUsed"`),
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
