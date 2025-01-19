package threadpool

import (
	"github.com/segmentio/kafka-go"
	"os"
	"sync"
)

type Worker struct {
	Id          int
	WorkChan    chan *os.File
	ReadyPool   chan chan *os.File
	Wg          *sync.WaitGroup
	WorkerQueue []int
	sharedChan  *chan []kafka.Message
}

type Pool struct {
	ReadyPool     chan chan *os.File
	Workers       []*Worker
	InternalQueue chan *os.File
	Wg            *sync.WaitGroup
	WorkSubmitWg  *sync.WaitGroup
	sharedChan    *chan []kafka.Message
	Pm            *sync.Mutex
}
