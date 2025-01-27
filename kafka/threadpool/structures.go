package threadpool

import (
	"sync"

	kafka "github.com/segmentio/kafka-go"
)

type Channels struct {
	WorkerFileChannel       chan any
	SharedChan              chan any
	ProcessFileChannel      chan any
	ChildWorkerReceiverChan chan MsgStruct
	ChildWorkerSenderChan   chan KafkaArray
}
type Pools struct {
	ProcessFilePool chan chan any
	kWriterPool     chan chan KafkaArray
}

type KafkaWorker struct {
	id          int
	kChannel    chan KafkaArray
	kWriterPool chan chan KafkaArray
	Wg          *sync.WaitGroup
}

type Worker struct {
	Id              int
	Wg              *sync.WaitGroup
	Channel         Channels
	pool            Pools
	miniWorkers     []Worker
	miniWorkerCount int
	GlobalIndex     int
	mu *sync.Mutex
}

type Pool struct {
	Workers      []*Worker
	KafkaWorkers []KafkaWorker
	Wg           *sync.WaitGroup
	WorkSubmitWg *sync.WaitGroup
	Channel      Channels
	pool         Pools
	ResultMap    map[int][]kafka.Message
}

type MsgStruct struct {
	parentWorkerId int
	id             int
	msg            []byte
}

type KafkaArray struct {
	parentWorkerId int
	id             int
	msg            []kafka.Message
}
