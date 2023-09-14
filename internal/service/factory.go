package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type factoryService struct {
}

var FactoryService = &factoryService{}

func (*factoryService) GetFactories(ctx context.Context, req vo.GetFactoriesRequest) (*vo.GetFactoriesResponse, error) {
	ctx, logger := log.With(ctx, "service", "GetFactories")
	err := vo.ValidateStruct(req)
	res := vo.GetFactoriesResponse{
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
	list, total, err := dao.FactoryDao.Pagination(ctx, client, req)
	if err != nil {
		return nil, err
	}
	res.TotalCount = total
	res.Records = make([]*vo.FactoryVo, len(list))
	if res.TotalCount > 0 {
		//
		for i, info := range list {
			_ = info
			factoryVo := &vo.FactoryVo{
				ID:           info.ID,
				AccountNum:   info.AccountNum,
				AccountNumD1: info.AccountNumD1,
			}
			res.Records[i] = factoryVo
		}
	}

	return &res, nil
}
