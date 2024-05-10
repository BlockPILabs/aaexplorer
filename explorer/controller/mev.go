package controller

import (
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/service"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameMEVBlock = "MEVBlock"

func MEVBlock(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(ctx)

	logger.Debug("start MEVBlock")

	req := vo.ListBlockMEVBundlersRequest{
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
	res, err := service.MevService.BlockMevList(ctx, req)
	if err != nil {
		logger.Error("MEVBlock error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}
