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

	client, err := entity.Client(ctx)
	if err != nil {
		return nil, err
	}
	//
	list, total, err := dao.UserOpDao.Pagination(ctx, client, req.Network, req.PaginationRequest)
	if err != nil {
		return nil, err
	}
	res.TotalCount = total

	//
	res.Records = make([]*vo.UserOpVo, len(list))
	for i, info := range list {
		res.Records[i] = &vo.UserOpVo{
			ID:                   info.ID,
			UserOperationHash:    info.UserOperationHash,
			TxHash:               info.TxHash,
			BlockNumber:          info.BlockNumber,
			Network:              info.Network,
			Sender:               info.Sender,
			Target:               info.Target,
			TxValue:              info.TxValue,
			Fee:                  info.Fee,
			Bundler:              info.Bundler,
			EntryPoint:           info.EntryPoint,
			Factory:              info.Factory,
			Paymaster:            info.Paymaster,
			PaymasterAndData:     info.PaymasterAndData,
			Signature:            info.Signature,
			Calldata:             info.Calldata,
			Nonce:                info.Nonce,
			CallGasLimit:         info.CallGasLimit,
			PreVerificationGas:   info.PreVerificationGas,
			VerificationGasLimit: info.VerificationGasLimit,
			MaxFeePerGas:         info.MaxFeePerGas,
			MaxPriorityFeePerGas: info.MaxPriorityFeePerGas,
			TxTime:               info.TxTime,
			TxTimeFormat:         info.TxTimeFormat,
			InitCode:             info.InitCode,
			Status:               info.Status,
			Source:               info.Source,
			ActualGasCost:        info.ActualGasCost,
			ActualGasUsed:        info.ActualGasUsed,
			CreateTime:           info.CreateTime,
		}
	}

	return &res, nil
}
