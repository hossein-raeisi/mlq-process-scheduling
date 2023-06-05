package main

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"

	"mlq/scheduler"
)

var (
	ProcessesNumber   = 5
	ProcessMaxCBT     = 16
	ProcessMaxATDelay = 10
)

func main() {
	run()
}

func run() {
	queues := []*scheduler.Queue{
		scheduler.NewQueue(time.Second*1, time.Second*2),
		scheduler.NewQueue(time.Second*2, time.Second*4),
		scheduler.NewQueue(time.Second*4, time.Second*8),
		scheduler.NewQueue(time.Second*8, time.Second*16),
	}
	mlq := scheduler.NewMultiLevelQueue(queues)
	processes := generateRandomProcesses(ProcessesNumber)
	ctx := context.Background()
	wg := &sync.WaitGroup{}

	doneChannel := make(chan *scheduler.CPUUsage, len(processes)*ProcessMaxCBT)
	updateChannel := make(chan scheduler.UpdateLog, len(processes)*(ProcessMaxCBT+1))

	wg.Add(1)
	go insertProcesses(ctx, wg, mlq, processes, updateChannel)
	go mlq.ScheduleCPU(ctx, wg, doneChannel, updateChannel)
	wg.Wait()

	scheduler.Display(updateChannel)
}

func generateRandomProcesses(number int) []*scheduler.Process {
	processes := make([]*scheduler.Process, number)

	for i := 0; i < number; i++ {
		processes[i] = scheduler.NewProcess(
			time.Second*time.Duration(rand.Intn(ProcessMaxCBT-1)+1),
			"p"+strconv.Itoa(i),
			time.Now().Add(time.Second*time.Duration(rand.Intn(ProcessMaxATDelay-1)+1)),
		)
	}

	sort.Slice(processes, func(i, j int) bool {
		return processes[i].AT.Before(processes[j].AT)
	})

	return processes
}

func insertProcesses(ctx context.Context, wg *sync.WaitGroup, mlq *scheduler.MultiLevelQueue,
	processes []*scheduler.Process, updateChannel chan scheduler.UpdateLog) {
	defer wg.Done()

	for _, process := range processes {
		wg.Add(1)
		if process.AT.After(time.Now()) {
			time.Sleep(process.AT.Sub(time.Now()))
		}
		_ = mlq.InsertProcess(process, updateChannel)

		go fmt.Printf("inserted process: %s \n", process.ToString())
	}
}
