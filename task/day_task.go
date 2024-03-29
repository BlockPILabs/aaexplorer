package task

import (
	"context"
	internalconfig "github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aauseropsinfo"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/bundlerstatisday"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/dailystatisticday"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/factorystatisday"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/paymasterstatisday"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/taskrecord"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/tokenpriceinfo"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/transactiondecode"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/transactionreceiptdecode"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/userassetinfo"
	"github.com/BlockPILabs/aaexplorer/service"
	"github.com/BlockPILabs/aaexplorer/third/moralis"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

const TimeLayout = "2006-01-02 15:04:05"

func InitDayStatis() {
	dayScheduler := chrono.NewDefaultTaskScheduler()
	_, err := dayScheduler.ScheduleWithCron(func(ctx context.Context) {
		doDayStatistic()
	}, "0 15 0 * * *")
	if err == nil {
		log.Print("dayStatistic has been scheduled")
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
		network := record.ID
		log.Printf("day-statistic start, network:%s", network)
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		taskRecords, err := client.TaskRecord.Query().Where(taskrecord.TaskTypeEQ("day"), taskrecord.NetworkEQ(network)).Limit(1).All(context.Background())
		if len(taskRecords) == 0 {
			continue
		}
		lastTime := taskRecords[0].LastTime
		now := time.Now()
		dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if lastTime.Add(24*time.Hour).Compare(dayStart) >= 0 {
			continue
		}
		startTime := time.Date(lastTime.Year(), lastTime.Month(), lastTime.Day()+1, 0, 0, 0, 0, lastTime.Location())
		endTime := time.Date(lastTime.Year(), lastTime.Month(), lastTime.Day()+2, 0, 0, 0, 0, lastTime.Location())
		for {
			if startTime.Compare(dayStart) >= 0 {
				break
			}
			opsInfos, err := client.AAUserOpsInfo.Query().
				Where(
					aauseropsinfo.TxTimeGTE(startTime.Unix()),
					aauseropsinfo.TxTimeLT(endTime.Unix())).
				All(context.Background())
			if err != nil {
				log.Println(err)
				break
			}

			txCount, err := client.TransactionDecode.Query().Where(transactiondecode.TimeGTE(startTime), transactiondecode.TimeLT(endTime)).Count(context.Background())

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

			bundlerMap := make(map[string]map[string][]*ent.AAUserOpsInfo)
			paymasterMap := make(map[string]map[string][]*ent.AAUserOpsInfo)
			factoryMap := make(map[string]map[string][]*ent.AAUserOpsInfo)
			txHashes := make(map[string]map[string]bool)

			for _, opsInfo := range opsInfos {
				addOpsInfo(opsInfo.Bundler, opsInfo, bundlerMap, "day")
				addOpsInfo(opsInfo.Paymaster, opsInfo, paymasterMap, "day")
				addOpsInfo(opsInfo.Factory, opsInfo, factoryMap, "day")

				addTxHash(opsInfo.TxHash, opsInfo, txHashes, "day")
			}

			dailyStatisticDays := calDailyStatistic(client, opsInfos, txHashes, network, txCount, startTime)

			bundlerList := calBundlerStatisDay(client, bundlerMap, earnMap, totalBundleMap, startTime, network)
			paymasterList := calPaymasterStatisDay(client, paymasterMap, startTime, network)
			factoryList := calFactoryStatisDay(client, factoryMap, startTime, network)

			bulkInsertBundlerStatsDay(context.Background(), client, bundlerList)
			bulkInsertPaymasterStatsDay(context.Background(), client, paymasterList)
			bulkInsertFactoryStatsDay(context.Background(), client, factoryList)
			bulkInsertDailyStatistic(context.Background(), client, dailyStatisticDays)

			//saveWhaleStatisticDay(context.Background(), client, startTime)
			client.TaskRecord.Update().SetLastTime(startTime).Where(taskrecord.IDEQ(taskRecords[0].ID)).Exec(context.Background())
			startTime = startTime.Add(24 * time.Hour)
			endTime = endTime.Add(24 * time.Hour)
			log.Printf("day task statistic success, day:%s", startTime.String())
		}
	}

}

func bulkInsertDailyStatistic(ctx context.Context, client *ent.Client, data []*ent.DailyStatisticDayCreate) error {
	if len(data) == 0 {
		return nil
	}

	for _, one := range data {
		mutation := one.Mutation()
		time, _ := mutation.StatisticTime()
		network, _ := mutation.Network()
		days, err := client.DailyStatisticDay.Query().Where(dailystatisticday.StatisticTimeEQ(time), dailystatisticday.NetworkEQ(network)).All(context.Background())
		if err != nil {
			continue
		}
		if len(days) != 0 {
			client.DailyStatisticDay.Delete().Where(dailystatisticday.IDEQ(days[0].ID)).Exec(context.Background())
		}
		one.Save(context.Background())
	}
	return nil
}

func addTxHash(hash string, info *ent.AAUserOpsInfo, hashMap map[string]map[string]bool, s string) {
	var startOf = ""
	if s == "day" {
		startOf = getDayStart(info.Time)
	} else {
		startOf = getHourStart(info.Time)
	}
	timeHash, timeHashOk := hashMap[startOf]
	if !timeHashOk {
		timeHash = make(map[string]bool)
	}
	timeHash[hash] = true
	hashMap[startOf] = timeHash
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
		if value.Cmp(decimal.NewFromInt(internalconfig.WhaleUsd)) < 0 {
			continue
		}
		addrCount += 1
		totalUsd = totalUsd.Add(value)

	}
	whaleDay := client.WhaleStatisticDay.Create().SetWhaleNum(int64(addrCount)).SetTotalUsd(totalUsd).SetNetwork("").SetStatisticTime(time)
	whaleDay.Save(ctx)

}

func calDailyStatistic(client *ent.Client, infos []*ent.AAUserOpsInfo, allTxHashes map[string]map[string]bool, network string, txCount int, startTime time.Time) []*ent.DailyStatisticDayCreate {
	if len(infos) == 0 {
		return nil
	}
	var resp []*ent.DailyStatisticDayCreate
	for statisticTime, txHashes := range allTxHashes {
		sTime, err := time.Parse(TimeLayout, statisticTime)
		if err != nil {
			continue
		}
		if len(txHashes) == 0 {
			continue
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
			spentGas = spentGas.Sub(GetReceiptGasRayDiv(receipt))
		}
		var totalGasFee = decimal.Zero
		var txMap = make(map[string]bool)
		var walletMap = make(map[string]bool)
		var paymasterFee = decimal.Zero
		for _, opsInfo := range infos {
			fee := RayDiv(decimal.NewFromInt(opsInfo.ActualGasCost))
			totalGasFee = totalGasFee.Add(fee)
			txMap[opsInfo.TxHash] = true
			walletMap[opsInfo.Sender] = true
			if len(opsInfo.Paymaster) > 2 {
				paymasterFee = paymasterFee.Add(fee)
			}
			spentGas = spentGas.Add(fee)
		}
		price := service.GetNativePrice(network)
		dailyStatistic := client.DailyStatisticDay.Create().
			SetNetwork(network).
			SetUserOpsNum(int64(len(infos))).
			SetStatisticTime(sTime.UnixMilli()).
			SetActiveWallet(int64(len(walletMap))).
			SetGasFee(totalGasFee).
			SetGasFeeUsd(price.Mul(totalGasFee)).
			SetBundlerGasProfit(spentGas).
			SetBundlerGasProfitUsd(price.Mul(spentGas)).
			SetPaymasterGasPaid(paymasterFee).
			SetPaymasterGasPaidUsd(price.Mul(paymasterFee)).
			SetAaTxNum(int64(len(txMap))).
			SetTxNum(int64(txCount))

		resp = append(resp, dailyStatistic)
	}

	return resp
}

func bulkInsertFactoryStatsDay(ctx context.Context, client *ent.Client, data []*ent.FactoryStatisDayCreate) error {
	if len(data) == 0 {
		return nil
	}

	for _, one := range data {
		mutation := one.Mutation()
		factory, _ := mutation.Factory()
		time, _ := mutation.StatisTime()
		network, _ := mutation.Network()
		factoryDays, err := client.FactoryStatisDay.Query().Where(factorystatisday.FactoryEqualFold(factory), factorystatisday.StatisTimeEQ(time), factorystatisday.NetworkEQ(network)).All(context.Background())
		if err != nil {
			continue
		}
		if len(factoryDays) != 0 {
			for _, old := range factoryDays {
				client.FactoryStatisDay.Delete().Where(factorystatisday.IDEQ(old.ID)).Exec(context.Background())
			}

		}
		one.Save(context.Background())
	}
	return nil
	/**
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	if _, err := client.FactoryStatisDay.CreateBulk(data...).Save(ctx); err != nil {
		tx.Rollback()
		log.Println(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	*/

}

func bulkInsertPaymasterStatsDay(ctx context.Context, client *ent.Client, data []*ent.PaymasterStatisDayCreate) error {
	if len(data) == 0 {
		return nil
	}

	for _, one := range data {
		mutation := one.Mutation()
		paymaster, _ := mutation.Paymaster()
		time, _ := mutation.StatisTime()
		network, _ := mutation.Network()
		paymasterDays, err := client.PaymasterStatisDay.Query().Where(paymasterstatisday.PaymasterEqualFold(paymaster), paymasterstatisday.StatisTimeEQ(time), paymasterstatisday.NetworkEQ(network)).All(context.Background())
		if err != nil {
			continue
		}
		if len(paymasterDays) != 0 {
			for _, old := range paymasterDays {
				client.PaymasterStatisDay.Delete().Where(paymasterstatisday.IDEQ(old.ID)).Exec(context.Background())
			}
		}
		one.Save(context.Background())
		log.Printf("paymaster-day-task statistic success, paymaster: %s, day:%s", paymaster, time.String())
	}
	return nil
}

func bulkInsertBundlerStatsDay(ctx context.Context, client *ent.Client, data []*ent.BundlerStatisDayCreate) error {
	if len(data) == 0 {
		return nil
	}

	for _, one := range data {
		mutation := one.Mutation()
		bundler, _ := mutation.Bundler()
		time, _ := mutation.StatisTime()
		network, _ := mutation.Network()
		bundlerDays, err := client.BundlerStatisDay.Query().Where(bundlerstatisday.BundlerEqualFold(bundler), bundlerstatisday.StatisTimeEQ(time), bundlerstatisday.NetworkEQ(network)).All(context.Background())
		if err != nil {
			continue
		}
		if len(bundlerDays) != 0 {
			for _, old := range bundlerDays {
				client.BundlerStatisDay.Delete().Where(bundlerstatisday.IDEQ(old.ID)).Exec(ctx)
			}
		}
		one.Save(context.Background())
		log.Printf("bundler-day-task statistic success, bundler: %s, day:%s", bundler, time.String())
	}
	return nil
}

func calBundlerStatisDay(client *ent.Client, bundlerMap map[string]map[string][]*ent.AAUserOpsInfo, earnMap map[string]decimal.Decimal, bundleMap map[string]map[string]int, startTime time.Time, network string) []*ent.BundlerStatisDayCreate {

	var bundlers []*ent.BundlerStatisDayCreate

	for key, allTimeUserOpsInfoList := range bundlerMap {
		if len(allTimeUserOpsInfoList) == 0 {
			continue
		}
		for statisticTime, userOpsInfoList := range allTimeUserOpsInfoList {
			sTime, err := time.Parse(TimeLayout, statisticTime)
			if err != nil {
				log.Println(err)
			}
			var successMumMap = make(map[string]bool)
			var failedNumMap = make(map[string]bool)
			var totalFee = decimal.Zero
			for _, userOpsInfo := range userOpsInfoList {
				totalFee = totalFee.Add(RayDiv(userOpsInfo.Fee))
				if userOpsInfo.Status == 1 {
					successMumMap[userOpsInfo.TxHash] = true
				} else {
					failedNumMap[userOpsInfo.TxHash] = true
				}
			}
			earnFee := earnMap[key]
			bundlers = append(bundlers, client.BundlerStatisDay.Create().
				SetBundler(key).
				SetNetwork(network).
				SetBundlesNum(int64(len(successMumMap))).
				SetGasCollected(totalFee).
				SetTotalNum(int64(len(successMumMap))+int64(len(failedNumMap))).
				SetFeeEarned(earnFee).
				SetUserOpsNum(int64(len(userOpsInfoList))).
				SetStatisTime(sTime).
				SetSuccessBundlesNum(int64(len(successMumMap))).
				SetFailedBundlesNum(int64(len(failedNumMap))),
			)
		}

	}

	return bundlers
}

func calPaymasterStatisDay(client *ent.Client, bundlerMap map[string]map[string][]*ent.AAUserOpsInfo, startTime time.Time, network string) []*ent.PaymasterStatisDayCreate {

	var paymasters []*ent.PaymasterStatisDayCreate
	price := service.GetNativePrice(network)
	for key, allTimeUserOpsInfoList := range bundlerMap {
		if len(key) <= 2 {
			continue
		}
		if len(allTimeUserOpsInfoList) == 0 {
			continue
		}
		for statisticTime, userOpsInfoList := range allTimeUserOpsInfoList {
			sTime, err := time.Parse(TimeLayout, statisticTime)
			if err != nil {
				log.Println(err)
			}

			var totalFee = decimal.Zero
			for _, userOpsInfo := range userOpsInfoList {
				totalFee = totalFee.Add(RayDiv(decimal.NewFromInt(userOpsInfo.ActualGasCost)))
			}
			nativeBalance := moralis.GetNativeTokenBalance(key, network)
			paymasters = append(paymasters, client.PaymasterStatisDay.Create().
				SetPaymaster(key).
				SetNetwork(network).
				SetUserOpsNum(int64(len(userOpsInfoList))).
				SetGasSponsored(totalFee).
				SetGasSponsoredUsd(totalFee.Mul(price)).
				SetReserve(nativeBalance).
				SetReserveUsd(price.Mul(nativeBalance)).
				SetStatisTime(sTime),
			)
		}

	}

	return paymasters
}

func calFactoryStatisDay(client *ent.Client, bundlerMap map[string]map[string][]*ent.AAUserOpsInfo, startTime time.Time, network string) []*ent.FactoryStatisDayCreate {

	var factories []*ent.FactoryStatisDayCreate
	for key, allTimeUserOpsInfoList := range bundlerMap {
		if len(key) <= 2 {
			continue
		}
		if len(allTimeUserOpsInfoList) == 0 {
			continue
		}
		for statisticTime, userOpsInfoList := range allTimeUserOpsInfoList {
			sTime, err := time.Parse(TimeLayout, statisticTime)
			if err != nil {
				log.Println(err)
			}
			accountDeployNum := 0
			accountMap := make(map[string]bool)
			for _, userOpsInfo := range userOpsInfoList {
				accountMap[userOpsInfo.Sender] = true
				if len(userOpsInfo.Factory) > 0 {
					accountDeployNum++
				}
			}
			factories = append(factories, client.FactoryStatisDay.Create().
				SetFactory(key).
				SetNetwork(network).
				SetStatisTime(sTime).
				SetAccountNum(int64(accountDeployNum)).
				SetAccountDeployNum(int64(accountDeployNum)),
			)
		}
	}

	return factories
}
