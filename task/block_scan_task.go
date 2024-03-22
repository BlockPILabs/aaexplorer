package task

import (
	"context"
	"github.com/BlockPILabs/aaexplorer/third/schedule"
	"log"
)

func init() {
	schedule.Add("scan_block", func(ctx context.Context) {
		BlockScanRun(ctx)
	})
	log.Println("BlockScanRun has been scheduled")
}

func BlockScanRun(ctx context.Context) {

}
