package threadpool

import (
	"sync"

	kafka "github.com/segmentio/kafka-go"
)

type Channels struct {
	WorkerFileChannel       chan any
	SharedChan              chan any
	ProcessFileChannel      chan any
	ChildWorkerReceiverChan chan any
	ChildWorkerSenderChan   chan any
}
type SharedPool struct {
	ReadyPool       chan chan any
	ProcessFilePool chan chan any
}
type Worker struct {
	Id          int
	Wg          *sync.WaitGroup
	WorkerQueue []int
	Channel     Channels
	pool        SharedPool
	miniWorkers []Worker
}

type Pool struct {
	Workers      []*Worker
	Wg           *sync.WaitGroup
	WorkSubmitWg *sync.WaitGroup
	Pm           *sync.Mutex
	Channel      Channels
	pool         SharedPool
	ResultMap    map[int][]kafka.Message
}

type MsgStruct struct {
	id  string
	msg string
}

type KafkaStruct struct {
	id  string
	msg kafka.Message
}
