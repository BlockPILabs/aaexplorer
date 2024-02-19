package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/transactionreceiptdecode"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
)

type transactionReceiptDao struct {
	baseDao
}

var TransactionReceiptDao = &transactionReceiptDao{}

type TransactionReceiptCondition struct {
	TxHash *string
}

func (dao *transactionReceiptDao) Pages(ctx context.Context, tx *ent.Client, page vo.PaginationRequest, condition TransactionReceiptCondition) (a []*ent.TransactionReceiptDecode, count int, err error) {
	query := tx.TransactionReceiptDecode.Query()

	if condition.TxHash != nil {
		query = query.Where(transactionreceiptdecode.ID(*condition.TxHash))
	}

	count = query.CountX(ctx)
	if count < 1 || page.GetOffset() > count {
		return
	}

	if page.Sort > 0 {
		query = query.Order(dao.orderPage(ctx, transactionreceiptdecode.Columns, page))
	} else {
		query = query.Order(transactionreceiptdecode.ByTime(sql.OrderDesc()))
	}

	query = query.Limit(page.GetPerPage()).Offset(page.GetOffset())
	a, err = query.All(ctx)
	return
}
