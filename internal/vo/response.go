package vo

import "github.com/gofiber/fiber/v2"

const vsn = "1.0"

type SetResponseOption func(*JsonResponse) *JsonResponse

type JsonResponse struct {
	Version string `json:"jsonrpc,omitempty"`
	Id      string `json:"id,omitempty"`
	Error   *Error `json:"error,omitempty"`
	Result  any    `json:"result,omitempty"`
}

func (r *JsonResponse) JSON(ctx *fiber.Ctx) error {
	return ctx.JSON(r)
}

func SetResponseId(id string) SetResponseOption {
	return func(r *JsonResponse) *JsonResponse {
		r.Id = id
		return r
	}
}

func SetResponseResult(result any) SetResponseOption {
	return func(r *JsonResponse) *JsonResponse {
		r.Result = result
		return r
	}
}

func SetResponseError(err error) SetResponseOption {
	return func(r *JsonResponse) *JsonResponse {
		if err == nil {
			switch e := err.(type) {
			case *Error:
				r.Error = e
			default:
				r.Error = ErrSystem.SetMessage(err.Error())
			}
		} else {
			r.Error = nil
		}

		return r
	}
}

func NewJsonResponse(sets ...SetResponseOption) *JsonResponse {
	r := &JsonResponse{Version: vsn}
	for _, set := range sets {
		r = set(r)
	}
	return r
}
func NewResultJsonResponse(result any, sets ...SetResponseOption) *JsonResponse {
	r := &JsonResponse{Version: vsn, Result: result}
	for _, set := range sets {
		r = set(r)
	}
	return r
}

func NewErrorJsonResponse(error *Error, sets ...SetResponseOption) *JsonResponse {
	r := &JsonResponse{Version: vsn, Error: error}
	for _, set := range sets {
		r = set(r)
	}
	return r
}
