package dao

import (
	"context"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/bundlerinfo"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type bundlerDao struct {
	baseDao
}

var BundlerDao = &bundlerDao{}

func (*bundlerDao) GetSortFields(ctx context.Context) []string {
	return []string{
		config.Default,
		bundlerinfo.FieldID,
		bundlerinfo.FieldBundlesNum,
		bundlerinfo.FieldBundlesNumD1,
		bundlerinfo.FieldBundlesNumD7,
		bundlerinfo.FieldBundlesNumD30,
	}
}
func (dao *bundlerDao) Sort(ctx context.Context, query *ent.BundlerInfoQuery, sort int, order int) *ent.BundlerInfoQuery {
	opts := dao.orderOptions(ctx, order)
	if len(opts) > 0 {
		switch dao.sortField(ctx, dao.GetSortFields(ctx), sort) {
		case bundlerinfo.FieldID:
			query.Order(bundlerinfo.ByID(opts...))
		//case bundlerinfo.FieldBundlesNum:
		//	query.Order(bundlerinfo.ByBundlesNum(opts...))
		case bundlerinfo.FieldBundlesNumD1:
			query.Order(bundlerinfo.ByBundlesNumD1(opts...))
		case bundlerinfo.FieldBundlesNumD7:
			query.Order(bundlerinfo.ByBundlesNumD7(opts...))
		case bundlerinfo.FieldBundlesNumD30:
			query.Order(bundlerinfo.ByBundlesNumD30(opts...))
		default:
			query.Order(bundlerinfo.ByBundlesNum(opts...))
		}
	}
	return query
}

func (dao *bundlerDao) Pagination(ctx context.Context, tx *ent.Client, network string, page vo.PaginationRequest) (list ent.BundlerInfos, total int, err error) {
	query := tx.BundlerInfo.Query().Where(
		bundlerinfo.NetworkEQ(network),
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