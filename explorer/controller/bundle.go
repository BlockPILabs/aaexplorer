package controller

import (
	"fmt"
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/gofiber/fiber/v2"
)

func GetBundles(fcx *fiber.Ctx) error {
	logger := log.Context(fcx.UserContext())

	//net, _ := dao.NetworkDao.ContextValue(fcx.UserContext())
	//ctx := fcx.UserContext()
	//networkFlag := fcx.Params("network")
	//
	//net, err := dao.NetworkDao.GetNetworkByNetwork(ctx, networkFlag)
	//if err != nil {
	//	return err
	//}
	req := vo.GetBundlersRequest{}
	err := fcx.ParamsParser(&req)
	if err != nil {
		logger.Warn("params parse error", "err", err)
	}
	err = fcx.QueryParser(&req)
	if err != nil {
		logger.Warn("query params parse error", "err", err, "network", req.Network)
	}
	err = vo.ValidateStruct(req)
	if err != nil {
		logger.Error("params error", "req", req, "err", err.Error())
		return vo.ErrParams.SetData(err)
	}
	fmt.Println(req)

	return vo.NewResultJsonResponse(req).JSON(fcx)
}
