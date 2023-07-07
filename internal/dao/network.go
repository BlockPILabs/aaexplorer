package dao

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/network"
	"github.com/BlockPILabs/aa-scan/internal/log"
)

type networkDao struct {
}

var NetworkDao = networkDao{}

func (networkDao) GetNetworks(ctx context.Context) ([]*ent.Network, error) {

	ctx, logger := log.With(ctx, "module", "network")
	db, err := entity.Client(ctx)
	if err != nil {
		return nil, err
	}
	networks, err := db.Network.Query().Where(
		network.DeleteTimeIsNil(),
	).All(ctx)
	if err != nil {
		logger.Warn("get networks error")
		return nil, err
	}
	return networks, err
}
