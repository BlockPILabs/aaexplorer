package controller

import (
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/service"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetDailyStatistic = "get_daily_statistic"
const NameGetAATxnDominance = "get_aa_txn_dominance"

func GetDailyStatistic(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get daily statistic")
	req := vo.DailyStatisticRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
	}

	res, err := service.GetDailyStatistic(ctx, req)
	if err != nil {
		logger.Error("get daily statistic error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}

func GetAATxnDominance(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get aa txn dominance")
	req := vo.AATxnDominanceRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
	}

	res, err := service.GetAATxnDominance(ctx, req)
	if err != nil {
		logger.Error("get aa txn dominance error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}

func GetAccountInfo(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get daily statistic")
	req := vo.DailyStatisticRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
	}

	res, err := service.GetDailyStatistic(ctx, req)
	if err != nil {
		logger.Error("get daily statistic error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}
