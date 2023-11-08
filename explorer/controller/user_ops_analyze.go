package controller

import (
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/service"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetUserOpType = "get_user_op_type"
const NameGetAAContractInteract = "get_aa_contract_interact"

func GetUserOpType(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get userop type")
	req := vo.UserOpsTypeRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
	}

	res, err := service.GetUserOpType(ctx, req)
	if err != nil {
		logger.Error("get userop type error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}

func GetAAContractInteract(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start get aa contract interact")
	req := vo.AAContractInteractRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
	}

	res, err := service.GetAAContractInteract(ctx, req)
	if err != nil {
		logger.Error("get aa contract interact error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}
