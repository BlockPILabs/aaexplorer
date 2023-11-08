package dao

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/tokenpriceinfo"
)

type tokenPriceInfoDao struct {
	baseDao
}

var TokenPriceInfoDao = &tokenPriceInfoDao{}

func (*tokenPriceInfoDao) GetBaseTokenPrice(ctx context.Context, tx *ent.Client) (block *ent.TokenPriceInfo, err error) {
	block, err = tx.TokenPriceInfo.Query().Where(
		tokenpriceinfo.TypeEQ(config.TokenTypeBase),
	).First(ctx)
	return
}
