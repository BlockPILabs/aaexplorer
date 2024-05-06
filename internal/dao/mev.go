package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/mevinfo"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
)

type mevDao struct {
	baseDao
}

var MevDao = &mevDao{}

func (*mevDao) GetSortFields() []string {
	return []string{
		config.Default,
		mevinfo.FieldProfit,
	}
}

func (dao *mevDao) Sort(ctx context.Context, query *ent.MevInfoQuery, sort int, order int) *ent.MevInfoQuery {
	opts := dao.orderOptions(ctx, order)
	if len(opts) > 0 {
		f := dao.sortField(ctx, dao.GetSortFields(), sort)
		switch f {
		case "", config.Default:
			query.Order(mevinfo.ByTxTime(opts...))
		default:
			query.Order(sql.OrderByField(f, opts...).ToFunc())
		}
	}
	return query
}

func (dao *mevDao) Pagination(ctx context.Context, tx *ent.Client, req vo.ListMEVBundlersRequest) (list ent.MevInfos, total int, err error) {
	query := tx.MevInfo.Query().Where(
		mevinfo.NetworkEQ(req.Network),
	)

	if req.TotalCount > 0 {
		total = req.TotalCount
	} else {
		total = query.CountX(ctx)
	}

	if total < 1 || req.GetOffset() > total {
		return
	}

	query = dao.Sort(ctx, query, req.Sort, req.Order)

	query = query.Offset(req.GetOffset()).Limit(req.GetPage())

	list, err = query.All(ctx)

	if err != nil {
		return
	}
	return list, total, err
}
