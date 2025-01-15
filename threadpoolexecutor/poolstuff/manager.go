package poolstuff

import (
	"fmt"
	"sync"
	"time"
)

type job struct {
	id int
	wg *sync.WaitGroup
}

func (j *job) Run() {
	time.Sleep(time.Second)
	fmt.Println("Performing the actual job", j.id)
	j.wg.Done()
}

type workerPool struct {
	readyPool     chan chan Work
	internalQueue chan Work
	workers       []*worker
	done          *sync.WaitGroup
}

func NewWorkerPool(noOfWorker int, capacity int) *workerPool {

	readyPool := make(chan chan Work, noOfWorker) // Think of it as workers adding their channel to the boss's channel
	workers := make([]*worker, noOfWorker)
	sharedStop := &sync.WaitGroup{} // We are creating a shared waitGroup (common among all workers they increment or decrement)

	for i := 0; i < noOfWorker; i++ {
		workers[i] = NewWorker(i, readyPool, sharedStop)
	}

	return &workerPool{
		readyPool:     readyPool,
		workers:       workers,
		internalQueue: make(chan Work),
		done:          sharedStop,
	}
}

func (p *workerPool) Start() {
	for i := 0; i < len(p.workers); i++ {
		p.workers[i].Start()
	}
	fmt.Println("Workers Ready !!!")
	jobs := []*job{}
	wg := &sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		jobs = append(jobs, &job{id: i, wg: wg})
	}
	p.done.Add(1)
	go p.Dispatch()
	p.EnqueueJobs(jobs)
	wg.Wait()
}

func (p *workerPool) Dispatch() {
	for {
		select {
		case work := <-p.internalQueue:
			fmt.Println("BOSS : Got Some work boys")
			workerChannel := <-p.readyPool
			workerChannel <- work
		}
	}
}

func (p *workerPool) EnqueueJobs(jobs []*job) {
	for _, actualJob := range jobs {
		fmt.Println("Enqueing Job", actualJob.id)
		p.internalQueue <- actualJob
	}
}
