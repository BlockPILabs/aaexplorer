package controller

import (
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/service"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/gofiber/fiber/v2"
)

func GetNetworks(fcx *fiber.Ctx) error {

	ctx := fcx.UserContext()

	log.Context(ctx).Debug("start get networks success")
	res := &vo.GetNetworksResponse{
		Pagination: vo.Pagination{
			Page:       1,
			TotalCount: 0,
		},
		Records: make([]*vo.NetworkVo, 0),
	}

	networks, err := service.NetworkService.GetNetworks(ctx)
	if err != nil {
		log.Context(ctx).Warn("get networks error", "err", err)
		return vo.NewResultJsonResponse(res).JSON(fcx)
	}

	// transfer to vo
	res.TotalCount = int64(len(networks))
	res.Records = make([]*vo.NetworkVo, int(res.TotalCount))

	for i, network := range networks {
		res.Records[i] = &vo.NetworkVo{
			Name:      network.Name,
			Network:   network.Network,
			Logo:      network.Logo,
			IsTestnet: network.IsTestnet,
		}
	}
	log.Context(ctx).Debug("get networks success", "totalCount", res.TotalCount)
	return vo.NewResultJsonResponse(res).JSON(fcx)
}
