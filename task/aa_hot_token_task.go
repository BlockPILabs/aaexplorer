package task

import (
	"context"
	"github.com/procyon-projects/chrono"
	"log"
)

func AAHotTokenTask() {
	d1Scheduler := chrono.NewDefaultTaskScheduler()
	_, err := d1Scheduler.ScheduleWithCron(func(ctx context.Context) {
		//day1TokenTask()
	}, "0 1 * * * *")
	if err != nil {
		log.Println(err)
	}

	dayScheduler := chrono.NewDefaultTaskScheduler()
	_, err = dayScheduler.ScheduleWithCron(func(ctx context.Context) {
		day7TokenTask()
		day30TokenTask()

	}, "0 5 0 * * *")

}

func day30TokenTask() {
	//doTokenTaskDay(30)

}

func day7TokenTask() {
	//doTokenTaskDay(7)
}

/**
func doTokenTaskDay(days int) {
	client, err := entity.Client(context.Background())
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
		target := calldata.
		count, targetOk := targetMap[target]
		if !targetOk {
			count = 0
		}
		count += 1
		targetMap[target] = count

	}
	network := opsCalldatas[0].Network
	var interactCreates []*ent.AAHotTokenStatisticCreate
	for target, count := range targetMap {
		interactCreate := client.AAHotTokenStatistic.Create().
			SetSymbol("").
			SetContractAddress(target).
			SetNetwork(opsCalldatas[0].Network).
			SetStatisticType("d" + string(days)).
			SetUsdAmount().Exec(context.Background())
		interactCreates = append(interactCreates, interactCreate)
	}
	client.AAContractInteract.Delete().
		Where(aacontractinteract.StatisticTypeEQ("d"+string(days)), aacontractinteract.NetworkEQ(network)).Exec(context.Background())
	err = client.AAContractInteract.CreateBulk(interactCreates...).Exec(context.Background())
	if err != nil {
		log.Println(err)
	}
}

func day1TokenTask() {
	client, err := entity.Client(context.Background())
	if err != nil {
		return
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
	for target, count := range targetMap {
		interactCreate := client.AAContractInteract.Create().
			SetStatisticType("d1").
			SetContractAddress(target).
			SetInteractNum(count).
			SetNetwork(opsCalldatas[0].Network)
		interactCreates = append(interactCreates, interactCreate)
	}
	client.AAContractInteract.Delete().
		Where(aacontractinteract.StatisticTypeEQ("d1"), aacontractinteract.NetworkEQ(opsCalldatas[0].Network)).Exec(context.Background())
	err = client.AAContractInteract.CreateBulk(interactCreates...).Exec(context.Background())
	if err != nil {
		log.Println(err)
	}
}

*/
