package task

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/useropsinfo"
	"github.com/procyon-projects/chrono"
	"log"
	"time"
)

func InitHourStatis() {
	doHourStatis()
	hourScheduler := chrono.NewDefaultTaskScheduler()

	_, err := hourScheduler.ScheduleWithCron(func(ctx context.Context) {
		doHourStatis()
	}, "0 5 * * * ?")

	if err == nil {
		log.Print("hourStatis has been scheduled")
	}
}

func doHourStatis() {
	client, err := entity.Client(context.Background())
	//cli, err := sql.Open("postgres", "postgres://postgres:root@127.0.0.1:5432/postgres?sslmode=disable")
	if err != nil {
		return
	}
	//client := ent.NewClient(ent.Driver(cli))
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day()-5, 0, 0, 0, 0, now.Location())
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	opsInfos, err := client.UserOpsInfo.
		Query().
		Where(
			useropsinfo.TxTimeGTE(startTime.Unix()),
			useropsinfo.TxTimeLT(endTime.Unix()),
		).
		All(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	bundlerMap := make(map[string][]*ent.UserOpsInfo)
	paymasterMap := make(map[string][]*ent.UserOpsInfo)
	factoryMap := make(map[string][]*ent.UserOpsInfo)

	for _, opsInfo := range opsInfos {
		addOpsInfo(opsInfo.Bundler, opsInfo, bundlerMap)
		addOpsInfo(opsInfo.Paymaster, opsInfo, paymasterMap)
		addOpsInfo(opsInfo.Factory, opsInfo, factoryMap)

	}

	bundlerList := calBundlerStatis(client, bundlerMap, startTime)
	paymasterList := calPaymasterStatis(client, paymasterMap, startTime)
	factoryList := calFactoryStatis(client, factoryMap, startTime)

	bulkInsertBundlerStatsHour(context.Background(), client, bundlerList)
	bulkInsertPaymasterStatsHour(context.Background(), client, paymasterList)
	bulkInsertFactoryStatsHour(context.Background(), client, factoryList)

}

func calBundlerStatis(client *ent.Client, bundlerMap map[string][]*ent.UserOpsInfo, startTime time.Time) []*ent.BundlerStatisHourCreate {
	totalCount := 0
	var totalFee float32 = 0.0

	var bundlers []*ent.BundlerStatisHourCreate
	for key, userOpsInfoList := range bundlerMap {
		totalCount += len(userOpsInfoList)
		txHashMap := make(map[string]bool)
		for _, userOpsInfo := range userOpsInfoList {
			totalFee += userOpsInfo.Fee
			txHashMap[userOpsInfo.TxHash] = true
		}
		bundlers = append(bundlers, client.BundlerStatisHour.Create().
			SetBundler(key).
			SetNetwork(userOpsInfoList[0].Network).
			SetBundlesNum(int64(len(txHashMap))).
			SetGasCollected(totalFee).
			SetUserOpsNum(int64(totalCount)).
			SetStatisTime(startTime),
		)
	}

	return bundlers
}

func calPaymasterStatis(client *ent.Client, bundlerMap map[string][]*ent.UserOpsInfo, startTime time.Time) []*ent.PaymasterStatisHourCreate {
	totalCount := 0
	var totalFee float32 = 0.0

	var paymasters []*ent.PaymasterStatisHourCreate
	for key, userOpsInfoList := range bundlerMap {
		totalCount += len(userOpsInfoList)
		for _, userOpsInfo := range userOpsInfoList {
			totalFee += float32(userOpsInfo.ActualGasCost) / 1e18
		}
		paymasters = append(paymasters, client.PaymasterStatisHour.Create().
			SetPaymaster(key).
			SetNetwork(userOpsInfoList[0].Network).
			SetUserOpsNum(int64(totalCount)).
			SetGasSponsored(totalFee).
			SetStatisTime(startTime),
		)
	}

	return paymasters
}

func calFactoryStatis(client *ent.Client, bundlerMap map[string][]*ent.UserOpsInfo, startTime time.Time) []*ent.FactoryStatisHourCreate {
	accountDeployNum := 0

	var factories []*ent.FactoryStatisHourCreate
	for key, userOpsInfoList := range bundlerMap {
		accountMap := make(map[string]bool)
		for _, userOpsInfo := range userOpsInfoList {
			accountMap[userOpsInfo.Sender] = true
			if userOpsInfo.Factory != "" {
				accountDeployNum++
			}
		}
		factories = append(factories, client.FactoryStatisHour.Create().
			SetFactory(key).
			SetNetwork(userOpsInfoList[0].Network).
			SetStatisTime(startTime).
			SetAccountNum(int64(len(accountMap))).
			SetAccountDeployNum(int64(accountDeployNum)),
		)
	}

	return factories
}

func bulkInsertBundlerStatsHour(ctx context.Context, client *ent.Client, data []*ent.BundlerStatisHourCreate) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	err = client.BundlerStatisHour.CreateBulk(data...).Exec(ctx)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func bulkInsertPaymasterStatsHour(ctx context.Context, client *ent.Client, data []*ent.PaymasterStatisHourCreate) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	if _, err := client.PaymasterStatisHour.CreateBulk(data...).Save(ctx); err != nil {
		tx.Rollback()
		log.Fatal(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func bulkInsertFactoryStatsHour(ctx context.Context, client *ent.Client, data []*ent.FactoryStatisHourCreate) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	if _, err := client.FactoryStatisHour.CreateBulk(data...).Save(ctx); err != nil {
		tx.Rollback()
		log.Fatal(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
