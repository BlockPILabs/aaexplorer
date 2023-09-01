package task

import (
	"context"
	"github.com/BlockPILabs/aa-scan/config"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/tokenpriceinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/userassetinfo"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/useropsinfo"
	"github.com/BlockPILabs/aa-scan/internal/utils"
	"github.com/procyon-projects/chrono"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

func InitDayStatis() {
	dayScheduler := chrono.NewDefaultTaskScheduler()

	_, err := dayScheduler.ScheduleWithCron(func(ctx context.Context) {
		doDayStatis()
	}, "0 15 0 * * ?")

	if err == nil {
		log.Print("dayStatis has been scheduled")
	}

}

func doDayStatis() {
	client, err := entity.Client(context.Background())
	if err != nil {
		return
	}
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

	dailyStatisticDay := calDailyStatistic(client, opsInfos, startTime)

	bundlerList := calBundlerStatisDay(client, bundlerMap, startTime)
	paymasterList := calPaymasterStatisDay(client, paymasterMap, startTime)
	factoryList := calFactoryStatisDay(client, factoryMap, startTime)

	bulkInsertBundlerStatsDay(context.Background(), client, bundlerList)
	bulkInsertPaymasterStatsDay(context.Background(), client, paymasterList)
	bulkInsertFactoryStatsDay(context.Background(), client, factoryList)
	dailyStatisticDay.Save(context.Background())

	saveWhaleStatisticDay(context.Background(), client, startTime)
}

func saveWhaleStatisticDay(ctx context.Context, client *ent.Client, time time.Time) {
	assetInfos, err := client.UserAssetInfo.Query().Where(userassetinfo.NetworkEQ("")).All(ctx)
	if err != nil {
		return
	}
	if len(assetInfos) == 0 {
		return
	}
	var contractMap = make(map[string]int)
	var contracts []string
	for _, asset := range assetInfos {

		_, ok := contractMap[asset.ContractAddress]
		if !ok {
			contractMap[asset.ContractAddress] = 1
			contracts = append(contracts, asset.ContractAddress)
		}
	}
	priceInfos, err := client.TokenPriceInfo.Query().Where(tokenpriceinfo.ContractAddressIn(contracts[:]...)).All(ctx)
	if err != nil {
		return
	}
	if len(priceInfos) == 0 {
		return
	}
	var priceMap = make(map[string]decimal.Decimal)
	for _, priceInfo := range priceInfos {
		priceMap[priceInfo.ContractAddress] = priceInfo.TokenPrice
	}

	var valueMap = make(map[string]decimal.Decimal)
	for _, asset := range assetInfos {
		preValue, accountOk := valueMap[asset.AccountAddress]
		if !accountOk {
			preValue = decimal.Zero
		}
		price, ok := priceMap[asset.ContractAddress]
		if !ok {
			valueMap[asset.AccountAddress] = preValue
		} else {
			valueMap[asset.AccountAddress] = asset.Amount.Mul(price).Add(preValue)
		}

	}

	var totalUsd = decimal.Zero
	var addrCount = 0
	for _, value := range valueMap {
		if value.Cmp(decimal.NewFromInt(config.WhaleUsd)) < 0 {
			continue
		}
		addrCount += 1
		totalUsd = totalUsd.Add(value)

	}
	whaleDay := client.WhaleStatisticDay.Create().SetWhaleNum(int64(addrCount)).SetTotalUsd(totalUsd).SetNetwork("").SetStatisticTime(time)
	whaleDay.Save(ctx)

}

func calDailyStatistic(client *ent.Client, infos []*ent.UserOpsInfo, startTime time.Time) *ent.DailyStatisticDayCreate {
	if len(infos) == 0 {
		return nil
	}
	var totalGasFee decimal.Decimal
	var txMap = make(map[string]bool)
	var walletMap = make(map[string]bool)
	for _, opsInfo := range infos {
		totalGasFee = opsInfo.Fee.Add(totalGasFee)
		txMap[opsInfo.TxHash] = true
		walletMap[opsInfo.Sender] = true
	}
	dailyStatistic := client.DailyStatisticDay.Create().
		SetNetwork(infos[0].Network).
		SetUserOpsNum(int64(len(infos))).
		SetStatisticTime(startTime).
		SetActiveWallet(int64(len(walletMap))).
		SetGasFee(totalGasFee).
		SetTxNum(int64(len(txMap)))

	return dailyStatistic
}

func bulkInsertFactoryStatsDay(ctx context.Context, client *ent.Client, data []*ent.FactoryStatisDayCreate) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	if _, err := client.FactoryStatisDay.CreateBulk(data...).Save(ctx); err != nil {
		tx.Rollback()
		log.Fatal(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func bulkInsertPaymasterStatsDay(ctx context.Context, client *ent.Client, data []*ent.PaymasterStatisDayCreate) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	if _, err := client.PaymasterStatisDay.CreateBulk(data...).Save(ctx); err != nil {
		tx.Rollback()
		log.Fatal(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func bulkInsertBundlerStatsDay(ctx context.Context, client *ent.Client, data []*ent.BundlerStatisDayCreate) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	err = client.BundlerStatisDay.CreateBulk(data...).Exec(ctx)
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

func calBundlerStatisDay(client *ent.Client, bundlerMap map[string][]*ent.UserOpsInfo, startTime time.Time) []*ent.BundlerStatisDayCreate {
	totalCount := 0
	var totalFee decimal.Decimal

	var bundlers []*ent.BundlerStatisDayCreate
	for key, userOpsInfoList := range bundlerMap {
		totalCount += len(userOpsInfoList)
		txHashMap := make(map[string]bool)
		for _, userOpsInfo := range userOpsInfoList {
			totalFee = totalFee.Add(userOpsInfo.Fee)
			txHashMap[userOpsInfo.TxHash] = true
		}
		bundlers = append(bundlers, client.BundlerStatisDay.Create().
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

func calPaymasterStatisDay(client *ent.Client, bundlerMap map[string][]*ent.UserOpsInfo, startTime time.Time) []*ent.PaymasterStatisDayCreate {
	totalCount := 0
	var totalFee decimal.Decimal

	var paymasters []*ent.PaymasterStatisDayCreate
	for key, userOpsInfoList := range bundlerMap {
		totalCount += len(userOpsInfoList)
		for _, userOpsInfo := range userOpsInfoList {
			totalFee = totalFee.Add(utils.DivRav(userOpsInfo.ActualGasCost))
		}
		paymasters = append(paymasters, client.PaymasterStatisDay.Create().
			SetPaymaster(key).
			SetNetwork(userOpsInfoList[0].Network).
			SetUserOpsNum(int64(totalCount)).
			SetGasSponsored(totalFee).
			SetStatisTime(startTime),
		)
	}

	return paymasters
}

func calFactoryStatisDay(client *ent.Client, bundlerMap map[string][]*ent.UserOpsInfo, startTime time.Time) []*ent.FactoryStatisDayCreate {
	accountDeployNum := 0

	var factories []*ent.FactoryStatisDayCreate
	for key, userOpsInfoList := range bundlerMap {
		accountMap := make(map[string]bool)
		for _, userOpsInfo := range userOpsInfoList {
			accountMap[userOpsInfo.Sender] = true
			if userOpsInfo.Factory != "" {
				accountDeployNum++
			}
		}
		factories = append(factories, client.FactoryStatisDay.Create().
			SetFactory(key).
			SetNetwork(userOpsInfoList[0].Network).
			SetStatisTime(startTime).
			SetAccountNum(int64(len(accountMap))).
			SetAccountDeployNum(int64(accountDeployNum)),
		)
	}

	return factories
}
