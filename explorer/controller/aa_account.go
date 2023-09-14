package controller

import (
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/service"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetAaAcountInfo = "GetAaAccountInfo"

func GetAaAccountInfo(fcx *fiber.Ctx) error {
	ctx := fcx.UserContext()
	logger := log.Context(fcx.UserContext())

	logger.Debug("start GetAaAccountInfo")
	req := vo.AaAccountRequestVo{}
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
	res, err := service.AaAccountService.GetAaAccountRecord(ctx, client, req)
	if err != nil {
		logger.Error("GetAaAccountInfo error", "err", err)
	}
	return vo.NewResultJsonResponse(res).JSON(fcx)

}
