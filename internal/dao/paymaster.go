package dao

import (
	"context"
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
		paymasterinfo.FieldID,
		paymasterinfo.FieldUserOpsNum,
	}
}
func (dao *paymasterDao) Sort(ctx context.Context, query *ent.PaymasterInfoQuery, sort int, order int) *ent.PaymasterInfoQuery {
	opts := dao.orderOptions(ctx, order)
	if len(opts) > 0 {
		switch dao.sortField(ctx, dao.GetSortFields(ctx), sort) {
		case paymasterinfo.FieldID:
			query.Order(paymasterinfo.ByID(opts...))
		default:
			query.Order(paymasterinfo.ByUserOpsNum(opts...))
		}
	}
	return query
}

func (dao *paymasterDao) Pagination(ctx context.Context, tx *ent.Client, network string, page vo.PaginationRequest) (list ent.PaymasterInfos, total int, err error) {
	query := tx.PaymasterInfo.Query().Where(
		paymasterinfo.NetworkEQ(network),
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
