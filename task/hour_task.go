package task

import (
	"context"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aauseropsinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/tokenpriceinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactionreceiptdecode"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/userassetinfo"
	"github.com/BlockPILabs/aa-scan/service"
	"github.com/BlockPILabs/aa-scan/third/moralis"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"log"
	"math"
	"time"
)

func InitHourStatis() {
	doHourStatistic()
	hourScheduler := chrono.NewDefaultTaskScheduler()

	_, err := hourScheduler.ScheduleWithCron(func(ctx context.Context) {
		doHourStatistic()
	}, "0 5 * * * ?")

	if err == nil {
		log.Print("hourStatis has been scheduled")
	}
}

func doHourStatistic() {
	cli, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	records, err := cli.Network.Query().All(context.Background())
	if len(records) == 0 {
		return
	}
	for _, record := range records {
		network := record.Name
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}

		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-10000, 0, 0, 0, now.Location())
		endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
		opsInfos, err := client.AAUserOpsInfo.
			Query().
			Where(
				aauseropsinfo.TxTimeGTE(startTime.Unix()),
				aauseropsinfo.TxTimeLT(endTime.Unix()),
			).
			All(context.Background())

		if err != nil {
			log.Println(err)
			continue
		}

		bundlerMap := make(map[string][]*ent.AAUserOpsInfo)
		paymasterMap := make(map[string][]*ent.AAUserOpsInfo)
		factoryMap := make(map[string][]*ent.AAUserOpsInfo)

		receiveMap := make(map[string]decimal.Decimal)
		totalBundleMap := make(map[string]map[string]int)
		txHashMap := make(map[string]bool)

		for _, opsInfo := range opsInfos {
			bundler := opsInfo.Bundler
			txHashMap[opsInfo.TxHash] = true

			bundle, bundleOk := totalBundleMap[bundler]
			if !bundleOk {
				bundle = make(map[string]int)
			}
			bundle[opsInfo.TxHash] = 1
			totalBundleMap[bundler] = bundle
			if opsInfo.Status == 0 {
				continue
			}
			receive, receiveOk := receiveMap[bundler]
			if !receiveOk {
				receive = decimal.Zero
			}
			receive = receive.Add(decimal.NewFromInt(opsInfo.ActualGasUsed).DivRound(decimal.NewFromFloat(math.Pow10(18)), 18))
			receiveMap[bundler] = receive

		}

		hashs := getKeySlice(txHashMap)
		receipts, err := client.TransactionReceiptDecode.Query().Where(transactionreceiptdecode.IDIn(hashs[:]...)).All(context.Background())
		costMap := getCostMap(receipts)
		//receiptMap := getReceiptMap(receipts)
		earnMap := getEarnMap(receiveMap, costMap)
		txHashes := make(map[string]bool)

		for _, opsInfo := range opsInfos {
			addOpsInfo(opsInfo.Bundler, opsInfo, bundlerMap)
			addOpsInfo(opsInfo.Paymaster, opsInfo, paymasterMap)
			addOpsInfo(opsInfo.Factory, opsInfo, factoryMap)
			txHashes[opsInfo.TxHash] = true
		}

		dailyStatisticHour := calHourStatistic(client, opsInfos, txHashes, network, startTime)

		bundlerList := calBundlerStatis(client, bundlerMap, startTime, earnMap, totalBundleMap, network)
		paymasterList := calPaymasterStatis(client, paymasterMap, startTime, network)
		factoryList := calFactoryStatis(client, factoryMap, startTime, network)

		bulkInsertBundlerStatsHour(context.Background(), client, bundlerList)
		bulkInsertPaymasterStatsHour(context.Background(), client, paymasterList)
		bulkInsertFactoryStatsHour(context.Background(), client, factoryList)
		dailyStatisticHour.Save(context.Background())
		//saveWhaleStatisticHour(context.Background(), client, startTime)
	}

}

func getBundleRateMap(bundleMap map[string]map[string]int) map[string]decimal.Decimal {
	var total = 0
	for _, value := range bundleMap {
		total += len(value)
	}
	var rateMap = make(map[string]decimal.Decimal)
	for key, value := range bundleMap {
		rate := decimal.NewFromInt(int64(len(value))).DivRound(decimal.NewFromInt(int64(total)), 4)
		rateMap[key] = rate
	}

	return rateMap
}

func getSuccessRateMap(sucMap map[string]map[string]int) map[string]decimal.Decimal {
	var rateMap = make(map[string]decimal.Decimal)
	for key, value := range sucMap {
		sucNum := 0
		totalNum := 0
		for _, status := range value {
			if status == 1 {
				sucNum += 1
			}
			totalNum += 1
		}
		rate := decimal.NewFromInt(int64(sucNum)).DivRound(decimal.NewFromInt(int64(totalNum)), 4)
		rateMap[key] = rate
	}
	return rateMap
}

func getEarnMap(receiveMap map[string]decimal.Decimal, costMap map[string]decimal.Decimal) map[string]decimal.Decimal {
	var earnMap = make(map[string]decimal.Decimal)
	for bundler, receive := range receiveMap {
		cost, costOk := costMap[bundler]
		if !costOk {
			cost = decimal.Zero
		}
		earnMap[bundler] = receive.Sub(cost)
	}

	return earnMap
}

func getCostMap(receipts []*ent.TransactionReceiptDecode) map[string]decimal.Decimal {
	if len(receipts) == 0 {
		return make(map[string]decimal.Decimal)
	}
	var receiptMap = make(map[string]decimal.Decimal)
	for _, receipt := range receipts {
		bundler := receipt.FromAddr
		cost, costOk := receiptMap[bundler]
		if !costOk {
			cost = decimal.Zero
		}
		cost = cost.Add(RayDiv(decimal.NewFromInt(receipt.CumulativeGasUsed)))
		receiptMap[bundler] = cost
	}
	return receiptMap
}

func getReceiptMap(receipts []*ent.TransactionReceiptDecode) map[string]*ent.TransactionReceiptDecode {
	if len(receipts) == 0 {
		return nil
	}
	var receiptMap = make(map[string]*ent.TransactionReceiptDecode)
	for _, receipt := range receipts {
		receiptMap[receipt.ID] = receipt
	}
	return receiptMap
}

func getKeySlice(maps map[string]bool) []string {
	if len(maps) == 0 {
		return nil
	}
	var keys []string
	for key, _ := range maps {
		keys = append(keys, key)
	}
	return keys
}

func saveWhaleStatisticHour(ctx context.Context, client *ent.Client, time time.Time) {
	assetInfos, err := client.UserAssetInfo.Query().Where(userassetinfo.NetworkEQ("")).All(ctx)
	if err != nil {
		return
	}
	if len(assetInfos) == 0 {
		return
	}
	var contractMap = make(map[string]int)
	var contracts []string
	for _, asset := range assetInfos {

		_, ok := contractMap[asset.ContractAddress]
		if !ok {
			contractMap[asset.ContractAddress] = 1
			contracts = append(contracts, asset.ContractAddress)
		}
	}
	priceInfos, err := client.TokenPriceInfo.Query().Where(tokenpriceinfo.ContractAddressIn(contracts[:]...)).All(ctx)
	if err != nil {
		return
	}
	if len(priceInfos) == 0 {
		return
	}
	var priceMap = make(map[string]decimal.Decimal)
	for _, priceInfo := range priceInfos {
		priceMap[priceInfo.ContractAddress] = priceInfo.TokenPrice
	}

	var valueMap = make(map[string]decimal.Decimal)
	for _, asset := range assetInfos {
		preValue, accountOk := valueMap[asset.AccountAddress]
		if !accountOk {
			preValue = decimal.Zero
		}
		price, ok := priceMap[asset.ContractAddress]
		if !ok {
			valueMap[asset.AccountAddress] = preValue
		} else {
			valueMap[asset.AccountAddress] = asset.Amount.Mul(price).Add(preValue)
		}

	}

	var totalUsd = decimal.Zero
	var addrCount = 0
	for _, value := range valueMap {
		if value.Cmp(decimal.NewFromInt(config.WhaleUsd)) < 0 {
			continue
		}
		addrCount += 1
		totalUsd = totalUsd.Add(value)

	}
	whaleHour := client.WhaleStatisticHour.Create().SetWhaleNum(int64(addrCount)).SetTotalUsd(totalUsd).SetNetwork("").SetStatisticTime(time)
	whaleHour.Save(ctx)

}

func calHourStatistic(client *ent.Client, infos []*ent.AAUserOpsInfo, txHashes map[string]bool, network string, startTime time.Time) *ent.DailyStatisticHourCreate {
	if len(infos) == 0 {
		return nil
	}

	var hashes []string
	for key, _ := range txHashes {
		hashes = append(hashes, key)
	}
	receipts, err := client.TransactionReceiptDecode.Query().Where(transactionreceiptdecode.IDIn(hashes[:]...)).All(context.Background())
	if err != nil {
		return nil
	}

	var spentGas = decimal.Zero
	for _, receipt := range receipts {
		if receipt.CumulativeGasUsed != 0 {
			spentGas = spentGas.Sub(RayDiv(decimal.NewFromInt(receipt.CumulativeGasUsed)))
		}
	}

	var totalGasFee decimal.Decimal
	var txMap = make(map[string]bool)
	var walletMap = make(map[string]bool)
	var paymasterFee = decimal.Zero
	for _, opsInfo := range infos {
		fee := RayDiv(opsInfo.Fee)
		totalGasFee = totalGasFee.Add(fee)
		txMap[opsInfo.TxHash] = true
		walletMap[opsInfo.Sender] = true
		if len(opsInfo.Paymaster) != 0 {
			paymasterFee = paymasterFee.Add(fee)
		}
		spentGas = spentGas.Add(fee)
	}

	price := service.GetNativePrice(network)
	dailyStatistic := client.DailyStatisticHour.Create().
		SetNetwork(infos[0].Network).
		SetUserOpsNum(int64(len(infos))).
		SetStatisticTime(startTime).
		SetActiveWallet(int64(len(walletMap))).
		SetGasFee(totalGasFee).
		SetBundlerGasProfit(spentGas).
		SetBundlerGasProfitUsd(price.Mul(spentGas)).
		SetPaymasterGasPaid(paymasterFee).
		SetPaymasterGasPaidUsd(price.Mul(paymasterFee)).
		SetTxNum(int64(len(txMap)))

	return dailyStatistic
}

func RayDiv(gas decimal.Decimal) decimal.Decimal {
	return gas.DivRound(decimal.NewFromFloat(math.Pow10(18)), 18)
}

func calBundlerStatis(client *ent.Client, bundlerMap map[string][]*ent.AAUserOpsInfo, startTime time.Time, earnMap map[string]decimal.Decimal, bundleMap map[string]map[string]int, network string) []*ent.BundlerStatisHourCreate {
	totalCount := 0
	var totalFee decimal.Decimal

	var bundlers []*ent.BundlerStatisHourCreate
	for key, userOpsInfoList := range bundlerMap {
		totalCount += len(userOpsInfoList)
		txHashMap := make(map[string]bool)
		for _, userOpsInfo := range userOpsInfoList {
			totalFee = totalFee.Add(userOpsInfo.Fee)
			txHashMap[userOpsInfo.TxHash] = true
		}
		//set properties
		totalBundleNum := len(bundleMap[key])
		earnFee := earnMap[key]
		bundlers = append(bundlers, client.BundlerStatisHour.Create().
			SetBundler(key).
			SetNetwork(network).
			SetBundlesNum(int64(len(txHashMap))).
			SetGasCollected(totalFee).
			SetTotalNum(int64(totalBundleNum)).
			SetFeeEarned(earnFee).
			SetUserOpsNum(int64(totalCount)).
			SetStatisTime(startTime),
		)
	}

	return bundlers
}

func calPaymasterStatis(client *ent.Client, bundlerMap map[string][]*ent.AAUserOpsInfo, startTime time.Time, network string) []*ent.PaymasterStatisHourCreate {
	totalCount := 0
	var totalFee decimal.Decimal

	var paymasters []*ent.PaymasterStatisHourCreate
	price := service.GetNativePrice(network)
	for key, userOpsInfoList := range bundlerMap {
		totalCount += len(userOpsInfoList)
		for _, userOpsInfo := range userOpsInfoList {
			totalFee = totalFee.Add(RayDiv(decimal.NewFromInt(userOpsInfo.ActualGasCost)))
		}
		nativeBalance := moralis.GetNativeTokenBalance(key, network)
		paymasters = append(paymasters, client.PaymasterStatisHour.Create().
			SetPaymaster(key).
			SetNetwork(network).
			SetUserOpsNum(int64(totalCount)).
			SetGasSponsored(totalFee).
			SetReserve(nativeBalance).
			SetReserveUsd(price.Mul(nativeBalance)).
			SetStatisTime(startTime),
		)
	}

	return paymasters
}

func calFactoryStatis(client *ent.Client, bundlerMap map[string][]*ent.AAUserOpsInfo, startTime time.Time, network string) []*ent.FactoryStatisHourCreate {
	accountDeployNum := 0

	var factories []*ent.FactoryStatisHourCreate
	for key, userOpsInfoList := range bundlerMap {
		accountMap := make(map[string]bool)
		for _, userOpsInfo := range userOpsInfoList {
			accountMap[userOpsInfo.Sender] = true
			if userOpsInfo.Factory != "" {
				accountDeployNum++
			}
		}
		factories = append(factories, client.FactoryStatisHour.Create().
			SetFactory(key).
			SetNetwork(network).
			SetStatisTime(startTime).
			SetAccountNum(int64(len(accountMap))).
			SetAccountDeployNum(int64(accountDeployNum)),
		)
	}

	return factories
}

func bulkInsertBundlerStatsHour(ctx context.Context, client *ent.Client, data []*ent.BundlerStatisHourCreate) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	err = client.BundlerStatisHour.CreateBulk(data...).Exec(ctx)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func bulkInsertPaymasterStatsHour(ctx context.Context, client *ent.Client, data []*ent.PaymasterStatisHourCreate) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	if _, err := client.PaymasterStatisHour.CreateBulk(data...).Save(ctx); err != nil {
		tx.Rollback()
		log.Fatal(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func bulkInsertFactoryStatsHour(ctx context.Context, client *ent.Client, data []*ent.FactoryStatisHourCreate) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	if _, err := client.FactoryStatisHour.CreateBulk(data...).Save(ctx); err != nil {
		tx.Rollback()
		log.Fatal(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
