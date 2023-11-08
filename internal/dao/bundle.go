package dao

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aatransactioninfo"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
)

type bundleDao struct {
	baseDao
}

var BundleDao = &bundleDao{}

func (*bundleDao) GetSortFields(ctx context.Context) []string {
	return []string{
		config.Default,
		aatransactioninfo.FieldID,
		aatransactioninfo.FieldTime,
	}
}
func (dao *bundleDao) Sort(ctx context.Context, query *ent.AaTransactionInfoQuery, sort int, order int) *ent.AaTransactionInfoQuery {
	opts := dao.orderOptions(ctx, order)
	if len(opts) > 0 {
		switch dao.sortField(ctx, dao.GetSortFields(ctx), sort) {
		case aatransactioninfo.FieldID:
			query.Order(aatransactioninfo.ByID(opts...))
		//case transactioninfo.FieldTxTime:
		//	query.Order(transactioninfo.ByTxTime(opts...))
		default:
			query.Order(aatransactioninfo.ByTime(opts...))
		}
	}
	return query
}

func (dao *bundleDao) Pagination(ctx context.Context, tx *ent.Client, network string, page vo.PaginationRequest) (list ent.AaTransactionInfos, total int, err error) {
	query := tx.AaTransactionInfo.Query().Where()
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
