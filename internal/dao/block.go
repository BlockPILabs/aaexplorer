package dao

import (
	"context"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/blockdatadecode"
	"github.com/BlockPILabs/aa-scan/internal/utils"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/shopspring/decimal"
	"strconv"
)

type blockDao struct {
	baseDao
}

var BlockDao = &blockDao{}

func (*blockDao) GetSortFields(ctx context.Context) []string {
	return []string{
		config.Default,
		blockdatadecode.FieldID,
		blockdatadecode.FieldCreateTime,
	}
}
func (dao *blockDao) Sort(ctx context.Context, query *ent.BlockDataDecodeQuery, sort int, order int) *ent.BlockDataDecodeQuery {
	opts := dao.orderOptions(ctx, order)
	if len(opts) > 0 {
		switch dao.sortField(ctx, dao.GetSortFields(ctx), sort) {
		case blockdatadecode.FieldID:
			query.Order(blockdatadecode.ByID(opts...))
		default:
			query.Order(blockdatadecode.ByID(opts...))
		}
	}
	return query
}

func (dao *blockDao) Pagination(ctx context.Context, tx *ent.Client, network string, page vo.PaginationRequest) (list ent.BlockDataDecodes, total int, err error) {
	query := tx.BlockDataDecode.Query()
	// sort
	query = dao.Sort(ctx, query, page.Sort, page.Order)

	total = query.CountX(ctx)

	if total < 1 || page.GetOffset() > total {
		return
	}

	// limit
	query = query.Select(page.Select...).
		Offset(page.GetOffset()).
		Limit(page.PerPage)

	list, err = query.All(ctx)
	return
}

func (*baseDao) GetByBlockNumber(ctx context.Context, tx *ent.Client, blockNumber decimal.Decimal) (block *ent.BlockDataDecode, err error) {
	//block, err = tx.BlockDataDecode.Query().Where(
	//	blockdatadecode.ID(blockNumber),
	//).First(ctx)
	return
}

func (*baseDao) GetByBlockHash(ctx context.Context, tx *ent.Client, hash string) (block *ent.BlockDataDecode, err error) {
	block, err = tx.BlockDataDecode.Query().Where(
		blockdatadecode.Hash(hash),
	).First(ctx)
	return
}

func (dao *baseDao) GetBlock(ctx context.Context, tx *ent.Client, blockNumberOrHash string) (block *ent.BlockDataDecode, err error) {
	if utils.Has0xPrefix(blockNumberOrHash) {
		return dao.GetByBlockHash(ctx, tx, blockNumberOrHash)
	} else {
		blockNumber, err := strconv.ParseInt(blockNumberOrHash, 10, 64)
		if err != nil {
			return nil, err
		}
		return dao.GetByBlockNumber(ctx, tx, decimal.NewFromInt(blockNumber))
	}
}
