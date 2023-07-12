package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type bundleService struct {
}

var BundleService = &bundleService{}

func (*bundleService) GetBundles(ctx context.Context, req vo.GetBundlesRequest) (*vo.GetBundlesResponse, error) {
	ctx, logger := log.With(ctx, "service", "GetBundlers")
	err := vo.ValidateStruct(req)
	res := vo.GetBundlesResponse{
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
	list, total, err := dao.BundleDao.Pagination(ctx, client, req.Network, req.PaginationRequest)
	if err != nil {
		return nil, err
	}
	res.TotalCount = total

	//
	res.Records = make([]*vo.BundleVo, len(list))
	for i, info := range list {
		res.Records[i] = &vo.BundleVo{
			ID:           info.ID,
			TxHash:       info.TxHash,
			BlockNumber:  info.BlockNumber,
			Network:      info.Network,
			Bundler:      info.Bundler,
			EntryPoint:   info.EntryPoint,
			UserOpsNum:   info.UserOpsNum,
			TxValue:      info.TxValue,
			Fee:          info.Fee,
			GasPrice:     info.GasPrice,
			GasLimit:     info.GasLimit,
			Status:       info.Status,
			TxTime:       info.TxTime,
			TxTimeFormat: info.TxTimeFormat,
			Beneficiary:  info.Beneficiary,
			CreateTime:   info.CreateTime,
		}
	}

	return &res, nil
}
