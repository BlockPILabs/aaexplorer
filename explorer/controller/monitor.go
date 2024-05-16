package controller

import (
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/service"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameAddMonitor = "add_monitor"
const NameRemoveMonitor = "remove_monitor"
const NameAssetDetail = "get_asset_detail"

func AddMonitor(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get add monitor")
	req := vo.AddMonitorRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.BodyParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
		return err
	}

	res, err := service.AddMonitor(ctx, req)
	if err != nil {
		logger.Error("get mev transaction error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}

func RemoveMonitor(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get mev transaction")
	req := vo.RemoveMonitorRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.BodyParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
		return err
	}

	res, err := service.RemoveMonitor(ctx, req)
	if err != nil {
		logger.Error("get mev transaction error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}

const NameListBundler = "list_mev_bundler"

func ListMEVBundlers(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(ctx)

	logger.Debug("start list mev bundlers", "")

	req := vo.ListMEVBundlersRequest{}
	res := &vo.ListMEVBundlersResponse{
		Pagination: vo.Pagination{
			TotalCount: 0,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
	}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
		return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(vo.ErrParams)).JSON(fcx)
	}

	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err)
		return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(vo.ErrParams)).JSON(fcx)
	}

	res, err = service.MevService.MevList(ctx, req)
	return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(err)).JSON(fcx)
}

const NameListWatchingAddress = "list_watching_address"

func ListWatchingAddress(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(ctx)
	logger.Debug("start list watching address", "")

	req := vo.ListWatchingAddressRequest{}
	res := &vo.ListWatchingAddressResponse{}

	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
		return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(vo.ErrParams)).JSON(fcx)
	}

	err = fcx.QueryParser(&req)

	if err != nil {
		logger.Warn("query params parse error", "err", err)
		return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(vo.ErrParams)).JSON(fcx)
	}

	res, err = service.ListWatchingAddress(ctx, req)
	return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(err)).JSON(fcx)
}

func GetAssetDetail(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get mev transaction")
	req := vo.AssetDetailRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
		return err
	}

	res, err := service.GetAssetDetail(ctx, req)
	if err != nil {
		logger.Error("get mev transaction error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}
