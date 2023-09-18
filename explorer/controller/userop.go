package controller

import (
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/service"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetUserOps = "get_user_ops"

func GetUserOps(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get userops")
	req := vo.GetUserOpsRequest{
		PaginationRequest: vo.NewDefaultPaginationRequest(),
	}
	res := &vo.GetUserOpsResponse{}
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

	res, err = service.UserOpService.GetUserOps(ctx, req)
	if err != nil {
		logger.Error("get userops error", "err", err)
	}
	return vo.NewResultJsonResponse(res, vo.SetResponseAutoDataError(err)).JSON(fcx)
}

const NameGetUserOpsAnalysis = "GetUserOpsAnalysis"

func GetUserOpsAnalysis(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start GetUserOpsAnalysis")
	req := vo.UserOpsAnalysisRequestVo{}
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
	res, err := service.UserOpService.GetUserOpsAnalysis(ctx, client, req)
	if err != nil {
		logger.Error("GetUserOpsAnalysis error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}

const NameGetUserOpsAnalysisList = "GetUserOpsAnalysisList"

func GetUserOpsAnalysisList(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start GetUserOpsAnalysisList")
	req := vo.UserOpsAnalysisListRequestVo{}
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
	res, err := service.UserOpService.GetUserOpsAnalysisList(ctx, client, req)
	if err != nil {
		logger.Error("GetUserOpsAnalysisList error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}
