package controller

import (
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/service"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetBundles = "get_bundles"

func GetBundles(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get bundles")
	req := vo.GetBundlesRequest{
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
	res, err := service.BundleService.GetBundles(ctx, req)
	if err != nil {
		logger.Error("get bundles error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}

const NameGetBundle = "get_bundle"

func GetBundle(fcx *fiber.Ctx) error {
	return vo.NewResultJsonResponse(fcx.Params("bundle")).JSON(fcx)
}
