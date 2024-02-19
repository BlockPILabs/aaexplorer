package controller

import (
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/service"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetUserBalance = "get_user_balance"

func GetUserBalance(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get user balance")
	req := vo.UserBalanceRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
	}

	res, err := service.GetUserBalance(ctx, req)
	if err != nil {
		logger.Error("get user balance error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}
