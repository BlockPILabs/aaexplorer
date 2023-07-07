package schema

import (
	"database/sql"
	"log"
	"time"
)

type BlockScanRecord struct {
	ID              int64     `db:"id"`
	Network         string    `db:"network"`
	LastBlockNumber int64     `db:"last_block_number"`
	LastScanTime    time.Time `db:"last_scan_time"`
	CreateTime      time.Time `db:"create_time"`
	UpdateTime      time.Time `db:"update_time"`
}

func InsertBlockScanRecord(record *BlockScanRecord) error {
	connStr := "user=postgres password=root dbname=postgres host=127.0.0.1 port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT INTO block_scan_record (network, last_block_number, last_scan_time)
		VALUES ($1, $2, $3)`,
		record.Network, record.LastBlockNumber, record.LastScanTime)
	if err != nil {
		return err
	}

	return nil
}

func UpdateBlockScanRecordByID(id int64, record *BlockScanRecord) error {
	connStr := "user=postgres password=root dbname=postgres host=127.0.0.1 port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		UPDATE block_scan_record
		SET last_block_number = $1, last_scan_time = $2, update_time = $3
		WHERE id = $4`,
		record.LastBlockNumber, record.LastScanTime, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}

func GetBlockScanRecordsByNetwork(network string) (*BlockScanRecord, error) {
	connStr := "user=postgres password=root dbname=postgres host=127.0.0.1 port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT id, network, last_block_number, last_scan_time, create_time
		FROM block_scan_record
		WHERE network = $1`

	rows, err := db.Query(query, network)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []BlockScanRecord

	for rows.Next() {
		var record BlockScanRecord
		err := rows.Scan(
			&record.ID,
			&record.Network,
			&record.LastBlockNumber,
			&record.LastScanTime,
			&record.CreateTime,
		)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}

	return &records[0], nil
}
