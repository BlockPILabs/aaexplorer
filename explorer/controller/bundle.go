package controller

import (
	"fmt"
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/gofiber/fiber/v2"
)

func GetBundles(fcx *fiber.Ctx) error {

	net, _ := dao.NetworkDao.ContextValue(fcx.UserContext())
	//ctx := fcx.UserContext()
	//networkFlag := fcx.Params("network")
	//
	//net, err := dao.NetworkDao.GetNetworkByNetwork(ctx, networkFlag)
	//if err != nil {
	//	return err
	//}
	fmt.Println(net)
	return nil
}
