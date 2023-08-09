package controller

import (
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/service"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetPaymasters = "get_paymasters"

func GetPaymasters(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get paymasters")
	req := vo.GetPaymastersRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
	}
	res, err := service.PaymasterService.GetPaymasters(ctx, req)
	if err != nil {
		logger.Error("get paymasters error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}
