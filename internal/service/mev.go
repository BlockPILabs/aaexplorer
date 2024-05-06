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

	res = &vo.ListMEVBundlersResponse{}
	list, total, err := dao.MevDao.Pagination(ctx, client, req)
	if err != nil {
		return nil, err
	}

	res.TotalCount = total

	for _, mev := range list {
		res.Records = append(res.Records, &vo.MEVBundlerAssets{
			Timestamp:   mev.TxTime,
			UserOpHash:  mev.TxHash,
			MevType:     mev.TxFrom, // undo
			Victim:      mev.TxFrom,
			Attacker:    mev.TxFrom,
			BundlerLoss: *mev.GasFee,
			MevProfits:  *mev.Profit,
		})
	}

	return res, nil
}
