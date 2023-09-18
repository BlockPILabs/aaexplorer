package commands

import (
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/memo"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

var FunctionSignatureCmd = &cobra.Command{
	Use:    "function_signature",
	Short:  "load solid function signature",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = entity.Start(logger.With("lib", "ent"), config)
		if err != nil {
			return
		}

		err = memo.Start(logger.With("lib", "memo"), config)
		if err != nil {
			return
		}

		nextUrl := "https://www.4byte.directory/api/v1/signatures/?format=json"
		for true {
			(func() {
				agent := fiber.AcquireAgent()
				defer fiber.ReleaseAgent(agent)

				request := agent.Request()
				request.Header.SetMethod(fiber.MethodGet)
				request.SetRequestURI(nextUrl)

				agent.Parse()

			})()
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(FunctionSignatureCmd)
}
