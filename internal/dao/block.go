package dao

import (
	"context"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/blockdata"
	"github.com/BlockPILabs/aa-scan/internal/utils"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"strconv"
)

type blockDao struct {
	baseDao
}

var BlockDao = &blockDao{}

func (*blockDao) GetSortFields(ctx context.Context) []string {
	return []string{
		config.Default,
		blockdata.FieldID,
		blockdata.FieldCreateTime,
	}
}
func (dao *blockDao) Sort(ctx context.Context, query *ent.BlockDataQuery, sort int, order int) *ent.BlockDataQuery {
	opts := dao.orderOptions(ctx, order)
	if len(opts) > 0 {
		switch dao.sortField(ctx, dao.GetSortFields(ctx), sort) {
		case blockdata.FieldID:
			query.Order(blockdata.ByID(opts...))
		default:
			query.Order(blockdata.ByID(opts...))
		}
	}
	return query
}

func (dao *blockDao) Pagination(ctx context.Context, tx *ent.Client, network string, page vo.PaginationRequest) (list ent.BlockDataSlice, total int, err error) {
	query := tx.BlockData.Query().Where()
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

func (*baseDao) GetByBlockNumber(ctx context.Context, tx *ent.Client, blockNumber int64) (block *ent.BlockData, err error) {
	block, err = tx.BlockData.Query().Where(
		blockdata.ID(blockNumber),
	).First(ctx)
	return
}

func (*baseDao) GetByBlockHash(ctx context.Context, tx *ent.Client, hash string) (block *ent.BlockData, err error) {
	block, err = tx.BlockData.Query().Where(
		blockdata.Hash(hash),
	).First(ctx)
	return
}

func (dao *baseDao) GetBlock(ctx context.Context, tx *ent.Client, blockNumberOrHash string) (block *ent.BlockData, err error) {
	if utils.Has0xPrefix(blockNumberOrHash) {
		return dao.GetByBlockHash(ctx, tx, blockNumberOrHash)
	} else {
		blockNumber, err := strconv.ParseInt(blockNumberOrHash, 10, 64)
		if err != nil {
			return nil, err
		}
		return dao.GetByBlockNumber(ctx, tx, blockNumber)
	}
}
