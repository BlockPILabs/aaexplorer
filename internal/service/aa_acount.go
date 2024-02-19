package service

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/internal/dao"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
)

type aaAccountService struct {
}

var AaAccountService = &aaAccountService{}

func (s *aaAccountService) GetAaAccountRecord(ctx context.Context, client *ent.Client, vo vo.AaAccountRequestVo) (*vo.AaAccountRecord, error) {
	ctx, logger := log.With(ctx, "service", "GetAaAccountRecord")
	logger.Info("GetAaAccountRecord ... ")

	address := vo.Address
	record, err := dao.AaAccountDao.GetAaAccountRecord(ctx, client, address)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (s *aaAccountService) GetAccounts(ctx context.Context, req vo.GetAccountsRequest) (res *vo.GetAccountsResponse, err error) {
	ctx, logger := log.With(ctx, "service", "GetAccounts")
	logger.Info("GetAccounts ... ")
	res = &vo.GetAccountsResponse{
		Pagination: vo.Pagination{
			TotalCount: req.TotalCount,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
		Records: []*vo.AaAccountDataVo{},
	}

	client, err := entity.Client(ctx, req.Network)
	if err != nil {
		return res, err
	}
	aaAccountDatas, total, err := dao.AaAccountDao.Pagination(ctx, client, req)
	if err != nil {
		return res, err
	}
	res.TotalCount = total
	for _, info := range aaAccountDatas {

		label := ""
		if info.Edges.Account != nil && info.Edges.Account.Label != nil {
			labels := []string{}
			info.Edges.Account.Label.AssignTo(&labels)
			if len(labels) > 0 {
				label = labels[0]
			}
		}

		a := &vo.AaAccountDataVo{
			ID:              info.ID,
			AddressLabel:    label,
			AaType:          info.AaType,
			Factory:         info.Factory,
			FactoryTime:     info.FactoryTime.UnixMilli(),
			UserOpsNum:      info.UserOpsNum,
			TotalBalanceUsd: info.TotalBalanceUsd,
			UpdateTime:      info.UpdateTime.UnixMilli(),
		}

		res.Records = append(res.Records, a)
	}

	return res, nil
}
