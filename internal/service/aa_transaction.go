package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/shopspring/decimal"
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

	aatxlist, _, err := dao.AaTransactionDao.Pagination(ctx, client, vo.PaginationRequest{
		PerPage: 1,
		Page:    1,
	}, dao.AaTransactionCondition{
		TxHash: &req.TxHash,
	})
	if err != nil {
		return nil, err
	}
	if len(aatxlist) != 1 {
		return nil, nil
	}
	aatx := aatxlist[0]

	txlist, _, err := dao.TransactionDao.Pages(ctx, client, vo.PaginationRequest{
		PerPage: 1,
		Page:    1,
	}, dao.TransactionCondition{
		TxHash: &req.TxHash,
	})

	if err != nil {
		return nil, err
	}
	if len(txlist) != 1 {
		return nil, nil
	}
	tx := txlist[0]

	txrlist, _, err := dao.TransactionReceiptDao.Pages(ctx, client, vo.PaginationRequest{
		PerPage: 1,
		Page:    1,
	}, dao.TransactionReceiptCondition{
		TxHash: &req.TxHash,
	})
	if err != nil {
		return nil, err
	}
	if len(txrlist) != 1 {
		return nil, nil
	}
	txr := txrlist[0]

	tokenPrice, err := dao.TokenPriceInfoDao.GetBaseTokenPrice(ctx, client)
	if err != nil {
		return nil, err
	}

	ret := &vo.AaTransactionRecord{
		Hash:                 aatx.ID,
		Time:                 aatx.Time.UnixMilli(),
		BlockHash:            aatx.BlockHash,
		BlockNumber:          aatx.BlockNumber,
		UseropCount:          aatx.UseropCount,
		IsMev:                aatx.IsMev,
		BundlerProfit:        aatx.BundlerProfit,
		BundlerProfitUsd:     aatx.BundlerProfitUsd,
		Nonce:                tx.Nonce,
		TransactionIndex:     tx.TransactionIndex,
		FromAddr:             tx.FromAddr,
		ToAddr:               tx.ToAddr,
		Value:                tx.Value,
		GasPrice:             tx.GasPrice,
		Gas:                  tx.Gas,
		Input:                tx.Input,
		R:                    tx.R,
		S:                    tx.S,
		V:                    tx.V,
		ChainID:              tx.ChainID,
		Type:                 tx.Type,
		MaxFeePerGas:         tx.MaxFeePerGas,
		MaxPriorityFeePerGas: tx.MaxPriorityFeePerGas,
		AccessList:           tx.AccessList,
		Method:               tx.Method,
		ContractAddress:      txr.ContractAddress,
		CumulativeGasUsed:    txr.CumulativeGasUsed,
		EffectiveGasPrice:    txr.EffectiveGasPrice,
		GasUsed:              txr.GasUsed,
		//Logs:                 txr.Logs,
		//LogsBloom:            txr.LogsBloom,
		Status: txr.Status,

		TokenPriceUsd: tokenPrice.TokenPrice,
		GasPriceUsd:   tx.GasPrice.Mul(tokenPrice.TokenPrice).Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))),
		ValueUsd:      tx.Value.Mul(tokenPrice.TokenPrice).Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))),
	}

	return ret, nil

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
	userOpsList, total, err := dao.AaTransactionDao.Pagination(ctx, client, req.PaginationRequest, condition)
	if err != nil {
		return nil, err
	}
	res.TotalCount = total
	for _, record := range userOpsList {
		txlist, _, err := dao.TransactionDao.Pages(ctx, client, vo.PaginationRequest{
			PerPage: 1,
			Page:    1,
		}, dao.TransactionCondition{
			TxHash: &req.TxHash,
		})

		if err != nil {
			return nil, err
		}
		if len(txlist) != 1 {
			return nil, nil
		}
		tx := txlist[0]
		ret := &vo.AaTransactionRecord{
			Hash:                 record.ID,
			Time:                 record.Time.UnixMilli(),
			BlockHash:            record.BlockHash,
			BlockNumber:          record.BlockNumber,
			UseropCount:          record.UseropCount,
			IsMev:                record.IsMev,
			BundlerProfit:        record.BundlerProfit,
			BundlerProfitUsd:     record.BundlerProfitUsd,
			Nonce:                tx.Nonce,
			TransactionIndex:     tx.TransactionIndex,
			FromAddr:             tx.FromAddr,
			ToAddr:               tx.ToAddr,
			Value:                tx.Value,
			GasPrice:             tx.GasPrice,
			Gas:                  tx.Gas,
			Input:                tx.Input,
			R:                    tx.R,
			S:                    tx.S,
			V:                    tx.V,
			ChainID:              tx.ChainID,
			Type:                 tx.Type,
			MaxFeePerGas:         tx.MaxFeePerGas,
			MaxPriorityFeePerGas: tx.MaxPriorityFeePerGas,
			AccessList:           tx.AccessList,
			Method:               tx.Method,
		}
		res.Records = append(res.Records, ret)
	}

	return &res, nil

}
