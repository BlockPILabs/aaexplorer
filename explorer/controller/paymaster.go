package controller

import (
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/service"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetPaymasters = "get_paymasters"
const NameGetPaymasterOverview = "get_paymaster_overview"

func GetPaymasters(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get paymasters")
	req := vo.GetPaymastersRequest{
		PaginationRequest: vo.NewDefaultPaginationRequest(),
	}
	res := &vo.GetPaymastersResponse{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)

		return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(vo.ErrParams)).JSON(fcx)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
		return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(vo.ErrParams)).JSON(fcx)
	}
	res, err = service.PaymasterService.GetPaymasters(ctx, req)
	if err != nil {
		logger.Error("get paymasters error", "err", err)
	}
	return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(err)).JSON(fcx)
}

func GetPaymasterOverview(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get paymaster overview")
	req := vo.GetPaymasterOverviewRequest{}
	res := &vo.GetPaymasterOverviewResponse{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)

		return vo.NewResultJsonResponse(res, vo.SetResponseError(vo.ErrParams)).JSON(fcx)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
		return vo.NewResultJsonResponse(res, vo.SetResponseError(vo.ErrParams)).JSON(fcx)
	}
	res, err = service.GetPaymasterOverview(ctx, req)
	if err != nil {
		logger.Error("get paymaster overview error", "err", err)
	}
	return vo.NewResultJsonResponse(res, vo.SetResponseError(err)).JSON(fcx)
}
