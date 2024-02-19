package service

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/internal/dao"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/factoryinfo"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
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
		Records: make([]*vo.FactoryVo, 0),
	}
	if err != nil {
		logger.Error("params error", "req", req, "err", err.Error())
		return &res, vo.ErrParams.SetData(err)
	}

	client, err := entity.Client(ctx, req.Network)
	if err != nil {
		return &res, err
	}
	//
	list, total, err := dao.FactoryDao.Pagination(ctx, client, req)
	if err != nil {
		return &res, err
	}
	res.TotalCount = total
	res.Records = make([]*vo.FactoryVo, len(list))
	if res.TotalCount > 0 {
		//
		for i, info := range list {
			label := ""
			if info.Edges.Account != nil && info.Edges.Account.Label != nil {
				labels := []string{}
				info.Edges.Account.Label.AssignTo(&labels)
				if len(labels) > 0 {
					label = labels[0]
				}
			}
			factoryVo := &vo.FactoryVo{
				ID:           info.ID,
				AccountNum:   info.AccountNum,
				AccountNumD1: info.AccountNumD1,
				Dominance:    info.Dominance,
				DominanceD1:  info.DominanceD1,
				FactoryLabel: label,
			}
			res.Records[i] = factoryVo
		}
	}

	return &res, nil
}

func (*factoryService) GetFactory(ctx context.Context, req vo.GetFactoryRequest) (res *vo.GetFactoryResponse, err error) {
	res = &vo.GetFactoryResponse{}
	client, err := entity.Client(ctx, req.Network)
	if err != nil {
		return
	}

	info, err := client.FactoryInfo.Get(ctx, req.Factory)
	if err != nil {
		return
	}

	acc, err := client.AaAccountData.Get(ctx, req.Factory)
	if err != nil {
		return
	}
	res = &vo.GetFactoryResponse{
		TotalAccountDeployNum: info.AccountDeployNum,
		AccountDeployNumD1:    info.AccountDeployNumD1,
		Dominance:             info.Dominance,
		UserOpsNum:            acc.UserOpsNum,
		Rank:                  0,
	}

	res.Rank = int64(
		client.FactoryInfo.Query().Where(
			factoryinfo.AccountDeployNumGT(info.AccountDeployNum),
		).CountX(ctx),
	) + 1
	res.TotalNumber = int64(
		client.FactoryInfo.Query().CountX(ctx),
	)

	addresses, _ := dao.AccountDao.GetAccountByAddresses(ctx, client, []string{req.Factory})
	if len(addresses) > 0 {
		addresses[0].Label.AssignTo(&res.Label)
	}

	return
}
