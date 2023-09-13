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
	list, total, err := dao.BundlerDao.Pagination(ctx, client, req)
	if err != nil {
		return nil, err
	}
	res.TotalCount = total

	//
	res.Records = make([]*vo.BundlersVo, len(list))
	for i, info := range list {

		res.Records[i] = &vo.BundlersVo{
			Bundler:        info.ID,
			BundlesNum:     info.BundlesNum,
			UserOpsNum:     info.UserOpsNum,
			SuccessRate:    info.SuccessRate,
			SuccessRateD1:  info.SuccessRateD1,
			BundlesNumD1:   info.BundlesNumD1,
			FeeEarnedD1:    info.FeeEarnedD1,
			FeeEarnedUsdD1: info.FeeEarnedUsdD1,
		}
		//account, err := dao.AccountDao.GetAbiByAddressWithMemo(ctx, client, info.Bundler)
		//if err == nil && account != nil && account.Label != nil {
		//	labels := []string{}
		//	account.Label.Scan(&labels)
		//	if len(labels) > 0 {
		//		res.Records[i].BundlerLabel = labels[0]
		//	}
		//}
	}

	return &res, nil
}
