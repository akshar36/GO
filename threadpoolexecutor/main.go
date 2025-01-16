package main

import (
	"sync"
	"threadpoolexecutor/poolstuff"
	"time"
	"fmt"
)

var Pool *poolstuff.WorkerPool

func testFunc() {
	wg := &sync.WaitGroup{}
	jobs := []*poolstuff.Job{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		jobs = append(jobs, &poolstuff.Job{Id: i, Wg: wg})
	}
	Pool.EnqueueJobs(jobs)
	wg.Wait()
	fmt.Println("🤖🤖🤖🤖 All jobs executed 🤖🤖🤖🤖")
}

func main() {
	fmt.Println("[BOSS 🗣️ ] : I am creating 4 workers and I can enqueue 5 jobs in one shot")
	Pool = poolstuff.NewWorkerPool(4, 5)
	Pool.Start()
	scheduler := poolstuff.NewScheduler()
	active := scheduler.Add(testFunc, 5*time.Second)
	timer := time.NewTicker(10 *time.Second)
	for range timer.C {
		fmt.Println("ENQUEUING JOBS EVERY 10 SECONDS")
		active <- true
	}

}
