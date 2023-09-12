package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type aaBlockService struct {
}

var AaBlockService = &aaBlockService{}

func (*aaBlockService) GetAaBlockInfo(ctx context.Context, client *ent.Client, req vo.GetAaBlocksRequest) (*vo.GetAaBlocksResponse, error) {
	ctx, logger := log.With(ctx, "service", "GetAaBlockInfo")
	logger.Info("GetAaBlockInfo ... ")
	res := vo.GetAaBlocksResponse{
		Pagination: vo.Pagination{
			TotalCount: 0,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
	}

	pages, total, err := dao.AaBlockDao.Pages(ctx, client, req.PaginationRequest, dao.AaBlockPagesCondition{})
	if err != nil {
		return nil, err
	}

	res.TotalCount = total

	res.Records = make([]*vo.AaBlocksVo, len(pages))
	for i, info := range pages {
		res.Records[i] = &vo.AaBlocksVo{
			Number:         info.ID,
			Time:           info.Time,
			Hash:           info.Hash,
			UseropCount:    info.UseropCount,
			UseropMevCount: info.UseropMevCount,
			BundlerProfit:  info.BundlerProfit,
			CreateTime:     info.CreateTime,
		}
	}

	return &res, nil
}