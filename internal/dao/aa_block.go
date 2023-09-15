package dao

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aablockinfo"
	"github.com/BlockPILabs/aa-scan/internal/utils"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type aaBlockDao struct {
	baseDao
}

var AaBlockDao = &aaBlockDao{}

type AaBlockPagesCondition struct {
	LatestBlockNumber int64
	HashTerm          string
}

func (dao *aaBlockDao) Pages(ctx context.Context, tx *ent.Client, page vo.PaginationRequest, condition AaBlockPagesCondition) (a []*ent.AaBlockInfo, count int, err error) {
	query := tx.AaBlockInfo.Query()
	if condition.LatestBlockNumber > 0 {
		query = query.Where(
			aablockinfo.IDGT(condition.LatestBlockNumber),
		)
	}

	if len(condition.HashTerm) > 0 && utils.IsHexSting(condition.HashTerm) {
		if utils.IsHashHex(condition.HashTerm) {
			query = query.Where(aablockinfo.HashEQ(utils.Fix0x(condition.HashTerm)))
		} else {
			query = query.Where(sql.FieldHasPrefix(aablockinfo.FieldHash, utils.Fix0x(condition.HashTerm)))
		}
	}

	if page.TotalCount > 0 {
		count = page.TotalCount
	} else {
		count = query.CountX(ctx)
	}
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

func (*aaBlockDao) GetByBlockNumber(ctx context.Context, tx *ent.Client, blockNumber int64) (block *ent.AaBlockInfo, err error) {
	block, err = tx.AaBlockInfo.Query().Where(
		aablockinfo.ID(blockNumber),
	).First(ctx)
	return
}
