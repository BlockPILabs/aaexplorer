package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/factoryinfo"
	"github.com/BlockPILabs/aa-scan/internal/vo"
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
	// sort
	query = dao.Sort(ctx, query, req.Sort, req.Order)

	total = query.CountX(ctx)

	if total < 1 || req.GetOffset() > total {
		return
	}

	// limit
	query = query.
		Offset(req.GetOffset()).
		Limit(req.PerPage)

	list, err = query.All(ctx)
	return
}
