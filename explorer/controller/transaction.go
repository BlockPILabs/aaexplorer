package controller

import (
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/service"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetTransaction = "GetTransaction"

func GetTransaction(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start GetTransaction")
	req := vo.TransactionRequestVo{}
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
	res, err := service.TransactionService.GetRecord(ctx, client, req)
	if err != nil {
		logger.Error("GetRecord error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}

const NameGetTransactionList = "GetTransactionList"

func GetTransactionList(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start GetTransactionList")
	req := vo.TransactionListRequestVo{}
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
	res, err := service.TransactionService.GetPages(ctx, client, req)
	if err != nil {
		logger.Error("GetPages error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)
}
