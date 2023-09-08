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
		case aauseropsinfo.FieldTxTime:
			query.Order(aauseropsinfo.ByTxTime(opts...))
		default:
			query.Order(aauseropsinfo.ByBlockNumber(opts...))
		}
	}
	return query
}

func (dao *userOpDao) Pagination(ctx context.Context, tx *ent.Client, req vo.GetUserOpsRequest) (list ent.AAUserOpsInfos, total int, err error) {
	query := tx.AAUserOpsInfo.Query().Where(
		aauseropsinfo.NetworkEQ(req.Network),
	)

	if req.LatestBlockNumber > 0 {
		query = query.Where(
			aauseropsinfo.BlockNumberGT(req.LatestBlockNumber),
		)
	}
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

type UserOpsCondition struct {
	UserOperationHash *string
	TxHash            *string
}

func (dao *userOpDao) Pages(ctx context.Context, tx *ent.Client, page vo.PaginationRequest, condition UserOpsCondition) (a []*ent.AAUserOpsInfo, count int, err error) {
	query := tx.AAUserOpsInfo.Query()

	if condition.UserOperationHash != nil {
		query = query.Where(aauseropsinfo.ID(*condition.UserOperationHash))
	}

	if condition.TxHash != nil {
		query = query.Where(aauseropsinfo.TxHash(*condition.TxHash))
	}

	if page.Sort > 0 {
		query = query.Order(dao.orderPage(ctx, aauseropsinfo.Columns, page))
	}

	count = query.CountX(ctx)
	if count < 1 || page.GetOffset() > count {
		return
	}

	query = query.Limit(page.GetPerPage()).Offset(page.GetOffset())
	a, err = query.All(ctx)
	return
}
