package controller

import (
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/service"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const (
	NameGetTopBundler   = "get_top_bundler"
	NameGetTopPaymaster = "get_top_paymaster"
	NameGetTopFactory   = "get_top_factory"
)

func GetTopBundler(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get top bundler")
	req := vo.TopBundlerRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
	}

	res, err := service.GetTopBundler(ctx, req)
	if err != nil {
		logger.Error("get top bundler error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}

func GetTopPaymaster(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get top paymaster")
	req := vo.TopPaymasterRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
	}

	res, err := service.GetTopPaymaster(ctx, req)
	if err != nil {
		logger.Error("get top paymaster error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}

func GetTopFactory(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get top factory")
	req := vo.TopFactoryRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
	}

	res, err := service.GetTopFactory(ctx, req)
	if err != nil {
		logger.Error("get top factory error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}
