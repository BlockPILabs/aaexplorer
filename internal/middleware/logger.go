package middleware

import (
	"github.com/BlockPILabs/aa-scan/internal/log"
	"github.com/gofiber/fiber/v2"
	"net/http"
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
			_logger.Error("api error", "err", err)
		} else if c.Response().StatusCode() != http.StatusOK {
			_logger.Warn("api warn")
		} else {
			_logger.Debug("api success")
		}
		return err
	}
}
