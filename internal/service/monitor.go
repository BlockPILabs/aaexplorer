package service

import (
	"context"
	"errors"
	"github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/dao"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaaccountdata"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaasset"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaassetdetail"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/mevtransaction"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/monitor"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/token"
	interlog "github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/internal/vo"
	"github.com/BlockPILabs/aaexplorer/util"
	"github.com/chenzhijie/go-web3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"math/big"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var logger = interlog.L()
var syncMap = map[string]int64{}

func SetLogger(lg interlog.Logger) {
	logger = lg
}

func AddMonitor(ctx context.Context, req vo.AddMonitorRequest) (*vo.AddMonitorResponse, error) {

	monitorAddress := req.MonitorAddress
	if len(monitorAddress) == 0 {
		return nil, nil
	}
	network := req.Network
	monitorAddress = strings.ToLower(monitorAddress)
	client, err := entity.Client(ctx)
	if err != nil {
		return nil, err
	}
	var resp = &vo.AddMonitorResponse{}

	userAddress, err := util.RecoverSignFromMetamask(req.Message, req.Sign)
	if err != nil {
		logger.Error("AddMonitor sign err ", "monitorAddress", monitorAddress, "network", network, "msg", err)
		return nil, err
	}
	userAddress = strings.ToLower(userAddress)

	olds, err := client.Monitor.Query().Where(monitor.MonitorAddressEqualFold(monitorAddress), monitor.UserAddressEqualFold(userAddress)).All(ctx)
	if len(olds) > 0 {
		return nil, vo.MonitorExist
	}

	nClient, err := entity.Client(ctx, network)
	if err != nil {
		return nil, nil
	}
	aaAccounts, err := nClient.AaAccountData.Query().Where(aaaccountdata.IDEqualFold(monitorAddress)).All(ctx)
	var monitorType = ""
	if len(aaAccounts) > 0 {
		monitorType = aaAccounts[0].AaType
	}

	monitor := client.Monitor.Create().SetCreateTime(time.Now()).SetMonitorAddress(monitorAddress).SetUserAddress(userAddress).SetMonitorAddressType(monitorType)
	_, err = monitor.Save(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func RemoveMonitor(ctx context.Context, req vo.RemoveMonitorRequest) (*vo.RemoveMonitorResponse, error) {

	monitorAddress := req.MonitorAddress
	if len(monitorAddress) == 0 {
		return nil, nil
	}
	network := req.Network
	monitorAddress = strings.ToLower(monitorAddress)
	client, err := entity.Client(ctx)
	if err != nil {
		return nil, err
	}
	var resp = &vo.RemoveMonitorResponse{}

	userAddress, err := util.RecoverSignFromMetamask(req.Message, req.Sign)
	if err != nil {
		logger.Error("RemoveMonitor sign err ", "monitorAddress", monitorAddress, "network", network, "msg", err)
		return nil, err
	}
	userAddress = strings.ToLower(userAddress)

	olds, err := client.Monitor.Query().Where(monitor.MonitorAddressEqualFold(monitorAddress), monitor.UserAddressEqualFold(userAddress)).All(ctx)
	if len(olds) == 0 {
		logger.Info("RemoveMonitor not exist ", "userAddress", userAddress, "monitorAddress", monitorAddress)
		return nil, vo.MonitorNotExist
	}

	_, err = client.Monitor.Delete().Where(monitor.UserAddressEqualFold(userAddress), monitor.MonitorAddressEqualFold(monitorAddress)).Exec(ctx)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetAssetDetail(ctx context.Context, req vo.AssetDetailRequest) (*vo.AssetDetailResponse, error) {
	userAddress := req.UserAddress
	network := req.Network
	if len(userAddress) == 0 {
		return nil, nil
	}
	userAddress = strings.ToLower(userAddress)
	client, err := entity.Client(ctx, network)
	if err != nil {
		return nil, err
	}
	var resp = &vo.AssetDetailResponse{}
	var assetDetails []vo.AssetDetail

	if syncMap[userAddress] == 1 {
		return resp, vo.NewUserErr
	}

	existUsers, _ := client.AaAssetDetail.Query().Where(aaassetdetail.UserAddressEqualFold(userAddress)).All(ctx)
	if len(existUsers) == 0 {
		syncMap[userAddress] = 1
		go AddOne(ctx, userAddress)
		return resp, vo.NewUserErr
	}

	details, err := client.AaAssetDetail.Query().Where(aaassetdetail.UserAddressEqualFold(userAddress), aaassetdetail.AssetAmountGT(decimal.Zero)).Order(ent.Desc(aaassetdetail.FieldAssetValue)).All(ctx)
	if len(details) == 0 {
		return resp, nil
	}

	//counts, err := client.AaAssetDetail.Query().Where(aaassetdetail.UserAddressEqualFold(userAddress), aaassetdetail.AssetAmountGT(decimal.Zero)).Count(ctx)

	totalUsd := decimal.Zero
	for _, detail := range details {
		totalUsd = totalUsd.Add(detail.AssetValue)
	}
	otherUsd := decimal.Zero
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
		assetDetail := vo.AssetDetail{
			Symbol:    detail.Symbol,
			Network:   network,
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
		otherDetail := vo.AssetDetail{
			Symbol:    "Other",
			Network:   network,
			AmountUsd: otherUsd,
			Percent:   percent,
		}
		otherDetail.TokenUrl = config.OtherCoinUrl
		assetDetails = append(assetDetails, otherDetail)
	}

	resp.AssetDetails = assetDetails
	resp.TotalAssetUsd = totalUsd.RoundDown(6)
	//resp.Pagination.TotalCount = counts

	return resp, nil
}

func ListWatchingAddress(ctx context.Context, req vo.ListWatchingAddressRequest) (*vo.ListWatchingAddressResponse, error) {
	if len(req.UserAddress) == 0 {
		return nil, errors.New("UserAddress Fields is nil")
	}
	req.UserAddress = strings.ToLower(req.UserAddress)
	list, total, err := dao.MonitorDao.ListMonitorDao(ctx, req)
	if err != nil {
		return nil, err
	}

	return &vo.ListWatchingAddressResponse{
		Pagination: vo.Pagination{
			TotalCount: total,
			PerPage:    req.GetPerPage(),
			Page:       req.GetPage(),
		},
		Monitors: list,
	}, nil
}

func GetMonitorMevInfo(ctx context.Context, req vo.MevInfoRequest) (*vo.MevInfoResponse, error) {
	client, err := entity.Client(ctx, req.Network)
	if err != nil {
		return nil, err
	}
	res := &vo.MevInfoResponse{}

	sum, err := client.MevTransaction.Query().Aggregate(ent.Sum(mevtransaction.FieldMevProfitUsd)).Float64(ctx)
	res.AttackerTotalProfit = sum
	if err != nil {
		return nil, err
	}

	mevTotal, err := client.MevTransaction.Query().Count(ctx)
	if err != nil {
		return nil, err
	}
	userOpsTotal, err := client.AAUserOpsInfo.Query().Count(ctx)

	if err != nil {
		return nil, err
	}

	if userOpsTotal == 0 {
		logger.Info("UserOpsTotal is zero")
		return nil, nil
	}

	res.MevUserOpsRatio = float64(mevTotal) / float64(userOpsTotal)

	attackerAccounts, err := client.MevTransaction.Query().GroupBy(mevtransaction.FieldAttacker).Strings(ctx)
	if err != nil {
		return nil, err
	}

	res.AttackerAccounts = len(attackerAccounts)

	return res, nil
}

func GetAccountType(ctx context.Context, req vo.AccountTypeRequest) (*vo.AccountTypeResponse, error) {
	network := req.Network
	address := req.AccountAddress
	if len(address) == 0 {
		return nil, nil
	}
	client, err := entity.Client(ctx, network)
	if err != nil {
		return nil, err
	}
	res := &vo.AccountTypeResponse{}
	address = strings.ToLower(address)

	accountDatas, err := client.AaAccountData.Query().Where(aaaccountdata.IDEqualFold(address)).All(ctx)
	if len(accountDatas) == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	typeStr := accountDatas[0].AaType
	if typeStr == "aa" {
		typeStr = "Account"
	}
	typeStr = capitalizeFirstLetter(typeStr)
	res.Type = typeStr
	res.AccountAddress = address
	return res, nil
}

func capitalizeFirstLetter(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

func AddOne(ctx context.Context, address string) {
	cli, err := entity.Client(ctx)
	if err != nil {
		syncMap[address] = 0
		logger.Error("AddOne err, ", "msg", err)
		return
	}
	networks, err := cli.Network.Query().All(ctx)
	if err != nil {
		syncMap[address] = 0
		return
	}
	if len(networks) == 0 {
		syncMap[address] = 0
		return
	}

	for _, net := range networks {
		network := net.ID
		client, err := entity.Client(ctx, network)
		if err != nil {
			continue
		}
		aas, err := client.AaAsset.Query().Where(aaasset.IDEqualFold(address)).All(ctx)
		if err != nil {
			logger.Error("AddOne query asset err ", "msg", err)
			continue
		}
		if len(aas) > 0 {
			continue
		}
		w3, err := web3.NewWeb3(net.HTTPRPC)
		if err != nil {
			logger.Error("AddOne newWeb3 err ", "msg", err)
			continue
		}
		w3.Eth.SetChainId(net.ChainID)
		tokens, err := client.Token.Query().All(ctx)
		if len(tokens) == 0 {
			continue
		}
		blockNum, err := w3.Eth.GetBlockNumber()
		if err != nil {
			logger.Error("AddOne blockNum err ", "msg", err)
			continue
		}
		doAddOne(ctx, client, tokens, w3, blockNum, network, address)
	}
	syncMap[address] = 0
}

func doAddOne(ctx context.Context, client *ent.Client, tokens []*ent.Token, w3 *web3.Web3, blockNum uint64, network string, userAddress string) {
	logger.Info("AddOne-doAddOne start.")
	totalValue := decimal.Zero
	for _, token := range tokens {
		contractAddress := token.ContractAddress
		if len(contractAddress) == 0 {
			continue
		}
		decimals := token.Decimals
		balance := getBalance(ctx, contractAddress, userAddress, decimals, w3)
		assetValue := balance.Mul(token.TokenPrice)
		addOrUpdateAssetDetail(ctx, contractAddress, userAddress, assetValue, client, network, token.Symbol, balance)
		totalValue = totalValue.Add(assetValue)
	}
	balance, err := w3.Eth.GetBalance(common.HexToAddress(userAddress), big.NewInt(int64(blockNum)))
	if err != nil {
		logger.Error("AddOne get balance err ", "user", userAddress, "network", network, "msg", err)
		return
	}
	nativeBalance := decimal.NewFromBigInt(balance, 0).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(config.DefaultDecimals)))
	nativeTokens, err := client.Token.Query().Where(token.TypeEQ("base"), token.NetworkEQ(network)).Limit(1).All(ctx)
	nativePrice := decimal.Zero
	nativeSymbol := ""
	if len(nativeTokens) > 0 {
		nativePrice = nativeTokens[0].TokenPrice
		nativeSymbol = nativeTokens[0].Symbol
	}
	nativeValue := nativeBalance.Mul(nativePrice)
	totalValue = totalValue.Add(nativeValue)
	client.AaAsset.Create().SetAssetValue(totalValue).
		SetLastTime(time.Now().UnixMilli()).SetCreateTime(time.Now()).SetUpdateTime(time.Now()).
		SetID(userAddress).SetNetwork(network).SetBalance(nativeBalance).Save(ctx)
	nativeDetails, err := client.AaAssetDetail.Query().Where(aaassetdetail.UserAddressEqualFold(userAddress), aaassetdetail.IsNativeEqualFold("true")).All(ctx)
	if err != nil {
		return
	}
	if len(nativeDetails) > 0 {
		client.AaAssetDetail.Update().SetAssetAmount(nativeBalance).SetAssetValue(nativeValue).SetLastTime(time.Now().UnixMilli()).Where(aaassetdetail.IDEQ(nativeDetails[0].ID)).Exec(ctx)
	} else {
		client.AaAssetDetail.Create().SetAssetValue(nativeValue).SetLastTime(time.Now().UnixMilli()).SetCreateTime(time.Now()).SetAssetAmount(nativeBalance).SetUserAddress(userAddress).SetContractAddress("").SetSymbol(nativeSymbol).SetNetwork(network).SetIsNative("true").SetUpdateTime(time.Now()).Save(ctx)
	}
	if nativeBalance.Cmp(decimal.Zero) > 0 {
		logger.Info("AddOne update balance success, ", "userAddress", userAddress, "network", network)
	} else {
		logger.Info("AddOne update balance success empty, ", "userAddress", userAddress, "network", network)
	}
}

func getBalance(ctx context.Context, tokenAddress string, userAddress string, decimals int64, web3 *web3.Web3) decimal.Decimal {

	contract, err := web3.Eth.NewContract(config.Abi, tokenAddress)

	res, err := contract.Call("balanceOf", common.HexToAddress(userAddress))
	if err != nil {
		logger.Error("GetBalance balanceOf err ", "msg", err)
		return decimal.Zero
	}
	value, ok := res.(*big.Int)
	if ok {
		balance := decimal.NewFromBigInt(value, 0).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(decimals)))
		return balance
	}
	return decimal.Zero
}

func addOrUpdateAssetDetail(ctx context.Context, contractAddress string, userAddress string, value decimal.Decimal, client *ent.Client, network string, symbol string, amount decimal.Decimal) {
	details, err := client.AaAssetDetail.Query().Where(aaassetdetail.UserAddressEqualFold(userAddress), aaassetdetail.ContractAddressEqualFold(contractAddress)).All(ctx)
	if err != nil {
		return
	}
	if len(details) > 0 {
		detail := details[0]
		client.AaAssetDetail.Update().SetAssetValue(value).SetLastTime(time.Now().UnixMilli()).Where(aaassetdetail.IDEQ(detail.ID)).Exec(ctx)
	} else {
		detail := client.AaAssetDetail.Create().SetAssetValue(value).SetLastTime(time.Now().UnixMilli()).SetCreateTime(time.Now()).SetUpdateTime(time.Now()).SetNetwork(network).SetContractAddress(contractAddress).SetSymbol(symbol).SetUserAddress(userAddress).SetAssetAmount(amount).SetIsNative("false")
		_, err := detail.Save(ctx)
		if err != nil {
			logger.Error("addOrUpdateAssetDetail err ", "symbol", symbol, "msg", err)
		}
	}
}
