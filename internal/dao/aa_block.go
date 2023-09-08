package dao

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aablockinfo"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type aaBlockDao struct {
	baseDao
}

var AaBlockDao = &aaBlockDao{}

type AaBlockPagesCondition struct {
}

func (dao *aaBlockDao) Pages(ctx context.Context, tx *ent.Client, page vo.PaginationRequest, condition AaBlockPagesCondition) (a []*ent.AaBlockInfo, count int, err error) {
	query := tx.AaBlockInfo.Query()

	if page.Sort > 0 {
		query = query.Order(dao.orderPage(ctx, aablockinfo.Columns, page))
	}

	count = query.CountX(ctx)
	if count < 1 || page.GetOffset() > count {
		return
	}

	query = query.Limit(page.GetPerPage()).Offset(page.GetOffset())
	a, err = query.All(ctx)
	return
}
