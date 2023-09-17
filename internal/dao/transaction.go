package dao

import (
	"context"
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
	}

	query = query.Limit(page.GetPerPage()).Offset(page.GetOffset())
	a, err = query.All(ctx)
	return
}

func (dao *transactionDao) PagesWithTxaa(ctx context.Context, tx *ent.Client, page vo.PaginationRequest, condition TransactionCondition) (a []*ent.TransactionDecode, count int, err error) {

	query := tx.TransactionDecode.Query()

	if condition.TxHash != nil && len(*(condition.TxHash)) > 0 {
		query = query.Where(transactiondecode.IDEQ(*condition.TxHash))
	}
	if condition.Address != nil && len(*(condition.Address)) > 0 {
		query = query.Where(
			transactiondecode.Or(
				transactiondecode.FromAddrEQ(*condition.Address),
				transactiondecode.ToAddrEQ(*condition.Address),
			),
		)
	}
	query = query.WithTxaa()

	if page.TotalCount > 0 {
		count = page.TotalCount
	} else {
		count = query.CountX(ctx)
	}
	if count < 1 || page.GetOffset() > count {
		return
	}
	if page.Sort > 0 {
		query = query.Order(dao.orderPage(ctx, transactiondecode.Columns, page))
	}
	query = query.Limit(page.GetPerPage()).Offset(page.GetOffset())
	a, err = query.All(ctx)
	return
}
