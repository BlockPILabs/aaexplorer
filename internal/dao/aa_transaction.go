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
	HASH                     *string          `json:"hash"`
	TIME                     *time.Time       `json:"time"`
	BLOCK_HASH               *string          `json:"block_hash"`
	BLOCK_NUMBER             *int64           `json:"block_number"`
	USEROP_COUNT             *int64           `json:"userop_count"`
	IS_MEV                   *bool            `json:"is_mev"`
	BUNDLER_PROFIT           *decimal.Decimal `json:"bundler_profit"`
	NONCE                    *decimal.Decimal `json:"nonce"`
	TRANSACTION_INDEX        *decimal.Decimal `json:"transaction_index"`
	FROM_ADDR                *string          `json:"from_addr"`
	TO_ADDR                  *string          `json:"to_addr"`
	VALUE                    *decimal.Decimal `json:"value"`
	GAS_PRICE                *decimal.Decimal `json:"gas_price"`
	GAS                      *decimal.Decimal `json:"gas"`
	INPUT                    *string          `json:"input"`
	R                        *string          `json:"r"`
	S                        *string          `json:"s"`
	V                        *decimal.Decimal `json:"v"`
	CHAIN_ID                 *int64           `json:"chain_id"`
	TYPE                     *string          `json:"type"`
	MAX_FEE_PER_GAS          *decimal.Decimal `json:"max_fee_per_gas"`
	MAX_PRIORITY_FEE_PER_GAS *decimal.Decimal `json:"max_priority_fee_per_gas"`
	ACCESS_LIST              *pgtype.JSONB    `json:"access_list"`
	Create_Time              *time.Time       `json:"create_time"`
	BundlerProfitUsd         decimal.Decimal  `json:"bundler_profit_usd,omitempty"`
}

func (dao *aaTransactionDao) Pages(ctx context.Context, tx *ent.Client, page vo.PaginationRequest, condition AaTransactionCondition) (a []*AaTransactionScan, count int, err error) {
	query := tx.AaTransactionInfo.Query().Modify(func(s *sql.Selector) {
		t := sql.Table(transactiondecode.Table)
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
	})
	err = query.Scan(ctx, &a)
	return
}
