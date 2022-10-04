package scheduler

import (
	schedulerjobs "Waffle/scheduler/scheduler_jobs"
	"time"
)

type SchedulerJob interface {
	//Gets executed every GetJobInterval() Seconds
	Execute()
	//Gets the Job Interval, i.e. how often to execute it, in seconds.
	GetJobInterval() uint32
}

var SchedulerReoccuringJobs []SchedulerJob

func InitializeJobScheduler() {
	SchedulerReoccuringJobs = append(SchedulerReoccuringJobs, schedulerjobs.SchedulerTestJob{})
}

func RunScheduler() {
	secondsPassed := uint32(1)

	for {
		for _, job := range SchedulerReoccuringJobs {
			if secondsPassed%job.GetJobInterval() == 0 {
				go job.Execute()
			}
		}

		time.Sleep(1000 * time.Millisecond)

		secondsPassed++
	}
}
