package task

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/paymasterinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/paymasterstatisday"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/paymasterstatishour"
	"github.com/BlockPILabs/aa-scan/service"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

func TopPaymaster() {
	paymasterScheduler := chrono.NewDefaultTaskScheduler()

	_, err := paymasterScheduler.ScheduleWithCron(func(ctx context.Context) {
		doTopPaymasterHour(1)
		//doTopPaymasterHour(7)
		//doTopPaymasterHour(30)
	}, "0 7 * * * *")

	paymasterSchedulerDay := chrono.NewDefaultTaskScheduler()

	_, err = paymasterSchedulerDay.ScheduleWithCron(func(ctx context.Context) {
		doTopPaymasterDay()
	}, "0 15 5 * * *")

	if err == nil {
		log.Print("TopPaymaster has been scheduled")
	}

}

func doTopPaymasterDay() {
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
		log.Printf("top paymaster day statistic start, network:%s", network)
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
		endTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		paymasterStatisDays, err := client.PaymasterStatisDay.
			Query().
			Where(
				paymasterstatisday.StatisTimeGTE(startTime),
				paymasterstatisday.StatisTimeLT(endTime),
			).
			All(context.Background())

		if err != nil {
			log.Println(err)
			continue
		}
		if len(paymasterStatisDays) == 0 {
			continue
		}
		price := service.GetNativePrice(network)
		paymasterInfoMap := make(map[string]*ent.PaymasterInfo)
		var repeatMap = make(map[string]bool)
		for _, paymasterStatisDay := range paymasterStatisDays {
			paymaster := paymasterStatisDay.Paymaster
			timeStr := paymasterStatisDay.StatisTime.String()
			_, exist := repeatMap[paymaster+timeStr]
			if exist {
				continue
			}
			repeatMap[paymaster+timeStr] = true
			paymasterInfo, paymasterInfoOk := paymasterInfoMap[paymaster]
			if paymasterInfoOk {
				paymasterInfo.UserOpsNum = paymasterInfo.UserOpsNum + paymasterStatisDay.UserOpsNum
				paymasterInfo.GasSponsored = paymasterInfo.GasSponsored.Add(paymasterStatisDay.GasSponsored)
				paymasterInfo.GasSponsoredUsd = paymasterInfo.GasSponsoredUsd.Add(paymasterStatisDay.GasSponsored.Mul(price))

			} else {
				paymasterInfo = &ent.PaymasterInfo{
					UserOpsNum:      paymasterStatisDay.UserOpsNum,
					GasSponsored:    paymasterStatisDay.GasSponsored,
					GasSponsoredUsd: paymasterStatisDay.GasSponsored.Mul(price),
				}
				paymasterInfoMap[paymaster] = paymasterInfo
			}
			paymasterInfo.ID = paymasterStatisDay.Paymaster
			paymasterInfo.Network = paymasterStatisDay.Network
			paymasterInfo.Reserve = paymasterStatisDay.Reserve
			paymasterInfo.ReserveUsd = paymasterStatisDay.ReserveUsd
			paymasterInfoMap[paymaster] = paymasterInfo

		}

		//price := service.GetNativePrice(network)
		for paymaster, paymasterInfo := range paymasterInfoMap {
			if len(paymaster) == 0 {
				continue
			}
			//nativeBalance := moralis.GetNativeTokenBalance(paymaster, network)
			//paymasterInfo.Reserve = nativeBalance
			//paymasterInfo.ReserveUsd = price.Mul(nativeBalance)
			saveOrUpdatePaymasterDay(client, paymaster, paymasterInfo)
		}
		log.Printf("top paymaster day statistic success, network:%s", network)
	}
}

func saveOrUpdatePaymasterDay(client *ent.Client, paymaster string, info *ent.PaymasterInfo) {
	paymasterInfos, err := client.PaymasterInfo.
		Query().
		Where(paymasterinfo.IDEQ(paymaster)).
		All(context.Background())
	if err != nil {
		log.Printf("saveOrUpdatePaymaster day err, %s, msg:{%s}\n", paymaster, err)
	}
	if paymasterInfos == nil || len(paymasterInfos) == 0 {

		_, err := client.PaymasterInfo.Create().
			SetID(info.ID).
			SetNetwork(info.Network).
			SetUserOpsNum(info.UserOpsNum).
			SetGasSponsored(info.GasSponsored).
			SetGasSponsoredUsd(info.GasSponsoredUsd).
			SetReserve(info.Reserve).
			SetReserveUsd(info.ReserveUsd).
			Save(context.Background())
		if err != nil {
			log.Printf("Save paymaster day err, %s\n", err)
		}
	} else {
		oldPaymaster := paymasterInfos[0]
		err = client.PaymasterInfo.UpdateOneID(oldPaymaster.ID).
			SetUserOpsNum(oldPaymaster.UserOpsNum + info.UserOpsNum).
			SetGasSponsored(oldPaymaster.GasSponsored.Add(info.GasSponsored)).
			SetGasSponsoredUsd(oldPaymaster.GasSponsoredUsd.Add(info.GasSponsoredUsd)).
			Exec(context.Background())
		if err != nil {
			log.Printf("Update paymaster day err, %s\n", err)
		}
	}
	log.Printf("top paymaster day, single statistic sync success, bundler:%s", info.ID)
}

func doTopPaymasterHour(timeRange int) {
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
		log.Printf("top paymaster hour statistic start, timeRange:%d network:%s", timeRange, network)
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-24*timeRange, 0, 0, 0, now.Location())
		endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
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
		price := service.GetNativePrice(network)
		paymasterInfoMap := make(map[string]*ent.PaymasterInfo)
		var repeatMap = make(map[string]bool)
		for _, paymasterStatisHour := range paymasterStatisHours {
			paymaster := paymasterStatisHour.Paymaster
			timeStr := paymasterStatisHour.StatisTime.String()
			_, exist := repeatMap[paymaster+timeStr]
			if exist {
				continue
			}
			repeatMap[paymaster+timeStr] = true
			paymasterInfo, paymasterInfoOk := paymasterInfoMap[paymaster]
			if paymasterInfoOk {
				if timeRange == 1 {
					paymasterInfo.UserOpsNumD1 = paymasterInfo.UserOpsNumD1 + paymasterStatisHour.UserOpsNum
					paymasterInfo.GasSponsoredD1 = paymasterInfo.GasSponsoredD1.Add(paymasterStatisHour.GasSponsored)
					paymasterInfo.GasSponsoredUsdD1 = paymasterInfo.GasSponsoredUsdD1.Add(paymasterStatisHour.GasSponsored.Mul(price))
				} else if timeRange == 7 {
					paymasterInfo.UserOpsNumD7 = paymasterInfo.UserOpsNumD7 + paymasterStatisHour.UserOpsNum
					paymasterInfo.GasSponsoredD7 = paymasterInfo.GasSponsoredD7.Add(paymasterStatisHour.GasSponsored)
					paymasterInfo.GasSponsoredUsdD7 = paymasterInfo.GasSponsoredUsdD7.Add(paymasterStatisHour.GasSponsored.Mul(price))
				} else if timeRange == 30 {
					paymasterInfo.UserOpsNumD30 = paymasterInfo.UserOpsNumD30 + paymasterStatisHour.UserOpsNum
					paymasterInfo.GasSponsoredD30 = paymasterInfo.GasSponsoredD30.Add(paymasterStatisHour.GasSponsored)
					paymasterInfo.GasSponsoredUsdD30 = paymasterInfo.GasSponsoredUsdD30.Add(paymasterStatisHour.GasSponsored.Mul(price))
				}
			} else {
				paymasterInfo = &ent.PaymasterInfo{}
				if timeRange == 1 {
					paymasterInfo.UserOpsNumD1 = paymasterStatisHour.UserOpsNum
					paymasterInfo.GasSponsoredD1 = paymasterStatisHour.GasSponsored
					paymasterInfo.GasSponsoredUsdD1 = paymasterStatisHour.GasSponsored.Mul(price)
				} else if timeRange == 7 {
					paymasterInfo.UserOpsNumD7 = paymasterStatisHour.UserOpsNum
					paymasterInfo.GasSponsoredD7 = paymasterStatisHour.GasSponsored
					paymasterInfo.GasSponsoredUsdD7 = paymasterStatisHour.GasSponsored.Mul(price)
				} else if timeRange == 30 {
					paymasterInfo.UserOpsNumD30 = paymasterStatisHour.UserOpsNum
					paymasterInfo.GasSponsoredD30 = paymasterStatisHour.GasSponsored
					paymasterInfo.GasSponsoredUsdD30 = paymasterStatisHour.GasSponsored.Mul(price)
				}
				paymasterInfoMap[paymaster] = paymasterInfo
			}
			paymasterInfo.ID = paymasterStatisHour.Paymaster
			paymasterInfo.Network = paymasterStatisHour.Network
			paymasterInfo.Reserve = paymasterStatisHour.Reserve
			paymasterInfo.ReserveUsd = paymasterStatisHour.ReserveUsd
			paymasterInfoMap[paymaster] = paymasterInfo

		}

		for paymaster, paymasterInfo := range paymasterInfoMap {
			if len(paymaster) == 0 {
				continue
			}
			//nativeBalance := moralis.GetNativeTokenBalance(paymaster, network)
			//paymasterInfo.Reserve = nativeBalance
			//paymasterInfo.ReserveUsd = price.Mul(nativeBalance)
			saveOrUpdatePaymaster(client, paymaster, paymasterInfo, timeRange)
		}

		allPaymaster, err := client.PaymasterInfo.Query().All(context.Background())
		if len(allPaymaster) > 0 {
			for _, paymaster := range allPaymaster {
				_, ok := paymasterInfoMap[paymaster.ID]
				if !ok {
					if timeRange == 1 {
						err = client.PaymasterInfo.UpdateOneID(paymaster.ID).
							SetUserOpsNumD1(int64(0)).
							SetGasSponsoredD1(decimal.Zero).
							SetGasSponsoredUsdD1(decimal.Zero).
							Exec(context.Background())
					} else if timeRange == 7 {
						err = client.PaymasterInfo.UpdateOneID(paymaster.ID).
							SetUserOpsNumD7(int64(0)).
							SetGasSponsoredD7(decimal.Zero).
							SetGasSponsoredUsdD30(decimal.Zero).
							Exec(context.Background())
					} else if timeRange == 30 {
						err = client.PaymasterInfo.UpdateOneID(paymaster.ID).
							SetUserOpsNumD30(int64(0)).
							SetGasSponsoredD30(decimal.Zero).
							SetGasSponsoredUsdD30(decimal.Zero).
							Exec(context.Background())
					}
				}
			}
		}
		log.Printf("top paymaster hour statistic success timeRange:%s, network:%s", string(timeRange), network)
	}

}

func saveOrUpdatePaymaster(client *ent.Client, paymaster string, info *ent.PaymasterInfo, timeRange int) {
	paymasterInfos, err := client.PaymasterInfo.
		Query().
		Where(paymasterinfo.IDEQ(paymaster)).
		All(context.Background())
	if err != nil {
		log.Printf("saveOrUpdatePaymaster err, %s, msg:{%s}\n", paymaster, err)
	}
	if paymasterInfos == nil || len(paymasterInfos) == 0 {

		newPaymaster := client.PaymasterInfo.Create().
			SetID(info.ID).
			SetNetwork(info.Network).
			SetReserve(info.Reserve).
			SetReserveUsd(info.ReserveUsd)

		if timeRange == 1 {
			newPaymaster.SetUserOpsNumD1(info.UserOpsNumD1).
				SetGasSponsoredD1(info.GasSponsoredD1).
				SetGasSponsoredUsdD1(info.GasSponsoredUsdD1)
		} else if timeRange == 7 {
			newPaymaster.SetUserOpsNumD7(info.UserOpsNumD7).
				SetGasSponsoredD7(info.GasSponsoredD7).
				SetGasSponsoredUsdD7(info.GasSponsoredUsdD7)
		} else if timeRange == 30 {
			newPaymaster.SetUserOpsNumD30(info.UserOpsNumD30).
				SetGasSponsoredD30(info.GasSponsoredD30).
				SetGasSponsoredUsdD30(info.GasSponsoredUsdD30)
		}
		_, err := newPaymaster.Save(context.Background())
		if err != nil {
			log.Printf("Save paymaster err, %s\n", err)
		}
	} else {
		oldPaymaster := paymasterInfos[0]
		if timeRange == 1 {
			err = client.PaymasterInfo.UpdateOneID(oldPaymaster.ID).
				SetUserOpsNumD1(info.UserOpsNumD1).
				SetGasSponsoredD1(info.GasSponsoredD1).
				SetGasSponsoredUsdD1(info.GasSponsoredUsdD1).
				Exec(context.Background())
		} else if timeRange == 7 {
			err = client.PaymasterInfo.UpdateOneID(oldPaymaster.ID).
				SetUserOpsNumD7(info.UserOpsNumD7).
				SetGasSponsoredD7(info.GasSponsoredD7).
				SetGasSponsoredUsdD30(info.GasSponsoredUsdD7).
				Exec(context.Background())
		} else if timeRange == 30 {
			err = client.PaymasterInfo.UpdateOneID(oldPaymaster.ID).
				SetUserOpsNumD30(info.UserOpsNumD30).
				SetGasSponsoredD30(info.GasSponsoredD30).
				SetGasSponsoredUsdD30(info.GasSponsoredUsdD30).
				Exec(context.Background())
		}

		if err != nil {
			log.Printf("Update paymaster err, %s\n", err)
		}
	}
	log.Printf("top paymaster hour, single statistic sync success, bundler:%s", info.ID)
}
