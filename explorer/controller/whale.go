package controller

import (
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/service"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetWhaleOverview = "get_whale_overview"
const NameGetWhaleChart = "get_whale_chart"

func GetWhaleOverview(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get whale overview")
	req := vo.WhaleOverviewRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
	}

	res, err := service.GetWhaleOverview(ctx, req)
	if err != nil {
		logger.Error("get whale overview error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}

func GetWhaleChart(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get whale chart")
	req := vo.WhaleChartRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
	}

	res, err := service.GetWhaleChart(ctx, req)
	if err != nil {
		logger.Error("get whale chart error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}
