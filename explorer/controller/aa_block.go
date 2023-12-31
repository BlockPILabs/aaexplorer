package controller

import (
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/service"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetAABlocksPage = "GetAABlocksPage"

func GetAABlocksPage(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start GetAABlocks")
	req := vo.GetAaBlocksRequest{
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
	client, err := entity.Client(ctx, req.Network)
	if err != nil {
		return err
	}
	res, err := service.AaBlockService.GetAaBlockInfo(ctx, client, req)
	if err != nil {
		logger.Error("GetAABlocks error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}
