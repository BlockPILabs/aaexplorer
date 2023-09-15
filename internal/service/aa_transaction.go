package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

type aaTransactionService struct {
}

var AaTransactionService = &aaTransactionService{}

func (*aaTransactionService) GetRecord(ctx context.Context, client *ent.Client, req vo.AaTransactionRequestVo) (*vo.AaTransactionRecord, error) {
	ctx, logger := log.With(ctx, "service", "GetRecord")
	err := vo.ValidateStruct(req)
	if err != nil {
		logger.Error("params error", "req", req, "err", err.Error())
		return nil, vo.ErrParams.SetData(err)
	}

	list, _, err := dao.AaTransactionDao.Pages(ctx, client, vo.PaginationRequest{
		PerPage: 1,
		Page:    1,
	}, dao.AaTransactionCondition{
		TxHash: &req.TxHash,
	})
	if err != nil {
		return nil, err
	}
	if len(list) != 1 {
		return nil, nil
	}
	record := list[0]

	return &vo.AaTransactionRecord{
		HASH:                     record.HASH,
		TIME:                     record.TIME,
		BLOCK_HASH:               record.BLOCK_HASH,
		BLOCK_NUMBER:             record.BLOCK_NUMBER,
		USEROP_COUNT:             record.USEROP_COUNT,
		IS_MEV:                   record.IS_MEV,
		BUNDLER_PROFIT:           record.BUNDLER_PROFIT,
		NONCE:                    record.NONCE,
		TRANSACTION_INDEX:        record.TRANSACTION_INDEX,
		FROM_ADDR:                record.FROM_ADDR,
		TO_ADDR:                  record.TO_ADDR,
		VALUE:                    record.VALUE,
		GAS_PRICE:                record.GAS_PRICE,
		GAS:                      record.GAS,
		INPUT:                    record.INPUT,
		R:                        record.R,
		S:                        record.S,
		V:                        record.V,
		CHAIN_ID:                 record.CHAIN_ID,
		TYPE:                     record.TYPE,
		MAX_FEE_PER_GAS:          record.MAX_FEE_PER_GAS,
		MAX_PRIORITY_FEE_PER_GAS: record.MAX_PRIORITY_FEE_PER_GAS,
		ACCESS_LIST:              record.ACCESS_LIST,
	}, nil

}

func (*aaTransactionService) GetPages(ctx context.Context, client *ent.Client, req vo.AaTransactionListRequestVo) (*vo.AaTransactionListResponse, error) {
	ctx, logger := log.With(ctx, "service", "GetPages")
	err := vo.ValidateStruct(req)
	if err != nil {
		logger.Error("params error", "req", req, "err", err.Error())
		return nil, vo.ErrParams.SetData(err)
	}

	res := vo.AaTransactionListResponse{
		Pagination: vo.Pagination{
			TotalCount: 0,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
	}

	condition := dao.AaTransactionCondition{}
	if req.TxHash != "" {
		condition.TxHash = &req.TxHash
	}
	userOpsList, total, err := dao.AaTransactionDao.Pages(ctx, client, req.PaginationRequest, condition)
	if err != nil {
		return nil, err
	}
	res.TotalCount = total
	for _, record := range userOpsList {
		res.Records = append(res.Records, &vo.AaTransactionRecord{
			HASH:                     record.HASH,
			TIME:                     record.TIME,
			BLOCK_HASH:               record.BLOCK_HASH,
			BLOCK_NUMBER:             record.BLOCK_NUMBER,
			USEROP_COUNT:             record.USEROP_COUNT,
			IS_MEV:                   record.IS_MEV,
			BUNDLER_PROFIT:           record.BUNDLER_PROFIT,
			NONCE:                    record.NONCE,
			TRANSACTION_INDEX:        record.TRANSACTION_INDEX,
			FROM_ADDR:                record.FROM_ADDR,
			TO_ADDR:                  record.TO_ADDR,
			VALUE:                    record.VALUE,
			GAS_PRICE:                record.GAS_PRICE,
			GAS:                      record.GAS,
			INPUT:                    record.INPUT,
			R:                        record.R,
			S:                        record.S,
			V:                        record.V,
			CHAIN_ID:                 record.CHAIN_ID,
			TYPE:                     record.TYPE,
			MAX_FEE_PER_GAS:          record.MAX_FEE_PER_GAS,
			MAX_PRIORITY_FEE_PER_GAS: record.MAX_PRIORITY_FEE_PER_GAS,
			ACCESS_LIST:              record.ACCESS_LIST,
		})
	}

	return &res, nil

}
