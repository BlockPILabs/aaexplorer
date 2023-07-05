package explorer

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
)

func Resister(router fiber.Router) {
	router.Use(recover2.New())
	router.Group("/v1")
}
