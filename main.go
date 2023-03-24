package main

import (
	"context"
	"fmt"
	"math/rand"
	scheduler "mlq/mlq_process_scheduling"
	"strconv"
	"sync"
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
		processes[i] = scheduler.NewProcess(rand.Intn(15)+1, "p"+strconv.Itoa(i))
	}

	return processes
}

func insertProcesses(ctx context.Context, wg *sync.WaitGroup, mlq *scheduler.MultiLevelQueue,
	processes []*scheduler.Process) {
	defer wg.Done()

	for _, process := range processes {
		wg.Add(1)
		_ = mlq.InsertProcess(process)

		fmt.Printf("inseting process: ")
		process.Print()
	}

	fmt.Printf("\n")
}
