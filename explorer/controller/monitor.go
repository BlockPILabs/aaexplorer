package controller

import (
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/service"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameAddMonitor = "add_monitor"
const NameRemoveMonitor = "remove_monitor"

func AddMonitor(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get add monitor")
	req := vo.AddMonitorRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
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
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
	}

	res, err := service.RemoveMonitor(ctx, req)
	if err != nil {
		logger.Error("get mev transaction error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}
