package controller

import (
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/service"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetBundlers = "get_bundlers"

func GetBundlers(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get bundlers")
	req := vo.GetBundlersRequest{
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
	res, err := service.BundlerService.GetBundlers(ctx, req)
	if err != nil {
		logger.Error("get bundlers error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}

const NameGetBundler = "get_bundler"

func GetBundler(fcx *fiber.Ctx) error {
	return vo.NewResultJsonResponse(fcx.Params("bundler")).JSON(fcx)
}
