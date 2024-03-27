package service

import (
	"context"
	"encoding/json"
	"entgo.io/ent/dialect/sql"
	"errors"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/functionsignature"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fastjson"
	"strings"
	"time"
)

type FunctionSignature struct {
	Id             int             `json:"id"`
	CreatedAt      time.Time       `json:"created_at"`
	TextSignature  string          `json:"text_signature"`
	Name           string          `json:"name"`
	HexSignature   string          `json:"hex_signature"`
	BytesSignature json.RawMessage `json:"bytes_signature"`
}

func GetSignature(signature string) (*FunctionSignature, error) {
	if !utils.IsHexSting(signature) && len(signature) < 8 {
		return nil, errors.New("not signature")
	}
	signature = utils.Fix0x(signature)
	//https://www.4byte.directory/api/v1/signatures/?format=json&hex_signature=0x4fa14b84
	agent := fiber.AcquireAgent()
	defer fiber.ReleaseAgent(agent)
	req := agent.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI("https://www.4byte.directory/api/v1/signatures/?format=json&hex_signature=" + signature)
	err := agent.Parse()
	if err != nil {
		return nil, err
	}
	code, bytes, errs := agent.Bytes()
	if errs != nil || code != fiber.StatusOK {
		return nil, errors.New("request error")
	}
	value, err := fastjson.ParseBytes(bytes)
	if err != nil {
		return nil, err
	}
	values := value.GetArray("results")
	if len(values) < 1 {
		return nil, errors.New("not found")
	}
	v := values[0]
	bytes = v.MarshalTo(nil)
	s := &FunctionSignature{}
	err = json.Unmarshal(bytes, s)
	if err != nil {
		return nil, err
	}
	if len(s.TextSignature) > 0 {
		ss := strings.Split(s.TextSignature, "(")
		s.Name = ss[0]
	}
	return s, nil
}

func ScanSignature(ctx context.Context) {
	logger := log.Context(ctx)
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
			client, _ := entity.Client(ctx)
			var ids []string
			var fs []*FunctionSignature
			for _, v := range values {
				bytes = v.MarshalTo(nil)
				f := &FunctionSignature{}
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
			).AllX(ctx)

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

			client.FunctionSignature.CreateBulk(bluk...).OnConflict(sql.DoNothing()).Exec(ctx)

		})()
	}
}
