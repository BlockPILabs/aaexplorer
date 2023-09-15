package dao

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aauseropscalldata"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

import (
	"entgo.io/ent/dialect/sql"
)

type userOpCallDataDao struct {
	baseDao
}

var UserOpCallDataDao = &userOpCallDataDao{}

func (dao *userOpCallDataDao) GetTargets(ctx context.Context, tx *ent.Client, userOpsHashIn []string) (lists map[string][]string, err error) {
	query := tx.AAUserOpsCalldata.Query().Where(
		aauseropscalldata.UserOpsHashIn(userOpsHashIn...),
	)
	// sort
	query = query.Order(
		aauseropscalldata.ByAaIndex(sql.OrderAsc()),
	)

	list, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	lists = map[string][]string{}
	for _, calldata := range list {
		if _, ok := lists[calldata.UserOpsHash]; ok {
			lists[calldata.UserOpsHash] = append(lists[calldata.UserOpsHash], calldata.Target)
		} else {
			lists[calldata.UserOpsHash] = []string{calldata.Target}
		}
	}

	return
}

type UserOpsCallDataCondition struct {
	UserOperationHash *string
}

func (dao *userOpCallDataDao) Pages(ctx context.Context, tx *ent.Client, page vo.PaginationRequest, condition UserOpsCallDataCondition) (a []*ent.AAUserOpsCalldata, count int, err error) {
	query := tx.AAUserOpsCalldata.Query()

	if condition.UserOperationHash != nil {
		query = query.Where(aauseropscalldata.UserOpsHash(*condition.UserOperationHash))
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
