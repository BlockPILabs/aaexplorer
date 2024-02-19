package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/factoryinfo"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
)

type factoryDao struct {
	baseDao
}

var FactoryDao = &factoryDao{}

func (*factoryDao) GetSortFields(ctx context.Context) []string {
	return []string{
		config.Default,
		factoryinfo.FieldAccountNum,
		factoryinfo.FieldAccountNumD1,
	}
}
func (dao *factoryDao) Sort(ctx context.Context, query *ent.FactoryInfoQuery, sort int, order int) *ent.FactoryInfoQuery {
	opts := dao.orderOptions(ctx, order)
	if len(opts) > 0 {
		f := dao.sortField(ctx, dao.GetSortFields(ctx), sort)
		switch f {
		case "", config.Default:
			query.Order(factoryinfo.ByAccountNum(opts...))
		default:
			query.Order(sql.OrderByField(f, opts...).ToFunc())
		}
	}
	return query
}

func (dao *factoryDao) Pagination(ctx context.Context, tx *ent.Client, req vo.GetFactoriesRequest) (list ent.FactoryInfos, total int, err error) {

	query := tx.FactoryInfo.Query().Where(
		factoryinfo.NetworkEQ(req.Network),
	)

	total = query.CountX(ctx)

	if total < 1 || req.GetOffset() > total {
		return
	}
	// sort
	query = dao.Sort(ctx, query, req.Sort, req.Order)

	list, err = query.WithAccount().
		Offset(req.GetOffset()).
		Limit(req.PerPage).All(ctx)
	return
}
