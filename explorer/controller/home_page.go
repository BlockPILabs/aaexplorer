package controller

import (
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/service"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetDailyStatistic = "get_daily_statistic"
const NameGetAATxnDominance = "get_aa_txn_dominance"
const NameGetLatestUserOps = "get_latest_user_ops"

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

func GetLatestUserOps(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get latest user ops")
	req := vo.LatestUserOpsRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
	}

	res, err := service.GetLatestUserOps(ctx, req)
	if err != nil {
		logger.Error("get latest user ops error", "err", err)
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
