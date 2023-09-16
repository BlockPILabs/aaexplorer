package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/bundlerinfo"
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
	list, total, err := dao.BundlerDao.Pagination(ctx, client, req)
	if err != nil {
		return nil, err
	}
	res.TotalCount = total

	//
	res.Records = make([]*vo.BundlersVo, len(list))
	for i, info := range list {
		label := ""
		if info.Edges.Account != nil && info.Edges.Account.Label != nil {
			labels := []string{}
			info.Edges.Account.Label.AssignTo(&labels)
			if len(labels) > 0 {
				label = labels[0]
			}
		}
		res.Records[i] = &vo.BundlersVo{
			Bundler:        info.ID,
			BundlesNum:     info.BundlesNum,
			UserOpsNum:     info.UserOpsNum,
			SuccessRate:    info.SuccessRate,
			SuccessRateD1:  info.SuccessRateD1,
			BundlesNumD1:   info.BundlesNumD1,
			FeeEarnedD1:    info.FeeEarnedD1,
			FeeEarnedUsdD1: info.FeeEarnedUsdD1,
			BundlerLabel:   label,
			BundleRate:     info.BundleRate,
		}
	}

	return &res, nil
}

func (*bundlerService) GetBundler(ctx context.Context, req vo.GetBundlerRequest) (res *vo.GetBundlerResponse, err error) {
	res = &vo.GetBundlerResponse{}
	client, err := entity.Client(ctx, req.Network)
	if err != nil {
		return
	}

	info, err := client.BundlerInfo.Get(ctx, req.Bundler)
	if err != nil {
		return
	}

	res = &vo.GetBundlerResponse{
		FeeEarnedUsdD1: info.FeeEarnedUsdD1,
		FeeEarnedUsd:   info.FeeEarnedUsd,
		SuccessRateD1:  info.SuccessRateD1,
		SuccessRate:    info.SuccessRate,
		BundleRate:     info.BundleRate,
		Rank:           999,
		TotalBundlers:  int64(client.BundlerInfo.Query().CountX(ctx)),
	}

	res.Rank = res.TotalBundlers
	res.Rank = int64(
		client.BundlerInfo.Query().Where(
			bundlerinfo.BundlesNumGT(info.BundlesNum),
		).CountX(ctx),
	)
	return
}
