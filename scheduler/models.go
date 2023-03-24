package scheduler

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	strftime "github.com/itchyny/timefmt-go"
)

type Process struct {
	CBT  int
	name string
	AT   time.Time
}

func NewProcess(CBT int, name string, AT time.Time) *Process {
	return &Process{
		CBT:  CBT,
		name: name,
		AT:   AT,
	}
}

func (p *Process) ToString() string {
	return fmt.Sprintf("name: %s, CBT: %d, AT: %s", p.name, p.CBT, strftime.Format(p.AT, "%M:%S"))
}

type CPUUsage struct {
	processName string
	start       time.Time
	end         time.Time
	duration    time.Duration
}

func NewCPUUsage(processName string, start time.Time, end time.Time) *CPUUsage {
	return &CPUUsage{
		processName: processName,
		start:       start,
		end:         end,
		duration:    end.Sub(start),
	}
}

type Queue struct {
	processes     chan *Process
	timeSlice     int
	MaxProcessCBT int
}

func NewQueue(timeSlice int, maxProcessCBT int) *Queue {
	return &Queue{
		processes:     make(chan *Process, 100),
		timeSlice:     timeSlice,
		MaxProcessCBT: maxProcessCBT,
	}
}

type MultiLevelQueue struct {
	queues []*Queue
}

func NewMultiLevelQueue(queues []*Queue) *MultiLevelQueue {
	return &MultiLevelQueue{
		queues: queues,
	}
}

func (mlq *MultiLevelQueue) InsertProcess(process *Process) error {
	for _, queue := range mlq.queues {
		if queue.MaxProcessCBT >= process.CBT {
			queue.processes <- process
			return nil
		}
	}

	return errors.New("couldn't find suitable queue")
}

func (mlq *MultiLevelQueue) ScheduleCPU(ctx context.Context, wg *sync.WaitGroup, doneChannel chan *CPUUsage) {
	for true {
		process, queue, err := mlq.getProcess()
		if err != nil {
			time.Sleep(time.Millisecond * time.Duration(5))
			continue
		}

		if process.CBT > queue.timeSlice {
			wg.Add(1)
		}

		start := time.Now()
		time.Sleep(time.Second * time.Duration(queue.timeSlice))
		end := time.Now()
		wg.Done()

		if process.CBT > queue.timeSlice {
			_ = mlq.InsertProcess(NewProcess(process.CBT-queue.timeSlice, process.name, process.AT))
		}

		if doneChannel != nil {
			doneChannel <- NewCPUUsage(process.name, start, end)
			go fmt.Printf("task: %s, start time: %s, end time %s \n", process.name, strftime.Format(start, "%M:%S"), strftime.Format(end, "%M:%S"))
		}
	}
}

func (mlq *MultiLevelQueue) getProcess() (*Process, *Queue, error) {
	for _, queue := range mlq.queues {
		if len(queue.processes) != 0 {

			process, ok := <-queue.processes
			if !ok {
				return nil, nil, errors.New("couldn't read from channel")
			}
			return process, queue, nil
		}
	}

	return nil, nil, errors.New("all channels are empty")
}