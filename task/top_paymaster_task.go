package task

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/paymasterinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/paymasterstatishour"
	"github.com/BlockPILabs/aa-scan/service"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

func TopPaymaster() {
	doTopPaymasterHour(1)
	doTopPaymasterHour(7)
	doTopPaymasterHour(30)
	doTopPaymasterDay()
	paymasterScheduler := chrono.NewDefaultTaskScheduler()

	_, err := paymasterScheduler.ScheduleWithCron(func(ctx context.Context) {
		doTopPaymasterHour(1)
		doTopPaymasterHour(7)
		doTopPaymasterHour(30)
	}, "0 5 * * * ?")

	paymasterSchedulerDay := chrono.NewDefaultTaskScheduler()

	_, err = paymasterSchedulerDay.ScheduleWithCron(func(ctx context.Context) {
		doTopPaymasterDay()
	}, "0 15 0 * * ?")

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
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day()-50, 0, 0, 0, 0, now.Location())
		endTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
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
		if price == nil {
			price = &decimal.Zero
		}
		paymasterInfoMap := make(map[string]*ent.PaymasterInfo)
		for _, paymasterStatisHour := range paymasterStatisHours {
			paymaster := paymasterStatisHour.Paymaster
			paymasterInfo, paymasterInfoOk := paymasterInfoMap[paymaster]
			if paymasterInfoOk {
				paymasterInfo.UserOpsNum = paymasterInfo.UserOpsNum + paymasterStatisHour.UserOpsNum
				paymasterInfo.GasSponsored = paymasterInfo.GasSponsored.Add(paymasterStatisHour.GasSponsored)
				paymasterInfo.GasSponsoredUsd = paymasterInfo.GasSponsoredUsd.Add(paymasterStatisHour.GasSponsored.Mul(*price))

			} else {
				paymasterInfo = &ent.PaymasterInfo{
					UserOpsNum:      paymasterStatisHour.UserOpsNum,
					GasSponsored:    paymasterStatisHour.GasSponsored,
					GasSponsoredUsd: paymasterStatisHour.GasSponsored.Mul(*price),
				}
			}
			paymasterInfo.Paymaster = paymasterStatisHour.Paymaster
			paymasterInfo.Network = paymasterStatisHour.Network
			paymasterInfo.Reserve = paymasterStatisHour.Reserve
			//paymasterInfo.ReserveUsd = paymasterStatisHour.ReserveUsd
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
	}
}

func saveOrUpdatePaymasterDay(client *ent.Client, paymaster string, info *ent.PaymasterInfo) {
	paymasterInfos, err := client.PaymasterInfo.
		Query().
		Where(paymasterinfo.PaymasterEQ(paymaster)).
		All(context.Background())
	if err != nil {
		log.Fatalf("saveOrUpdatePaymaster day err, %s, msg:{%s}\n", paymaster, err)
	}
	if paymasterInfos == nil || len(paymasterInfos) == 0 {

		_, err := client.PaymasterInfo.Create().
			SetPaymaster(info.Paymaster).
			SetNetwork(info.Network).
			SetUserOpsNum(info.UserOpsNum).
			SetGasSponsored(info.GasSponsored).
			SetGasSponsoredUsd(info.GasSponsoredUsdD1).
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
			Exec(context.Background())
		if err != nil {
			log.Printf("Update paymaster day err, %s\n", err)
		}
	}
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
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day()-50, now.Hour()-720, 0, 0, 0, now.Location())
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
		price := service.GetNativePrice(network)
		if price == nil {
			price = &decimal.Zero
		}
		paymasterInfoMap := make(map[string]*ent.PaymasterInfo)
		for _, paymasterStatisHour := range paymasterStatisHours {
			paymaster := paymasterStatisHour.Paymaster
			paymasterInfo, paymasterInfoOk := paymasterInfoMap[paymaster]
			if paymasterInfoOk {
				if timeRange == 1 {
					paymasterInfo.UserOpsNumD1 = paymasterInfo.UserOpsNumD1 + paymasterStatisHour.UserOpsNum
					paymasterInfo.GasSponsoredD1 = paymasterInfo.GasSponsoredD1.Add(paymasterStatisHour.GasSponsored)
					paymasterInfo.GasSponsoredUsdD1 = paymasterInfo.GasSponsoredUsdD1.Add(paymasterStatisHour.GasSponsored.Mul(*price))
				} else if timeRange == 7 {
					paymasterInfo.UserOpsNumD7 = paymasterInfo.UserOpsNumD7 + paymasterStatisHour.UserOpsNum
					paymasterInfo.GasSponsoredD7 = paymasterInfo.GasSponsoredD7.Add(paymasterStatisHour.GasSponsored)
					paymasterInfo.GasSponsoredUsdD7 = paymasterInfo.GasSponsoredUsdD7.Add(paymasterStatisHour.GasSponsored.Mul(*price))
				} else if timeRange == 30 {
					paymasterInfo.UserOpsNumD30 = paymasterInfo.UserOpsNumD30 + paymasterStatisHour.UserOpsNum
					paymasterInfo.GasSponsoredD30 = paymasterInfo.GasSponsoredD30.Add(paymasterStatisHour.GasSponsored)
					paymasterInfo.GasSponsoredUsdD30 = paymasterInfo.GasSponsoredUsdD30.Add(paymasterStatisHour.GasSponsored.Mul(*price))
				}
			} else {
				paymasterInfo = &ent.PaymasterInfo{
					UserOpsNumD1:       paymasterStatisHour.UserOpsNum,
					GasSponsoredD1:     paymasterStatisHour.GasSponsored,
					GasSponsoredUsdD1:  paymasterStatisHour.GasSponsored.Mul(*price),
					UserOpsNumD7:       paymasterStatisHour.UserOpsNum,
					GasSponsoredD7:     paymasterStatisHour.GasSponsored,
					GasSponsoredUsdD7:  paymasterStatisHour.GasSponsored.Mul(*price),
					UserOpsNumD30:      paymasterStatisHour.UserOpsNum,
					GasSponsoredD30:    paymasterStatisHour.GasSponsored,
					GasSponsoredUsdD30: paymasterStatisHour.GasSponsored.Mul(*price),
				}
			}
			paymasterInfo.Paymaster = paymasterStatisHour.Paymaster
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
	}

}

func saveOrUpdatePaymaster(client *ent.Client, paymaster string, info *ent.PaymasterInfo, timeRange int) {
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
			SetUserOpsNumD1(info.UserOpsNumD1).
			SetGasSponsoredD1(info.GasSponsoredD1).
			SetGasSponsoredUsdD1(info.GasSponsoredUsdD1).
			SetUserOpsNumD7(info.UserOpsNumD7).
			SetGasSponsoredD7(info.GasSponsoredD7).
			SetGasSponsoredUsdD7(info.GasSponsoredUsdD7).
			SetUserOpsNumD30(info.UserOpsNumD30).
			SetGasSponsoredD30(info.GasSponsoredD30).
			SetGasSponsoredUsdD30(info.GasSponsoredUsdD30).
			SetReserve(info.Reserve).
			SetReserveUsd(info.ReserveUsd).
			Save(context.Background())
		if err != nil {
			log.Printf("Save paymaster err, %s\n", err)
		}
	} else {
		oldPaymaster := paymasterInfos[0]
		if timeRange == 1 {
			err = client.PaymasterInfo.UpdateOneID(oldPaymaster.ID).
				SetUserOpsNumD1(info.UserOpsNumD1).
				SetGasSponsoredD1(info.GasSponsoredD1).
				SetGasSponsoredD1(info.GasSponsoredUsdD1).
				Exec(context.Background())
		} else if timeRange == 7 {
			err = client.PaymasterInfo.UpdateOneID(oldPaymaster.ID).
				SetUserOpsNumD7(info.UserOpsNumD7).
				SetGasSponsoredD7(info.GasSponsoredD7).
				SetGasSponsoredD7(info.GasSponsoredUsdD7).
				Exec(context.Background())
		} else if timeRange == 30 {
			err = client.PaymasterInfo.UpdateOneID(oldPaymaster.ID).
				SetUserOpsNumD30(info.UserOpsNumD30).
				SetGasSponsoredD30(info.GasSponsoredD30).
				SetGasSponsoredD30(info.GasSponsoredUsdD30).
				Exec(context.Background())
		}

		if err != nil {
			log.Printf("Update paymaster err, %s\n", err)
		}
	}
}
