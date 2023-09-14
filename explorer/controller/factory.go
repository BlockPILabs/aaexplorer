package controller

import (
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/service"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetFactories = "get_factories"

func GetFactories(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get userops")
	req := vo.GetFactoriesRequest{
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

	res, err := service.FactoryService.GetFactories(ctx, req)
	if err != nil {
		logger.Error("get userops error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}
