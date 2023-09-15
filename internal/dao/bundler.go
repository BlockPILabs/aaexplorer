package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/bundlerinfo"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type bundlerDao struct {
	baseDao
}

var BundlerDao = &bundlerDao{}

func (*bundlerDao) GetSortFields(ctx context.Context) []string {
	return []string{
		config.Default,
		bundlerinfo.FieldBundlesNum,
		bundlerinfo.FieldSuccessRate,
		bundlerinfo.FieldBundleRateD1,
		bundlerinfo.FieldUserOpsNum,
		bundlerinfo.FieldBundlesNumD1,
		bundlerinfo.FieldFeeEarnedD1,
	}
}
func (dao *bundlerDao) Sort(ctx context.Context, query *ent.BundlerInfoQuery, sort int, order int) *ent.BundlerInfoQuery {
	opts := dao.orderOptions(ctx, order)
	if len(opts) > 0 {
		f := dao.sortField(ctx, dao.GetSortFields(ctx), sort)
		switch f {
		case "", config.Default:
			query.Order(bundlerinfo.ByBundlesNum(opts...))
		default:
			query.Order(sql.OrderByField(f, opts...).ToFunc())
		}
	}
	return query
}

func (dao *bundlerDao) Pagination(ctx context.Context, tx *ent.Client, req vo.GetBundlersRequest) (list ent.BundlerInfos, total int, err error) {
	query := tx.BundlerInfo.Query().Where(
		bundlerinfo.NetworkEQ(req.Network),
	)

	total = query.CountX(ctx)

	if total < 1 || req.GetOffset() > total {
		return
	}

	// sort
	query = dao.Sort(ctx, query, req.Sort, req.Order)
	// limit
	query = query.WithAccount().
		Offset(req.GetOffset()).
		Limit(req.PerPage)

	list, err = query.All(ctx)
	return
}
