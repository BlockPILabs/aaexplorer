package schema

import (
	"database/sql"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

type UserOpsInfo struct {
	ID                   int64           `db:"id"`
	UserOperationHash    string          `db:"user_operation_hash"`
	TxHash               string          `db:"tx_hash"`
	BlockNumber          int64           `db:"block_number"`
	Network              string          `db:"network"`
	Sender               string          `db:"sender"`
	Target               string          `db:"target"`
	TxValue              decimal.Decimal `db:"tx_value"`
	Fee                  decimal.Decimal `db:"fee"`
	Bundler              string          `db:"bundler"`
	EntryPoint           string          `db:"entry_point"`
	Factory              string          `db:"factory"`
	Paymaster            string          `db:"paymaster"`
	PaymasterAndData     string          `db:"paymaster_and_data"`
	Signature            string          `db:"signature"`
	Calldata             string          `db:"calldata"`
	Nonce                int64           `db:"nonce"`
	CallGasLimit         int64           `db:"call_gas_limit"`
	PreVerificationGas   int64           `db:"pre_verification_gas"`
	VerificationGasLimit int64           `db:"verification_gas_limit"`
	MaxFeePerGas         int64           `db:"max_fee_per_gas"`
	MaxPriorityFeePerGas int64           `db:"max_priority_fee_per_gas"`
	TxTime               int64           `db:"tx_time"`
	TxTimeFormat         string          `db:"tx_time_format"`
	InitCode             string          `db:"init_code"`
	Status               int             `db:"status"`
	Source               string          `db:"source"`
	ActualGasCost        int64           `db:"actual_gas_cost"`
	ActualGasUsed        int64           `db:"actual_gas_used"`
	CreateTime           time.Time       `db:"create_time"`
	ent.Schema
}

// Fields of the UserOpsInfo.
func (UserOpsInfo) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			StructTag(`json:"id"`),
		field.String("user_operation_hash").
			StructTag(`json:"userOperationHash"`),
		field.String("tx_hash").
			StructTag(`json:"txHash"`),
		field.Int64("block_number").
			StructTag(`json:"blockNumber"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.String("sender").
			StructTag(`json:"sender"`),
		field.String("target").
			StructTag(`json:"target"`),
		field.Int64("tx_value").
			StructTag(`json:"txValue"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.Int64("fee").
			StructTag(`json:"fee"`).GoType(decimal.Zero).SchemaType(map[string]string{dialect.Postgres: "numeric(50, 20)"}),
		field.String("bundler").
			StructTag(`json:"bundler"`),
		field.String("entry_point").
			StructTag(`json:"entryPoint"`),
		field.String("factory").
			StructTag(`json:"factory"`),
		field.String("paymaster").
			StructTag(`json:"paymaster"`),
		field.String("paymaster_and_data").
			StructTag(`json:"paymasterAndData"`),
		field.String("signature").
			StructTag(`json:"signature"`),
		field.String("calldata").
			StructTag(`json:"calldata"`),
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
			StructTag(`json:"source"`),
		field.Int64("actual_gas_cost").
			StructTag(`json:"actualGasCost"`),
		field.Int64("actual_gas_used").
			StructTag(`json:"actualGasUsed"`),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
	}
}

func (UserOpsInfo) Edges() []ent.Edge {
	return nil
}

func (UserOpsInfo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "user_ops_info"},
	}
}

func BulkInsertUserOpsInfo(userOpsInfo []UserOpsInfo) error {
	connStr := "user=postgres password=root dbname=postgres host=127.0.0.1 port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}

	stmt, err := tx.Prepare(pq.CopyIn("user_ops_info",
		"user_operation_hash", "tx_hash", "block_number", "network", "sender", "target", "fee", "tx_value", "bundler", "entry_point",
		"factory", "paymaster", "calldata", "nonce", "call_gas_limit", "pre_verification_gas", "verification_gas_limit",
		"max_fee_per_gas", "max_priority_fee_per_gas", "tx_time", "tx_time_format", "init_code", "status", "source", "actual_gas_cost", "actual_gas_used",
		"paymaster_and_data", "signature"))
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, u := range userOpsInfo {
		_, err = stmt.Exec(
			u.UserOperationHash, u.TxHash, u.BlockNumber, u.Network, u.Sender, u.Target, u.Fee, u.TxValue, u.Bundler, u.EntryPoint,
			u.Factory, u.Paymaster, u.Calldata, u.Nonce, u.CallGasLimit, u.PreVerificationGas, u.VerificationGasLimit,
			u.MaxFeePerGas, u.MaxPriorityFeePerGas, u.TxTime, u.TxTimeFormat, u.InitCode, u.Status, u.Source, u.ActualGasCost, u.ActualGasUsed,
			u.PaymasterAndData, u.Signature,
		)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
