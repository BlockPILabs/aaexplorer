package service

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
)

type networkService struct {
}

var NetworkService = networkService{}

func (networkService) GetNetworks(ctx context.Context) ([]*ent.Network, error) {

	return dao.NetworkDao.GetNetworks(ctx)
}
