package schema

import (
	"database/sql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"fmt"
	"github.com/lib/pq"
	"log"
	"time"
)

type TransactionInfo struct {
	ID           int64     `db:"id"`
	TxHash       string    `db:"tx_hash"`
	BlockNumber  int64     `db:"block_number"`
	Network      string    `db:"network"`
	Bundler      string    `db:"bundler"`
	EntryPoint   string    `db:"entry_point"`
	UserOpsNum   int64     `db:"user_ops_num"`
	TxValue      float64   `db:"tx_value"`
	Fee          float64   `db:"fee"`
	GasPrice     string    `db:"gas_price"`
	GasLimit     int64     `db:"gas_limit"`
	Status       int       `db:"status"`
	TxTime       int64     `db:"tx_time"`
	TxTimeFormat string    `db:"tx_time_format"`
	Beneficiary  string    `db:"beneficiary"`
	CreateTime   time.Time `db:"create_time"`
	ent.Schema
}

// Fields of the TransactionInfo.
func (TransactionInfo) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			StructTag(`json:"id"`),
		field.String("tx_hash").
			MaxLen(255).
			StructTag(`json:"txHash"`),
		field.Int64("block_number").
			StructTag(`json:"blockNumber"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.String("bundler").
			MaxLen(255).
			StructTag(`json:"bundler"`),
		field.String("entry_point").
			MaxLen(255).
			StructTag(`json:"entryPoint"`),
		field.Int64("user_ops_num").
			StructTag(`json:"userOpsNum"`),
		field.Float32("tx_value").
			StructTag(`json:"txValue"`),
		field.Float32("fee").
			StructTag(`json:"fee"`),
		field.String("gas_price").
			StructTag(`json:"gasPrice"`),
		field.Int64("gas_limit").
			StructTag(`json:"gasLimit"`),
		field.Int("status").
			StructTag(`json:"status"`),
		field.Int64("tx_time").
			StructTag(`json:"txTime"`),
		field.String("tx_time_format").
			StructTag(`json:"txTimeFormat"`),
		field.String("beneficiary").
			MaxLen(255).
			StructTag(`json:"beneficiary"`),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`),
	}
}

func (TransactionInfo) Edges() []ent.Edge {
	return nil
}

// Annotations of the TransactionInfo.
func (TransactionInfo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "transaction_info"},
	}
}

func BulkInsertTransactions(transactions []TransactionInfo) error {
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

	stmt, err := tx.Prepare(pq.CopyIn("transaction_info", "tx_hash", "block_number", "network", "bundler", "entry_point", "user_ops_num", "fee", "tx_value", "gas_price", "gas_limit", "status", "tx_time", "tx_time_format", "beneficiary"))
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, t := range transactions {
		_, err = stmt.Exec(t.TxHash, t.BlockNumber, t.Network, t.Bundler, t.EntryPoint, t.UserOpsNum, t.Fee, t.TxValue, t.GasPrice, t.GasLimit, t.Status, t.TxTime, t.TxTimeFormat, t.Beneficiary)
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

func FindTransactionByTxHash(txHash string) (*TransactionInfo, error) {
	db, err := sql.Open("postgres", "your-database-connection-string")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT id, tx_hash, block_number, network, bundler, entry_point, user_ops_num, fee, tx_value, gas_price, gas_limit, status, tx_time, tx_time_format, beneficiary, create_time FROM transaction_info WHERE tx_hash = $1`
	row := db.QueryRow(query, txHash)

	transaction := &TransactionInfo{}

	err = row.Scan(
		&transaction.ID,
		&transaction.TxHash,
		&transaction.BlockNumber,
		&transaction.Network,
		&transaction.Bundler,
		&transaction.EntryPoint,
		&transaction.UserOpsNum,
		&transaction.Fee,
		&transaction.TxValue,
		&transaction.GasPrice,
		&transaction.GasLimit,
		&transaction.Status,
		&transaction.TxTime,
		&transaction.TxTimeFormat,
		&transaction.Beneficiary,
		&transaction.CreateTime,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("")
		}
		return nil, err
	}

	return transaction, nil
}
