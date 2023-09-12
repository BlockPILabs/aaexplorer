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
	TxHash *string
}

func (dao *transactionDao) Pages(ctx context.Context, tx *ent.Client, page vo.PaginationRequest, condition TransactionCondition) (a []*ent.TransactionDecode, count int, err error) {
	query := tx.TransactionDecode.Query()

	if condition.TxHash != nil {
		query = query.Where(transactiondecode.ID(*condition.TxHash))
	}

	if page.Sort > 0 {
		query = query.Order(dao.orderPage(ctx, transactiondecode.Columns, page))
	}

	count = query.CountX(ctx)
	if count < 1 || page.GetOffset() > count {
		return
	}

	query = query.Limit(page.GetPerPage()).Offset(page.GetOffset())
	a, err = query.All(ctx)
	return
}
