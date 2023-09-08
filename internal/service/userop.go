package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type userOpService struct {
}

var UserOpService = &userOpService{}

func (*userOpService) GetUserOps(ctx context.Context, req vo.GetUserOpsRequest) (*vo.GetUserOpsResponse, error) {
	ctx, logger := log.With(ctx, "service", "GetUserOps")
	err := vo.ValidateStruct(req)
	res := vo.GetUserOpsResponse{
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
	list, total, err := dao.UserOpDao.Pagination(ctx, client, req)
	if err != nil {
		return nil, err
	}
	res.TotalCount = total
	res.Records = make([]*vo.UserOpVo, len(list))
	if res.TotalCount > 0 {

		var lists = map[string][]string{}
		var userOpsHashIn = []string{}
		for _, info := range list {
			lists[info.ID] = []string{info.Target}
			if info.TargetsCount > 0 {
				userOpsHashIn = append(userOpsHashIn, info.ID)
			}
		}
		if len(userOpsHashIn) > 0 {
			getTargets, _ := dao.UserOpCallDataDao.GetTargets(ctx, client, userOpsHashIn)
			for id, targets := range getTargets {
				lists[id] = targets
			}
		}

		//
		for i, info := range list {
			res.Records[i] = &vo.UserOpVo{
				Time:              info.Time.Unix(),
				UserOperationHash: info.ID,
				TxHash:            info.TxHash,
				BlockNumber:       info.BlockNumber,
				Network:           info.Network,
				Sender:            info.Sender,
				Target:            info.Target,
				TxValue:           info.TxValue,
				Fee:               info.Fee,
				InitCode:          info.InitCode,
				Status:            info.Status,
				Source:            info.Source,
				Targets:           lists[info.ID],
				TargetsCount:      info.TargetsCount,
			}
		}
	}

	return &res, nil
}
