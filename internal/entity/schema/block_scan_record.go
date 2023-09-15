package schema

import (
	"database/sql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
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
	ent.Schema
}

func (BlockScanRecord) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive().
			Unique().
			StructTag(`json:"id"`),
		field.String("network").
			MaxLen(255).
			StructTag(`json:"network"`),
		field.Int64("last_block_number").
			StructTag(`json:"lastBlockNumber"`),
		field.Time("last_scan_time").
			StructTag(`json:"lastScanTime"`),
		field.Time("create_time").
			Default(time.Now).
			StructTag(`json:"createTime"`).
			Immutable(),
		field.Time("update_time").
			Default(time.Now).
			UpdateDefault(time.Now).
			StructTag(`json:"updateTime"`).
			Immutable(),
	}
}

func (BlockScanRecord) Edges() []ent.Edge {
	return nil
}

func (BlockScanRecord) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "block_scan_record"},
	}
}

func InsertBlockScanRecord(record *BlockScanRecord) error {
	connStr := "user=postgres password=root dbname=postgres host=127.0.0.1 port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
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
		log.Println(err)
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
