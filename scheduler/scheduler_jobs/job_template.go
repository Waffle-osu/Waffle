package schedulerjobs

import "fmt"

type SchedulerTemplateJob struct{}

func (job SchedulerTemplateJob) Execute() {
	fmt.Printf("Test\n")
}

func (job SchedulerTemplateJob) GetJobInterval() uint32 {
	return 2
}
