package task

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/factoryinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/factorystatishour"
	"github.com/procyon-projects/chrono"
	"log"
	"time"
)

func TopFactories() {
	doTopFactoryHour(1)
	doTopFactoryHour(7)
	doTopFactoryHour(30)
	doTopFactoryDay()
	factoryScheduler := chrono.NewDefaultTaskScheduler()

	_, err := factoryScheduler.ScheduleWithCron(func(ctx context.Context) {
		doTopFactoryHour(1)
		doTopFactoryHour(7)
		doTopFactoryHour(30)
	}, "0 5 * * * ?")

	factorySchedulerDay := chrono.NewDefaultTaskScheduler()

	_, err = factorySchedulerDay.ScheduleWithCron(func(ctx context.Context) {
		doTopFactoryDay()
	}, "0 20 0 * * ?")

	if err == nil {
		log.Print("TopFactory has been scheduled")
	}

}

func doTopFactoryDay() {
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
			} else {
				factoryInfo = &ent.FactoryInfo{
					AccountDeployNum: int(factory.AccountDeployNum),
					AccountNum:       int(factory.AccountNum),
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
			saveOrUpdateFactoryDay(client, factory, factoryInfo)
		}
	}

}

func saveOrUpdateFactoryDay(client *ent.Client, factory string, info *ent.FactoryInfo) {
	factoryInfos, err := client.FactoryInfo.
		Query().
		Where(factoryinfo.FactoryEQ(factory)).
		All(context.Background())
	if err != nil {
		log.Fatalf("saveOrUpdateFactory day err, %s, msg:{%s}\n", factory, err)
	}
	if len(factoryInfos) == 0 {

		_, err := client.FactoryInfo.Create().
			SetFactory(info.Factory).
			SetNetwork(info.Network).
			SetAccountNum(info.AccountNum).
			SetAccountDeployNum(info.AccountDeployNum).
			Save(context.Background())
		if err != nil {
			log.Printf("Save factory day err, %s\n", err)
		}
	} else {
		oldFactory := factoryInfos[0]
		err = client.FactoryInfo.UpdateOneID(oldFactory.ID).
			SetAccountDeployNum(oldFactory.AccountDeployNum + info.AccountDeployNum).
			SetAccountNum(oldFactory.AccountNum + info.AccountNum).
			Exec(context.Background())
		if err != nil {
			log.Printf("Update factory day err, %s\n", err)
		}
	}
}

func doTopFactoryHour(timeRange int) {
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
				if timeRange == 1 {
					factoryInfo.AccountNumD1 = factoryInfo.AccountNumD1 + int(factory.AccountNum)
					factoryInfo.AccountDeployNumD1 = factoryInfo.AccountDeployNumD1 + int(factory.AccountDeployNum)
				} else if timeRange == 7 {
					factoryInfo.AccountNumD7 = factoryInfo.AccountNumD7 + int(factory.AccountNum)
					factoryInfo.AccountDeployNumD7 = factoryInfo.AccountDeployNumD7 + int(factory.AccountDeployNum)
				} else if timeRange == 30 {
					factoryInfo.AccountNumD30 = factoryInfo.AccountNumD30 + int(factory.AccountNum)
					factoryInfo.AccountDeployNumD30 = factoryInfo.AccountDeployNumD30 + int(factory.AccountDeployNum)
				}

			} else {
				factoryInfo = &ent.FactoryInfo{
					AccountNumD1:        int(factory.AccountNum),
					AccountDeployNumD1:  int(factory.AccountDeployNum),
					AccountNumD7:        int(factory.AccountNum),
					AccountDeployNumD7:  int(factory.AccountDeployNum),
					AccountNumD30:       int(factory.AccountNum),
					AccountDeployNumD30: int(factory.AccountDeployNum),
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
			saveOrUpdateFactory(client, factory, factoryInfo, timeRange)
		}
	}

}

func saveOrUpdateFactory(client *ent.Client, factory string, info *ent.FactoryInfo, timeRange int) {
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
			SetAccountNumD1(info.AccountNumD1).
			SetAccountDeployNumD1(info.AccountDeployNumD1).
			SetAccountNumD7(info.AccountNumD7).
			SetAccountDeployNumD7(info.AccountDeployNumD7).
			SetAccountNumD30(info.AccountNumD30).
			SetAccountDeployNumD30(info.AccountDeployNumD30).
			Save(context.Background())
		if err != nil {
			log.Printf("Save factory err, %s\n", err)
		}
	} else {
		oldFactory := factoryInfos[0]
		if timeRange == 1 {
			err = client.FactoryInfo.UpdateOneID(oldFactory.ID).
				SetAccountDeployNumD1(info.AccountDeployNumD1).
				SetAccountNumD1(info.AccountNumD1).
				Exec(context.Background())
		} else if timeRange == 7 {
			err = client.FactoryInfo.UpdateOneID(oldFactory.ID).
				SetAccountDeployNumD7(info.AccountDeployNumD7).
				SetAccountNumD7(info.AccountNumD7).
				Exec(context.Background())
		} else if timeRange == 30 {
			err = client.FactoryInfo.UpdateOneID(oldFactory.ID).
				SetAccountDeployNumD30(info.AccountDeployNumD30).
				SetAccountNumD30(info.AccountNumD30).
				Exec(context.Background())
		}

		if err != nil {
			log.Printf("Update factory err, %s\n", err)
		}
	}
}
