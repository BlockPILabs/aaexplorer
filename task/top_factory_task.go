package task

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/factoryinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/factorystatisday"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/factorystatishour"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

func TopFactories() {
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
		startTime := time.Date(now.Year(), now.Month(), now.Day()-70, 0, 0, 0, 0, now.Location())
		endTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		factoryStatisDays, err := client.FactoryStatisDay.
			Query().
			Where(
				factorystatisday.StatisTimeGTE(startTime),
				factorystatisday.StatisTimeLT(endTime),
			).
			All(context.Background())

		if err != nil {
			log.Println(err)
			continue
		}
		if len(factoryStatisDays) == 0 {
			continue
		}

		factoryInfoMap := make(map[string]*ent.FactoryInfo)
		var totalNum = 0
		var repeatMap = make(map[string]bool)
		for _, factory := range factoryStatisDays {
			factoryAddr := factory.Factory
			timeStr := string(factory.StatisTime.UnixMilli())
			_, exist := repeatMap[factoryAddr+timeStr]
			if exist {
				continue
			}
			repeatMap[factoryAddr+timeStr] = true
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
			factoryInfo.ID = factory.Factory
			factoryInfo.Network = factory.Network
			factoryInfoMap[factoryAddr] = factoryInfo
			totalNum += factoryInfo.AccountNum
		}

		for factory, factoryInfo := range factoryInfoMap {
			if len(factory) == 0 {
				continue
			}
			factoryInfo.Dominance = decimal.NewFromInt(int64(factoryInfo.AccountDeployNum)).DivRound(decimal.NewFromInt(int64(totalNum)), 4)
			saveOrUpdateFactoryDay(client, factory, factoryInfo)
		}
	}

}

func saveOrUpdateFactoryDay(client *ent.Client, factory string, info *ent.FactoryInfo) {
	factoryInfos, err := client.FactoryInfo.
		Query().
		Where(factoryinfo.IDEQ(factory)).
		All(context.Background())
	if err != nil {
		log.Printf("saveOrUpdateFactory day err, %s, msg:{%s}\n", factory, err)
	}
	if len(factoryInfos) == 0 {

		_, err := client.FactoryInfo.Create().
			SetID(info.ID).
			SetNetwork(info.Network).
			SetAccountNum(info.AccountNum).
			SetAccountDeployNum(info.AccountDeployNum).
			SetDominance(info.Dominance).
			Save(context.Background())
		if err != nil {
			log.Printf("Save factory day err, %s\n", err)
		}
	} else {
		oldFactory := factoryInfos[0]
		err = client.FactoryInfo.UpdateOneID(oldFactory.ID).
			SetAccountDeployNum(oldFactory.AccountDeployNum + info.AccountDeployNum).
			SetAccountNum(oldFactory.AccountNum + info.AccountNum).
			SetDominance(info.Dominance).
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
		startTime := time.Date(now.Year(), now.Month(), now.Day()-70, now.Hour()-720, 0, 0, 0, now.Location())
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
		var totalNum = 0
		var repeatMap = make(map[string]bool)
		for _, factory := range factoryStatisHours {
			factoryAddr := factory.Factory
			timeStr := string(factory.StatisTime.UnixMilli())
			_, exist := repeatMap[factoryAddr+timeStr]
			if exist {
				continue
			}
			repeatMap[factoryAddr+timeStr] = true
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
				factoryInfo = &ent.FactoryInfo{}
				if timeRange == 1 {
					factoryInfo.AccountNumD1 = int(factory.AccountNum)
					factoryInfo.AccountDeployNumD1 = int(factory.AccountDeployNum)
				} else if timeRange == 7 {
					factoryInfo.AccountNumD7 = int(factory.AccountNum)
					factoryInfo.AccountDeployNumD7 = int(factory.AccountDeployNum)
				} else if timeRange == 30 {
					factoryInfo.AccountNumD30 = int(factory.AccountNum)
					factoryInfo.AccountDeployNumD30 = int(factory.AccountDeployNum)
				}
			}
			factoryInfo.ID = factory.Factory
			factoryInfo.Network = factory.Network
			factoryInfoMap[factoryAddr] = factoryInfo
			if timeRange == 1 {
				totalNum += factoryInfo.AccountNumD1
			} else if timeRange == 7 {
				totalNum += factoryInfo.AccountNumD7
			} else if timeRange == 30 {
				totalNum += factoryInfo.AccountNumD30
			}

		}

		for factory, factoryInfo := range factoryInfoMap {
			if len(factory) == 0 {
				continue
			}
			if timeRange == 1 {
				factoryInfo.DominanceD1 = getSingleRate(int64(factoryInfo.AccountDeployNumD1), int64(totalNum))
			} else if timeRange == 7 {
				factoryInfo.DominanceD7 = getSingleRate(int64(factoryInfo.AccountDeployNumD7), int64(totalNum))
			} else if timeRange == 30 {
				factoryInfo.DominanceD30 = getSingleRate(int64(factoryInfo.AccountDeployNumD30), int64(totalNum))
			}
			saveOrUpdateFactory(client, factory, factoryInfo, timeRange)
		}
	}

}

func saveOrUpdateFactory(client *ent.Client, factory string, info *ent.FactoryInfo, timeRange int) {
	factoryInfos, err := client.FactoryInfo.
		Query().
		Where(factoryinfo.IDEQ(factory)).
		All(context.Background())
	if err != nil {
		log.Printf("saveOrUpdateFactory err, %s, msg:{%s}\n", factory, err)
	}
	if len(factoryInfos) == 0 {

		newFactory := client.FactoryInfo.Create().
			SetID(info.ID).
			SetNetwork(info.Network)

		if timeRange == 1 {
			newFactory.SetAccountNumD1(info.AccountNumD1).
				SetAccountDeployNumD1(info.AccountDeployNumD1).
				SetDominanceD1(info.DominanceD1)
		} else if timeRange == 7 {
			newFactory.SetAccountNumD7(info.AccountNumD7).
				SetAccountDeployNumD7(info.AccountDeployNumD7).
				SetDominanceD7(info.DominanceD7)
		} else if timeRange == 30 {
			newFactory.SetAccountNumD30(info.AccountNumD30).
				SetAccountDeployNumD30(info.AccountDeployNumD30).
				SetDominanceD30(info.DominanceD30)
		}
		_, err := newFactory.Save(context.Background())
		if err != nil {
			log.Printf("Save factory err, %s\n", err)
		}
	} else {
		oldFactory := factoryInfos[0]
		if timeRange == 1 {
			err = client.FactoryInfo.UpdateOneID(oldFactory.ID).
				SetAccountDeployNumD1(info.AccountDeployNumD1).
				SetAccountNumD1(info.AccountNumD1).
				SetDominanceD1(info.DominanceD1).
				Exec(context.Background())
		} else if timeRange == 7 {
			err = client.FactoryInfo.UpdateOneID(oldFactory.ID).
				SetAccountDeployNumD7(info.AccountDeployNumD7).
				SetAccountNumD7(info.AccountNumD7).
				SetDominanceD7(info.DominanceD7).
				Exec(context.Background())
		} else if timeRange == 30 {
			err = client.FactoryInfo.UpdateOneID(oldFactory.ID).
				SetAccountDeployNumD30(info.AccountDeployNumD30).
				SetAccountNumD30(info.AccountNumD30).
				SetDominanceD30(info.DominanceD30).
				Exec(context.Background())
		}

		if err != nil {
			log.Printf("Update factory err, %s\n", err)
		}
	}
}
