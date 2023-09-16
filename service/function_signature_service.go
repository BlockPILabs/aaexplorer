package service

import (
	"encoding/json"
	"errors"
	"github.com/BlockPILabs/aa-scan/internal/utils"
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
