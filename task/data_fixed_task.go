package task

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/bundlerinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/bundlerstatisday"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/bundlerstatishour"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/factoryinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/factorystatisday"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/factorystatishour"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/paymasterinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/paymasterstatisday"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/paymasterstatishour"
	"github.com/BlockPILabs/aa-scan/service"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

//manu trigger

func DataFixedTask() {
	go FixBundlerDayAll()
	go FixPaymasterDayAll()
	go FixFactoryDayAll()
	go FixBundlerDay1()
	go FixPaymasterDay1()
	go FixFactoryDay1()
}

func FixBundlerDayAll() {
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
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		bundlerInfos, err := client.BundlerInfo.Query().All(context.Background())
		if len(bundlerInfos) == 0 {
			continue
		}
		var totalBundles = int64(0)
		now := time.Now()
		dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		for _, bundler := range bundlerInfos {
			statisticDays, _ := client.BundlerStatisDay.Query().Where(bundlerstatisday.BundlerEQ(bundler.ID)).All(context.Background())
			if len(statisticDays) == 0 {
				continue
			}
			var bundlesNum = int64(0)
			var successNum = int64(0)
			var failedNum = int64(0)
			var userOpsNum = int64(0)
			var gasCollected = decimal.Zero
			var feeEarned = decimal.Zero
			for _, day := range statisticDays {
				if day.StatisTime.Compare(dayStart) >= 0 {
					continue
				}
				bundlesNum += day.BundlesNum
				successNum += day.SuccessBundlesNum
				failedNum += day.FailedBundlesNum
				userOpsNum += day.UserOpsNum
				gasCollected = gasCollected.Add(day.GasCollected)
				feeEarned = feeEarned.Add(day.FeeEarned)
			}
			price := service.GetNativePrice(network)
			bundler.FeeEarnedUsd = bundler.FeeEarned.Mul(price)
			bundler.BundlesNum = successNum + failedNum
			bundler.SuccessBundlesNum = successNum
			bundler.FailedBundlesNum = failedNum
			bundler.UserOpsNum = userOpsNum
			bundler.GasCollected = gasCollected
			bundler.FeeEarned = feeEarned
			totalBundles += successNum + failedNum

			bundler.SuccessRate = getSingleRate(bundler.SuccessBundlesNum, bundler.SuccessBundlesNum+bundler.FailedBundlesNum)
			client.BundlerInfo.Update().
				SetSuccessRate(bundler.SuccessRate).
				SetBundlesNum(bundler.BundlesNum).
				SetSuccessBundlesNum(bundler.SuccessBundlesNum).
				SetFailedBundlesNum(bundler.FailedBundlesNum).
				SetUserOpsNum(bundler.UserOpsNum).
				SetGasCollected(bundler.GasCollected).
				SetFeeEarned(bundler.FeeEarned).
				SetFeeEarnedUsd(bundler.FeeEarnedUsd).
				Where(bundlerinfo.IDEQ(bundler.ID)).Exec(context.Background())
			log.Printf("bundler all fixed success, %s", bundler.ID)
		}

		for _, bundler := range bundlerInfos {
			bundleRate := getSingleRate(bundler.BundlesNum, totalBundles)
			client.BundlerInfo.Update().
				SetBundleRate(bundleRate).
				Where(bundlerinfo.IDEQ(bundler.ID)).Exec(context.Background())

			log.Printf("bundler rate fixed success, %s", bundler.ID)
		}
	}
}

func FixBundlerDay1() {
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
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		bundlerInfos, err := client.BundlerInfo.Query().All(context.Background())
		if len(bundlerInfos) == 0 {
			continue
		}
		var totalBundles = int64(0)
		now := time.Now()
		hourStart := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
		startTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-24, 0, 0, 0, now.Location())
		for _, bundler := range bundlerInfos {
			statisticDays, _ := client.BundlerStatisHour.Query().Where(bundlerstatishour.BundlerEQ(bundler.ID), bundlerstatishour.StatisTimeGTE(startTime)).All(context.Background())
			if len(statisticDays) == 0 {
				continue
			}
			var bundlesNumD1 = int64(0)
			var successNumD1 = int64(0)
			var failedNumD1 = int64(0)
			var userOpsNumD1 = int64(0)
			var gasCollectedD1 = decimal.Zero
			var feeEarnedD1 = decimal.Zero
			for _, day := range statisticDays {
				if day.StatisTime.Compare(hourStart) >= 0 {
					continue
				}
				bundlesNumD1 += day.BundlesNum
				successNumD1 += day.SuccessBundlesNum
				failedNumD1 += day.FailedBundlesNum
				userOpsNumD1 += day.UserOpsNum
				gasCollectedD1 = gasCollectedD1.Add(day.GasCollected)
				feeEarnedD1 = feeEarnedD1.Add(day.FeeEarned)
			}
			price := service.GetNativePrice(network)
			bundler.FeeEarnedUsd = bundler.FeeEarned.Mul(price)
			bundler.BundlesNumD1 = successNumD1 + failedNumD1
			bundler.UserOpsNumD1 = userOpsNumD1
			bundler.GasCollectedD1 = gasCollectedD1
			bundler.FeeEarnedD1 = feeEarnedD1
			totalBundles += successNumD1 + failedNumD1

			bundler.SuccessRateD1 = getSingleRate(successNumD1, successNumD1+failedNumD1)
			client.BundlerInfo.Update().
				SetSuccessRateD1(bundler.SuccessRateD1).
				SetBundlesNumD1(bundler.BundlesNumD1).
				SetUserOpsNumD1(bundler.UserOpsNumD1).
				SetGasCollectedD1(bundler.GasCollectedD1).
				SetFeeEarnedD1(bundler.FeeEarnedD1).
				SetFeeEarnedUsdD1(bundler.FeeEarnedUsdD1).
				Where(bundlerinfo.IDEQ(bundler.ID)).Exec(context.Background())
			log.Printf("bundler d1 all fixed success, %s", bundler.ID)
		}

		for _, bundler := range bundlerInfos {
			bundleRate := getSingleRate(bundler.BundlesNumD1, totalBundles)
			client.BundlerInfo.Update().
				SetBundleRateD1(bundleRate).
				Where(bundlerinfo.IDEQ(bundler.ID)).Exec(context.Background())

			log.Printf("bundler d1 rate fixed success, %s", bundler.ID)
		}
	}
}

func FixPaymasterDayAll() {
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
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		paymasterInfos, err := client.PaymasterInfo.Query().All(context.Background())
		if len(paymasterInfos) == 0 {
			continue
		}
		now := time.Now()
		dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		price := service.GetNativePrice(network)
		for _, paymaster := range paymasterInfos {
			statisticDays, _ := client.PaymasterStatisDay.Query().Where(paymasterstatisday.PaymasterEQ(paymaster.ID)).All(context.Background())
			if len(statisticDays) == 0 {
				continue
			}
			var gasSponsored = decimal.Zero
			var gasSponsoredUsd = decimal.Zero
			var userOpsNum = int64(0)
			for _, day := range statisticDays {
				if day.StatisTime.Compare(dayStart) >= 0 {
					continue
				}
				userOpsNum += day.UserOpsNum
				gasSponsored = gasSponsored.Add(day.GasSponsored)
				gasSponsoredUsd = gasSponsoredUsd.Add(day.GasSponsoredUsd)
			}
			paymaster.UserOpsNum = userOpsNum
			paymaster.GasSponsored = gasSponsored
			paymaster.GasSponsoredUsd = gasSponsored.Mul(price)
			paymaster.ReserveUsd = statisticDays[0].ReserveUsd

			client.PaymasterInfo.Update().
				SetUserOpsNum(paymaster.UserOpsNum).
				SetReserveUsd(paymaster.ReserveUsd).
				SetGasSponsoredUsd(paymaster.GasSponsoredUsd).
				SetGasSponsored(paymaster.GasSponsored).
				Where(paymasterinfo.IDEQ(paymaster.ID)).Exec(context.Background())
			log.Printf("paymaster all fixed success, %s", paymaster.ID)
		}

	}
}

func FixPaymasterDay1() {
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
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		paymasterInfos, err := client.PaymasterInfo.Query().All(context.Background())
		if len(paymasterInfos) == 0 {
			continue
		}
		now := time.Now()
		hourStart := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
		startTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-24, 0, 0, 0, now.Location())

		price := service.GetNativePrice(network)
		for _, paymaster := range paymasterInfos {
			statisticHours, _ := client.PaymasterStatisHour.Query().Where(paymasterstatishour.PaymasterEQ(paymaster.ID), paymasterstatishour.StatisTimeGTE(startTime)).All(context.Background())
			if len(statisticHours) == 0 {
				continue
			}
			var gasSponsoredD1 = decimal.Zero
			var gasSponsoredUsdD1 = decimal.Zero
			var userOpsNumD1 = int64(0)
			for _, hour := range statisticHours {
				if hour.StatisTime.Compare(hourStart) >= 0 {
					continue
				}
				userOpsNumD1 += hour.UserOpsNum
				gasSponsoredD1 = gasSponsoredD1.Add(hour.GasSponsored)
				gasSponsoredUsdD1 = gasSponsoredUsdD1.Add(hour.GasSponsoredUsd)
			}
			paymaster.UserOpsNumD1 = userOpsNumD1
			paymaster.GasSponsoredD1 = gasSponsoredD1
			paymaster.GasSponsoredUsdD1 = gasSponsoredD1.Mul(price)
			paymaster.ReserveUsd = statisticHours[0].ReserveUsd

			client.PaymasterInfo.Update().
				SetUserOpsNumD1(paymaster.UserOpsNumD1).
				//SetReserveUsd(paymaster.ReserveUsd).
				SetGasSponsoredUsdD1(paymaster.GasSponsoredUsdD1).
				SetGasSponsoredD1(paymaster.GasSponsoredD1).
				Where(paymasterinfo.IDEQ(paymaster.ID)).Exec(context.Background())
			log.Printf("paymaster d1 all fixed success, %s", paymaster.ID)
		}

	}
}

func FixFactoryDayAll() {
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
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		factoryInfos, err := client.FactoryInfo.Query().All(context.Background())
		if len(factoryInfos) == 0 {
			continue
		}
		now := time.Now()
		dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		var totalNum = 0
		for _, factory := range factoryInfos {
			statisticDays, _ := client.FactoryStatisDay.Query().Where(factorystatisday.FactoryEQ(factory.ID)).All(context.Background())
			if len(statisticDays) == 0 {
				continue
			}
			var deployedAccountNum = int64(0)
			var accountNum = int64(0)
			for _, day := range statisticDays {
				if day.StatisTime.Compare(dayStart) >= 0 {
					continue
				}
				accountNum += day.AccountNum
				deployedAccountNum += day.AccountDeployNum
			}
			factory.AccountNum = int(accountNum)
			factory.AccountDeployNum = int(deployedAccountNum)
			totalNum += factory.AccountNum

			client.FactoryInfo.Update().
				SetAccountNum(factory.AccountNum).
				SetAccountDeployNum(factory.AccountDeployNum).
				Where(factoryinfo.IDEQ(factory.ID)).Exec(context.Background())
			log.Printf("factory all fixed success, %s", factory.ID)
		}

		for _, factory := range factoryInfos {
			dominance := getSingleRate(int64(factory.AccountNum), int64(totalNum))
			factory.Dominance = dominance
			client.FactoryInfo.Update().SetDominance(dominance).Where(factoryinfo.IDEQ(factory.ID)).Exec(context.Background())
			log.Printf("factory dominance fixed success, %s", factory.ID)
		}
	}
}

func FixFactoryDay1() {
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
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		factoryInfos, err := client.FactoryInfo.Query().All(context.Background())
		if len(factoryInfos) == 0 {
			continue
		}
		now := time.Now()
		hourStart := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
		startTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-24, 0, 0, 0, now.Location())
		var totalNum = 0
		for _, factory := range factoryInfos {
			statisticHours, _ := client.FactoryStatisHour.Query().Where(factorystatishour.FactoryEQ(factory.ID), factorystatishour.StatisTimeGTE(startTime)).All(context.Background())
			if len(statisticHours) == 0 {
				continue
			}
			var deployedAccountNum = int64(0)
			var accountNum = int64(0)
			for _, hour := range statisticHours {
				if hour.StatisTime.Compare(hourStart) >= 0 {
					continue
				}
				accountNum += hour.AccountNum
				deployedAccountNum += hour.AccountDeployNum
			}
			factory.AccountNum = int(accountNum)
			factory.AccountDeployNum = int(deployedAccountNum)
			totalNum += factory.AccountNum

			client.FactoryInfo.Update().
				SetAccountNumD1(factory.AccountNumD1).
				SetAccountDeployNumD1(factory.AccountDeployNumD1).
				Where(factoryinfo.IDEQ(factory.ID)).Exec(context.Background())
			log.Printf("factory d1 all fixed success, %s", factory.ID)
		}

		for _, factory := range factoryInfos {
			dominanceD1 := getSingleRate(int64(factory.AccountNumD1), int64(totalNum))
			factory.Dominance = dominanceD1
			client.FactoryInfo.Update().SetDominanceD1(dominanceD1).Where(factoryinfo.IDEQ(factory.ID)).Exec(context.Background())
			log.Printf("factory d1 dominance fixed success, %s", factory.ID)
		}
	}
}
