package schedulerjobs

import "fmt"

type SchedulerTestJob struct{}

func (job SchedulerTestJob) Execute() {
	fmt.Printf("Test\n")
}

func (job SchedulerTestJob) GetJobInterval() uint32 {
	return 2
}
