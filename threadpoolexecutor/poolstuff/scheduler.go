package poolstuff

import (
	"fmt"
	"sync"
	"time"
)

type Scheduler struct {
	Wg     *sync.WaitGroup
	Active chan bool
}

type jobFunc func()

func NewScheduler() *Scheduler {
	return &Scheduler{
		Wg:     new(sync.WaitGroup),
		Active: make(chan bool),
	}
}
func (s *Scheduler) Add(j jobFunc, interval time.Duration) chan bool {
	activeChannel := make(chan bool)
	s.Wg.Add(1)
	go s.Process(j, interval, activeChannel)
	return activeChannel
}

func (s *Scheduler) Process(j jobFunc, interval time.Duration, activeChannel chan bool) {
	ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				fmt.Println("[TEST] Tick Every ", interval)
			case <-activeChannel:
				j()
			}
		}
}
