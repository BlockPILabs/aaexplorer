package controller

import (
	"fmt"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/bundlerinfo"
	"github.com/BlockPILabs/aa-scan/internal/vo"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

const (
	Bundler   = "bundler"
	Paymaster = "paymaster"
	Factory   = "factory"
)

func GetTopList(fcx *fiber.Ctx, req vo.TopRequest) error {
	context := fcx.Context()
	reqType := req.Type
	client, err := entity.Client(context)
	if err != nil {
		return nil
	}
	if reqType == Bundler {
		bundlerInfos, err := client.BundlerInfo.Query().Order(ent.Desc(bundlerinfo.FieldBundlesNumD1)).Limit(10).All(context)
		if err != nil {
			return nil
		}
		toBundlerResponse(bundlerInfos)
		return vo.NewResultJsonResponse(bundlerInfos).JSON(fcx)
	} else if reqType == Paymaster {

	} else if reqType == Factory {

	}

	return nil
}

func toBundlerResponse(infos []*ent.BundlerInfo) []*vo.TopBundlerResponse {
	if len(infos) == 0 {
		return nil
	}
	var resp []*vo.TopBundlerResponse
	for _, bundler := range infos {
		top := vo.TopBundlerResponse{
			Address:    bundler.Bundler,
			Bundles:    bundler.BundlesNumD1,
			Success24H: decimal.Zero,
		}
		fmt.Println(top)
	}
	return resp
}
