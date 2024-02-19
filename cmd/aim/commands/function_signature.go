package commands

import (
	"encoding/json"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/functionsignature"
	"github.com/BlockPILabs/aaexplorer/internal/memo"
	out_service "github.com/BlockPILabs/aaexplorer/service"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	"github.com/valyala/fastjson"
	"strings"
	"time"
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
		for len(nextUrl) > 0 {
			(func() {
				agent := fiber.AcquireAgent()
				defer fiber.ReleaseAgent(agent)

				request := agent.Request()
				request.Header.SetMethod(fiber.MethodGet)
				request.SetRequestURI(nextUrl)
				err := agent.Parse()
				nextUrl = ""
				if err != nil {
					logger.Error("agent parse error", "err", err)
					return
				}
				code, bytes, errs := agent.Bytes()
				if errs != nil || code != fiber.StatusOK {
					logger.Error("agent request error", "status", code, "errs", errs)
					return
				}
				value, err := fastjson.ParseBytes(bytes)
				if err != nil {
					logger.Error("json parse error", "err", err)
					return
				}
				nextUrl = strings.Replace(string(value.GetStringBytes("next")), "http:", "https:", 1)
				values := value.GetArray("results")
				if len(values) < 1 {
					logger.Error("data not found")
					return
				}
				client, _ := entity.Client(cmd.Context())
				var ids []string
				var fs []*out_service.FunctionSignature
				for _, v := range values {
					bytes = v.MarshalTo(nil)
					f := &out_service.FunctionSignature{}
					err = json.Unmarshal(bytes, f)
					if err != nil {
						logger.Error("json parse error 1", "err", err)
						continue
					}
					if len(f.TextSignature) > 0 {
						ss := strings.Split(f.TextSignature, "(")
						f.Name = ss[0]
						fs = append(fs, f)
						ids = append(ids, f.HexSignature)
					}

				}

				if len(ids) < 1 {
					return
				}

				functionSignatures := client.FunctionSignature.Query().Where(
					functionsignature.IDIn(ids...),
				).AllX(cmd.Context())

				sfm := map[string]*ent.FunctionSignature{}
				for _, sf := range functionSignatures {
					sfm[sf.ID] = sf
				}

				bluk := []*ent.FunctionSignatureCreate{}
				for _, f := range fs {
					if _, ok := sfm[f.HexSignature]; ok {
						continue
					}
					bluk = append(bluk,
						client.FunctionSignature.Create().
							SetID(f.HexSignature).
							SetName(f.Name).
							SetText(f.TextSignature).
							SetBytes(f.BytesSignature).
							SetCreateTime(time.Now()),
					)
				}

				client.FunctionSignature.CreateBulk(bluk...).ExecX(cmd.Context())

			})()
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(FunctionSignatureCmd)
}
