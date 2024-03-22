package task

import (
	"context"
	internalconfig "github.com/BlockPILabs/aaexplorer/config"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aacontractinteract"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aauseropscalldata"
	"github.com/procyon-projects/chrono"
	"log"
	"time"
)

func AAContractInteractTask() {
	d1Scheduler := chrono.NewDefaultTaskScheduler()
	_, err := d1Scheduler.ScheduleWithCron(func(ctx context.Context) {
		day1InteractTask()
	}, "0 1 * * * *")
	if err != nil {
		log.Println(err)
	}

	dayScheduler := chrono.NewDefaultTaskScheduler()
	_, err = dayScheduler.ScheduleWithCron(func(ctx context.Context) {
		day7InteractTask()
		day30InteractTask()

	}, "0 5 0 * * *")

}

func day30InteractTask() {
	doInteractTaskDay(30)

}

func day7InteractTask() {
	doInteractTaskDay(7)
}

func doInteractTaskDay(days int) {
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
			return
		}
		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day()-days, 0, 0, 0, 0, now.Location())
		endTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		opsCalldatas, err := client.AAUserOpsCalldata.
			Query().
			Where(
				aauseropscalldata.TxTimeGTE(startTime.Unix()),
				aauseropscalldata.TxTimeLT(endTime.Unix()),
			).
			All(context.Background())

		if err != nil {
			log.Println(err)
		}
		if len(opsCalldatas) == 0 {
			return
		}

		targetMap := make(map[string]int64)
		for _, calldata := range opsCalldatas {
			target := calldata.Target
			count, targetOk := targetMap[target]
			if !targetOk {
				count = 0
			}
			count += 1
			targetMap[target] = count

		}
		var interactCreates []*ent.AAContractInteractCreate
		var statisticType = internalconfig.RangeD7
		if days == 30 {
			statisticType = internalconfig.RangeD30
		}
		for target, count := range targetMap {
			interactCreate := client.AAContractInteract.Create().
				SetStatisticType(statisticType).
				SetContractAddress(target).
				SetInteractNum(count).
				SetNetwork(network)
			interactCreates = append(interactCreates, interactCreate)
		}
		client.AAContractInteract.Delete().
			Where(aacontractinteract.StatisticTypeEQ(statisticType), aacontractinteract.NetworkEQ(network)).Exec(context.Background())
		err = client.AAContractInteract.CreateBulk(interactCreates...).Exec(context.Background())
		if err != nil {
			log.Println(err)
		}

		log.Printf("aa_contract_interact task d30 success, network:%s", network)
	}

}

func day1InteractTask() {
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
		startTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-24, 0, 0, 0, now.Location())
		endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
		opsCalldatas, err := client.AAUserOpsCalldata.
			Query().
			Where(
				aauseropscalldata.TxTimeGTE(startTime.Unix()),
				aauseropscalldata.TxTimeLT(endTime.Unix()),
			).
			All(context.Background())

		if err != nil {
			log.Println(err)
		}
		if len(opsCalldatas) == 0 {
			continue
		}

		targetMap := make(map[string]int64)
		for _, calldata := range opsCalldatas {
			target := calldata.Target
			count, targetOk := targetMap[target]
			if !targetOk {
				count = 0
			}
			count += 1
			targetMap[target] = count

		}
		var interactCreates []*ent.AAContractInteractCreate
		for target, count := range targetMap {
			interactCreate := client.AAContractInteract.Create().
				SetStatisticType(internalconfig.RangeH24).
				SetContractAddress(target).
				SetInteractNum(count).
				SetNetwork(network)
			interactCreates = append(interactCreates, interactCreate)
		}
		client.AAContractInteract.Delete().
			Where(aacontractinteract.StatisticTypeEQ(internalconfig.RangeH24), aacontractinteract.NetworkEQ(opsCalldatas[0].Network)).Exec(context.Background())
		_, err = client.AAContractInteract.CreateBulk(interactCreates...).Save(context.Background())
		if err != nil {
			log.Println(err)
		}
		log.Printf("aa_contract_interact task d1 success, network:%s", network)
	}

}
