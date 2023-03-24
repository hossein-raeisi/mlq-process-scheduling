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

func main() {
	_ = run()
}

func run() error {
	queues := []*scheduler.Queue{
		scheduler.NewQueue(1, 2),
		scheduler.NewQueue(2, 4),
		scheduler.NewQueue(4, 8),
		scheduler.NewQueue(8, 16),
	}
	mlq := scheduler.NewMultiLevelQueue(queues)
	processes := generateRandomProcesses(5)
	ctx := context.Background()
	wg := &sync.WaitGroup{}

	doneChannel := make(chan *scheduler.CPUUsage, len(processes)*16+1)

	wg.Add(1)
	go insertProcesses(ctx, wg, mlq, processes)
	go mlq.ScheduleCPU(ctx, wg, doneChannel)
	wg.Wait()

	return nil
}

func generateRandomProcesses(number int) []*scheduler.Process {
	processes := make([]*scheduler.Process, number)

	for i := 0; i < number; i++ {
		processes[i] = scheduler.NewProcess(
			rand.Intn(15)+1,
			"p"+strconv.Itoa(i),
			time.Now().Add(time.Second*time.Duration(rand.Intn(10)+1)),
		)
	}

	sort.Slice(processes, func(i, j int) bool {
		return processes[i].AT.Before(processes[j].AT)
	})

	return processes
}

func insertProcesses(ctx context.Context, wg *sync.WaitGroup, mlq *scheduler.MultiLevelQueue,
	processes []*scheduler.Process) {
	defer wg.Done()

	for _, process := range processes {
		wg.Add(1)
		if process.AT.After(time.Now()) {
			time.Sleep(process.AT.Sub(time.Now()))
		}
		_ = mlq.InsertProcess(process)

		go fmt.Printf("inserted process: %s \n", process.ToString())
	}
}
