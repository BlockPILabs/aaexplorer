package dao

import (
	"context"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactioninfo"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type bundleDao struct {
	baseDao
}

var BundleDao = &bundleDao{}

func (*bundleDao) GetSortFields(ctx context.Context) []string {
	return []string{
		config.Default,
		transactioninfo.FieldID,
		transactioninfo.FieldTxTime,
	}
}
func (dao *bundleDao) Sort(ctx context.Context, query *ent.TransactionInfoQuery, sort int, order int) *ent.TransactionInfoQuery {
	opts := dao.orderOptions(ctx, order)
	if len(opts) > 0 {
		switch dao.sortField(ctx, dao.GetSortFields(ctx), sort) {
		case transactioninfo.FieldID:
			query.Order(transactioninfo.ByID(opts...))
		//case transactioninfo.FieldTxTime:
		//	query.Order(transactioninfo.ByTxTime(opts...))
		default:
			query.Order(transactioninfo.ByTxTime(opts...))
		}
	}
	return query
}

func (dao *bundleDao) Pagination(ctx context.Context, tx *ent.Client, network string, page vo.PaginationRequest) (list ent.TransactionInfos, total int, err error) {
	query := tx.TransactionInfo.Query().Where(
		transactioninfo.NetworkEQ(network),
	)
	// sort
	query = dao.Sort(ctx, query, page.Sort, page.Order)

	total = query.CountX(ctx)

	if total < 1 || page.GetOffset() > total {
		return
	}

	// limit
	query = query.
		Offset(page.GetOffset()).
		Limit(page.PerPage)

	list, err = query.All(ctx)
	return
}
