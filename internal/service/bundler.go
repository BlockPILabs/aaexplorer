package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type bundlerService struct {
}

var BundlerService = &bundlerService{}

func (*bundlerService) GetBundlers(ctx context.Context, req vo.GetBundlersRequest) (*vo.GetBundlersResponse, error) {
	ctx, logger := log.With(ctx, "service", "GetBundlers")
	err := vo.ValidateStruct(req)
	res := vo.GetBundlersResponse{
		Pagination: vo.Pagination{
			TotalCount: 0,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
	}
	if err != nil {
		logger.Error("params error", "req", req, "err", err.Error())
		return &res, vo.ErrParams.SetData(err)
	}

	client, err := entity.Client(ctx, req.Network)
	if err != nil {
		return nil, err
	}
	//
	list, total, err := dao.BundlerDao.Pagination(ctx, client, req.Network, req.PaginationRequest)
	if err != nil {
		return nil, err
	}
	res.TotalCount = total

	//
	res.Records = make([]*vo.BundlersVo, len(list))
	for i, info := range list {
		res.Records[i] = &vo.BundlersVo{
			Bundler:      info.ID,
			Network:      info.Network,
			UserOpsNum:   info.UserOpsNum,
			BundlesNum:   info.BundlesNum,
			GasCollected: info.GasCollected,
		}
	}

	return &res, nil
}
