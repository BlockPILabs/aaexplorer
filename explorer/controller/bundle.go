package controller

import (
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/gofiber/fiber/v2"
)

func GetBundles(ctx *fiber.Ctx) error {

	return vo.NewResultJsonResponse(fiber.Map{
		"params": ctx.Params("network"),
	}).JSON(ctx)
}
