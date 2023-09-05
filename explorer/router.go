package explorer

import (
	"github.com/BlockPILabs/aa-scan/explorer/controller"
	"github.com/gofiber/fiber/v2"
	fiber_recover "github.com/gofiber/fiber/v2/middleware/recover"
)

func Resister(router fiber.Router) {
	router.Use(fiber_recover.New())
	v1 := router.Group("/v1")
	// networks

	v1.Get("/networks", controller.GetNetworks).Name(controller.NameGetNetworks)

	networksV1 := v1.Group("/network/:network<regex(^[a-z0-9]{1,}$)}>?")
	networksV1.Use(controller.NetworkMiddleware())
	networksV1.Get("/", controller.GetNetwork)
	// Bundlers
	networksV1.Get("/bundlers", controller.GetBundlers).Name(controller.NameGetBundlers)
	networksV1.Get("/bundler/:bundler<regex(0x[a-z0-9]{40}$)}>", controller.GetBundler).Name(controller.NameGetBundler)

	// bundles
	networksV1.Get("/bundles", controller.GetBundles).Name(controller.NameGetBundles)
	networksV1.Get("/bundle/:bundle<regex(0x[a-z0-9]{40}$)}>", controller.GetBundle).Name(controller.NameGetBundle)

	// userops
	networksV1.Get("/userops", controller.GetUserOps).Name(controller.NameGetUserOps)

	// paymasters
	networksV1.Get("/paymasters", controller.GetPaymasters).Name(controller.NameGetPaymasters)
	// blocks
	networksV1.Get("/blocks", controller.GetBlocks).Name(controller.NameGetBlocks)
	networksV1.Get("/block/:block<regex((^0x[a-z0-9]{64}$|^\\d+$))}>", controller.GetBlock).Name(controller.NameGetBlock)

	//user
	//networksV1.Get("/user")
	//networksV1.Get("/user")

	//home page
	networksV1.Get("/dailyStatistic", controller.GetDailyStatistic).Name(controller.NameGetDailyStatistic)
}

func Error() {

}
