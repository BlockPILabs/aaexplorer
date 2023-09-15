package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aatransactioninfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactiondecode"
	"github.com/BlockPILabs/aa-scan/internal/utils"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type aaTransactionDao struct {
	baseDao
}

var AaTransactionDao = &aaTransactionDao{}

type AATransactionCondition struct {
	TxHashTerm string
}

func (dao *aaTransactionDao) Pagination(ctx context.Context, tx *ent.Client, page vo.PaginationRequest, condition AATransactionCondition) (a ent.AaTransactionInfos, count int, err error) {
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
		query = query.Order(dao.orderPage(ctx, transactiondecode.Columns, page))
	}

	query = query.Limit(page.GetPerPage()).Offset(page.GetOffset())
	a, err = query.All(ctx)
	return
}
