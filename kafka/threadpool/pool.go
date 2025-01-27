package threadpool

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"syscall"
)

func NewPool(noOfWorkers int, queueCapacity int) *Pool {
	c := Channels{
		SharedChan:         make(chan any, noOfWorkers),
		ProcessFileChannel: make(chan any, noOfWorkers),
	}
	p := Pools{
		ProcessFilePool: make(chan chan any, 200),
		kWriterPool:     make(chan chan KafkaArray, 200),
	}
	workers := make([]*Worker, noOfWorkers)
	kworkers := make([]KafkaWorker, noOfWorkers)

	workersWg := &sync.WaitGroup{}
	kworkersWg := &sync.WaitGroup{}

	for i := 0; i < noOfWorkers; i++ {
		workers[i] = NewWorker(i, workersWg, c, p)
	}

	for i := 0; i < noOfWorkers; i++ {
		kworkers[i] = NewKafkaWorker(i, kworkersWg, c, p)
	}
	return &Pool{
		Workers:      workers,
		Wg:           workersWg,
		KafkaWorkers: kworkers,
		WorkSubmitWg: &sync.WaitGroup{},
		Channel:      c,
		pool:         p,
	}
}

func (p *Pool) SubmitWorkToBoss() {
	for i := 0; i < 3; i++ {
		p.WorkSubmitWg.Add(1)
		fileName := "newfile" + strconv.Itoa(i+1) + ".txt"
		file, err := os.OpenFile(fileName, syscall.O_RDONLY, 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Sending File ", file.Name())
		p.Channel.ProcessFileChannel <- file
	}
}

func (p *Pool) StartWorkers() {
	for i := 0; i < len(p.KafkaWorkers); i++ {
		go p.KafkaWorkers[i].kafkaWorkerStart()
	}
	for i := range p.Workers {
		go p.Workers[i].Start()
	}

}

func (p *Pool) Start() {
	for {
		select {
		case work := <-p.Channel.ProcessFileChannel:
			workerFileChan := <-p.pool.ProcessFilePool
			workerFileChan <- work
		}
	}

}
