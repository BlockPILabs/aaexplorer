package service

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/internal/dao"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
)

type transactionService struct {
}

var TransactionService = &transactionService{}

func (*transactionService) GetRecord(ctx context.Context, client *ent.Client, req vo.TransactionRequestVo) (*vo.TransactionRecord, error) {
	ctx, logger := log.With(ctx, "service", "GetRecord")
	err := vo.ValidateStruct(req)
	if err != nil {
		logger.Error("params error", "req", req, "err", err.Error())
		return nil, vo.ErrParams.SetData(err)
	}

	list, _, err := dao.TransactionDao.Pages(ctx, client, vo.PaginationRequest{
		PerPage: 1,
		Page:    1,
	}, dao.TransactionCondition{
		TxHash: &req.TxHash,
	})
	if err != nil {
		return nil, err
	}
	if len(list) != 1 {
		return nil, nil
	}
	record := list[0]

	ret := &vo.TransactionRecord{
		Hash:                 record.ID,
		Time:                 record.Time.UnixMilli(),
		CreateTime:           record.CreateTime.UnixMilli(),
		BlockHash:            record.BlockHash,
		BlockNumber:          record.BlockNumber,
		Nonce:                record.Nonce,
		TransactionIndex:     record.TransactionIndex,
		FromAddr:             record.FromAddr,
		ToAddr:               record.ToAddr,
		Value:                record.Value,
		GasPrice:             record.GasPrice,
		Gas:                  record.Gas,
		Input:                record.Input,
		R:                    record.R,
		S:                    record.S,
		V:                    record.V,
		ChainID:              record.ChainID,
		Type:                 record.Type,
		MaxFeePerGas:         record.MaxFeePerGas,
		MaxPriorityFeePerGas: record.MaxPriorityFeePerGas,
		AccessList:           record.AccessList,
	}

	return ret, nil

}

func (*transactionService) GetPages(ctx context.Context, client *ent.Client, req vo.TransactionListRequestVo) (*vo.TransactionListResponse, error) {
	ctx, logger := log.With(ctx, "service", "GetPages")
	err := vo.ValidateStruct(req)
	if err != nil {
		logger.Error("params error", "req", req, "err", err.Error())
		return nil, vo.ErrParams.SetData(err)
	}

	res := vo.TransactionListResponse{
		Pagination: vo.Pagination{
			TotalCount: 0,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
	}

	condition := dao.TransactionCondition{}
	if req.TxHash != "" {
		condition.TxHash = &req.TxHash
	}
	userOpsList, total, err := dao.TransactionDao.Pages(ctx, client, req.PaginationRequest, condition)
	if err != nil {
		return nil, err
	}
	res.TotalCount = total
	for _, record := range userOpsList {
		res.Records = append(res.Records, &vo.TransactionRecord{
			Hash:                 record.ID,
			Time:                 record.Time.UnixMilli(),
			CreateTime:           record.CreateTime.UnixMilli(),
			BlockHash:            record.BlockHash,
			BlockNumber:          record.BlockNumber,
			Nonce:                record.Nonce,
			TransactionIndex:     record.TransactionIndex,
			FromAddr:             record.FromAddr,
			ToAddr:               record.ToAddr,
			Value:                record.Value,
			GasPrice:             record.GasPrice,
			Gas:                  record.Gas,
			Input:                record.Input,
			R:                    record.R,
			S:                    record.S,
			V:                    record.V,
			ChainID:              record.ChainID,
			Type:                 record.Type,
			MaxFeePerGas:         record.MaxFeePerGas,
			MaxPriorityFeePerGas: record.MaxPriorityFeePerGas,
			AccessList:           record.AccessList,
		})
	}

	return &res, nil

}
