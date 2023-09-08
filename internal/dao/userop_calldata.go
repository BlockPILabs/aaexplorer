package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aauseropscalldata"
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
