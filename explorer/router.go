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

	v1.Get("/networks", controller.GetNetworks).Name("get_networks")

	// Bundles
	v1.Get("/networks/:network<regex(^[a-z0-9]{1,}$)}>?/bundles", controller.GetBundles).Name("get_bundles")
}

func Error() {

}
