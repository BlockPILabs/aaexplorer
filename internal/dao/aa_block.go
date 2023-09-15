package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aablockinfo"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type aaBlockDao struct {
	baseDao
}

var AaBlockDao = &aaBlockDao{}

type AaBlockPagesCondition struct {
	LatestBlockNumber int64
}

func (dao *aaBlockDao) Pages(ctx context.Context, tx *ent.Client, page vo.PaginationRequest, condition AaBlockPagesCondition) (a []*ent.AaBlockInfo, count int, err error) {
	query := tx.AaBlockInfo.Query()
	if condition.LatestBlockNumber > 0 {
		query = query.Where(
			aablockinfo.IDGT(condition.LatestBlockNumber),
		)
	}
	count = query.CountX(ctx)
	if count < 1 || page.GetOffset() > count {
		return
	}

	if page.Sort > 0 {
		query = query.Order(dao.orderPage(ctx, aablockinfo.Columns, page))
	}

	query = query.Limit(page.GetPerPage()).Offset(page.GetOffset())
	a, err = query.All(ctx)
	return
}

func (dao *aaBlockDao) GetLatestBlock(ctx context.Context, tx *ent.Client) (a *ent.AaBlockInfo, err error) {
	query := tx.AaBlockInfo.Query()
	query = query.Order(aablockinfo.ByID(sql.OrderDesc()))
	return query.First(ctx)
}
