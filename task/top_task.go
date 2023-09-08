package task

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/bundlerinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/bundlerstatishour"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/factoryinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/factorystatishour"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/paymasterinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/paymasterstatishour"
	"github.com/BlockPILabs/aa-scan/service"
	"github.com/BlockPILabs/aa-scan/third/moralis"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

func InitTask() {

	//hour statistics
	InitHourStatis()

	//day statistics
	InitDayStatis()

	TopBundlers()

	TopPaymaster()

	TopFactories()

	UserOpTypeTask()

	AAContractInteractTask()

}

func addOpsInfo(key string, opsInfo *ent.AAUserOpsInfo, bundlerMap map[string][]*ent.AAUserOpsInfo) {
	bundlerOps, bundlerOk := bundlerMap[key]
	if !bundlerOk {
		bundlerOps = []*ent.AAUserOpsInfo{}
	}

	bundlerOps = append(bundlerOps, opsInfo)
	bundlerMap[key] = bundlerOps
}

func TopFactories() {
	doTopFactoryHour()
	factoryScheduler := chrono.NewDefaultTaskScheduler()

	_, err := factoryScheduler.ScheduleWithCron(func(ctx context.Context) {
		doTopFactoryHour()
	}, "0 7 0 * * ?")

	if err == nil {
		log.Print("TopFactory has been scheduled")
	}

}

func doTopFactoryHour() {
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
		startTime := time.Date(now.Year(), now.Month(), now.Day()-500, now.Hour()-23, 0, 0, 0, now.Location())
		endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
		factoryStatisHours, err := client.FactoryStatisHour.
			Query().
			Where(
				factorystatishour.StatisTimeGTE(startTime),
				factorystatishour.StatisTimeLT(endTime),
			).
			All(context.Background())

		if err != nil {
			log.Println(err)
			continue
		}
		if len(factoryStatisHours) == 0 {
			continue
		}

		factoryInfoMap := make(map[string]*ent.FactoryInfo)
		for _, factory := range factoryStatisHours {
			factoryAddr := factory.Factory
			factoryInfo, bundlerInfoOk := factoryInfoMap[factoryAddr]
			if bundlerInfoOk {

				factoryInfo.AccountDeployNum = factoryInfo.AccountNum + int(factory.AccountNum)
				factoryInfo.AccountNum = factoryInfo.AccountDeployNum + int(factory.AccountDeployNum)
				factoryInfo.AccountNumD1 = factoryInfo.AccountNumD1 + int(factory.AccountNum)
				factoryInfo.AccountDeployNumD1 = factoryInfo.AccountDeployNumD1 + int(factory.AccountDeployNum)
			} else {
				factoryInfo = &ent.FactoryInfo{
					AccountDeployNum:   int(factory.AccountDeployNum),
					AccountNum:         int(factory.AccountNum),
					AccountNumD1:       int(factory.AccountNum),
					AccountDeployNumD1: int(factory.AccountDeployNum),
				}
			}
			factoryInfo.Factory = factory.Factory
			factoryInfo.Network = factory.Network
			factoryInfoMap[factoryAddr] = factoryInfo

		}

		for factory, factoryInfo := range factoryInfoMap {
			if len(factory) == 0 {
				continue
			}
			saveOrUpdateFactory(client, factory, factoryInfo)
		}
	}

}

func TopPaymaster() {
	doTopPaymasterHour()
	paymasterScheduler := chrono.NewDefaultTaskScheduler()

	_, err := paymasterScheduler.ScheduleWithCron(func(ctx context.Context) {
		doTopPaymasterHour()
	}, "0 6 0 * * ?")

	if err == nil {
		log.Print("TopPaymaster has been scheduled")
	}

}

func doTopPaymasterHour() {
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
		startTime := time.Date(now.Year(), now.Month(), now.Day()-500, now.Hour()-23, 0, 0, 0, now.Location())
		endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
		paymasterStatisHours, err := client.PaymasterStatisHour.
			Query().
			Where(
				paymasterstatishour.StatisTimeGTE(startTime),
				paymasterstatishour.StatisTimeLT(endTime),
			).
			All(context.Background())

		if err != nil {
			log.Println(err)
			continue
		}
		if len(paymasterStatisHours) == 0 {
			continue
		}

		paymasterInfoMap := make(map[string]*ent.PaymasterInfo)
		for _, paymasterStatisHour := range paymasterStatisHours {
			paymaster := paymasterStatisHour.Paymaster
			paymasterInfo, paymasterInfoOk := paymasterInfoMap[paymaster]
			if paymasterInfoOk {
				paymasterInfo.UserOpsNum = paymasterInfo.UserOpsNum + paymasterStatisHour.UserOpsNum
				paymasterInfo.GasSponsored = paymasterInfo.GasSponsored.Add(paymasterStatisHour.GasSponsored)
				paymasterInfo.UserOpsNumD1 = paymasterInfo.UserOpsNumD1 + paymasterStatisHour.UserOpsNum
				paymasterInfo.GasSponsoredD1 = paymasterInfo.GasSponsoredD1.Add(paymasterStatisHour.GasSponsored)
			} else {
				paymasterInfo = &ent.PaymasterInfo{
					UserOpsNum:     paymasterStatisHour.UserOpsNum,
					GasSponsored:   paymasterStatisHour.GasSponsored,
					UserOpsNumD1:   paymasterStatisHour.UserOpsNum,
					GasSponsoredD1: paymasterStatisHour.GasSponsored,
				}
			}
			paymasterInfo.Paymaster = paymasterStatisHour.Paymaster
			paymasterInfo.Network = paymasterStatisHour.Network
			paymasterInfo.Reserve = paymasterStatisHour.Reserve
			paymasterInfo.ReserveUsd = paymasterStatisHour.ReserveUsd
			paymasterInfoMap[paymaster] = paymasterInfo

		}

		price := service.GetNativePrice(network)
		for paymaster, paymasterInfo := range paymasterInfoMap {
			if len(paymaster) == 0 {
				continue
			}
			nativeBalance := moralis.GetNativeTokenBalance(paymaster, network)
			paymasterInfo.Reserve = nativeBalance
			paymasterInfo.ReserveUsd = price.Mul(nativeBalance)
			saveOrUpdatePaymaster(client, paymaster, paymasterInfo)
		}
	}

}

func TopBundlers() {
	doTopBundlersHour()
	bundlerScheduler := chrono.NewDefaultTaskScheduler()

	_, err := bundlerScheduler.ScheduleWithCron(func(ctx context.Context) {
		doTopBundlersHour()
	}, "0 5 0 * * ?")

	if err == nil {
		log.Print("TopBundlers has been scheduled")
	}

}

func doTopBundlersHour() {
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
		startTime := time.Date(now.Year(), now.Month(), now.Day()-500, now.Hour()-23, 0, 0, 0, now.Location())
		endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
		bundlerStatisHours, err := client.BundlerStatisHour.
			Query().
			Where(
				bundlerstatishour.StatisTimeGTE(startTime),
				bundlerstatishour.StatisTimeLT(endTime),
			).
			All(context.Background())

		if err != nil {
			log.Println(err)
			continue
		}
		if len(bundlerStatisHours) == 0 {
			continue
		}

		var totalBundleNum = int64(0)
		var bundleNumMap = make(map[string]int64)
		var totalNumMap = make(map[string]int64)
		var feeEarnedMap = make(map[string]decimal.Decimal)
		for _, bundlerStatisHour := range bundlerStatisHours {
			bundler := bundlerStatisHour.Bundler
			feeEarned, feeOk := feeEarnedMap[bundler]
			if !feeOk {
				feeEarned = decimal.Zero
			}
			feeEarnedMap[bundler] = feeEarned.Add(*bundlerStatisHour.FeeEarned)

			totalBundleNum += bundlerStatisHour.BundlesNum
			bundleNum, ok := bundleNumMap[bundlerStatisHour.Bundler]
			if !ok {
				bundleNum = 0
			}
			bundleNumMap[bundlerStatisHour.Bundler] = int64(bundleNum) + bundlerStatisHour.BundlesNum
			totalNum, totalOk := totalNumMap[bundlerStatisHour.Bundler]
			if !totalOk {
				totalNum = 0
			}
			totalNumMap[bundlerStatisHour.Bundler] = int64(totalNum) + bundlerStatisHour.TotalNum
		}
		bundleRateMap, sucRateMap := getRate(totalBundleNum, bundleNumMap, totalNumMap)
		price := service.GetNativePrice(network)
		bundlerInfoMap := make(map[string]*ent.BundlerInfo)
		for _, bundlerStatisHour := range bundlerStatisHours {
			bundler := bundlerStatisHour.Bundler
			bundlerInfo, bundlerInfoOk := bundlerInfoMap[bundler]

			if bundlerInfoOk {
				bundlerInfo.UserOpsNum = bundlerInfo.UserOpsNum + bundlerStatisHour.UserOpsNum
				bundlerInfo.BundlesNum = bundlerInfo.BundlesNum + bundlerStatisHour.BundlesNum
				bundlerInfo.GasCollected = bundlerInfo.GasCollected.Add(bundlerStatisHour.GasCollected)
				bundlerInfo.UserOpsNumD1 = bundlerInfo.UserOpsNumD1 + bundlerStatisHour.UserOpsNum
				bundlerInfo.BundlesNumD1 = bundlerInfo.BundlesNumD1 + bundlerStatisHour.BundlesNum
				bundlerInfo.GasCollectedD1 = bundlerInfo.GasCollectedD1.Add(bundlerStatisHour.GasCollected)
			} else {
				bundlerInfo = &ent.BundlerInfo{
					UserOpsNum:     bundlerStatisHour.UserOpsNum,
					BundlesNum:     bundlerStatisHour.BundlesNum,
					GasCollected:   bundlerStatisHour.GasCollected,
					UserOpsNumD1:   bundlerStatisHour.UserOpsNum,
					BundlesNumD1:   bundlerStatisHour.BundlesNum,
					GasCollectedD1: bundlerStatisHour.GasCollected,
				}
			}

			if bundleRateMap != nil {
				bundlerInfo.BundleRateD1 = bundleRateMap[bundler]
			}

			if sucRateMap != nil {
				bundlerInfo.SuccessRateD1 = sucRateMap[bundler]
			}

			if feeEarnedMap != nil {
				feeEarnUsd := price.Mul(feeEarnedMap[bundler])
				bundlerInfo.FeeEarned = bundlerInfo.FeeEarned.Add(feeEarnedMap[bundler])
				bundlerInfo.FeeEarnedD1 = feeEarnedMap[bundler]
				bundlerInfo.FeeEarnedUsd = bundlerInfo.FeeEarnedUsd.Add(feeEarnUsd)
				bundlerInfo.FeeEarnedUsdD1 = feeEarnUsd
			}

			bundlerInfo.Bundler = bundlerStatisHour.Bundler
			bundlerInfo.Network = bundlerStatisHour.Network
			bundlerInfoMap[bundler] = bundlerInfo

		}
		for bundler, bundlerInfo := range bundlerInfoMap {
			if len(bundler) == 0 {
				continue
			}
			saveOrUpdateBundler(client, bundler, bundlerInfo)
		}
	}

}

func getRate(num int64, bundleNumMap map[string]int64, totalNumMap map[string]int64) (map[string]decimal.Decimal, map[string]decimal.Decimal) {
	if len(bundleNumMap) == 0 || len(totalNumMap) == 0 {
		return nil, nil
	}
	var bundleRateMap = make(map[string]decimal.Decimal)
	var sucRateMap = make(map[string]decimal.Decimal)
	for bundler, singleNum := range bundleNumMap {
		totalNum := totalNumMap[bundler]
		if singleNum == 0 || totalNum == 0 {
			bundleRateMap[bundler] = decimal.Zero
			sucRateMap[bundler] = decimal.Zero
			continue
		}
		sucRate := decimal.NewFromInt(singleNum).DivRound(decimal.NewFromInt(totalNum), 4)
		bundleRate := decimal.NewFromInt(singleNum).DivRound(decimal.NewFromInt(num), 4)
		bundleRateMap[bundler] = bundleRate
		sucRateMap[bundler] = sucRate
	}
	return bundleRateMap, sucRateMap
}

func saveOrUpdateBundler(client *ent.Client, bundler string, info *ent.BundlerInfo) {
	bundlerInfos, err := client.BundlerInfo.
		Query().
		Where(bundlerinfo.BundlerEQ(bundler)).
		All(context.Background())
	if err != nil {
		log.Fatalf("saveOrUpdateBundler err, %s, msg:{%s}\n", bundler, err)
	}
	if len(bundlerInfos) == 0 {
		_, err := client.BundlerInfo.Create().
			SetBundler(info.Bundler).
			SetNetwork(info.Network).
			SetGasCollectedD1(info.GasCollectedD1).
			SetUserOpsNum(info.UserOpsNum).
			SetBundlesNum(info.BundlesNum).
			SetGasCollected(info.GasCollected).
			SetUserOpsNumD1(info.UserOpsNumD1).
			SetBundlesNumD1(info.BundlesNumD1).
			SetFeeEarned(info.FeeEarned).
			SetFeeEarnedUsd(info.FeeEarnedUsd).
			SetFeeEarnedD1(info.FeeEarnedD1).
			SetFeeEarnedUsdD1(info.FeeEarnedUsdD1).
			SetFeeEarnedD7(info.FeeEarnedD7).
			SetFeeEarnedUsdD7(info.FeeEarnedUsdD7).
			SetFeeEarnedD30(info.FeeEarnedD30).
			SetFeeEarnedUsdD30(info.FeeEarnedUsdD30).
			SetSuccessRate(info.SuccessRate).
			SetSuccessRateD1(info.SuccessRateD1).
			SetSuccessRateD7(info.SuccessRateD7).
			SetSuccessRateD30(info.SuccessRateD30).
			SetBundleRate(info.BundleRate).
			SetBundleRateD7(info.BundleRateD7).
			SetBundleRateD1(info.BundleRateD1).
			SetBundleRateD30(info.BundleRateD30).
			Save(context.Background())
		if err != nil {
			log.Printf("Save bundler err, %s\n", err)
		}
	} else {
		oldBundler := bundlerInfos[0]
		err = client.BundlerInfo.UpdateOneID(oldBundler.ID).
			SetBundlesNum(oldBundler.BundlesNum + info.BundlesNum).
			SetUserOpsNum(oldBundler.UserOpsNum + info.UserOpsNum).
			SetGasCollected(oldBundler.GasCollected.Add(info.GasCollected)).
			SetBundlesNumD1(oldBundler.BundlesNumD1 + info.BundlesNumD1).
			SetUserOpsNumD1(oldBundler.UserOpsNumD1 + info.UserOpsNumD1).
			SetGasCollectedD1(oldBundler.GasCollectedD1.Add(info.GasCollectedD1)).Exec(context.Background())
		if err != nil {
			log.Printf("Update bundler err, %s\n", err)
		}
	}
}

func saveOrUpdatePaymaster(client *ent.Client, paymaster string, info *ent.PaymasterInfo) {
	paymasterInfos, err := client.PaymasterInfo.
		Query().
		Where(paymasterinfo.PaymasterEQ(paymaster)).
		All(context.Background())
	if err != nil {
		log.Fatalf("saveOrUpdatePaymaster err, %s, msg:{%s}\n", paymaster, err)
	}
	if paymasterInfos == nil || len(paymasterInfos) == 0 {

		_, err := client.PaymasterInfo.Create().
			SetPaymaster(info.Paymaster).
			SetNetwork(info.Network).
			SetUserOpsNum(info.UserOpsNum).
			SetGasSponsored(info.GasSponsored).
			SetUserOpsNumD1(info.UserOpsNumD1).
			SetGasSponsoredD1(info.GasSponsoredD1).
			SetGasSponsoredUsd(info.GasSponsoredUsdD1).
			SetGasSponsoredUsdD1(info.GasSponsoredUsdD1).
			SetReserve(info.Reserve).
			SetReserveUsd(info.ReserveUsd).
			Save(context.Background())
		if err != nil {
			log.Printf("Save paymaster err, %s\n", err)
		}
	} else {
		oldPaymaster := paymasterInfos[0]
		err = client.PaymasterInfo.UpdateOneID(oldPaymaster.ID).
			SetUserOpsNum(oldPaymaster.UserOpsNum + info.UserOpsNum).
			SetGasSponsored(oldPaymaster.GasSponsored.Add(info.GasSponsored)).
			SetUserOpsNumD1(oldPaymaster.UserOpsNumD1 + info.UserOpsNumD1).
			SetGasSponsoredD1(oldPaymaster.GasSponsoredD1.Add(info.GasSponsoredD1)).Exec(context.Background())
		if err != nil {
			log.Printf("Update paymaster err, %s\n", err)
		}
	}
}

func saveOrUpdateFactory(client *ent.Client, factory string, info *ent.FactoryInfo) {
	factoryInfos, err := client.FactoryInfo.
		Query().
		Where(factoryinfo.FactoryEQ(factory)).
		All(context.Background())
	if err != nil {
		log.Fatalf("saveOrUpdateFactory err, %s, msg:{%s}\n", factory, err)
	}
	if len(factoryInfos) == 0 {

		_, err := client.FactoryInfo.Create().
			SetFactory(info.Factory).
			SetNetwork(info.Network).
			SetAccountNum(info.AccountNum).
			SetAccountDeployNum(info.AccountDeployNum).
			SetAccountNumD1(info.AccountNumD1).
			SetAccountDeployNumD1(info.AccountDeployNumD1).
			Save(context.Background())
		if err != nil {
			log.Printf("Save factory err, %s\n", err)
		}
	} else {
		oldFactory := factoryInfos[0]
		err = client.FactoryInfo.UpdateOneID(oldFactory.ID).
			SetAccountDeployNum(oldFactory.AccountDeployNum + info.AccountDeployNum).
			SetAccountNum(oldFactory.AccountNum + info.AccountNum).
			SetAccountDeployNumD1(oldFactory.AccountDeployNumD1 + info.AccountDeployNumD1).
			SetAccountNumD1(oldFactory.AccountNumD1 + info.AccountNumD1).Exec(context.Background())
		if err != nil {
			log.Printf("Update factory err, %s\n", err)
		}
	}
}
