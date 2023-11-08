package controller

import (
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/service"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
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
	res := &vo.GetBundlersResponse{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
		return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(vo.ErrParams)).JSON(fcx)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
		return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(vo.ErrParams)).JSON(fcx)
	}
	res, err = service.BundlerService.GetBundlers(ctx, req)
	if err != nil {
		logger.Error("get bundlers error", "err", err)
	}
	return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(err)).JSON(fcx)
}

const NameGetBundler = "get_bundler"

func GetBundler(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(ctx)

	logger.Debug("start get bundler", "")
	req := vo.GetBundlerRequest{}
	res := &vo.GetBundlerResponse{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
		return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(vo.ErrParams)).JSON(fcx)
	}

	ctx, logger = log.With(ctx, "bundler", req.Bundler, "network", req.Network)

	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err)
		return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(vo.ErrParams)).JSON(fcx)
	}
	res, err = service.BundlerService.GetBundler(ctx, req)
	return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(err)).JSON(fcx)
}
