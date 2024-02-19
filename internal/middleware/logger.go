package middleware

import (
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/gofiber/fiber/v2"
	"time"
)

func Logger() fiber.Handler {

	return func(c *fiber.Ctx) error {

		start := time.Now()
		err := c.Next()
		end := time.Now()

		_logger := log.Context(c.UserContext()).With(
			"status", c.Response().StatusCode(),
			"duration", end.Sub(start).Round(time.Millisecond),
		)

		if err != nil {
			switch e := err.(type) {
			case *fiber.Error:
				_logger.Error("api error", "err", e.Message, "status", e.Code)
			default:
				_logger.Error("api error", "err", err)
			}
		} else if c.Response().StatusCode() != fiber.StatusOK {
			_logger.Warn("api warn")
		} else {
			_logger.Debug("api success")
		}
		return err
	}
}
