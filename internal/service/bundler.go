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

	client, err := entity.Client(ctx)
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
	res.Records = make([]*vo.BundlerVo, len(list))
	for i, info := range list {
		res.Records[i] = &vo.BundlerVo{
			ID:              info.ID,
			Bundler:         info.Bundler,
			Network:         info.Network,
			UserOpsNum:      info.UserOpsNum,
			BundlesNum:      info.BundlesNum,
			GasCollected:    info.GasCollected,
			UserOpsNumD1:    info.UserOpsNumD1,
			BundlesNumD1:    info.BundlesNumD1,
			GasCollectedD1:  info.GasCollectedD1,
			UserOpsNumD7:    info.UserOpsNumD7,
			BundlesNumD7:    info.BundlesNumD7,
			GasCollectedD7:  info.GasCollectedD7,
			UserOpsNumD30:   info.UserOpsNumD30,
			BundlesNumD30:   info.BundlesNumD30,
			GasCollectedD30: info.GasCollectedD30,
		}
	}

	return &res, nil
}
