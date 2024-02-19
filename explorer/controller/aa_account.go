package controller

import (
	"github.com/BlockPILabs/aaexplorer/internal/dao"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/service"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetAaAcountInfo = "GetAaAccountInfo"

func GetAaAccountInfo(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start GetAaAccountInfo")
	req := vo.AaAccountRequestVo{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
	}
	client, err := entity.Client(ctx, req.Network)
	if err != nil {
		return err
	}
	res, err := service.AaAccountService.GetAaAccountRecord(ctx, client, req)
	if err != nil {
		logger.Error("GetAaAccountInfo error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)

}

const NameGetAaAccountChains = "GetAaAccountChains"

func GetAaAccountChains(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start GetAaAccountChains")

	req := vo.AaAccountNetworkRequestVo{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "addr", req.Address)
	}

	networks, err := service.NetworkService.GetNetworks(ctx)
	if err != nil {
		return err
	}

	var ret []string
	for _, network := range networks {
		client, err := entity.NetworkClient(ctx, network)
		if err != nil {
			return err
		}
		exists := dao.AaAccountDao.AaAccountExists(ctx, client, req.Address)
		if exists {
			ret = append(ret, network.ID)
		}
	}
	result := vo.AaAccountNetworkResponseVo{Chains: ret}
	return vo.NewResultJsonResponse(result).JSON(fcx)
}
