package controller

import (
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/service"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
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
	req := vo.TopBundlerRequest{
		PaginationRequest: vo.NewDefaultPaginationRequest(),
	}
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
	req := vo.TopPaymasterRequest{
		PaginationRequest: vo.NewDefaultPaginationRequest(),
	}
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
	req := vo.TopFactoryRequest{
		PaginationRequest: vo.NewDefaultPaginationRequest(),
	}
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
