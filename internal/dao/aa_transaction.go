package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aatransactioninfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactiondecode"
	"github.com/BlockPILabs/aa-scan/internal/utils"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
	"time"
)

type aaTransactionDao struct {
	baseDao
}

var AaTransactionDao = &aaTransactionDao{}

type AaTransactionCondition struct {
	TxHashTerm string
	TxHash     *string
	Address    *string
}

func (dao *aaTransactionDao) Pagination(ctx context.Context, tx *ent.Client, page vo.PaginationRequest, condition AaTransactionCondition) (a ent.AaTransactionInfos, count int, err error) {
	query := tx.AaTransactionInfo.Query()

	if len(condition.TxHashTerm) > 0 && utils.IsHexSting(condition.TxHashTerm) {
		if utils.IsHashHex(condition.TxHashTerm) {
			query = query.Where(aatransactioninfo.IDEQ(utils.Fix0x(condition.TxHashTerm)))
		} else {
			query = query.Where(sql.FieldHasPrefix(aatransactioninfo.FieldID, utils.Fix0x(condition.TxHashTerm)))
		}
	}
	if page.TotalCount > 0 {
		count = page.TotalCount
	} else {
		count = query.CountX(ctx)
	}
	if count < 1 || page.GetOffset() > count {
		return
	}

	if page.Sort > 0 {
		query = query.Order(dao.orderPage(ctx, aatransactioninfo.Columns, page))
	}

	query = query.Limit(page.GetPerPage()).Offset(page.GetOffset())
	a, err = query.All(ctx)
	return
}

type AaTransactionScan struct {
	Hash             string          `json:"hash"`
	Time             time.Time       `json:"time"`
	BlockHash        string          `json:"block_hash"`
	BlockNumber      int64           `json:"block_number"`
	UseropCount      int64           `json:"userop_count"`
	IsMev            bool            `json:"is_mev"`
	BundlerProfit    decimal.Decimal `json:"bundler_profit"`
	CreateTime       time.Time       `json:"create_time"`
	BundlerProfitUsd decimal.Decimal `json:"bundler_profit_usd"`

	Nonce                decimal.Decimal  `json:"nonce"`
	TransactionIndex     decimal.Decimal  `json:"transaction_index"`
	FromAddr             string           `json:"from_addr"`
	ToAddr               string           `json:"to_addr"`
	Value                decimal.Decimal  `json:"value"`
	GasPrice             decimal.Decimal  `json:"gas_price"`
	Gas                  decimal.Decimal  `json:"gas"`
	Input                string           `json:"input"`
	R                    string           `json:"r"`
	S                    string           `json:"s"`
	V                    decimal.Decimal  `json:"v"`
	ChainID              int64            `json:"chain_id"`
	Type                 string           `json:"type"`
	MaxFeePerGas         *decimal.Decimal `json:"max_fee_per_gas"`
	MaxPriorityFeePerGas *decimal.Decimal `json:"max_priority_fee_per_gas"`
	AccessList           *pgtype.JSONB    `json:"access_list"`
	Method               string           `json:"method"`
}

func (dao *aaTransactionDao) Pages(ctx context.Context, tx *ent.Client, page vo.PaginationRequest, condition AaTransactionCondition) (a []*AaTransactionScan, count int, err error) {
	query := tx.AaTransactionInfo.Query().Modify(func(s *sql.Selector) {
		t := sql.Table(transactiondecode.Table).As(transactiondecode.Table)
		s.LeftJoin(t).On(s.C(aatransactioninfo.FieldID), t.C(transactiondecode.FieldID))
		if len(condition.TxHashTerm) > 0 && utils.IsHexSting(condition.TxHashTerm) {
			s.Where(sql.HasPrefix(aatransactioninfo.FieldID, utils.Fix0x(condition.TxHashTerm)))
		}
		if condition.TxHash != nil && len(*(condition.TxHash)) > 0 {
			sql.FieldEQ(aatransactioninfo.FieldID, *condition.TxHash)(s)
		}
		if condition.Address != nil && len(*(condition.Address)) > 0 {
			s.Where(sql.Or(
				sql.EQ(transactiondecode.FieldFromAddr, *condition.Address),
				sql.EQ(transactiondecode.FieldToAddr, *condition.Address),
			),
			)

		}
	})

	if page.TotalCount > 0 {
		count = page.TotalCount
	} else {
		count = query.CountX(ctx)
	}
	if count < 1 || page.GetOffset() > count {
		return
	}

	query = query.Modify(func(s *sql.Selector) {
		if page.Sort > 0 {
			dao.orderPage(ctx, aatransactioninfo.Columns, page)(s)
		}
		s.Limit(page.GetPerPage()).Offset(page.GetOffset())
		a := sql.Dialect(s.Dialect()).Table(aatransactioninfo.Table)
		t := sql.Dialect(s.Dialect()).Table(transactiondecode.Table)
		s.Select(a.C("*"), t.C("*"))
	})
	err = query.Scan(ctx, &a)
	return
}
