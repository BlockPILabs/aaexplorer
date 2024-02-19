package commands

import (
	"fmt"
	"github.com/BlockPILabs/aaexplorer/explorer"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/memo"
	"github.com/BlockPILabs/aaexplorer/internal/middleware"
	aimos "github.com/BlockPILabs/aaexplorer/internal/os"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	fiber_recover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"os"
	"runtime/debug"
	"strings"
	"time"
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
			logger = logger.With(
				"child", fiber.IsChild(),
				"pid", os.Getpid(),
			)
			// db start
			err := entity.Start(logger.With("lib", "ent"), config)
			if err != nil {
				return err
			}

			err = memo.Start(logger.With("lib", "memo"), config)
			if err != nil {
				return err
			}

			app := fiber.New(fiber.Config{
				Prefork:     config.Api.Prefork,
				BodyLimit:   int(config.Api.MaxBodyBytes),
				Concurrency: config.Api.MaxOpenConnections,
				ErrorHandler: func(ctx *fiber.Ctx, err error) error {
					switch e := err.(type) {
					case *vo.Error:
						return vo.NewErrorJsonResponse(e).JSON(ctx)
					case *fiber.Error:
						return vo.NewErrorJsonResponse(&vo.Error{Code: e.Code, Message: e.Message}).JSON(ctx)
					default:
						return vo.NewErrorJsonResponse(vo.ErrSystem).JSON(ctx)
					}
				},
			})
			// Use middleware
			app.Use(fiber.Handler(func(ctx *fiber.Ctx) error {
				_logger := logger.With(
					"module", "api",
					"method", string(ctx.Request().Header.Method()),
					"requestUri", string(ctx.Request().RequestURI()),
					"remoteIp", ctx.IP(),
					"requestId", uuid.NewString(),
				)
				ctx.SetUserContext(log.WithContext(ctx.UserContext(), _logger))
				return ctx.Next()
			}))
			app.Use(favicon.New())

			// logger
			app.Use(middleware.Logger())

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
				StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
					//fmt.Println(c, e)
					_, _ = os.Stderr.WriteString(fmt.Sprintf("panic: %v\n%s\n", e, debug.Stack())) //nolint:errcheck // This will never fail
					//log.Context(c.UserContext()).Error("request panic", "debug", string(debug.Stack()))
					c.Next()
				},
			}))

			// register router
			explorer.Resister(app.Group("/explorer"))

			app.Get("/", func(ftx *fiber.Ctx) error {
				_, err := ftx.WriteString(time.Now().String())
				return err
			})

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
