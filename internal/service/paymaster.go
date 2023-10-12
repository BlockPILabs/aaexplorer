package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/paymasterinfo"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/BlockPILabs/aa-scan/service"
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
		label := ""
		if info.Edges.Account != nil && info.Edges.Account.Label != nil {
			labels := []string{}
			info.Edges.Account.Label.AssignTo(&labels)
			if len(labels) > 0 {
				label = labels[0]
			}
		}
		res.Records[i] = &vo.PaymastersVo{
			Paymaster:       info.ID,
			UserOpsNum:      info.UserOpsNum,
			UserOpsNumD1:    info.UserOpsNumD1,
			Reserve:         info.Reserve,
			GasSponsored:    info.GasSponsored,
			GasSponsoredUsd: info.GasSponsoredUsd,
			PaymasterLabel:  label,
		}
	}

	return &res, nil
}

/*func GetReservePage(ctx context.Context, req vo.GetReserveRequest) (*vo.GetReserveResponse, error) {
	ctx, logger := log.With(ctx, "service", "GetReservePage")
	err := vo.ValidateStruct(req)
	res := vo.GetReserveResponse{
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

	return &res, nil
}*/

func GetPaymasterOverview(ctx context.Context, req vo.GetPaymasterOverviewRequest) (*vo.GetPaymasterOverviewResponse, error) {
	ctx, logger := log.With(ctx, "service", "GetPaymasterOverview")
	err := vo.ValidateStruct(req)
	res := vo.GetPaymasterOverviewResponse{}
	if err != nil {
		logger.Error("params error", "req", req, "err", err.Error())
		return &res, vo.ErrParams.SetData(err)
	}

	client, err := entity.Client(ctx, req.Network)
	if err != nil {
		return nil, err
	}
	paymaster := req.Paymaster
	paymasterInfos, err := client.PaymasterInfo.Query().Where(paymasterinfo.IDEQ(paymaster)).All(context.Background())
	if err != nil {
		return nil, err
	}
	if len(paymasterInfos) == 0 {
		return nil, nil
	}

	allUserOpsNum, err := client.PaymasterInfo.Query().Aggregate(ent.Sum(paymasterinfo.FieldUserOpsNum)).Ints(context.Background())

	info := paymasterInfos[0]
	res.SponsorGasFeeUsdTotal = info.GasSponsoredUsd.RoundDown(2)
	res.SponsorGasFeeUsd24h = info.GasSponsoredUsdD1.RoundDown(2)
	res.UserOpsNumTotal = info.UserOpsNum
	res.UserOpsNum24h = info.UserOpsNumD1
	allNum := allUserOpsNum[0]
	if allNum != 0 {
		res.Dominance = getRate(int64(allNum), res.UserOpsNumTotal)
	}
	highCount, err := client.PaymasterInfo.Query().Where(paymasterinfo.UserOpsNumGT(res.UserOpsNumTotal)).Count(context.Background())
	totalBalance := service.GetTotalBalance(paymaster, req.Network)
	res.Rank = highCount + 1
	res.AccountBalance = totalBalance

	addresses, _ := dao.AccountDao.GetAccountByAddresses(ctx, client, []string{req.Paymaster})
	if len(addresses) > 0 {
		addresses[0].Label.AssignTo(&res.Label)
	}

	res.TotalNumber = int64(
		client.PaymasterInfo.Query().CountX(ctx),
	)
	return &res, nil
}
