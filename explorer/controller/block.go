package controller

import (
	"github.com/BlockPILabs/aa-scan/internal/dao"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/gofiber/fiber/v2"
)

const NameGetBlocks = "get_blocks"
const NameGetBlock = "get_block"

func GetBlocks(fcx *fiber.Ctx) error {

	network, _ := dao.NetworkDao.ContextValue(fcx.UserContext())

	return vo.NewResultJsonResponse(network).JSON(fcx)
}

func GetBlock(fcx *fiber.Ctx) error {

	network, _ := dao.NetworkDao.ContextValue(fcx.UserContext())
	return vo.NewResultJsonResponse(fiber.Map{
		"network": network,
		"block":   fcx.Params("block"),
	}).JSON(fcx)
}
