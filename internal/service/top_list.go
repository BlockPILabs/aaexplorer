package service

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/bundlerinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/factoryinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/paymasterinfo"
	"github.com/BlockPILabs/aa-scan/internal/vo"
)

func GetTopBundler(ctx context.Context, req vo.TopBundlerRequest) (*vo.TopBundlerResponse, error) {
	network := req.Network
	client, err := entity.Client(ctx, network)
	if err != nil {
		return nil, err
	}
	var resp = &vo.TopBundlerResponse{
		Pagination: vo.Pagination{
			TotalCount: 0,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
	}

	bundlerInfos, err := client.BundlerInfo.Query().Order(bundlerinfo.ByFeeEarnedUsdD1(sql.OrderDesc())).Offset(req.GetOffset()).Limit(req.GetPerPage()).All(ctx)
	if len(bundlerInfos) == 0 {
		return nil, nil
	}
	var bundlerDetails []*vo.BundlerDetail
	for _, info := range bundlerInfos {
		detail := &vo.BundlerDetail{
			Address:         info.Bundler,
			Bundles:         info.BundlesNumD1,
			Success24H:      info.SuccessRateD1,
			FeeEarned24H:    info.FeeEarnedD1.Round(2),
			FeeEarnedUsd24H: info.FeeEarnedUsdD1.Round(2),
		}
		bundlerDetails = append(bundlerDetails, detail)
	}
	resp.BundlerDetails = bundlerDetails

	return resp, nil
}

func GetTopPaymaster(ctx context.Context, req vo.TopPaymasterRequest) (*vo.TopPaymasterResponse, error) {
	network := req.Network
	client, err := entity.Client(ctx, network)
	if err != nil {
		return nil, err
	}
	var resp = &vo.TopPaymasterResponse{
		Pagination: vo.Pagination{
			TotalCount: 0,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
	}

	paymasterInfos, err := client.PaymasterInfo.Query().Order(paymasterinfo.ByGasSponsoredUsdD1(sql.OrderDesc())).Offset(req.GetOffset()).Limit(req.GetPerPage()).All(ctx)
	if len(paymasterInfos) == 0 {
		return nil, nil
	}
	var paymasterDetails []*vo.PaymasterDetail
	for _, info := range paymasterInfos {
		detail := &vo.PaymasterDetail{
			Address:         info.Paymaster,
			ReserveUsd:      info.ReserveUsd,
			GasSponsored:    info.GasSponsoredD1.Round(2),
			GasSponsoredUsd: info.GasSponsoredUsdD1.Round(2),
		}
		paymasterDetails = append(paymasterDetails, detail)
	}
	resp.PaymasterDetails = paymasterDetails

	return resp, nil
}

func GetTopFactory(ctx context.Context, req vo.TopFactoryRequest) (*vo.TopFactoryResponse, error) {
	network := req.Network
	client, err := entity.Client(ctx, network)
	if err != nil {
		return nil, err
	}
	var resp = &vo.TopFactoryResponse{
		Pagination: vo.Pagination{
			TotalCount: 0,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
	}

	factoryInfos, err := client.FactoryInfo.Query().Order(factoryinfo.ByAccountNumD1(sql.OrderDesc())).Offset(req.GetOffset()).Limit(req.GetPerPage()).All(ctx)
	if len(factoryInfos) == 0 {
		return nil, nil
	}
	var factoryDetails []*vo.FactoryDetail
	for _, info := range factoryInfos {
		detail := &vo.FactoryDetail{
			Address:       info.Factory,
			ActiveAccount: int64(info.AccountDeployNumD1),
			TotalAccount:  int64(info.AccountDeployNumD1),
		}
		factoryDetails = append(factoryDetails, detail)
	}
	resp.FactoryDetails = factoryDetails

	return resp, nil
}
