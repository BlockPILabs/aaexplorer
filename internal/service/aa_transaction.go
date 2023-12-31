package service

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/internal/dao"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
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

	/*	txlist, _, err := dao.TransactionDao.Pages(ctx, client, vo.PaginationRequest{
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
		txr := txrlist[0]*/

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
		Nonce:                aatx.Nonce,
		TransactionIndex:     aatx.TransactionIndex,
		FromAddr:             aatx.FromAddr,
		ToAddr:               aatx.ToAddr,
		Value:                aatx.Value,
		GasPrice:             aatx.GasPrice,
		Gas:                  aatx.Gas,
		Input:                aatx.Input,
		R:                    aatx.R,
		S:                    aatx.S,
		V:                    aatx.V,
		ChainID:              aatx.ChainID,
		Type:                 aatx.Type,
		MaxFeePerGas:         aatx.MaxFeePerGas,
		MaxPriorityFeePerGas: aatx.MaxPriorityFeePerGas,
		AccessList:           aatx.AccessList,
		Method:               *aatx.Method,
		ContractAddress:      *aatx.ContractAddress,
		CumulativeGasUsed:    *aatx.CumulativeGasUsed,
		EffectiveGasPrice:    *aatx.EffectiveGasPrice,
		GasUsed:              *aatx.GasUsed,
		//Logs:                 txr.Logs,
		//LogsBloom:            txr.LogsBloom,
		Status: *aatx.Status,

		TokenPriceUsd: tokenPrice.TokenPrice,
		GasPriceUsd:   aatx.GasPrice.Mul(tokenPrice.TokenPrice).Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))),
		ValueUsd:      aatx.Value.Mul(tokenPrice.TokenPrice).Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))),
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
			TotalCount: 10,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
	}

	condition := dao.AaTransactionCondition{}
	if req.TxHash != "" {
		condition.TxHash = &req.TxHash
	}

	if req.Address != "" {
		condition.Address = &req.Address
	}

	list, total, err := dao.AaTransactionDao.Pagination(ctx, client, req.PaginationRequest, condition)
	if err != nil {
		return nil, err
	}
	res.TotalCount = total
	for _, record := range list {
		//txlist, _, err := dao.TransactionDao.Pages(ctx, client, vo.PaginationRequest{
		//	PerPage: 1,
		//	Page:    1,
		//}, dao.TransactionCondition{
		//	TxHash: &req.TxHash,
		//})
		//if err != nil {
		//	return nil, err
		//}
		//tx := ent.TransactionDecode{}
		//if len(txlist) < 1 {
		//	tx = *txlist[0]
		//}
		//aa := record.Edges.Txaa
		//if aa == nil {
		//	aa = &ent.AaTransactionInfo{}
		//}

		ret := &vo.AaTransactionRecord{
			Hash:                 record.ID,
			Time:                 record.Time.UnixMilli(),
			BlockHash:            record.BlockHash,
			BlockNumber:          record.BlockNumber,
			UseropCount:          record.UseropCount,
			IsMev:                record.IsMev,
			BundlerProfit:        record.BundlerProfit,
			BundlerProfitUsd:     record.BundlerProfitUsd,
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
			Method:               *record.Method,
			ContractAddress:      *record.ContractAddress,
			CumulativeGasUsed:    *record.CumulativeGasUsed,
			EffectiveGasPrice:    *record.EffectiveGasPrice,
			GasUsed:              *record.GasUsed,
			Status:               *record.Status,
		}
		res.Records = append(res.Records, ret)
	}

	return &res, nil

}
