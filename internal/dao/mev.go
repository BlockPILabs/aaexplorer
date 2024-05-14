package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/mevtransaction"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
)

type mevDao struct {
	baseDao
}

var MevDao = &mevDao{}

func (*mevDao) GetSortFields() []string {
	return []string{
		config.Default,
		mevtransaction.FieldBundlerLoss,
		mevtransaction.FieldMevProfit,
	}
}

func (dao *mevDao) Sort(ctx context.Context, query *ent.MevTransactionQuery, sort int, order int) *ent.MevTransactionQuery {
	opts := dao.orderOptions(ctx, order)
	if len(opts) > 0 {
		f := dao.sortField(ctx, dao.GetSortFields(), sort)
		switch f {
		case "", config.Default:
			query.Order(mevtransaction.ByTime(opts...))
		default:
			query.Order(sql.OrderByField(f, opts...).ToFunc())
		}
	}
	return query
}

func (dao *mevDao) Pagination(ctx context.Context, tx *ent.Client, req vo.ListMEVBundlersRequest) (list ent.MevTransactions, total int, err error) {
	query := tx.MevTransaction.Query()

	if req.TotalCount > 0 {
		total = req.TotalCount
	} else {
		total = query.CountX(ctx)
	}

	if total < 1 || req.GetOffset() > total {
		return
	}

	query = dao.Sort(ctx, query, req.Sort, req.Order)

	query = query.Offset(req.GetOffset()).Limit(req.PerPage)

	list, err = query.All(ctx)

	if err != nil {
		return
	}
	return list, total, err
}

func (dao *mevDao) BlockMevPagination(ctx context.Context, tx *ent.Client, req vo.ListBlockMEVBundlersRequest) (list ent.MevTransactions, total int, err error) {
	query := tx.MevTransaction.Query().Where(mevtransaction.BlockNumberEQ(req.BlockNumber))

	if req.TotalCount > 0 {
		total = req.TotalCount
	} else {
		total = query.CountX(ctx)
	}

	if total < 1 || req.GetOffset() > total {
		return
	}

	query = dao.Sort(ctx, query, req.Sort, req.Order)

	query = query.Offset(req.GetOffset()).Limit(req.PerPage)

	list, err = query.All(ctx)

	if err != nil {
		return
	}
	return list, total, err
}

func (dao *mevDao) BundlerMevPagination(ctx context.Context, tx *ent.Client, req vo.ListBundlerMEVBundlersRequest) (list ent.MevTransactions, total int, err error) {
	query := tx.MevTransaction.Query().Where(mevtransaction.VictimEQ(req.Bundler))

	if req.TotalCount > 0 {
		total = req.TotalCount
	} else {
		total = query.CountX(ctx)
	}

	if total < 1 || req.GetOffset() > total {
		return
	}

	query = dao.Sort(ctx, query, req.Sort, req.Order)

	query = query.Offset(req.GetOffset()).Limit(req.PerPage)

	list, err = query.All(ctx)

	if err != nil {
		return
	}
	return list, total, err
}
