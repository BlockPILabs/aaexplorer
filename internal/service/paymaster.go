package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type paymasterService struct {
}

var PaymasterService = &paymasterService{}

func (*paymasterService) GetPaymasters(ctx context.Context, req vo.GetPaymastersRequest) (*vo.GetPaymastersResponse, error) {
	ctx, logger := log.With(ctx, "service", "GetPaymasters")
	err := vo.ValidateStruct(req)
	res := vo.GetPaymastersResponse{
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

	list, total, err := dao.PaymasterDao.Pagination(ctx, client, req)
	if err != nil {
		return nil, err
	}
	res.TotalCount = total

	//
	res.Records = make([]*vo.PaymastersVo, len(list))
	for i, info := range list {
		res.Records[i] = &vo.PaymastersVo{
			Paymaster:       info.ID,
			UserOpsNum:      info.UserOpsNum,
			UserOpsNumD1:    info.UserOpsNumD1,
			Reserve:         info.Reserve,
			GasSponsored:    info.GasSponsored,
			GasSponsoredUsd: info.GasSponsoredUsd,
		}
	}

	return &res, nil
}
