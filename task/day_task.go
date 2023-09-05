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
	"github.com/BlockPILabs/aa-scan/parser"
	"github.com/BlockPILabs/aa-scan/service"
	"github.com/BlockPILabs/aa-scan/third/moralis"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

func InitDayStatis() {
	dayScheduler := chrono.NewDefaultTaskScheduler()
	doDayStatistic()
	_, err := dayScheduler.ScheduleWithCron(func(ctx context.Context) {
		doDayStatistic()
	}, "0 15 0 * * ?")

	if err == nil {
		log.Print("dayStatis has been scheduled")
	}

}

func doDayStatistic() {
	cli, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	records, err := cli.Network.Query().All(context.Background())
	if len(records) == 0 {
		return
	}
	for _, record := range records {
		network := record.Network
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day()-100, 0, 0, 0, 0, now.Location())
		endTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
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
			receive = receive.Add(RayDiv(decimal.NewFromInt(opsInfo.ActualGasCost)))
			receiveMap[bundler] = receive

		}

		hashs := getKeySlice(txHashMap)
		receipts, err := client.TransactionReceiptDecode.Query().Where(transactionreceiptdecode.IDIn(hashs[:]...)).All(context.Background())
		costMap := getCostMap(receipts)
		earnMap := getEarnMap(receiveMap, costMap)

		bundlerMap := make(map[string][]*ent.AAUserOpsInfo)
		paymasterMap := make(map[string][]*ent.AAUserOpsInfo)
		factoryMap := make(map[string][]*ent.AAUserOpsInfo)
		txHashes := make(map[string]bool)

		for _, opsInfo := range opsInfos {
			addOpsInfo(opsInfo.Bundler, opsInfo, bundlerMap)
			addOpsInfo(opsInfo.Paymaster, opsInfo, paymasterMap)
			addOpsInfo(opsInfo.Factory, opsInfo, factoryMap)
			txHashes[opsInfo.TxHash] = true
		}

		dailyStatisticDay := calDailyStatistic(client, opsInfos, txHashes, network, startTime)

		bundlerList := calBundlerStatisDay(client, bundlerMap, earnMap, totalBundleMap, startTime, network)
		paymasterList := calPaymasterStatisDay(client, paymasterMap, startTime, network)
		factoryList := calFactoryStatisDay(client, factoryMap, startTime, network)

		bulkInsertBundlerStatsDay(context.Background(), client, bundlerList)
		bulkInsertPaymasterStatsDay(context.Background(), client, paymasterList)
		bulkInsertFactoryStatsDay(context.Background(), client, factoryList)
		dailyStatisticDay.Save(context.Background())

		//saveWhaleStatisticDay(context.Background(), client, startTime)
	}

}

func saveWhaleStatisticDay(ctx context.Context, client *ent.Client, time time.Time) {
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
	whaleDay := client.WhaleStatisticDay.Create().SetWhaleNum(int64(addrCount)).SetTotalUsd(totalUsd).SetNetwork("").SetStatisticTime(time)
	whaleDay.Save(ctx)

}

func calDailyStatistic(client *ent.Client, infos []*ent.AAUserOpsInfo, txHashes map[string]bool, network string, startTime time.Time) *ent.DailyStatisticDayCreate {
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
		if receipt.CumulativeGasUsed != nil {
			spentGas = spentGas.Sub(RayDiv(*receipt.CumulativeGasUsed))
		}
	}

	var totalGasFee = decimal.Zero
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
	dailyStatistic := client.DailyStatisticDay.Create().
		SetNetwork(network).
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

func bulkInsertFactoryStatsDay(ctx context.Context, client *ent.Client, data []*ent.FactoryStatisDayCreate) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	if _, err := client.FactoryStatisDay.CreateBulk(data...).Save(ctx); err != nil {
		tx.Rollback()
		log.Fatal(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func bulkInsertPaymasterStatsDay(ctx context.Context, client *ent.Client, data []*ent.PaymasterStatisDayCreate) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	if _, err := client.PaymasterStatisDay.CreateBulk(data...).Save(ctx); err != nil {
		tx.Rollback()
		log.Fatal(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func bulkInsertBundlerStatsDay(ctx context.Context, client *ent.Client, data []*ent.BundlerStatisDayCreate) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	err = client.BundlerStatisDay.CreateBulk(data...).Exec(ctx)
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

func calBundlerStatisDay(client *ent.Client, bundlerMap map[string][]*ent.AAUserOpsInfo, earnMap map[string]decimal.Decimal, bundleMap map[string]map[string]int, startTime time.Time, network string) []*ent.BundlerStatisDayCreate {
	totalCount := 0
	var totalFee decimal.Decimal

	var bundlers []*ent.BundlerStatisDayCreate
	for key, userOpsInfoList := range bundlerMap {
		totalCount += len(userOpsInfoList)
		txHashMap := make(map[string]bool)
		for _, userOpsInfo := range userOpsInfoList {
			totalFee = totalFee.Add(RayDiv(userOpsInfo.Fee))
			txHashMap[userOpsInfo.TxHash] = true
		}
		totalBundleNum := len(bundleMap[key])
		earnFee := earnMap[key]
		bundlers = append(bundlers, client.BundlerStatisDay.Create().
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

func calPaymasterStatisDay(client *ent.Client, bundlerMap map[string][]*ent.AAUserOpsInfo, startTime time.Time, network string) []*ent.PaymasterStatisDayCreate {
	totalCount := 0
	var totalFee decimal.Decimal

	var paymasters []*ent.PaymasterStatisDayCreate
	price := service.GetNativePrice(network)
	for key, userOpsInfoList := range bundlerMap {
		totalCount += len(userOpsInfoList)
		for _, userOpsInfo := range userOpsInfoList {
			totalFee = totalFee.Add(parser.DivRav(userOpsInfo.ActualGasCost))
		}
		nativeBalance := moralis.GetNativeTokenBalance(key, network)
		paymasters = append(paymasters, client.PaymasterStatisDay.Create().
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

func calFactoryStatisDay(client *ent.Client, bundlerMap map[string][]*ent.AAUserOpsInfo, startTime time.Time, network string) []*ent.FactoryStatisDayCreate {
	accountDeployNum := 0

	var factories []*ent.FactoryStatisDayCreate
	for key, userOpsInfoList := range bundlerMap {
		accountMap := make(map[string]bool)
		for _, userOpsInfo := range userOpsInfoList {
			accountMap[userOpsInfo.Sender] = true
			if userOpsInfo.Factory != "" {
				accountDeployNum++
			}
		}
		factories = append(factories, client.FactoryStatisDay.Create().
			SetFactory(key).
			SetNetwork(network).
			SetStatisTime(startTime).
			SetAccountNum(int64(len(accountMap))).
			SetAccountDeployNum(int64(accountDeployNum)),
		)
	}

	return factories
}
