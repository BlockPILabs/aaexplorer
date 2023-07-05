package middleware

import (
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/gofiber/fiber/v2"
	"time"
)

func Logger(logger log.Logger) fiber.Handler {

	return func(c *fiber.Ctx) error {

		start := time.Now()
		err := c.Next()
		end := time.Now()

		_logger := logger.With(
			"status", c.Response().StatusCode(),
			"ip", c.IP(),
			"method", c.Method(),
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
