package poolstuff

import (
	"fmt"
	"sync"
	"time"
)

type worker struct {
	id        int
	readyPool chan chan Work
	done      *sync.WaitGroup
	work      chan Work // You will put this chan Work in to readypool and that's why readyPool is of type chan chan Work
	quit      chan bool
}

func NewWorker(id int, readyPool chan chan Work, waitGroup *sync.WaitGroup) *worker {
	return &worker{
		id:        id,
		readyPool: readyPool,
		done:      waitGroup,
		work:      make(chan Work),
		quit:      make(chan bool),
	}
}

func (w *worker) Process(work Work){
	time.Sleep(time.Second)
	work.Run()
}

func (w *worker) Start() {
	fmt.Println("Starting Worker 🤖 ", w.id)
	go func() { // All you need is the keyword 'go' to spawn a thread :)
		for {
			w.readyPool <- w.work // Submit the worker's channel to boss
			select {
				case work := <- w.work: // If the boss has assigned work for you
					fmt.Println("[WORKER 🤖] : Got your work boss", w.id)
					w.Process(work)
			}
		}
	}()
}
