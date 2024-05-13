package service

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/internal/dao"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
)

type mevService struct {
}

var MevService = &mevService{}

func (*mevService) MevList(ctx context.Context, req vo.ListMEVBundlersRequest) (res *vo.ListMEVBundlersResponse, err error) {
	client, err := entity.Client(ctx, req.Network)
	if err != nil {
		return nil, err
	}

	res = &vo.ListMEVBundlersResponse{
		Pagination: vo.Pagination{
			TotalCount: 0,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
	}
	list, total, err := dao.MevDao.Pagination(ctx, client, req)
	if err != nil {
		return nil, err
	}

	res.TotalCount = total

	for _, mev := range list {
		res.Records = append(res.Records, &vo.MEVBundlerAssets{
			Timestamp:        mev.Time,
			UserOpHash:       mev.VictimTxHash,
			MevType:          mev.MevType,
			Victim:           mev.Victim,
			Attacker:         mev.Attacker,
			BundlerLoss:      mev.BundlerLoss,
			BundlerLossInUsd: mev.BundlerLossUsd,
			MevProfits:       mev.MevProfit,
			MevProfitsInUsd:  mev.MevProfitUsd,
		})
	}

	return res, nil
}

func (*mevService) BlockMevList(ctx context.Context, req vo.ListBlockMEVBundlersRequest) (res *vo.ListMEVBundlersResponse, err error) {
	client, err := entity.Client(ctx, req.Network)
	if err != nil {
		return nil, err
	}
	res = &vo.ListMEVBundlersResponse{
		Pagination: vo.Pagination{
			TotalCount: 0,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
	}
	list, total, err := dao.MevDao.BlockMevPagination(ctx, client, req)
	if err != nil {
		return nil, err
	}

	res.TotalCount = total

	for _, mev := range list {
		res.Records = append(res.Records, &vo.MEVBundlerAssets{
			Timestamp:        mev.Time,
			UserOpHash:       mev.VictimTxHash,
			MevType:          mev.MevType,
			Victim:           mev.Victim,
			Attacker:         mev.Attacker,
			BundlerLoss:      mev.BundlerLoss,
			BundlerLossInUsd: mev.BundlerLossUsd,
			MevProfits:       mev.MevProfit,
			MevProfitsInUsd:  mev.MevProfitUsd,
		})
	}

	return res, nil
}

func (*mevService) BundlerMevList(ctx context.Context, req vo.ListBundlerMEVBundlersRequest) (res *vo.ListMEVBundlersResponse, err error) {
	client, err := entity.Client(ctx, req.Network)
	if err != nil {
		return nil, err
	}
	res = &vo.ListMEVBundlersResponse{
		Pagination: vo.Pagination{
			TotalCount: 0,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
	}
	list, total, err := dao.MevDao.BundlerMevPagination(ctx, client, req)
	if err != nil {
		return nil, err
	}

	res.TotalCount = total

	for _, mev := range list {
		res.Records = append(res.Records, &vo.MEVBundlerAssets{
			Timestamp:        mev.Time,
			UserOpHash:       mev.VictimTxHash,
			MevType:          mev.MevType,
			Victim:           mev.Victim,
			Attacker:         mev.Attacker,
			BundlerLoss:      mev.BundlerLoss,
			BundlerLossInUsd: mev.BundlerLossUsd,
			MevProfits:       mev.MevProfit,
			MevProfitsInUsd:  mev.MevProfitUsd,
		})
	}

	return res, nil
}
