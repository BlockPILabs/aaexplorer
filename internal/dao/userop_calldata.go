package dao

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aauseropscalldata"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type userOpCallDataDao struct {
	baseDao
}

var UserOpCallDataDao = &userOpCallDataDao{}

type UserOpsCallDataCondition struct {
	UserOperationHash *string
}

func (dao *userOpCallDataDao) Pages(ctx context.Context, tx *ent.Client, page vo.PaginationRequest, condition UserOpsCallDataCondition) (a []*ent.AAUserOpsCalldata, count int, err error) {
	query := tx.AAUserOpsCalldata.Query()

	if condition.UserOperationHash != nil {
		query = query.Where(aauseropscalldata.ID(*condition.UserOperationHash))
	}

	if page.Sort > 0 {
		query = query.Order(dao.orderPage(ctx, aauseropscalldata.Columns, page))
	}

	count = query.CountX(ctx)
	if count < 1 || page.GetOffset() > count {
		return
	}

	query = query.Limit(page.GetPerPage()).Offset(page.GetOffset())
	a, err = query.All(ctx)
	return
}
