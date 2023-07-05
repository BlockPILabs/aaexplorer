package commands

import (
	"github.com/BlockPILabs/aa-scan/internal/middleware"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	fiber_recover "github.com/gofiber/fiber/v2/middleware/recover"
	"strings"

	aimos "github.com/BlockPILabs/aa-scan/internal/os"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

// AddFlags exposes some common configuration options on the command-line
// These are exposed for convenience of commands embedding
func AddFlags(cmd *cobra.Command) {

	// rpc flags
	cmd.Flags().String("api.laddr", config.Api.ListenAddress, "api listen address. Port required")
	cmd.Flags().Bool("api.unsafe", config.Api.Unsafe, "enabled unsafe api methods")
	cmd.Flags().String("api.pprof_prefix", config.Api.PprofPrefix, "pprof path (https://golang.org/pkg/net/http/pprof)")

}

// NewStartCmd returns the command that allows the CLI to start a node.
func NewStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"node", "run"},
		Short:   "Run the aim api",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := fiber.New(fiber.Config{
				Prefork:     config.Api.Prefork,
				BodyLimit:   int(config.Api.MaxBodyBytes),
				Concurrency: config.Api.MaxOpenConnections,
				ErrorHandler: func(ctx *fiber.Ctx, err error) error {
					return err
				},
			})
			// Use middleware
			// logger
			app.Use(middleware.Logger(logger.With("module", "api")))

			// cros
			app.Use(cors.New(cors.Config{
				AllowOrigins:     strings.Join(config.Api.CORSAllowedOrigins, ", "),
				AllowMethods:     strings.Join(config.Api.CORSAllowedMethods, ", "),
				AllowHeaders:     strings.Join(config.Api.CORSAllowedHeaders, ", "),
				AllowCredentials: config.Api.CORSAllowedCredentials,
				MaxAge:           config.Api.CORSAMaxAge,
			}))

			// pprof
			if len(config.Api.PprofPrefix) > 0 {
				app.Use(pprof.New(pprof.Config{Prefix: config.Api.PprofPrefix}))
			}

			// error recover
			app.Use(fiber_recover.New(fiber_recover.Config{
				EnableStackTrace: true,
			}))

			go func() {
				err := app.Listen(config.Api.ListenAddress)
				if err != nil {
					aimos.Exit(err.Error())
					return
				}
			}()

			logger.Info("start api")

			// Stop upon receiving SIGTERM or CTRL-C.
			aimos.TrapSignal(logger, func() {
				err := app.Shutdown()
				if err != nil {
					logger.Error("stop api error", "err", err)
				}
			})

			// Run forever.
			select {}
		},
	}

	AddFlags(cmd)
	return cmd
}
