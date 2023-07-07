package schema

import (
	"database/sql"
	"entgo.io/ent"
	"fmt"
	"github.com/lib/pq"
	"log"
	"time"
)

type TransactionInfo struct {
	ID          int64     `db:"id"`
	TxHash      string    `db:"tx_hash"`
	BlockNumber int64     `db:"block_number"`
	Network     string    `db:"network"`
	Bundler     string    `db:"bundler"`
	EntryPoint  string    `db:"entry_point"`
	UserOpsNum  int64     `db:"user_ops_num"`
	TxValue     float64   `db:"tx_value"`
	Fee         float64   `db:"fee"`
	GasPrice    string    `db:"gas_price"`
	GasLimit    int64     `db:"gas_limit"`
	Status      int       `db:"status"`
	TxTime      int64     `db:"tx_time"`
	Beneficiary string    `db:"beneficiary"`
	CreateTime  time.Time `db:"create_time"`
	ent.Schema
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

	stmt, err := tx.Prepare(pq.CopyIn("transaction_info", "tx_hash", "block_number", "network", "bundler", "entry_point", "user_ops_num", "fee", "tx_value", "gas_price", "gas_limit", "status", "tx_time", "beneficiary"))
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, t := range transactions {
		_, err = stmt.Exec(t.TxHash, t.BlockNumber, t.Network, t.Bundler, t.EntryPoint, t.UserOpsNum, t.Fee, t.TxValue, t.GasPrice, t.GasLimit, t.Status, t.TxTime, t.Beneficiary)
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

	query := `SELECT id, tx_hash, block_number, network, bundler, entry_point, user_ops_num, fee, tx_value, gas_price, gas_limit, status, tx_time, beneficiary, create_time FROM transaction_info WHERE tx_hash = $1`
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
