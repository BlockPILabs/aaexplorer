package dao

import (
	"context"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aauseropsinfo"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type userOpDao struct {
	baseDao
}

var UserOpDao = &userOpDao{}

func (*userOpDao) GetSortFields(ctx context.Context) []string {
	return []string{
		config.Default,
		aauseropsinfo.FieldID,
		aauseropsinfo.FieldTime,
	}
}
func (dao *userOpDao) Sort(ctx context.Context, query *ent.AAUserOpsInfoQuery, sort int, order int) *ent.AAUserOpsInfoQuery {
	opts := dao.orderOptions(ctx, order)
	if len(opts) > 0 {
		switch dao.sortField(ctx, dao.GetSortFields(ctx), sort) {
		case aauseropsinfo.FieldID:
			query.Order(aauseropsinfo.ByID(opts...))
		//case useropsinfo.FieldTxTime:
		//	query.Order(useropsinfo.ByTxTime(opts...))
		default:
			query.Order(aauseropsinfo.ByTime(opts...))
		}
	}
	return query
}

func (dao *userOpDao) Pagination(ctx context.Context, tx *ent.Client, network string, page vo.PaginationRequest) (list ent.AAUserOpsInfos, total int, err error) {
	query := tx.AAUserOpsInfo.Query().Where(
		aauseropsinfo.NetworkEQ(network),
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
