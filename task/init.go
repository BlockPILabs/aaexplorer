package task

import (
	"context"
	interConfig "github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/dao"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	interlog "github.com/BlockPILabs/aaexplorer/internal/log"
)

var logger = interlog.L()

func SetLogger(lg interlog.Logger) {
	logger = lg
}

var config = interConfig.DefaultConfig()

func SetConfig(cfg *interConfig.Config) {
	config = cfg
}

func getTaskNetworks(ctx context.Context, tx *ent.Client) (ent.Networks, error) {
	if len(config.Task.Networks) > 0 {
		return dao.NetworkDao.GetNetworksByNetworks(ctx, tx, config.Task.Networks...)
	}
	return dao.NetworkDao.GetNetworks(ctx, tx)
}
