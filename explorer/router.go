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

	networksV1 := v1.Group("/networks/:network<regex(^[a-z0-9]{1,}$)}>?")
	networksV1.Use(controller.NetworkMiddleware())
	// Bundles
	networksV1.Get("/bundlers", controller.GetBundlers).Name(controller.NameGetBundlers)
	networksV1.Get("/userops", controller.GetUserOps).Name(controller.NameGetUserOps)
}

func Error() {

}
