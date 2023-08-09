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
	"github.com/procyon-projects/chrono"
	"log"
	"time"
)

func InitTask() {

	//hour statistics
	InitHourStatis()

	//day statistics
	InitDayStatis()

	topBundlers()

	topPaymaster()

	topFactories()

}

func addOpsInfo(key string, opsInfo *ent.UserOpsInfo, bundlerMap map[string][]*ent.UserOpsInfo) {
	bundlerOps, bundlerOk := bundlerMap[key]
	if !bundlerOk {
		bundlerOps = []*ent.UserOpsInfo{}
	}

	bundlerOps = append(bundlerOps, opsInfo)
	bundlerMap[key] = bundlerOps
}

func topFactories() {
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
	client, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day()-7, now.Hour()-23, 0, 0, 0, now.Location())
	endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
	factoryStatisHours, err := client.FactoryStatisHour.
		Query().
		Where(
			factorystatishour.StatisTimeGTE(startTime),
			factorystatishour.StatisTimeLT(endTime),
		).
		All(context.Background())

	if err != nil {
		return
	}
	if len(factoryStatisHours) == 0 {
		return
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
		if factory == "" {
			continue
		}
		saveOrUpdateFactory(client, factory, factoryInfo)
	}
}

func topPaymaster() {
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
	client, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day()-7, now.Hour()-23, 0, 0, 0, now.Location())
	endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
	paymasterStatisHours, err := client.PaymasterStatisHour.
		Query().
		Where(
			paymasterstatishour.StatisTimeGTE(startTime),
			paymasterstatishour.StatisTimeLT(endTime),
		).
		All(context.Background())

	if err != nil {
		return
	}
	if len(paymasterStatisHours) == 0 {
		return
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
		paymasterInfoMap[paymaster] = paymasterInfo

	}

	for paymaster, paymasterInfo := range paymasterInfoMap {
		if paymaster == "" {
			continue
		}
		saveOrUpdatePaymaster(client, paymaster, paymasterInfo)
	}

}

func topBundlers() {
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
	client, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day()-7, now.Hour()-23, 0, 0, 0, now.Location())
	endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
	bundlerStatisHours, err := client.BundlerStatisHour.
		Query().
		Where(
			bundlerstatishour.StatisTimeGTE(startTime),
			bundlerstatishour.StatisTimeLT(endTime),
		).
		All(context.Background())

	if err != nil {
		return
	}
	if len(bundlerStatisHours) == 0 {
		return
	}

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
		bundlerInfo.Bundler = bundlerStatisHour.Bundler
		bundlerInfo.Network = bundlerStatisHour.Network
		bundlerInfoMap[bundler] = bundlerInfo

	}

	for bundler, bundlerInfo := range bundlerInfoMap {
		if bundler == "" {
			continue
		}
		saveOrUpdateBundler(client, bundler, bundlerInfo)
	}
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
