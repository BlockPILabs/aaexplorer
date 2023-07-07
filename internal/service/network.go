package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/log"
)

type networkService struct {
}

var NetworkService = networkService{}

func (networkService) GetNetworks(ctx context.Context) ([]*ent.Network, error) {
	ctx, _ = log.With(ctx, "service", "network")
	return dao.NetworkDao.GetNetworks(ctx)
}
