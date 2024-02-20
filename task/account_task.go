package task

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/internal/entity"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aaaccountdata"
	"github.com/BlockPILabs/aaexplorer/internal/entity/ent/aauseropsinfo"
	interlog "github.com/BlockPILabs/aaexplorer/internal/log"
	"github.com/BlockPILabs/aaexplorer/service"
	"github.com/procyon-projects/chrono"
	"log"
	"os"
	"time"
)

var logger = interlog.NewTMLogger(interlog.NewSyncWriter(os.Stdout))

func AccountTask() {
	logger.Info("account-task init ")
	go doAccountTask()
	d1Scheduler := chrono.NewDefaultTaskScheduler()
	_, err := d1Scheduler.ScheduleWithCron(func(ctx context.Context) {
		doAccountTask()
	}, "0 30 0 * * *")
	if err != nil {
		log.Println(err)
	}

}

func doAccountTask() {
	logger.Info("account-task start ")
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
	inner:
		for {
			hour6Ago := time.Now().UnixMilli() - 6*3600*1000
			accountDatas, err := client.AaAccountData.Query().Where(aaaccountdata.LastTimeLT(hour6Ago), aaaccountdata.AaTypeEQ("aa")).Limit(1000).All(context.Background())

			if err != nil {
				log.Println(err)
				break inner
			}
			if len(accountDatas) == 0 {
				break inner
			}
			for _, accountData := range accountDatas {
				address := accountData.ID
				totalBalance := service.GetTotalBalance(address, network)
				userOpsNum, _ := client.AAUserOpsInfo.Query().Where(aauseropsinfo.SenderEqualFold(address)).Count(context.Background())
				err := client.AaAccountData.Update().SetLastTime(time.Now().UnixMilli()).SetTotalBalanceUsd(totalBalance).SetUserOpsNum(int64(userOpsNum)).Where(aaaccountdata.IDEQ(address)).Exec(context.Background())
				if err != nil {
					log.Printf("update data err %s", err)
				}
			}
			time.Sleep(1 * time.Second)
		}
		logger.Info("account-task update success  ", "network", network)
	}
}
