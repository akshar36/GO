package poolstuff

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	Id int
	Wg *sync.WaitGroup
}

func (j *Job) Run() {
	time.Sleep(time.Second)
	fmt.Println("[✅ Job] ID ", j.Id)
	j.Wg.Done()
}

type WorkerPool struct {
	readyPool     chan chan Work
	internalQueue chan Work
	workers       []*worker
	done          *sync.WaitGroup
}

func NewWorkerPool(noOfWorker int, capacity int) *WorkerPool {

	readyPool := make(chan chan Work, noOfWorker) // Think of it as workers adding their channel to the boss's channel
	workers := make([]*worker, noOfWorker)
	sharedStop := &sync.WaitGroup{} // We are creating a shared waitGroup (common among all workers they increment or decrement)

	for i := 0; i < noOfWorker; i++ {
		workers[i] = NewWorker(i, readyPool, sharedStop)
	}

	return &WorkerPool{
		readyPool:     readyPool,
		workers:       workers,
		internalQueue: make(chan Work, capacity),
		done:          sharedStop,
	}
}

func (p *WorkerPool) Start() {
	for i := 0; i < len(p.workers); i++ {
		p.workers[i].Start()
	}
	fmt.Println("🚀🚀🚀🚀 Workers Ready 🚀🚀🚀🚀")
	go p.Dispatch()
}

func (p *WorkerPool) Dispatch() {
	for {
		select {
		case work := <-p.internalQueue:
			fmt.Println("[BOSS 🗣️] : Got Some work boys")
			workerChannel := <-p.readyPool
			workerChannel <- work
		}
	}
}

func (p *WorkerPool) EnqueueJobs(jobs []*Job) {
	go func() {
		for _, actualJob := range jobs {
			fmt.Println("[BOSS 🗣️] : I am getting a job", actualJob.Id)
			p.internalQueue <- actualJob
		}
	}()
}
