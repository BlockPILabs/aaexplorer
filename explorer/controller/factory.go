package controller

import (
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/service"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetFactories = "get_factories"
const NameGetFactoryAccounts = "get_factory_accounts"

func GetFactories(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get userops")
	req := vo.GetFactoriesRequest{
		PaginationRequest: vo.NewDefaultPaginationRequest(),
	}
	res := &vo.GetFactoriesResponse{}
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

	res, err = service.FactoryService.GetFactories(ctx, req)
	if err != nil {
		logger.Error("get factories error", "err", err)
	}
	return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(err)).JSON(fcx)
}

func GetFactoryAccounts(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get userops")
	req := vo.GetFactoryAccountsRequest{
		PaginationRequest: vo.NewDefaultPaginationRequest(),
	}
	res := &vo.GetAccountsResponse{}
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

	res, err = service.AaAccountService.GetAccounts(ctx, vo.GetAccountsRequest{
		PaginationRequest: req.PaginationRequest,
		Network:           req.Network,
		Factory:           req.Factory,
	})
	if err != nil {
		logger.Error("get factories error", "err", err)
	}
	return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(err)).JSON(fcx)
}

const NameGetFactory = "GetFactory"

func GetFactory(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start GetFactory")
	req := vo.GetFactoryRequest{}
	res := &vo.GetFactoryResponse{}
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

	res, err = service.FactoryService.GetFactory(ctx, req)
	if err != nil {
		logger.Error("GetFactory error", "err", err)
	}
	return vo.NewResultJsonResponse(res, vo.SetResponseError(err)).JSON(fcx)
}
