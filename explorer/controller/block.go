package controller

import (
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/service"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetBlocks = "get_blocks"
const NameGetBlock = "get_block"

func GetBlocks(fcx *fiber.Ctx) error {

	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get bundles")
	req := vo.GetBlocksRequest{
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

	res, err := service.BlockService.GetBlocks(ctx, req)
	if err != nil {
		logger.Error("get bundles error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}

func GetBlock(fcx *fiber.Ctx) error {

	ctx := fcx.UserContext()
	logger := log.Context(ctx)
	req := vo.GetBlockRequest{
		Block: "",
	}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
	}
	response, err := service.BlockService.GetBlock(ctx, req)
	if err != nil {
		return err
	}
	return vo.NewResultJsonResponse(response).JSON(fcx)
}
