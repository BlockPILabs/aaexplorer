package task

import (
	"context"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aauseropsinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/bundlerstatishour"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/dailystatistichour"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/factorystatishour"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/paymasterstatishour"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/taskrecord"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/tokenpriceinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactiondecode"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/transactionreceiptdecode"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/userassetinfo"
	"github.com/BlockPILabs/aa-scan/service"
	"github.com/BlockPILabs/aa-scan/third/moralis"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"log"
	"math"
	"math/big"
	"time"
)

func InitHourStatis() {
	hourScheduler := chrono.NewDefaultTaskScheduler()
	_, err := hourScheduler.ScheduleWithCron(func(ctx context.Context) {
		doHourStatistic()
	}, "0 5 * * * *")

	if err == nil {
		log.Print("hourStatistic has been scheduled")
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
		network := record.ID
		log.Printf("hour-statistic start, network:%s", network)
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}

		taskRecords, err := client.TaskRecord.Query().Where(taskrecord.TaskTypeEQ("hour"), taskrecord.NetworkEQ(network)).Limit(1).All(context.Background())
		if len(taskRecords) == 0 {
			continue
		}
		lastTime := taskRecords[0].LastTime
		if lastTime.Add(1*time.Hour).Compare(time.Now()) > 0 {
			continue
		}
		startTime := time.Date(lastTime.Year(), lastTime.Month(), lastTime.Day(), lastTime.Hour()+1, 0, 0, 0, lastTime.Location())
		endTime := time.Date(lastTime.Year(), lastTime.Month(), lastTime.Day(), lastTime.Hour()+2, 0, 0, 0, lastTime.Location())
		now := time.Now()
		for {
			if startTime.Compare(now) > 0 {
				break
			}
			opsInfos, err := client.AAUserOpsInfo.Query().
				Where(
					aauseropsinfo.TimeGTE(startTime),
					aauseropsinfo.TimeLT(endTime)).
				All(context.Background())

			if err != nil {
				log.Println(err)
				continue
			}

			txCount, err := client.TransactionDecode.Query().Where(transactiondecode.TimeGTE(startTime), transactiondecode.TimeLT(endTime)).Count(context.Background())

			bundlerMap := make(map[string]map[string][]*ent.AAUserOpsInfo)
			paymasterMap := make(map[string]map[string][]*ent.AAUserOpsInfo)
			factoryMap := make(map[string]map[string][]*ent.AAUserOpsInfo)

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
			//receiptMap := getReceiptMap(receipts)
			earnMap := getEarnMap(receiveMap, costMap)
			txHashes := make(map[string]map[string]bool)

			for _, opsInfo := range opsInfos {
				addOpsInfo(opsInfo.Bundler, opsInfo, bundlerMap, "hour")
				addOpsInfo(opsInfo.Paymaster, opsInfo, paymasterMap, "hour")
				addOpsInfo(opsInfo.Factory, opsInfo, factoryMap, "hour")
				addTxHash(opsInfo.TxHash, opsInfo, txHashes, "hour")
			}

			dailyStatisticHours := calHourStatistic(client, opsInfos, txHashes, network, txCount, startTime)

			bundlerList := calBundlerStatistic(client, bundlerMap, startTime, earnMap, totalBundleMap, network)
			paymasterList := calPaymasterStatistic(client, paymasterMap, startTime, network)
			factoryList := calFactoryStatis(client, factoryMap, startTime, network)

			bulkInsertBundlerStatsHour(context.Background(), client, bundlerList)
			bulkInsertPaymasterStatsHour(context.Background(), client, paymasterList)
			bulkInsertFactoryStatsHour(context.Background(), client, factoryList)
			bulkInsertDailyStatisticHour(context.Background(), client, dailyStatisticHours)
			//saveWhaleStatisticHour(context.Background(), client, startTime)
			client.TaskRecord.Update().SetLastTime(startTime).Where(taskrecord.IDEQ(taskRecords[0].ID)).Exec(context.Background())
			startTime = startTime.Add(1 * time.Hour)
			endTime = endTime.Add(1 * time.Hour)
			log.Printf("hour task statistic success, day:%s", startTime.String())
		}
	}

}

func bulkInsertDailyStatisticHour(ctx context.Context, client *ent.Client, data []*ent.DailyStatisticHourCreate) error {
	if len(data) == 0 {
		return nil
	}

	for _, one := range data {
		mutation := one.Mutation()
		time, _ := mutation.StatisticTime()
		network, _ := mutation.Network()
		hours, err := client.DailyStatisticHour.Query().Where(dailystatistichour.StatisticTimeEQ(time), dailystatistichour.NetworkEQ(network)).All(context.Background())
		if err != nil {
			continue
		}
		if len(hours) != 0 {
			client.DailyStatisticHour.Delete().Where(dailystatistichour.IDEQ(hours[0].ID)).Exec(context.Background())
		}
		one.Save(context.Background())
	}
	return nil
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
		cost = cost.Add(GetReceiptGasRayDiv(receipt))
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

func GetReceiptGasRayDiv(receipt *ent.TransactionReceiptDecode) decimal.Decimal {
	var gasPrice big.Int
	_, success := gasPrice.SetString(receipt.EffectiveGasPrice, 0)
	if !success {
		log.Printf("GetReceiptGasRayDiv convert err, %s", receipt.ID)
		return decimal.Zero
	}
	return receipt.GasUsed.Mul(RayDiv(decimal.NewFromInt(gasPrice.Int64())))
}
func GetReceiptGas(receipt *ent.TransactionReceiptDecode) decimal.Decimal {
	var gasPrice big.Int
	_, success := gasPrice.SetString(receipt.EffectiveGasPrice, 0)
	if !success {
		log.Printf("GetReceiptGasRayDiv convert err, %s", receipt.ID)
		return decimal.Zero
	}
	return receipt.GasUsed.Mul(decimal.NewFromInt(gasPrice.Int64()))
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

func calHourStatistic(client *ent.Client, infos []*ent.AAUserOpsInfo, allTxHashes map[string]map[string]bool, network string, txCount int, startTime time.Time) []*ent.DailyStatisticHourCreate {
	if len(infos) == 0 {
		return nil
	}

	var resp []*ent.DailyStatisticHourCreate
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
		dailyStatistic := client.DailyStatisticHour.Create().
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

func RayDiv(gas decimal.Decimal) decimal.Decimal {
	return gas.DivRound(decimal.NewFromFloat(math.Pow10(18)), 18)
}

func calBundlerStatistic(client *ent.Client, bundlerMap map[string]map[string][]*ent.AAUserOpsInfo, startTime time.Time, earnMap map[string]decimal.Decimal, bundleMap map[string]map[string]int, network string) []*ent.BundlerStatisHourCreate {

	var bundlers []*ent.BundlerStatisHourCreate
	for key, allTimeUserOpsInfoList := range bundlerMap {
		if len(allTimeUserOpsInfoList) == 0 {
			continue
		}
		for statisticTime, userOpsInfoList := range allTimeUserOpsInfoList {
			sTime, err := time.Parse(TimeLayout, statisticTime)
			if err != nil {
				log.Println(err)
			}
			txHashMap := make(map[string]bool)
			var successMumMap = make(map[string]bool)
			var failedNumMap = make(map[string]bool)
			var totalFee = decimal.Zero
			for _, userOpsInfo := range userOpsInfoList {
				totalFee = totalFee.Add(userOpsInfo.Fee)
				if userOpsInfo.Status == 1 {
					successMumMap[userOpsInfo.TxHash] = true
					txHashMap[userOpsInfo.TxHash] = true
				} else {
					failedNumMap[userOpsInfo.TxHash] = true
				}
			}
			//set properties
			earnFee := earnMap[key]
			bundlers = append(bundlers, client.BundlerStatisHour.Create().
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

func calPaymasterStatistic(client *ent.Client, bundlerMap map[string]map[string][]*ent.AAUserOpsInfo, startTime time.Time, network string) []*ent.PaymasterStatisHourCreate {

	var paymasters []*ent.PaymasterStatisHourCreate
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
			paymasters = append(paymasters, client.PaymasterStatisHour.Create().
				SetPaymaster(key).
				SetNetwork(network).
				SetUserOpsNum(int64(len(userOpsInfoList))).
				SetGasSponsored(totalFee).
				SetReserve(nativeBalance).
				SetReserveUsd(price.Mul(nativeBalance)).
				SetStatisTime(sTime),
			)
		}

	}

	return paymasters
}

func calFactoryStatis(client *ent.Client, factoryMap map[string]map[string][]*ent.AAUserOpsInfo, startTime time.Time, network string) []*ent.FactoryStatisHourCreate {

	var factories []*ent.FactoryStatisHourCreate
	for key, allTimeUserOpsInfoList := range factoryMap {
		if len(key) == 0 {
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
				if userOpsInfo.Factory != "" {
					accountDeployNum++
				}
			}
			factories = append(factories, client.FactoryStatisHour.Create().
				SetFactory(key).
				SetNetwork(network).
				SetStatisTime(sTime).
				SetAccountNum(int64(len(accountMap))).
				SetAccountDeployNum(int64(accountDeployNum)),
			)
		}

	}

	return factories
}

func bulkInsertBundlerStatsHour(ctx context.Context, client *ent.Client, data []*ent.BundlerStatisHourCreate) error {
	if len(data) == 0 {
		return nil
	}

	for _, one := range data {
		mutation := one.Mutation()
		bundler, _ := mutation.Bundler()
		time, _ := mutation.StatisTime()
		network, _ := mutation.Network()
		bundlerDays, err := client.BundlerStatisHour.Query().Where(bundlerstatishour.BundlerEqualFold(bundler), bundlerstatishour.StatisTimeEQ(time), bundlerstatishour.NetworkEQ(network)).All(context.Background())
		if err != nil {
			continue
		}
		if len(bundlerDays) != 0 {
			client.BundlerStatisHour.Delete().Where(bundlerstatishour.IDEQ(bundlerDays[0].ID)).Exec(context.Background())
		}
		one.Save(context.Background())
	}
	return nil
}

func bulkInsertPaymasterStatsHour(ctx context.Context, client *ent.Client, data []*ent.PaymasterStatisHourCreate) error {
	if len(data) == 0 {
		return nil
	}

	for _, one := range data {
		mutation := one.Mutation()
		paymaster, _ := mutation.Paymaster()
		time, _ := mutation.StatisTime()
		network, _ := mutation.Network()
		bundlerDays, err := client.PaymasterStatisHour.Query().Where(paymasterstatishour.PaymasterEqualFold(paymaster), paymasterstatishour.StatisTimeEQ(time), paymasterstatishour.NetworkEQ(network)).All(context.Background())
		if err != nil {
			continue
		}
		if len(bundlerDays) != 0 {
			client.PaymasterStatisHour.Delete().Where(paymasterstatishour.IDEQ(bundlerDays[0].ID)).Exec(context.Background())
		}
		one.Save(context.Background())
	}
	return nil
}

func bulkInsertFactoryStatsHour(ctx context.Context, client *ent.Client, data []*ent.FactoryStatisHourCreate) error {
	if len(data) == 0 {
		return nil
	}

	if len(data) == 0 {
		return nil
	}

	for _, one := range data {
		mutation := one.Mutation()
		factory, _ := mutation.Factory()
		time, _ := mutation.StatisTime()
		network, _ := mutation.Network()
		bundlerDays, err := client.FactoryStatisHour.Query().Where(factorystatishour.FactoryEqualFold(factory), factorystatishour.StatisTimeEQ(time), factorystatishour.NetworkEQ(network)).All(context.Background())
		if err != nil {
			continue
		}
		if len(bundlerDays) != 0 {
			client.FactoryStatisHour.Delete().Where(factorystatishour.IDEQ(bundlerDays[0].ID)).Exec(context.Background())
		}
		one.Save(context.Background())
	}
	return nil
}
