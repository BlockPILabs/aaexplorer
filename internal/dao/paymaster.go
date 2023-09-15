package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/paymasterinfo"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type paymasterDao struct {
	baseDao
}

var PaymasterDao = &paymasterDao{}

func (*paymasterDao) GetSortFields(ctx context.Context) []string {
	return []string{
		config.Default,
		paymasterinfo.FieldUserOpsNum,
		paymasterinfo.FieldUserOpsNumD1,
		paymasterinfo.FieldReserve,
		paymasterinfo.FieldGasSponsored,
	}
}
func (dao *paymasterDao) Sort(ctx context.Context, query *ent.PaymasterInfoQuery, sort int, order int) *ent.PaymasterInfoQuery {
	opts := dao.orderOptions(ctx, order)
	if len(opts) > 0 {
		f := dao.sortField(ctx, dao.GetSortFields(ctx), sort)
		switch f {
		case "", config.Default:
			query.Order(paymasterinfo.ByUserOpsNum(opts...))
		default:
			query.Order(sql.OrderByField(f, opts...).ToFunc())
		}
	}
	return query
}

func (dao *paymasterDao) Pagination(ctx context.Context, tx *ent.Client, req vo.GetPaymastersRequest) (list ent.PaymasterInfos, total int, err error) {
	query := tx.PaymasterInfo.Query().Where(
		paymasterinfo.NetworkEQ(req.Network),
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
