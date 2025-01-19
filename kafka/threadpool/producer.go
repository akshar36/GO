package threadpool

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	kafka "github.com/segmentio/kafka-go"
)

func NewPool(noOfWorkers int, queueCapacity int) *Pool {
	workers := make([]*Worker, noOfWorkers)
	workersWg := &sync.WaitGroup{}
	readyPool := make(chan chan *os.File, noOfWorkers)
	sharedChan := make(chan []kafka.Message, noOfWorkers)
	for i := 0; i < noOfWorkers; i++ {
		workers[i] = NewWorker(i, workersWg, readyPool, &sharedChan)
	}
	return &Pool{
		ReadyPool:     readyPool,
		InternalQueue: make(chan *os.File, queueCapacity),
		Workers:       workers,
		Wg:            workersWg,
		WorkSubmitWg:  &sync.WaitGroup{},
		sharedChan:    &sharedChan,
		Pm: &sync.Mutex{},
	}
}

func (p *Pool) Start() {
	for {
		select {
		case work := <-p.InternalQueue:
			workerChan := <-p.ReadyPool
			workerChan <- work
		case send := <-*p.sharedChan:
			p.WorkSubmitWg.Add(1)
			fmt.Println("Received from worker")
			p.SendKafka(send)

		}
	}
}

func (p *Pool) SendKafka(km []kafka.Message) {
	fmt.Printf("Worker Received %v messages \n", len(km))
	defer p.Wg.Done()
	w := &kafka.Writer{
		Addr:                   kafka.TCP("localhost:19092"),
		AllowAutoTopicCreation: true,
	}

	err := w.WriteMessages(context.Background(), km...)
	if err != nil {
		panic(err.Error())
	}
	defer w.Close()
	fmt.Println("Written successfully")
}

func (p *Pool) StartWorkers() {
	for i := range p.Workers {
		go p.Workers[i].Start()
	}
}

func (w *Worker) Start() {
		w.ReadyPool <- w.WorkChan
		for work := range w.WorkChan {
			w.Process(work)
			w.ReadyPool <- w.WorkChan
		}
}

func (w *Worker) Process(file *os.File) {
	fileName := strings.Split(file.Name(), ".")[0]
	defer file.Close()
	scanner := bufio.NewScanner(file)
	fmt.Println("Processing file: ", fileName)
	km := []kafka.Message{}
	for scanner.Scan() {
		keyval := strings.Split(scanner.Text(), ":")
		km = append(km, kafka.Message{
			Topic: fileName,
			Key:   []byte(fileName),
			Value: []byte(keyval[1]),
		})
	}
	fmt.Printf("Worker %v Finished Processing File %v Sending Messages of Length %v Back To Boss \n", w.Id, fileName, len(km))
	*w.sharedChan <- km

}

func NewWorker(id int, wg *sync.WaitGroup, rp chan chan *os.File, sm *chan []kafka.Message) *Worker {
	return &Worker{
		Id:         id,
		WorkChan:   make(chan *os.File),
		ReadyPool:  rp,
		Wg:         wg,
		sharedChan: sm,
	}
}

func (p *Pool) SubmitWorkToBoss() {
	for i := 0; i < 3; i++ {
		p.Wg.Add(1)
		fileName := "topic" + strconv.Itoa(i+1) + ".txt"
		file, err := os.OpenFile(fileName, syscall.O_RDONLY, 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Sending File ", file.Name())
		p.InternalQueue <- file
	}
}
