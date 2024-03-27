package task

import (
	"github.com/BlockPILabs/aaexplorer/service"
	"github.com/BlockPILabs/aaexplorer/third/schedule"
	"time"
)

func init() {
	schedule.Add("function_signature", service.ScanSignature).ScheduleWithFixedDelay(time.Hour * 24)
}
