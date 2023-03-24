package mlq_process_scheduling

import (
	"time"
)

type Process struct {
	cbt  int
	name string
}

func NewProcess(cbt int, name string) *Process {
	return &Process{
		cbt:  cbt,
		name: name,
	}
}

type CPUUsage struct {
	process  Process
	start    time.Time
	end      time.Time
	duration time.Duration
}

func NewCPUUsage(process Process, start time.Time, end time.Time) *CPUUsage {
	return &CPUUsage{
		process:  process,
		start:    start,
		end:      end,
		duration: end.Sub(start),
	}
}

type Queue struct {
	processes     chan Process
	timeSlice     int
	MaxProcessCBT int
}

func NewQueue(timeSlice int, maxProcessCBT int) *Queue {
	return &Queue{
		processes:     make(chan Process, 100),
		timeSlice:     timeSlice,
		MaxProcessCBT: maxProcessCBT,
	}
}

type MultiLevelQueue struct {
	queues []Queue
}

func NewMultiLevelQueue(queues []Queue) *MultiLevelQueue {
	return &MultiLevelQueue{
		queues: queues,
	}
}
