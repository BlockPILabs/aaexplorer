package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactiondecode"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type transactionDao struct {
	baseDao
}

var TransactionDao = &transactionDao{}

type TransactionCondition struct {
	TxHash  *string
	Address *string
}

func (dao *transactionDao) Pages(ctx context.Context, tx *ent.Client, page vo.PaginationRequest, condition TransactionCondition) (a []*ent.TransactionDecode, count int, err error) {
	query := tx.TransactionDecode.Query()

	if condition.TxHash != nil {
		query = query.Where(transactiondecode.ID(*condition.TxHash))
	}

	count = query.CountX(ctx)
	if count < 1 || page.GetOffset() > count {
		return
	}

	if page.Sort > 0 {
		query = query.Order(dao.orderPage(ctx, transactiondecode.Columns, page))
	} else {
		query = query.Order(transactiondecode.ByTime(sql.OrderDesc()))
	}

	query = query.Limit(page.GetPerPage()).Offset(page.GetOffset())
	a, err = query.All(ctx)
	return
}

/*type TxPages struct {
	ID string `json:"hash"`
	// Time holds the value of the "time" field.
	Time time.Time `json:"time"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"createTime"`
	// BlockHash holds the value of the "block_hash" field.
	BlockHash string `json:"blockHash"`
	// BlockNumber holds the value of the "block_number" field.
	BlockNumber int64 `json:"blockNumber"`
	// Nonce holds the value of the "nonce" field.
	Nonce decimal.Decimal `json:"nonce"`
	// TransactionIndex holds the value of the "transaction_index" field.
	TransactionIndex decimal.Decimal `json:"transactionIndex"`
	// FromAddr holds the value of the "from_addr" field.
	FromAddr string `json:"from_addr"`
	// ToAddr holds the value of the "to_addr" field.
	ToAddr string `json:"to_addr"`
	// Value holds the value of the "value" field.
	Value decimal.Decimal `json:"value"`
	// GasPrice holds the value of the "gas_price" field.
	GasPrice decimal.Decimal `json:"gasPrice"`
	// Gas holds the value of the "gas" field.
	Gas decimal.Decimal `json:"gas"`
	// Input holds the value of the "input" field.
	Input string `json:"input"`
	// R holds the value of the "r" field.
	R string `json:"r"`
	// S holds the value of the "s" field.
	S string `json:"s"`
	// V holds the value of the "v" field.
	V decimal.Decimal `json:"v"`
	// ChainID holds the value of the "chain_id" field.
	ChainID int64 `json:"chainId"`
	// Type holds the value of the "type" field.
	Type string `json:"type"`
	// MaxFeePerGas holds the value of the "max_fee_per_gas" field.
	MaxFeePerGas *decimal.Decimal `json:"maxFeePerGas"`
	// MaxPriorityFeePerGas holds the value of the "max_priority_fee_per_gas" field.
	MaxPriorityFeePerGas *decimal.Decimal `json:"maxPriorityFeePerGas"`
	// AccessList holds the value of the "access_list" field.
	AccessList *pgtype.JSONB `json:"accessList"`
	// Method holds the value of the "method" field.
	Method string `json:"method,omitempty"`

	UseropCount int64 `json:"useropCount"`
	// IsMev holds the value of the "is_mev" field.
	IsMev bool `json:"isMev"`
	// BundlerProfit holds the value of the "bundler_profit" field.
	BundlerProfit decimal.Decimal `json:"bundlerProfit"`
	// BundlerProfitUsd holds the value of the "bundler_profit_usd" field.
	BundlerProfitUsd decimal.Decimal `json:"bundler_profit_usd,omitempty"`
}
*/
/*func (dao *transactionDao) PagesWithTxaa(ctx context.Context, tx *ent.Client, page vo.PaginationRequest, condition TransactionCondition) (a []*TxPages, count int, err error) {
	query := tx.AaTransactionInfo.Query().
		Modify(func(s *sql.Selector) {
			d := sql.Dialect(s.Dialect())
			s.LeftJoin(d.Table(transactiondecode.Table).As(transactiondecode.Table)).OnP(sql.ColumnsEQ(d.Table(transactiondecode.Table).C(transactiondecode.FieldID), d.Table(aatransactioninfo.Table).C(aatransactioninfo.FieldID)))

		})
	if condition.TxHash != nil && len(*(condition.TxHash)) > 0 {
		query = query.Modify(func(s *sql.Selector) {
			s.Where(sql.EQ(transactiondecode.FieldID, *condition.TxHash))
		})
	}
	if condition.Address != nil && len(*(condition.Address)) > 0 {

		query = query.Modify(func(s *sql.Selector) {
			s.Where(
				sql.Or(
					sql.EQ(transactiondecode.FieldFromAddr, *condition.Address),
					sql.EQ(transactiondecode.FieldToAddr, *condition.Address),
				),
			)
		})
		//query = query.Where(
		//	transactiondecode.Or(
		//		transactiondecode.FromAddrEQ(*condition.Address),
		//		transactiondecode.ToAddrEQ(*condition.Address),
		//	),
		//)
	}
	//query = query.WithTxaa()

	if page.TotalCount > 0 {
		count = page.TotalCount
	} else {
		count = query.CountX(ctx)
	}

	if count < 1 || page.GetOffset() > count {
		return
	}
	if page.Sort > 0 {

		query = query.Modify(func(s *sql.Selector) {
			dao.orderPage(ctx, transactiondecode.Columns, page)(s)
		})

		//query = query.Order(dao.orderPage(ctx, transactiondecode.Columns, page))
	} else {
		query = query.Modify(func(s *sql.Selector) {
			sql.OrderByField(transactiondecode.FieldTime, sql.OrderDesc())
		})
		//query = query.Order(transactiondecode.ByTime(sql.OrderDesc()))
	}
	//query := query.Limit(page.GetPerPage()).Offset(page.GetOffset())
	query = query.Modify(func(s *sql.Selector) {
		s.Limit(page.GetPerPage())
		s.Offset(page.GetOffset())
	})
	err = query.Scan(ctx, a)
	return
}*/
