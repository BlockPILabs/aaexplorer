package controller

import (
	"github.com/BlockPILabs/aaexplorer/internal/dao"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/service"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetNetworks = "get_networks"

func GetNetworks(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()

	log.Context(ctx).Debug("start get networks")
	res := &vo.GetNetworksResponse{
		Pagination: vo.Pagination{
			Page:       1,
			TotalCount: 0,
			PerPage:    0,
		},
		Records: make([]*vo.NetworkVo, 0),
	}

	networks, err := service.NetworkService.GetNetworks(ctx)
	if err != nil {
		log.Context(ctx).Warn("get networks error", "err", err)
		return vo.NewResultJsonResponse(res).JSON(fcx)
	}

	// transfer to vo
	res.TotalCount = len(networks)
	res.Records = make([]*vo.NetworkVo, int(res.TotalCount))

	for i, network := range networks {
		res.Records[i] = &vo.NetworkVo{
			Name:        network.Name,
			ChainName:   network.ChainName,
			Network:     network.ID,
			ChainID:     network.ChainID,
			IsTestnet:   network.IsTestnet,
			Scan:        network.Scan,
			ScanTx:      network.ScanTx,
			ScanBlock:   network.ScanBlock,
			ScanAddress: network.ScanAddress,
			ScanName:    network.ScanName,
		}
	}
	log.Context(ctx).Debug("get networks success", "totalCount", res.TotalCount)
	return vo.NewResultJsonResponse(res).JSON(fcx)
}
func GetNetwork(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	network, _ := dao.NetworkDao.ContextValue(ctx)
	return vo.NewResultJsonResponse(vo.NetworkVo{
		Name:        network.Name,
		ChainName:   network.ChainName,
		Network:     network.ID,
		ChainID:     network.ChainID,
		IsTestnet:   network.IsTestnet,
		Scan:        network.Scan,
		ScanTx:      network.ScanTx,
		ScanBlock:   network.ScanBlock,
		ScanAddress: network.ScanAddress,
		ScanName:    network.ScanName,
	}).JSON(fcx)
}

// NetworkMiddleware check network params
func NetworkMiddleware() fiber.Handler {
	return func(fcx *fiber.Ctx) error {
		ctx := fcx.UserContext()
		networkFlag := fcx.Params("network")

		nw, err := dao.NetworkDao.GetNetworkByNetwork(ctx, networkFlag)
		if err != nil {
			return vo.ErrNetworkNotFound
		}
		// set value
		fcx.SetUserContext(
			dao.NetworkDao.WithContext(ctx, nw),
		)
		return fcx.Next()
	}
}
