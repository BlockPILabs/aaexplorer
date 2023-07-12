package dao

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/useropsinfo"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type userOpDao struct {
	baseDao
}

var UserOpDao = &userOpDao{}

func (*userOpDao) GetSortFields(ctx context.Context) []string {
	return []string{
		useropsinfo.FieldID,
	}
}
func (dao *userOpDao) Sort(ctx context.Context, query *ent.UserOpsInfoQuery, sort int, order int) *ent.UserOpsInfoQuery {
	opts := dao.orderOptions(ctx, order)
	if len(opts) > 0 {
		switch dao.sortField(ctx, dao.GetSortFields(ctx), sort) {
		case useropsinfo.FieldID:
			query.Order(useropsinfo.ByID(opts...))
		default:
			break
		}
	}
	return query
}

func (dao *userOpDao) Pagination(ctx context.Context, tx *ent.Client, network string, page vo.PaginationRequest) (list ent.UserOpsInfos, total int, err error) {
	query := tx.UserOpsInfo.Query().Where(
		useropsinfo.NetworkEQ(network),
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
