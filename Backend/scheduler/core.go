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
	CBT  time.Duration `json:"CBT"`
	Name string        `json:"name"`
	AT   time.Time     `json:"AT"`
	QI   int           `json:"QI"`
}

func NewProcess(CBT time.Duration, name string, AT time.Time) *Process {
	return &Process{
		CBT:  CBT,
		Name: name,
		AT:   AT,
	}
}

func (proc *Process) ToString() string {
	return fmt.Sprintf("name: %s, CBT: %d, AT: %s", proc.Name, int(proc.CBT.Seconds()), strftime.Format(proc.AT, "%M:%S"))
}

type CPUUsage struct {
	ProcessName string
	Start       time.Time
	End         time.Time
	QI          int
}

func NewCPUUsage(processName string, start time.Time, end time.Time, qi int) *CPUUsage {
	return &CPUUsage{
		ProcessName: processName,
		Start:       start,
		End:         end,
		QI:          qi,
	}
}

type Queue struct {
	processes     chan *Process
	timeSlice     time.Duration
	MaxProcessCBT time.Duration
}

func NewQueue(timeSlice time.Duration, maxProcessCBT time.Duration) *Queue {
	return &Queue{
		processes:     make(chan *Process, 100),
		timeSlice:     timeSlice,
		MaxProcessCBT: maxProcessCBT,
	}
}

func (q *Queue) ToString() string {
	return fmt.Sprintf("time slice: %d, max CBT: %d", int(q.timeSlice.Seconds()), int(q.MaxProcessCBT.Seconds()))
}

type MultiLevelQueue struct {
	queues []*Queue
}

func NewMultiLevelQueue(queues []*Queue) *MultiLevelQueue {
	return &MultiLevelQueue{
		queues: queues,
	}
}

func (mlq *MultiLevelQueue) InsertProcess(process *Process, updateChannel chan UpdateLog) error {
	for i, queue := range mlq.queues {
		if queue.MaxProcessCBT >= process.CBT {
			queue.processes <- process
			process.QI = i
			updateChannel <- process.toUpdate()
			return nil
		}
	}

	return errors.New("couldn't find suitable queue")
}

func (mlq *MultiLevelQueue) ScheduleCPU(ctx context.Context, wg *sync.WaitGroup, updateChannel chan UpdateLog) {
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
		time.Sleep(minDuration(queue.timeSlice, process.CBT))
		end := time.Now()
		wg.Done()

		if updateChannel != nil {
			cu := NewCPUUsage(process.Name, start, end, process.QI)
			updateChannel <- cu.toUpdate()
			go fmt.Printf(
				"task: %s from queue with %s | start time: %s, end time %s \n",
				process.Name, queue.ToString(),
				strftime.Format(start, "%M:%S"), strftime.Format(end, "%M:%S"),
			)
		}

		if process.CBT > queue.timeSlice {
			newProcess := NewProcess(process.CBT-queue.timeSlice, process.Name, process.AT)
			queue.processes <- newProcess
			updateChannel <- process.toUpdate()
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

func minDuration(a, b time.Duration) time.Duration {
	if a <= b {
		return a
	}
	return b
}
