package dao

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/network"
	"github.com/BlockPILabs/aaexplorer/internal/log"
)

type networkDao struct {
}

type networkCtxKey struct {
}

var NetworkDao = &networkDao{}

func (*networkDao) GetNetworks(ctx context.Context, tx *ent.Client) (ent.Networks, error) {

	ctx, logger := log.With(ctx, "module", "network")

	networks, err := tx.Network.Query().Where(
		network.DeleteTimeIsNil(),
	).All(ctx)
	if err != nil {
		logger.Warn("get networks error", "err", err)
		return nil, err
	}
	return networks, err
}

func (*networkDao) GetNetworkByNetwork(ctx context.Context, network_ string) (*ent.Network, error) {

	ctx, logger := log.With(ctx, "module", "network")
	db, err := entity.Client(ctx)
	if err != nil {
		return nil, err
	}
	net, err := db.Network.Query().Where(
		network.IDEQ(network_),
		network.DeleteTimeIsNil(),
	).First(ctx)
	if err != nil {
		logger.Warn("get networks error", "err", err)
		return nil, err
	}
	return net, err
}

func (*networkDao) GetNetworksByNetworks(ctx context.Context, tx *ent.Client, networks ...string) (ent.Networks, error) {

	ctx, logger := log.With(ctx, "module", "network")
	db, err := entity.Client(ctx)
	if err != nil {
		return nil, err
	}
	nets, err := db.Network.Query().Where(
		network.IDIn(networks...),
		network.DeleteTimeIsNil(),
	).All(ctx)
	if err != nil {
		logger.Warn("get networks error", "err", err)
		return nil, err
	}
	return nets, err
}

func (*networkDao) WithContext(ctx context.Context, net *ent.Network) context.Context {
	if net == nil {
		return ctx
	}
	return context.WithValue(ctx, networkCtxKey{}, net)
}

func (*networkDao) ContextValue(ctx context.Context) (*ent.Network, bool) {
	v := ctx.Value(networkCtxKey{})
	if v == nil {
		return nil, false
	}
	net, ok := v.(*ent.Network)
	return net, ok
}
