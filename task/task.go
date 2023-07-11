package task

import (
	"context"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/bundlerstatishour"
	"time"
)

type BundlerList struct {
	Address      string
	UserOpsNum   int
	BundlesNum   int64
	GasCollected float32
}

func InitTask() {

	//hour statistics
	InitHourStatis()

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

}

func topPaymaster() {

}

func topBundlers() {
	//get bundler list

	//
	//client := ent.NewClient(ent.Driver(config.NewDriver()))
	client, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-23, 0, 0, 0, now.Location())
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
	/*
		bundlerInfoMap := make(map[string]*ent.BundlerInfo)
		for _, bundlerStatisHour := range bundlerStatisHours {
			bundlerInfo = &ent.BundlerInfo{
				Bundler:        bundlerStatisHour.Bundler,
				Network:        bundlerStatisHour.Network,
				UserOpsNum:     bundlerStatisHour.UserOpsNum,
				BundlesNum:     bundlerStatisHour.BundlesNum,
				GasCollected:   bundlerStatisHour.GasCollected,
				UserOpsNumD1:   bundlerStatisHour.UserOpsNum,
				BundlesNumD1:   bundlerStatisHour.BundlesNum,
				GasCollectedD1: bundlerStatisHour.GasCollected,
			}
		}

	*/

}
