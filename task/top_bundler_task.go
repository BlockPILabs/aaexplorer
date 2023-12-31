package task

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/bundlerinfo"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/bundlerstatisday"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/bundlerstatishour"
	"github.com/BlockPILabs/aaexplorer/service"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

func InitTask() {

	AccountTask()
	InitDayStatis()
	InitHourStatis()
	TopBundlers()
	TopPaymaster()
	TopFactories()
	AAContractInteractTask()
	UserOpTypeTask()
	AssetTask()
	go AATransactionFix()
	//temp
	//DataFixedTask()

}

func addOpsInfo(key string, opsInfo *ent.AAUserOpsInfo, bundlerMap map[string]map[string][]*ent.AAUserOpsInfo, startType string) {
	bundlerOps, bundlerOk := bundlerMap[key]
	if !bundlerOk {
		bundlerOps = make(map[string][]*ent.AAUserOpsInfo)
	}
	var startOf = ""
	if startType == "day" {
		startOf = getDayStart(opsInfo.Time)
	} else {
		startOf = getHourStart(opsInfo.Time)
	}
	timeOps, timeOpOk := bundlerOps[startOf]
	if !timeOpOk {
		timeOps = []*ent.AAUserOpsInfo{}
	}
	timeOps = append(timeOps, opsInfo)
	bundlerOps[startOf] = timeOps
	bundlerMap[key] = bundlerOps
}

func getDayStart(t time.Time) string {
	startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Format(TimeLayout)
	return startOfDay
}

func getHourStart(t time.Time) string {
	startOfHour := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.Local).Format(TimeLayout)
	return startOfHour
}

func TopBundlers() {
	bundlerScheduler := chrono.NewDefaultTaskScheduler()
	_, err := bundlerScheduler.ScheduleWithCron(func(ctx context.Context) {
		doTopBundlersHour(1)
		//doTopBundlersHour(7)
		//doTopBundlersHour(30)
	}, "0 7 * * * *")

	bundlerSchedulerDay := chrono.NewDefaultTaskScheduler()

	_, err = bundlerSchedulerDay.ScheduleWithCron(func(ctx context.Context) {
		doTopBundlersDay()
	}, "0 10 5 * * *")
	if err == nil {
		log.Print("TopBundlers has been scheduled")
	}

}

func doTopBundlersDay() {
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
		log.Printf("top bundler day statistic start, network:%s", network)
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
		endTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		bundlerStatisDays, err := client.BundlerStatisDay.
			Query().
			Where(
				bundlerstatisday.StatisTimeGTE(startTime),
				bundlerstatisday.StatisTimeLT(endTime),
			).
			All(context.Background())

		if err != nil {
			log.Println(err)
			continue
		}
		if len(bundlerStatisDays) == 0 {
			continue
		}

		var totalBundleNum = int64(0)
		var bundleNumMap = make(map[string]int64)
		var totalNumMap = make(map[string]int64)
		var feeEarnedMap = make(map[string]decimal.Decimal)
		var repeatMap = make(map[string]bool)
		var allBundlerNum = int64(0)
		for _, bundlerStatisDay := range bundlerStatisDays {
			bundler := bundlerStatisDay.Bundler
			timeStr := bundlerStatisDay.StatisTime.String()
			_, exist := repeatMap[bundler+timeStr]
			if exist {
				continue
			}
			repeatMap[bundler+timeStr] = true
			feeEarned, feeOk := feeEarnedMap[bundler]
			if !feeOk {
				feeEarned = decimal.Zero
			}
			feeEarnedMap[bundler] = feeEarned.Add(bundlerStatisDay.FeeEarned)

			totalBundleNum += bundlerStatisDay.BundlesNum
			bundleNum, ok := bundleNumMap[bundlerStatisDay.Bundler]
			if !ok {
				bundleNum = 0
			}
			bundleNumMap[bundlerStatisDay.Bundler] = bundleNum + bundlerStatisDay.TotalNum
			totalNum, totalOk := totalNumMap[bundlerStatisDay.Bundler]
			if !totalOk {
				totalNum = 0
			}
			totalNumMap[bundlerStatisDay.Bundler] = totalNum + bundlerStatisDay.TotalNum
			allBundlerNum += bundlerStatisDay.TotalNum
		}
		price := service.GetNativePrice(network)
		bundlerInfoMap := make(map[string]*ent.BundlerInfo)
		for _, bundlerStatisDay := range bundlerStatisDays {
			bundler := bundlerStatisDay.Bundler
			bundlerInfo, bundlerInfoOk := bundlerInfoMap[bundler]

			if bundlerInfoOk {
				bundlerInfo.UserOpsNum = bundlerInfo.UserOpsNum + bundlerStatisDay.UserOpsNum
				bundlerInfo.BundlesNum = bundlerInfo.BundlesNum + bundlerStatisDay.TotalNum
				bundlerInfo.GasCollected = bundlerInfo.GasCollected.Add(bundlerStatisDay.GasCollected)
				bundlerInfo.SuccessBundlesNum = bundlerInfo.SuccessBundlesNum + bundlerStatisDay.SuccessBundlesNum
				bundlerInfo.FailedBundlesNum = bundlerInfo.FailedBundlesNum + bundlerStatisDay.FailedBundlesNum

			} else {
				bundlerInfo = &ent.BundlerInfo{
					UserOpsNum:        bundlerStatisDay.UserOpsNum,
					BundlesNum:        bundlerStatisDay.TotalNum,
					GasCollected:      bundlerStatisDay.GasCollected,
					SuccessBundlesNum: bundlerStatisDay.SuccessBundlesNum,
					FailedBundlesNum:  bundlerStatisDay.FailedBundlesNum,
				}
				bundlerInfoMap[bundler] = bundlerInfo
			}

			if feeEarnedMap != nil {
				feeEarnUsd := price.Mul(feeEarnedMap[bundler])
				bundlerInfo.FeeEarned = bundlerInfo.FeeEarned.Add(feeEarnedMap[bundler])
				bundlerInfo.FeeEarnedUsd = bundlerInfo.FeeEarnedUsd.Add(feeEarnUsd)
			}

			bundlerInfo.ID = bundlerStatisDay.Bundler
			bundlerInfo.Network = bundlerStatisDay.Network
			bundlerInfo.SuccessRate = getSingleRate(bundlerInfo.SuccessBundlesNum, bundlerInfo.SuccessBundlesNum+bundlerInfo.FailedBundlesNum)
			bundlerInfoMap[bundler] = bundlerInfo

		}
		for bundler, bundlerInfo := range bundlerInfoMap {
			if len(bundler) == 0 {
				continue
			}
			bundleRate := getSingleRate(bundlerInfo.BundlesNum, allBundlerNum)
			bundlerInfo.BundleRate = bundleRate
			saveOrUpdateBundlerDay(client, bundler, bundlerInfo)
		}
		now1 := time.Now()
		log.Printf("top bundler day statistic success, network:%s, spent:%d", network, now1.Second()-now.Second())

	}
}

func getSingleRate(num int64, num2 int64) decimal.Decimal {
	if num == 0 || num2 == 0 {
		return decimal.Zero
	}
	return decimal.NewFromInt(num).DivRound(decimal.NewFromInt(num2), 4)
}

func saveOrUpdateBundlerDay(client *ent.Client, bundler string, info *ent.BundlerInfo) {
	bundlerInfos, err := client.BundlerInfo.
		Query().
		Where(bundlerinfo.IDEQ(bundler)).
		All(context.Background())
	if err != nil {
		log.Printf("saveOrUpdateBundler day err, %s, msg:{%s}\n", bundler, err)
	}
	if len(bundlerInfos) == 0 {
		_, err := client.BundlerInfo.Create().
			SetID(info.ID).
			SetNetwork(info.Network).
			SetUserOpsNum(info.UserOpsNum).
			SetBundlesNum(info.BundlesNum).
			SetGasCollected(info.GasCollected).
			SetFeeEarned(info.FeeEarned).
			SetFeeEarnedUsd(info.FeeEarnedUsd).
			SetSuccessRate(info.SuccessRate).
			SetBundleRate(info.BundleRate).
			SetSuccessBundlesNum(info.SuccessBundlesNum).
			SetFailedBundlesNum(info.FailedBundlesNum).
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
			SetFeeEarned(oldBundler.FeeEarned.Add(info.FeeEarned)).
			SetFeeEarnedUsd(oldBundler.FeeEarnedUsd.Add(info.FeeEarnedUsd)).
			SetSuccessRate(info.SuccessRate).
			SetBundleRate(info.BundleRate).
			SetSuccessBundlesNum(info.SuccessBundlesNum + info.SuccessBundlesNum).
			SetFailedBundlesNum(info.FailedBundlesNum + info.FailedBundlesNum).
			Exec(context.Background())
		if err != nil {
			log.Printf("Update bundler day err, %s\n", err)
		}
	}

	log.Printf("top bundler day, single statistic sync success, bundler:%s", info.ID)

}

func doTopBundlersHour(timeRange int) {
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
		log.Printf("top bundler hour statistic start timeRange:%d, network:%s", timeRange, network)
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-24*timeRange, 0, 0, 0, now.Location())
		endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
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
		var successBundleNumMap = make(map[string]int64)
		var totalNumMap = make(map[string]int64)
		var feeEarnedMap = make(map[string]decimal.Decimal)
		var repeatMap = make(map[string]bool)
		for _, bundlerStatisHour := range bundlerStatisHours {
			bundler := bundlerStatisHour.Bundler
			timeStr := bundlerStatisHour.StatisTime.String()
			_, exist := repeatMap[bundler+timeStr]
			if exist {
				continue
			}
			repeatMap[bundler+timeStr] = true
			feeEarned, feeOk := feeEarnedMap[bundler]
			if !feeOk {
				feeEarned = decimal.Zero
			}
			feeEarnedMap[bundler] = feeEarned.Add(bundlerStatisHour.FeeEarned)

			bundleNum, ok := successBundleNumMap[bundlerStatisHour.Bundler]
			if !ok {
				bundleNum = 0
			}
			successBundleNumMap[bundler] = bundleNum + bundlerStatisHour.SuccessBundlesNum

			totalNum, totalOk := totalNumMap[bundlerStatisHour.Bundler]
			if !totalOk {
				totalNum = 0
			}
			totalNumMap[bundler] = totalNum + bundlerStatisHour.TotalNum
			totalBundleNum += bundlerStatisHour.TotalNum
		}
		bundleRateMap, sucRateMap := getRate(totalBundleNum, successBundleNumMap, totalNumMap)
		price := service.GetNativePrice(network)
		bundlerInfoMap := make(map[string]*ent.BundlerInfo)
		for _, bundlerStatisHour := range bundlerStatisHours {
			bundler := bundlerStatisHour.Bundler
			bundlerInfo, bundlerInfoOk := bundlerInfoMap[bundler]

			if bundlerInfoOk {
				if timeRange == 1 {
					bundlerInfo.UserOpsNumD1 = bundlerInfo.UserOpsNumD1 + bundlerStatisHour.UserOpsNum
					bundlerInfo.BundlesNumD1 = bundlerInfo.BundlesNumD1 + bundlerStatisHour.TotalNum
					bundlerInfo.GasCollectedD1 = bundlerInfo.GasCollectedD1.Add(bundlerStatisHour.GasCollected)
				} else if timeRange == 7 {
					bundlerInfo.UserOpsNumD7 = bundlerInfo.UserOpsNumD7 + bundlerStatisHour.UserOpsNum
					bundlerInfo.BundlesNumD7 = bundlerInfo.BundlesNumD7 + bundlerStatisHour.TotalNum
					bundlerInfo.GasCollectedD7 = bundlerInfo.GasCollectedD7.Add(bundlerStatisHour.GasCollected)
				} else if timeRange == 30 {
					bundlerInfo.UserOpsNumD30 = bundlerInfo.UserOpsNumD30 + bundlerStatisHour.UserOpsNum
					bundlerInfo.BundlesNumD30 = bundlerInfo.BundlesNumD30 + bundlerStatisHour.TotalNum
					bundlerInfo.GasCollectedD30 = bundlerInfo.GasCollectedD30.Add(bundlerStatisHour.GasCollected)
				}

			} else {
				bundlerInfo = &ent.BundlerInfo{}
				if timeRange == 1 {
					bundlerInfo.UserOpsNumD1 = bundlerStatisHour.UserOpsNum
					bundlerInfo.BundlesNumD1 = bundlerStatisHour.TotalNum
					bundlerInfo.GasCollectedD1 = bundlerStatisHour.GasCollected
				} else if timeRange == 7 {
					bundlerInfo.UserOpsNumD7 = bundlerStatisHour.UserOpsNum
					bundlerInfo.BundlesNumD7 = bundlerStatisHour.TotalNum
					bundlerInfo.GasCollectedD7 = bundlerStatisHour.GasCollected
				} else if timeRange == 30 {
					bundlerInfo.UserOpsNumD30 = bundlerStatisHour.UserOpsNum
					bundlerInfo.BundlesNumD30 = bundlerStatisHour.TotalNum
					bundlerInfo.GasCollectedD30 = bundlerStatisHour.GasCollected
				}
				bundlerInfoMap[bundler] = bundlerInfo
			}

			if bundleRateMap != nil {
				if timeRange == 1 {
					bundlerInfo.BundleRateD1 = bundleRateMap[bundler]
				} else if timeRange == 7 {
					bundlerInfo.BundleRateD7 = bundleRateMap[bundler]
				} else if timeRange == 30 {
					bundlerInfo.BundleRateD30 = bundleRateMap[bundler]
				}

			}

			if sucRateMap != nil {
				if timeRange == 1 {
					bundlerInfo.SuccessRateD1 = sucRateMap[bundler]
				} else if timeRange == 7 {
					bundlerInfo.SuccessRateD7 = sucRateMap[bundler]
				} else if timeRange == 30 {
					bundlerInfo.SuccessRateD30 = sucRateMap[bundler]
				}
			}

			if feeEarnedMap != nil {
				feeEarnUsd := price.Mul(feeEarnedMap[bundler])
				if timeRange == 1 {
					bundlerInfo.FeeEarnedD1 = feeEarnedMap[bundler]
					bundlerInfo.FeeEarnedUsdD1 = feeEarnUsd
				} else if timeRange == 7 {
					bundlerInfo.FeeEarnedD7 = feeEarnedMap[bundler]
					bundlerInfo.FeeEarnedUsdD7 = feeEarnUsd
				} else if timeRange == 30 {
					bundlerInfo.FeeEarnedD30 = feeEarnedMap[bundler]
					bundlerInfo.FeeEarnedUsdD30 = feeEarnUsd
				}

			}

			bundlerInfo.ID = bundlerStatisHour.Bundler
			bundlerInfo.Network = bundlerStatisHour.Network
			bundlerInfoMap[bundler] = bundlerInfo

		}
		allBundlers, err := client.BundlerInfo.Query().All(context.Background())

		for bundler, bundlerInfo := range bundlerInfoMap {
			if len(bundler) == 0 {
				continue
			}
			saveOrUpdateBundler(client, bundler, bundlerInfo, timeRange)
		}
		if len(allBundlers) > 0 {
			for _, bundler := range allBundlers {
				_, ok := bundlerInfoMap[bundler.ID]
				if !ok {
					if timeRange == 1 {
						err = client.BundlerInfo.UpdateOneID(bundler.ID).
							SetSuccessRateD1(decimal.Zero).
							SetBundleRateD1(decimal.Zero).
							SetFeeEarnedD1(decimal.Zero).
							SetFeeEarnedUsdD1(decimal.Zero).
							SetBundlesNumD1(int64(0)).
							SetUserOpsNumD1(int64(0)).
							SetGasCollectedD1(decimal.Zero).Exec(context.Background())
					} else if timeRange == 7 {
						err = client.BundlerInfo.UpdateOneID(bundler.ID).
							SetSuccessRateD7(decimal.Zero).
							SetBundleRateD7(decimal.Zero).
							SetFeeEarnedD7(decimal.Zero).
							SetFeeEarnedUsdD7(decimal.Zero).
							SetBundlesNumD7(int64(0)).
							SetUserOpsNumD7(int64(0)).
							SetGasCollectedD7(decimal.Zero).Exec(context.Background())
					} else if timeRange == 30 {
						err = client.BundlerInfo.UpdateOneID(bundler.ID).
							SetSuccessRateD30(decimal.Zero).
							SetBundleRateD30(decimal.Zero).
							SetFeeEarnedD30(decimal.Zero).
							SetFeeEarnedUsdD30(decimal.Zero).
							SetBundlesNumD30(int64(0)).
							SetUserOpsNumD30(int64(0)).
							SetGasCollectedD30(decimal.Zero).Exec(context.Background())
					}
				}
			}
		}
		log.Printf("top bundler hour statistic success timeRange:%s, network:%s", string(timeRange), network)
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

func saveOrUpdateBundler(client *ent.Client, bundler string, info *ent.BundlerInfo, timeRange int) {
	bundlerInfos, err := client.BundlerInfo.
		Query().
		Where(bundlerinfo.IDEQ(bundler)).
		All(context.Background())
	if err != nil {
		log.Printf("saveOrUpdateBundler err, %s, msg:{%s}\n", bundler, err)
	}
	if len(bundlerInfos) == 0 {
		newBundler := client.BundlerInfo.Create().
			SetID(info.ID).
			SetNetwork(info.Network)
		if timeRange == 1 {
			newBundler.SetSuccessRateD1(info.SuccessRateD1).
				SetBundleRateD1(info.BundleRateD1).
				SetFeeEarnedD1(info.FeeEarnedD1).
				SetFeeEarnedUsdD1(info.FeeEarnedUsdD1).
				SetBundlesNumD1(info.BundlesNumD1).
				SetUserOpsNumD1(info.UserOpsNumD1).
				SetGasCollectedD1(info.GasCollectedD1)
		} else if timeRange == 7 {
			newBundler.SetSuccessRateD7(info.SuccessRateD7).
				SetBundleRateD7(info.BundleRateD7).
				SetFeeEarnedD7(info.FeeEarnedD7).
				SetFeeEarnedUsdD7(info.FeeEarnedUsdD7).
				SetBundlesNumD7(info.BundlesNumD7).
				SetUserOpsNumD7(info.UserOpsNumD7).
				SetGasCollectedD7(info.GasCollectedD7)
		} else if timeRange == 30 {
			newBundler.SetSuccessRateD30(info.SuccessRateD30).
				SetBundleRateD30(info.BundleRateD30).
				SetFeeEarnedD30(info.FeeEarnedD30).
				SetFeeEarnedUsdD30(info.FeeEarnedUsdD30).
				SetBundlesNumD30(info.BundlesNumD30).
				SetUserOpsNumD30(info.UserOpsNumD30).
				SetGasCollectedD30(info.GasCollectedD30)
		}
		newBundler.Save(context.Background())

		if err != nil {
			log.Printf("Save bundler err, %s\n", err)
		}
	} else {
		oldBundler := bundlerInfos[0]
		if timeRange == 1 {
			err = client.BundlerInfo.UpdateOneID(oldBundler.ID).
				SetSuccessRateD1(info.SuccessRateD1).
				SetBundleRateD1(info.BundleRateD1).
				SetFeeEarnedD1(info.FeeEarnedD1).
				SetFeeEarnedUsdD1(info.FeeEarnedUsdD1).
				SetBundlesNumD1(info.BundlesNumD1).
				SetUserOpsNumD1(info.UserOpsNumD1).
				SetGasCollectedD1(info.GasCollectedD1).Exec(context.Background())
		} else if timeRange == 7 {
			err = client.BundlerInfo.UpdateOneID(oldBundler.ID).
				SetSuccessRateD7(info.SuccessRateD7).
				SetBundleRateD7(info.BundleRateD7).
				SetFeeEarnedD7(info.FeeEarnedD7).
				SetFeeEarnedUsdD7(info.FeeEarnedUsdD7).
				SetBundlesNumD7(info.BundlesNumD7).
				SetUserOpsNumD7(info.UserOpsNumD7).
				SetGasCollectedD7(info.GasCollectedD7).Exec(context.Background())
		} else if timeRange == 30 {
			err = client.BundlerInfo.UpdateOneID(oldBundler.ID).
				SetSuccessRateD30(info.SuccessRateD30).
				SetBundleRateD30(info.BundleRateD30).
				SetFeeEarnedD30(info.FeeEarnedD30).
				SetFeeEarnedUsdD30(info.FeeEarnedUsdD30).
				SetBundlesNumD30(info.BundlesNumD30).
				SetUserOpsNumD30(info.UserOpsNumD30).
				SetGasCollectedD30(info.GasCollectedD30).Exec(context.Background())
		}

		if err != nil {
			log.Printf("Update bundler err, %s\n", err)
		}
	}
	log.Printf("top bundler hour, single statistic sync success, bundler:%s", info.ID)

}
