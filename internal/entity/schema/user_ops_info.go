package schema

import (
	"database/sql"
	"entgo.io/ent"
	"github.com/lib/pq"
	"log"
	"time"
)

type UserOpsInfo struct {
	ID                   int64     `db:"id"`
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
	Nonce                int64     `db:"nonce"`
	CallGasLimit         int64     `db:"call_gas_limit"`
	PreVerificationGas   int64     `db:"pre_verification_gas"`
	VerificationGasLimit int64     `db:"verification_gas_limit"`
	MaxFeePerGas         int64     `db:"max_fee_per_gas"`
	MaxPriorityFeePerGas int64     `db:"max_priority_fee_per_gas"`
	TxTime               int64     `db:"tx_time"`
	InitCode             string    `db:"init_code"`
	Status               int       `db:"status"`
	Source               string    `db:"source"`
	ActualGasCost        int64     `db:"actual_gas_cost"`
	ActualGasUsed        int64     `db:"actual_gas_used"`
	CreateTime           time.Time `db:"create_time"`
	ent.Schema
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
		"max_fee_per_gas", "max_priority_fee_per_gas", "tx_time", "init_code", "status", "source", "actual_gas_cost", "actual_gas_used",
		"paymaster_and_data", "signature"))
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, u := range userOpsInfo {
		_, err = stmt.Exec(
			u.UserOperationHash, u.TxHash, u.BlockNumber, u.Network, u.Sender, u.Target, u.Fee, u.TxValue, u.Bundler, u.EntryPoint,
			u.Factory, u.Paymaster, u.Calldata, u.Nonce, u.CallGasLimit, u.PreVerificationGas, u.VerificationGasLimit,
			u.MaxFeePerGas, u.MaxPriorityFeePerGas, u.TxTime, u.InitCode, u.Status, u.Source, u.ActualGasCost, u.ActualGasUsed,
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
