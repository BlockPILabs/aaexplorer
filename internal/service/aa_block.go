package service

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/internal/dao"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/mevtransaction"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/shopspring/decimal"
)

type aaBlockService struct {
}

var AaBlockService = &aaBlockService{}

func (*aaBlockService) GetAaBlockInfo(ctx context.Context, client *ent.Client, req vo.GetAaBlocksRequest) (*vo.GetAaBlocksResponse, error) {
	ctx, logger := log.With(ctx, "service", "GetAaBlockInfo")
	logger.Info("GetAaBlockInfo ... ")
	res := vo.GetAaBlocksResponse{
		Pagination: vo.Pagination{
			TotalCount: 0,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
	}

	pages, total, err := dao.AaBlockDao.Pages(ctx, client, req.PaginationRequest, dao.AaBlockPagesCondition{LatestBlockNumber: req.LatestBlockNumber})
	if err != nil {
		return nil, err
	}

	res.TotalCount = total

	res.Records = make([]*vo.AaBlocksVo, len(pages))
	for i, info := range pages {
		mevTxs, _ := client.MevTransaction.Query().Where(mevtransaction.BlockNumberEQ(info.ID)).All(ctx)
		bundlerLoss := decimal.Zero
		bundlerLossUsd := decimal.Zero
		mevProfit := decimal.Zero
		mevProfitUsd := decimal.Zero
		if len(mevTxs) > 0 {
			for _, tx := range mevTxs {
				bundlerLoss = bundlerLoss.Add(tx.BundlerLoss)
				bundlerLossUsd = bundlerLossUsd.Add(tx.BundlerLossUsd)
				mevProfit = mevProfit.Add(tx.MevProfit)
				mevProfitUsd = mevProfitUsd.Add(tx.MevProfitUsd)
			}
		}
		res.Records[i] = &vo.AaBlocksVo{
			Number:           info.ID,
			Time:             info.Time.UnixMilli(),
			Hash:             info.Hash,
			UseropCount:      info.UseropCount,
			UseropMevCount:   info.UseropMevCount,
			BundlerProfit:    info.BundlerProfit,
			BundlerProfitUsd: info.BundlerProfitUsd,
			CreateTime:       info.CreateTime.UnixMilli(),
			MevCount:         int64(len(mevTxs)),
			MevProfits:       mevProfit,
			MevProfitsUsd:    mevProfitUsd,
			BundlerLoss:      bundlerLoss,
			BundlerLossUsd:   bundlerLossUsd,
		}

	}

	return &res, nil
}
