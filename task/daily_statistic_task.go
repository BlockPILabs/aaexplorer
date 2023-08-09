package task

import (
	"context"
	"github.com/procyon-projects/chrono"
	"log"
)

func ExecuteTask() {
	dayScheduler := chrono.NewDefaultTaskScheduler()

	_, err := dayScheduler.ScheduleWithCron(func(ctx context.Context) {
		dailyStatisticTask()
	}, "0 15 0 * * ?")

	if err == nil {
		log.Print("dayStatisticTask has been scheduled")
	}
}

func dailyStatisticTask() {
	day1()

	day7()

	day30()

}

func day30() {

}

func day7() {

}

func day1() {

}
