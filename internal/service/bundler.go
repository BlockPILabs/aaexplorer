package service

import (
	"context"
	"strings"
	"time"

	"github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/dao"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaassetdetail"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/bundlerinfo"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/token"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/transfertransaction"
	"github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/shopspring/decimal"
)

type bundlerService struct {
}

var BundlerService = &bundlerService{}

func (*bundlerService) GetBundlers(ctx context.Context, req vo.GetBundlersRequest) (*vo.GetBundlersResponse, error) {
	ctx, logger := log.With(ctx, "service", "GetBundlers")
	err := vo.ValidateStruct(req)
	res := vo.GetBundlersResponse{
		Pagination: vo.Pagination{
			TotalCount: 0,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
	}
	if err != nil {
		logger.Error("params error", "req", req, "err", err.Error())
		return &res, vo.ErrParams.SetData(err)
	}

	client, err := entity.Client(ctx, req.Network)
	if err != nil {
		return nil, err
	}
	//
	list, total, err := dao.BundlerDao.Pagination(ctx, client, req)
	if err != nil {
		return nil, err
	}
	res.TotalCount = total

	//
	res.Records = make([]*vo.BundlersVo, len(list))
	for i, info := range list {
		label := ""
		if info.Edges.Account != nil && info.Edges.Account.Label != nil {
			labels := []string{}
			info.Edges.Account.Label.AssignTo(&labels)
			if len(labels) > 0 {
				label = labels[0]
			}
		}
		res.Records[i] = &vo.BundlersVo{
			Bundler:        info.ID,
			BundlesNum:     info.BundlesNum,
			UserOpsNum:     info.UserOpsNum,
			SuccessRate:    info.SuccessRate,
			SuccessRateD1:  info.SuccessRateD1,
			BundlesNumD1:   info.BundlesNumD1,
			FeeEarnedD1:    info.FeeEarnedD1,
			FeeEarnedUsdD1: info.FeeEarnedUsdD1,
			BundlerLabel:   label,
			BundleRate:     info.BundleRate,
		}
	}

	return &res, nil
}

func (*bundlerService) GetBundler(ctx context.Context, req vo.GetBundlerRequest) (res *vo.GetBundlerResponse, err error) {
	res = &vo.GetBundlerResponse{}
	client, err := entity.Client(ctx, req.Network)
	if err != nil {
		return
	}

	info, err := client.BundlerInfo.Get(ctx, req.Bundler)
	if err != nil {
		return
	}

	res = &vo.GetBundlerResponse{
		FeeEarnedUsdD1: info.FeeEarnedUsdD1,
		FeeEarnedUsd:   info.FeeEarnedUsd,
		SuccessRateD1:  info.SuccessRateD1,
		SuccessRate:    info.SuccessRate,
		BundleRate:     info.BundleRate,
		Rank:           999,
		TotalBundlers:  int64(client.BundlerInfo.Query().CountX(ctx)),
	}

	res.Rank = res.TotalBundlers
	res.Rank = int64(
		client.BundlerInfo.Query().Where(
			bundlerinfo.BundlesNumGT(info.BundlesNum),
		).CountX(ctx),
	) + 1
	addresses, _ := dao.AccountDao.GetAccountByAddresses(ctx, client, []string{req.Bundler})
	if len(addresses) > 0 {
		addresses[0].Label.AssignTo(&res.Label)
	}
	return
}

func (*bundlerService) GetBundlerTransfers(ctx context.Context, req vo.BundlerTransferRequest) (res *vo.BundlerTransferResponse, err error) {

	client, err := entity.Client(ctx, req.Network)
	if err != nil {
		return nil, err
	}

	bundler := req.Address
	if len(bundler) == 0 {
		return nil, nil
	}
	bundler = strings.ToLower(bundler)

	res = &vo.BundlerTransferResponse{
		Pagination: vo.Pagination{
			TotalCount: 0,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
	}

	startTime := time.UnixMilli(time.Now().UnixMilli() - 24*3600*1000*180)

	transferTxs, err := client.TransferTransaction.Query().Where(transfertransaction.FromAddr(bundler), transfertransaction.TokenSymbolNEQ(""), transfertransaction.TimeGTE(startTime)).Order(ent.Desc(transfertransaction.FieldTime)).Offset(req.GetOffset()).Limit(req.PerPage).All(ctx)
	totalCount, err := client.TransferTransaction.Query().Where(transfertransaction.FromAddr(bundler), transfertransaction.TokenSymbolNEQ(""), transfertransaction.TimeGTE(startTime)).Count(ctx)

	if len(transferTxs) == 0 {
		return nil, nil
	}
	var details []vo.TransferInfo
	for _, tx := range transferTxs {
		tokenUrl := ""
		if len(tx.TokenURL) > 0 {
			tokenUrl = config.UrlPrefix + tx.TokenURL
		}
		info := vo.TransferInfo{
			Id:          tx.ID,
			TxnHash:     tx.TxHash,
			Source:      "Transfer",
			Timestamp:   tx.Time.UnixMilli(),
			From:        tx.FromAddr,
			To:          tx.ToAddr,
			Value:       tx.TransferValue,
			TokenSymbol: tx.TokenSymbol,
			TokenImage:  tokenUrl,
		}
		details = append(details, info)
	}
	res.TransferList = details
	res.TotalCount = totalCount

	return res, nil
}

func (*bundlerService) GetBundlerBalance(ctx context.Context, req vo.BundlerBalanceRequest) (res *vo.BundlerBalanceResponse, err error) {

	client, err := entity.Client(ctx, req.Network)
	if err != nil {
		return nil, err
	}

	bundler := req.Address
	if len(bundler) == 0 {
		return nil, nil
	}
	bundler = strings.ToLower(bundler)

	res = &vo.BundlerBalanceResponse{}

	details, err := client.AaAssetDetail.Query().Where(aaassetdetail.UserAddressEqualFold(bundler), aaassetdetail.AssetAmountGT(decimal.Zero)).Order(ent.Desc(aaassetdetail.FieldAssetValue)).All(ctx)
	if len(details) == 0 {
		return res, nil
	}

	totalUsd := decimal.Zero
	for _, detail := range details {
		totalUsd = totalUsd.Add(detail.AssetValue)
	}
	otherUsd := decimal.Zero
	var assetDetails []vo.AssetInfo
	for idx, detail := range details {
		if idx >= 7 {
			otherUsd = otherUsd.Add(detail.AssetValue.RoundDown(6))
			continue
		}
		tokens, _ := client.Token.Query().Where(token.SymbolEqualFold(detail.Symbol)).All(ctx)
		tokenUrl := ""
		if len(tokens) > 0 {
			tokenUrl = config.UrlPrefix + tokens[0].ImageURL
		}
		assetDetail := vo.AssetInfo{
			Symbol:    detail.Symbol,
			Amount:    detail.AssetAmount,
			AmountUsd: detail.AssetValue.RoundDown(6),
			TokenUrl:  tokenUrl,
		}
		percent := decimal.Zero
		if totalUsd.Cmp(decimal.Zero) > 0 {
			percent = detail.AssetValue.DivRound(totalUsd, 4)
		}
		assetDetail.Percent = percent
		assetDetails = append(assetDetails, assetDetail)
	}
	if otherUsd.Cmp(decimal.Zero) > 0 {
		percent := otherUsd.DivRound(totalUsd, 4)
		otherDetail := vo.AssetInfo{
			Symbol:    "Other",
			AmountUsd: otherUsd,
			Percent:   percent,
		}
		otherDetail.TokenUrl = config.OtherCoinUrl
		assetDetails = append(assetDetails, otherDetail)
	}

	res.AssetDetails = assetDetails
	res.TotalUsd = totalUsd.RoundDown(6)

	return res, nil
}
