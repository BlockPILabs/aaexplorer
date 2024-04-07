package service

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/bundlerinfo"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/factoryinfo"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/paymasterinfo"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/shopspring/decimal"
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
	total, err := client.BundlerInfo.Query().Count(ctx)
	bundlerInfos, err := client.BundlerInfo.Query().Order(bundlerinfo.ByBundlesNumD1(sql.OrderDesc())).Offset(req.GetOffset()).Limit(req.GetPerPage()).All(ctx)
	if len(bundlerInfos) == 0 {
		return nil, nil
	}
	var bundlerDetails []*vo.BundlerDetail
	for _, info := range bundlerInfos {
		detail := &vo.BundlerDetail{
			Address:         info.ID,
			Bundles:         info.BundlesNumD1,
			Success24H:      info.SuccessRateD1,
			FeeEarned24H:    info.FeeEarnedD1.Round(config.FeePrecision),
			FeeEarnedUsd24H: info.FeeEarnedUsdD1.Round(config.FeePrecision),
		}
		bundlerDetails = append(bundlerDetails, detail)
	}
	resp.BundlerDetails = bundlerDetails
	resp.TotalCount = total

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
	total, err := client.PaymasterInfo.Query().Count(ctx)
	paymasterInfos, err := client.PaymasterInfo.Query().Order(paymasterinfo.ByGasSponsoredUsdD1(sql.OrderDesc())).Offset(req.GetOffset()).Limit(req.GetPerPage()).All(ctx)
	if len(paymasterInfos) == 0 {
		return nil, nil
	}
	var paymasterDetails []*vo.PaymasterDetail
	for _, info := range paymasterInfos {
		detail := &vo.PaymasterDetail{
			Address:         info.ID,
			ReserveUsd:      info.ReserveUsd,
			GasSponsored:    info.GasSponsoredD1.Round(config.FeePrecision),
			GasSponsoredUsd: info.GasSponsoredUsdD1.Round(config.FeePrecision),
		}
		paymasterDetails = append(paymasterDetails, detail)
	}
	resp.PaymasterDetails = paymasterDetails
	resp.TotalCount = total
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

	total, err := client.FactoryInfo.Query().Count(ctx)
	factoryInfos, err := client.FactoryInfo.Query().Order(factoryinfo.ByAccountNumD1(sql.OrderDesc())).Offset(req.GetOffset()).Limit(req.GetPerPage()).All(ctx)
	if len(factoryInfos) == 0 {
		return nil, nil
	}
	var factoryDetails []*vo.FactoryDetail
	for _, info := range factoryInfos {
		detail := &vo.FactoryDetail{
			Address:       info.ID,
			ActiveAccount: int64(info.AccountDeployNumD1),
			TotalAccount:  int64(info.AccountDeployNumD1),
		}
		factoryDetails = append(factoryDetails, detail)
	}
	resp.FactoryDetails = factoryDetails
	resp.TotalCount = total

	return resp, nil
}

func GetTopWhale(ctx context.Context, req vo.TopWhaleRequest) (*vo.TopWhaleResponse, error) {
	network := req.Network
	client, err := entity.Client(ctx, network)
	if err != nil {
		return nil, err
	}

	var resp = &vo.TopWhaleResponse{
		Pagination: vo.Pagination{
			TotalCount: 0,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
	}

	total := req.RankLimit
	if total == 0 {
		return nil, nil
	}
	aaAssets, err := client.QueryContext(ctx, `select row_number() over(order by asset_value desc) as row_number, rank_tab.user_address, rank_tab.asset_value from (select * from aa_asset order by asset_value DESC offset $1 limit $2) rank_tab offset $3 limit $4`, 0, req.RankLimit, req.GetOffset(), req.GetPerPage())
	defer aaAssets.Close()
	if err != nil {
		return nil, err
	}

	var aaAssetsRankList []*vo.TopWhaleInfo
	queryContext, err := client.QueryContext(ctx, "select sum(rank_tab.asset_value) from (select asset_value from aa_asset order by asset_value DESC offset $1 limit $2) rank_tab", 0, req.RankLimit)
	defer queryContext.Close()
	if err != nil {
		return nil, err
	}
	var totalValue float64
	if queryContext.Next() {
		queryContext.Scan(&totalValue)
	}
	if totalValue == 0 {
		return nil, nil
	}
	for aaAssets.Next() {
		var rank int
		var address string
		var assetValue decimal.Decimal
		aaAssets.Scan(&rank, &address, &assetValue)
		assetValueForFloat64, _ := assetValue.Float64()
		whale := &vo.TopWhaleInfo{
			Rank:       rank,
			Address:    address,
			Balance:    assetValue,
			Percentage: assetValueForFloat64 / totalValue,
		}

		aaAssetsRankList = append(aaAssetsRankList, whale)
	}

	resp.TopWhaleRankList = aaAssetsRankList
	resp.TotalCount = total

	return resp, nil
}
