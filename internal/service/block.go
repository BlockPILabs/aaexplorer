package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type blockService struct {
}

var BlockService = &blockService{}

func (*blockService) GetBlocks(ctx context.Context, req vo.GetBlocksRequest) (*vo.GetBlocksResponse, error) {
	ctx, logger := log.With(ctx, "service", "GetBlocks")
	err := vo.ValidateStruct(req)
	res := vo.GetBlocksResponse{
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
		return &res, err
	}

	list, total, err := dao.BlockDao.Pagination(ctx, client, req.Network, req.PaginationRequest)
	if err != nil {
		return &res, err
	}
	res.TotalCount = total

	//
	res.Records = make([]*vo.BlocksVo, len(list))
	for i, info := range list {
		res.Records[i] = &vo.BlocksVo{
			Time:             info.Time,
			BlockNum:         info.ID,
			CreateTime:       info.CreateTime,
			Hash:             info.Hash,
			Size:             info.Size,
			Miner:            info.Miner,
			Nonce:            info.Nonce,
			Number:           info.Number,
			Uncles:           info.Uncles,
			GasUsed:          info.GasUsed,
			MixHash:          info.MixHash,
			GasLimit:         info.GasLimit,
			ExtraData:        info.ExtraData,
			LogsBloom:        info.LogsBloom,
			StateRoot:        info.StateRoot,
			Timestamp:        info.Timestamp,
			Difficulty:       info.Difficulty,
			ParentHash:       info.ParentHash,
			Sha3Uncles:       info.Sha3Uncles,
			Withdrawals:      info.Withdrawals,
			ReceiptsRoot:     info.ReceiptsRoot,
			BaseFeePerGas:    info.BaseFeePerGas,
			TotalDifficulty:  info.TotalDifficulty,
			WithdrawalsRoot:  info.WithdrawalsRoot,
			TransactionsRoot: info.TransactionsRoot,
		}
	}

	return &res, nil
}
func (*blockService) GetBlock(ctx context.Context, req vo.GetBlockRequest) (*vo.GetBlockResponse, error) {
	ctx, logger := log.With(ctx, "service", "GetBlock")
	err := vo.ValidateStruct(req)
	if err != nil {
		logger.Error("params error", "req", req, "err", err.Error())
		return nil, vo.ErrParams.SetData(err)
	}
	res := vo.GetBlockResponse{
		Block: nil,
	}

	client, err := entity.Client(ctx)
	if err != nil {
		return &res, err
	}
	block, err := dao.BlockDao.GetBlock(ctx, client, req.Block)
	if err != nil {
		return nil, err
	}
	res.Block = &vo.BlocksVo{
		BlockNum:   block.ID,
		CreateTime: block.CreateTime,
		Hash:       block.Hash,
	}

	return &res, nil
}
