package dao

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/network"
)

type networkDao struct {
}

var NetworkDao = networkDao{}

func (networkDao) GetNetworks(ctx context.Context) ([]*ent.Network, error) {
	db, err := entity.Client()
	if err != nil {
		return nil, err
	}
	networks, err := db.Network.Query().Where(
		network.DeleteTimeIsNil(),
	).All(ctx)
	if err != nil {
		return nil, err
	}
	return networks, err
}
