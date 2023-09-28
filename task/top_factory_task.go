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
	doTopFactoryHour(1)
	factoryScheduler := chrono.NewDefaultTaskScheduler()
	_, err := factoryScheduler.ScheduleWithCron(func(ctx context.Context) {
		doTopFactoryHour(1)
		//doTopFactoryHour(7)
		//doTopFactoryHour(30)
	}, "0 7 * * * *")

	factorySchedulerDay := chrono.NewDefaultTaskScheduler()

	_, err = factorySchedulerDay.ScheduleWithCron(func(ctx context.Context) {
		doTopFactoryDay()
	}, "0 20 0 * * *")

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
		log.Printf("top factory day statistic start, network:%s", network)
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
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
		var totalNum = int64(0)
		var repeatMap = make(map[string]bool)
		for _, factory := range factoryStatisDays {
			factoryAddr := factory.Factory
			timeStr := factory.StatisTime.String()
			_, exist := repeatMap[factoryAddr+timeStr]
			if exist {
				continue
			}
			repeatMap[factoryAddr+timeStr] = true
			factoryInfo, bundlerInfoOk := factoryInfoMap[factoryAddr]
			if bundlerInfoOk {

				factoryInfo.AccountDeployNum = factoryInfo.AccountDeployNum + int(factory.AccountNum)
				factoryInfo.AccountNum = factoryInfo.AccountNum + int(factory.AccountDeployNum)
			} else {
				factoryInfo = &ent.FactoryInfo{
					AccountDeployNum: int(factory.AccountDeployNum),
					AccountNum:       int(factory.AccountNum),
				}
				factoryInfoMap[factoryAddr] = factoryInfo
			}
			factoryInfo.ID = factory.Factory
			factoryInfo.Network = factory.Network
			factoryInfoMap[factoryAddr] = factoryInfo
			totalNum += factory.AccountDeployNum
		}

		for factory, factoryInfo := range factoryInfoMap {
			if len(factory) == 0 {
				continue
			}
			factoryInfo.Dominance = decimal.NewFromInt(int64(factoryInfo.AccountDeployNum)).DivRound(decimal.NewFromInt(totalNum), 4)
			saveOrUpdateFactoryDay(client, factory, factoryInfo)
		}
		now1 := time.Now()
		log.Printf("top factory hour statistic success, network:%s, spent:%d", network, now1.Second()-now.Second())

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
	log.Printf("top factory day, single statistic sync success, factory:%s", info.ID)
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
		log.Printf("top factory hour statistic start timeRange:%d, network:%s", timeRange, network)
		client, err := entity.Client(context.Background(), network)
		if err != nil {
			continue
		}
		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-24*timeRange, 0, 0, 0, now.Location())
		endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
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
		var totalNum = int64(0)
		var repeatMap = make(map[string]bool)
		for _, factory := range factoryStatisHours {
			factoryAddr := factory.Factory
			timeStr := factory.StatisTime.String()
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
				factoryInfoMap[factoryAddr] = factoryInfo
			}
			factoryInfo.ID = factory.Factory
			factoryInfo.Network = factory.Network
			factoryInfoMap[factoryAddr] = factoryInfo
			if timeRange == 1 {
				totalNum += factory.AccountNum
			} else if timeRange == 7 {
				totalNum += factory.AccountNum
			} else if timeRange == 30 {
				totalNum += factory.AccountNum
			}

		}

		for factory, factoryInfo := range factoryInfoMap {
			if len(factory) == 0 {
				continue
			}
			if timeRange == 1 {
				factoryInfo.DominanceD1 = getSingleRate(int64(factoryInfo.AccountNumD1), totalNum)
			} else if timeRange == 7 {
				factoryInfo.DominanceD7 = getSingleRate(int64(factoryInfo.AccountNumD7), totalNum)
			} else if timeRange == 30 {
				factoryInfo.DominanceD30 = getSingleRate(int64(factoryInfo.AccountNumD30), totalNum)
			}
			saveOrUpdateFactory(client, factory, factoryInfo, timeRange)
		}

		allFactory, err := client.FactoryInfo.Query().All(context.Background())
		if len(allFactory) > 0 {
			for _, factory := range allFactory {
				_, ok := factoryInfoMap[factory.ID]
				if !ok {
					if timeRange == 1 {
						err = client.FactoryInfo.UpdateOneID(factory.ID).
							SetAccountDeployNumD1(0).
							SetAccountNumD1(0).
							SetDominanceD1(decimal.Zero).
							Exec(context.Background())
					} else if timeRange == 7 {
						err = client.FactoryInfo.UpdateOneID(factory.ID).
							SetAccountDeployNumD7(0).
							SetAccountNumD7(0).
							SetDominanceD7(decimal.Zero).
							Exec(context.Background())
					} else if timeRange == 30 {
						err = client.FactoryInfo.UpdateOneID(factory.ID).
							SetAccountDeployNumD30(0).
							SetAccountNumD30(0).
							SetDominanceD30(decimal.Zero).
							Exec(context.Background())
					}
				}
			}
		}
		log.Printf("top factory hour statistic success timeRange:%s, network:%s", string(timeRange), network)
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
	log.Printf("top factory hour, single statistic sync success, factory:%s", info.ID)
}
