package task

import (
	"context"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/aauseropscalldata"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/useroptypestatistic"
	"github.com/procyon-projects/chrono"
	"log"
	"time"
)

func UserOpTypeTask() {
	day1Task()
	d1Scheduler := chrono.NewDefaultTaskScheduler()
	_, err := d1Scheduler.ScheduleWithCron(func(ctx context.Context) {
		day1Task()
	}, "0 1 * * * ?")
	if err != nil {
		log.Println(err)
	}

	day7Task()
	day30Task()
	dayScheduler := chrono.NewDefaultTaskScheduler()
	_, err = dayScheduler.ScheduleWithCron(func(ctx context.Context) {
		day7Task()
		day30Task()

	}, "0 5 0 * * ?")

}

func day30Task() {
	doTaskDay(30)

}

func day7Task() {
	doTaskDay(7)

}

func doTaskDay(days int) {
	cli, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	records, err := cli.BlockScanRecord.Query().All(context.Background())
	if len(records) == 0 {
		return
	}
	for _, record := range records {
		network := record.Network
		client, err := entity.Client(context.Background())
		if err != nil {
			return
		}
		now := time.Now()
		//-days
		startTime := time.Date(now.Year(), now.Month(), now.Day()-200, 0, 0, 0, 0, now.Location())
		endTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		opsCalldatas, err := client.AAUserOpsCalldata.
			Query().
			Where(
				aauseropscalldata.TxTimeGTE(startTime.Unix()),
				aauseropscalldata.TxTimeLT(endTime.Unix()),
			).
			All(context.Background())

		if err != nil {
			log.Fatal(err)
		}
		if len(opsCalldatas) == 0 {
			return
		}

		sourceMap := make(map[string]int64)
		for _, calldata := range opsCalldatas {
			source := calldata.Source
			count, sourceOk := sourceMap[source]
			if !sourceOk {
				count = 0
			}
			count += 1
			sourceMap[source] = count

		}
		var statisticType = config.RangeD7
		if days == 30 {
			statisticType = config.RangeD30
		}
		var userOpCreates []*ent.UserOpTypeStatisticCreate
		for source, count := range sourceMap {
			userOpCreate := client.UserOpTypeStatistic.Create().
				SetStatisticType(statisticType).
				SetUserOpType(source).
				SetOpNum(count).
				SetUserOpSign(source).
				SetNetwork(network)
			userOpCreates = append(userOpCreates, userOpCreate)

		}
		client.UserOpTypeStatistic.Delete().
			Where(useroptypestatistic.StatisticTypeEQ(statisticType), useroptypestatistic.NetworkEQ(network)).Exec(context.Background())
		err = client.UserOpTypeStatistic.CreateBulk(userOpCreates...).Exec(context.Background())
		if err != nil {
			log.Println(err)
		}
	}

}

func day1Task() {
	cli, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	records, err := cli.BlockScanRecord.Query().All(context.Background())
	if len(records) == 0 {
		return
	}
	for _, record := range records {
		network := record.Network
		client, err := entity.Client(context.Background())
		if err != nil {
			return
		}
		now := time.Now()
		//-24
		startTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-10000, 0, 0, 0, now.Location())
		endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
		opsCalldatas, err := client.AAUserOpsCalldata.
			Query().
			Where(
				aauseropscalldata.TxTimeGTE(startTime.Unix()),
				aauseropscalldata.TxTimeLT(endTime.Unix()),
			).
			All(context.Background())

		if err != nil {
			log.Fatal(err)
		}
		if len(opsCalldatas) == 0 {
			return
		}

		sourceMap := make(map[string]int64)
		for _, calldata := range opsCalldatas {
			source := calldata.Source
			count, sourceOk := sourceMap[source]
			if !sourceOk {
				count = 0
			}
			count += 1
			sourceMap[source] = count

		}
		var userOpCreates []*ent.UserOpTypeStatisticCreate
		for source, count := range sourceMap {
			userOpCreate := client.UserOpTypeStatistic.Create().
				SetStatisticType(config.RangeH24).
				SetUserOpType(source).
				SetOpNum(count).
				SetUserOpSign(source).
				SetNetwork(network)
			userOpCreates = append(userOpCreates, userOpCreate)
		}
		client.UserOpTypeStatistic.Delete().
			Where(useroptypestatistic.StatisticTypeEQ(config.RangeH24), useroptypestatistic.NetworkEQ(opsCalldatas[0].Network)).Exec(context.Background())
		_, err = client.UserOpTypeStatistic.CreateBulk(userOpCreates...).Save(context.Background())
		if err != nil {
			log.Println(err)

		}
	}

}
