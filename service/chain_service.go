package service

import (
	"context"
	"fmt"
	"github.com/BlockPILabs/aa-scan/internal/entity"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/blockdata"
	"github.com/BlockPILabs/aa-scan/internal/entity/ent/blockscanrecord"
)

//1.scan block, get active account, contract
//2.add trace table
//3.mev detect

func ScanBlock() {
	client, err := entity.Client(context.Background())
	if err != nil {
		return
	}
	records, err := client.BlockScanRecord.Query().All(context.Background())
	if len(records) == 0 {
		return
	}
	for _, record := range records {
		network := record.Network
		last := record.LastBlockNumber
		fmt.Println(network)
		//1.get max block num by network
		//2.
		blockData, err := client.BlockData.Query().Order(ent.Desc(blockdata.FieldBlockNum)).Limit(1).All(context.Background())
		if err != nil {
			continue
		}
		if len(blockData) == 0 {
			continue
		}
		for i := last + 1; i <= blockData[0].BlockNum; i++ {
			//do biz

			client.BlockScanRecord.Update().Where(blockscanrecord.IDEQ(record.ID)).SetLastBlockNumber(i).Exec(context.Background())
		}
	}
}
