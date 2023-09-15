package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"strings"
)

type userOpService struct {
}

var UserOpService = &userOpService{}

func (*userOpService) GetUserOps(ctx context.Context, req vo.GetUserOpsRequest) (*vo.GetUserOpsResponse, error) {
	ctx, logger := log.With(ctx, "service", "GetFactories")
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
		var accountsMap = map[string]struct{}{}
		var accounts = []string{}
		for _, info := range list {
			lists[info.ID] = []string{info.Target}
			if info.TargetsCount > 0 {
				if len(info.Targets.Elements) == 1 {
					lists[info.ID] = strings.Split(strings.Trim(strings.Repeat(info.Target+"-", info.TargetsCount), "-"), "-")
				} else {
					userOpsHashIn = append(userOpsHashIn, info.ID)
				}
			}

			if _, ok := accountsMap[info.Sender]; !ok {
				accountsMap[info.Sender] = struct{}{}
				accounts = append(accounts, info.Sender)
			}

			if _, ok := accountsMap[info.Paymaster]; !ok {
				accountsMap[info.Paymaster] = struct{}{}
				accounts = append(accounts, info.Paymaster)
			}

			if _, ok := accountsMap[info.Bundler]; !ok {
				accountsMap[info.Bundler] = struct{}{}
				accounts = append(accounts, info.Bundler)
			}

			//if _, ok := accountsMap[info.Factory]; !ok {
			//	accountsMap[info.Factory] = struct{}{}
			//	accounts = append(accounts, info.Factory)
			//}

		}
		if len(userOpsHashIn) > 0 {
			getTargets, _ := dao.UserOpCallDataDao.GetTargets(ctx, client, userOpsHashIn)
			for id, targets := range getTargets {
				lists[id] = targets
			}
		}

		var labelMap = map[string][]string{}
		if len(accounts) > 0 {
			addresses, err := dao.AccountDao.GetAccountByAddresses(ctx, client, accounts)
			if err == nil {
				for i, account := range addresses {
					var labels []string
					addresses[i].Label.Scan(&labels)
					if len(labels) > 0 {
						labelMap[account.ID] = labels
					}

				}
			}
		}

		//
		for i, info := range list {
			userOpVo := &vo.UserOpVo{
				Time:              info.Time.UnixMilli(),
				UserOperationHash: info.ID,
				TxHash:            info.TxHash,
				BlockNumber:       info.BlockNumber,
				Network:           info.Network,
				Sender:            info.Sender,
				SenderLabel:       "",
				Target:            info.Target,
				TargetLabel:       "",
				TxValue:           info.TxValue,
				Fee:               info.Fee,
				InitCode:          info.InitCode,
				Status:            info.Status,
				Source:            info.Source,
				Targets:           lists[info.ID],
				TargetsCount:      info.TargetsCount,
				Bundler:           info.Bundler,
				Paymaster:         info.Paymaster,
			}

			if a, ok := labelMap[userOpVo.Sender]; ok {
				userOpVo.SenderLabel = a[0]
			}
			if a, ok := labelMap[userOpVo.Target]; ok {
				userOpVo.TargetLabel = a[0]
			}
			if a, ok := labelMap[userOpVo.Bundler]; ok {
				userOpVo.BundlerLabel = a[0]
			}
			if a, ok := labelMap[userOpVo.Paymaster]; ok {
				userOpVo.Paymaster = a[0]
			}
			res.Records[i] = userOpVo
		}
	}

	return &res, nil
}

func (*userOpService) GetUserOpsAnalysis(ctx context.Context, client *ent.Client, req vo.UserOpsAnalysisRequestVo) (*vo.UserOpsAnalysisRecord, error) {
	ctx, logger := log.With(ctx, "service", "GetUserOpsAnalysis")
	err := vo.ValidateStruct(req)
	if err != nil {
		logger.Error("params error", "req", req, "err", err.Error())
		return nil, vo.ErrParams.SetData(err)
	}

	userOpsList, _, err := dao.UserOpDao.Pages(ctx, client, vo.PaginationRequest{
		PerPage: 1,
		Page:    1,
	}, dao.UserOpsCondition{
		UserOperationHash: &req.UserOperationHash,
	})
	if err != nil {
		return nil, err
	}
	if len(userOpsList) != 1 {
		return nil, nil
	}
	userOps := userOpsList[0]

	callDataList, _, err := dao.UserOpCallDataDao.Pages(ctx, client, vo.PaginationRequest{
		PerPage: 1,
		Page:    1,
	}, dao.UserOpsCallDataCondition{
		UserOperationHash: &req.UserOperationHash,
	})
	if err != nil {
		return nil, err
	}

	var callData []vo.CallDataInfo
	for _, info := range callDataList {
		data := vo.CallDataInfo{
			Time:        info.Time,
			UserOpsHash: info.UserOpsHash,
			TxHash:      info.TxHash,
			BlockNumber: info.BlockNumber,
			Network:     info.Network,
			Sender:      info.Sender,
			Target:      info.Target,
			TxValue:     *info.TxValue,
			Source:      info.Source,
			Calldata:    info.Calldata,
			TxTime:      info.TxTime,
			CreateTime:  info.CreateTime,
			UpdateTime:  info.UpdateTime,
			AaIndex:     info.AaIndex,
		}
		callData = append(callData, data)
	}

	return &vo.UserOpsAnalysisRecord{
		UserOperationHash:    userOps.ID,
		Time:                 userOps.Time,
		TxHash:               userOps.TxHash,
		BlockNumber:          userOps.BlockNumber,
		Network:              userOps.Network,
		Sender:               userOps.Sender,
		Target:               userOps.Target,
		Targets:              userOps.Targets,
		TargetsCount:         userOps.TargetsCount,
		TxValue:              userOps.TxValue,
		Fee:                  userOps.Fee,
		Bundler:              userOps.Bundler,
		EntryPoint:           userOps.EntryPoint,
		Factory:              userOps.Factory,
		Paymaster:            userOps.Paymaster,
		PaymasterAndData:     userOps.PaymasterAndData,
		Signature:            userOps.Signature,
		Calldata:             userOps.Calldata,
		CalldataContract:     userOps.CalldataContract,
		Nonce:                userOps.Nonce,
		CallGasLimit:         userOps.CallGasLimit,
		PreVerificationGas:   userOps.PreVerificationGas,
		VerificationGasLimit: userOps.VerificationGasLimit,
		MaxFeePerGas:         userOps.MaxFeePerGas,
		MaxPriorityFeePerGas: userOps.MaxPriorityFeePerGas,
		TxTime:               userOps.TxTime,
		InitCode:             userOps.InitCode,
		Status:               userOps.Status,
		Source:               userOps.Source,
		ActualGasCost:        userOps.ActualGasCost,
		ActualGasUsed:        userOps.ActualGasUsed,
		CreateTime:           userOps.CreateTime,
		UpdateTime:           userOps.UpdateTime,
		UsdAmount:            userOps.UsdAmount,
		CallData:             callData,
	}, nil

}

func (*userOpService) GetUserOpsAnalysisList(ctx context.Context, client *ent.Client, req vo.UserOpsAnalysisListRequestVo) (*vo.UserOpsAnalysisListResponse, error) {
	ctx, logger := log.With(ctx, "service", "GetUserOpsAnalysis")
	err := vo.ValidateStruct(req)
	if err != nil {
		logger.Error("params error", "req", req, "err", err.Error())
		return nil, vo.ErrParams.SetData(err)
	}

	res := vo.UserOpsAnalysisListResponse{
		Pagination: vo.Pagination{
			TotalCount: 0,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
	}
	userOpsList, total, err := dao.UserOpDao.Pages(ctx, client, req.PaginationRequest, dao.UserOpsCondition{
		TxHash: &req.TxHash,
	})
	if err != nil {
		return nil, err
	}
	res.TotalCount = total
	for _, info := range userOpsList {
		res.Records = append(res.Records, &vo.UserOpsAnalysisRecord{
			UserOperationHash:    info.ID,
			Time:                 info.Time,
			TxHash:               info.TxHash,
			BlockNumber:          info.BlockNumber,
			Network:              info.Network,
			Sender:               info.Sender,
			Target:               info.Target,
			Targets:              info.Targets,
			TargetsCount:         info.TargetsCount,
			TxValue:              info.TxValue,
			Fee:                  info.Fee,
			Bundler:              info.Bundler,
			EntryPoint:           info.EntryPoint,
			Factory:              info.Factory,
			Paymaster:            info.Paymaster,
			PaymasterAndData:     info.PaymasterAndData,
			Signature:            info.Signature,
			Calldata:             info.Calldata,
			CalldataContract:     info.CalldataContract,
			Nonce:                info.Nonce,
			CallGasLimit:         info.CallGasLimit,
			PreVerificationGas:   info.PreVerificationGas,
			VerificationGasLimit: info.VerificationGasLimit,
			MaxFeePerGas:         info.MaxFeePerGas,
			MaxPriorityFeePerGas: info.MaxPriorityFeePerGas,
			TxTime:               info.TxTime,
			InitCode:             info.InitCode,
			Status:               info.Status,
			Source:               info.Source,
			ActualGasCost:        info.ActualGasCost,
			ActualGasUsed:        info.ActualGasUsed,
			CreateTime:           info.CreateTime,
			UpdateTime:           info.UpdateTime,
			UsdAmount:            info.UsdAmount,
		})
	}

	return &res, nil

}
