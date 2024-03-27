package schedule

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	Add("test", func(ctx context.Context) {
		fmt.Println(time.Now())
	}).ScheduleWithCron("0 0 9 */1 * *")
	Schedule(context.Background())
	select {}
}
