package threadpool

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func NewWorker(id int, wg *sync.WaitGroup, c Channels, p Pools) *Worker {
	return &Worker{
		Id: id,
		pool: Pools{
			ProcessFilePool: p.ProcessFilePool,
			kWriterPool:     p.kWriterPool,
		},
		Channel: Channels{
			WorkerFileChannel: make(chan any, 100000),
		},
		Wg:              wg,
		miniWorkerCount: 10,
		mu:              &sync.Mutex{},
	}
}

func NewKafkaWorker(id int, wg *sync.WaitGroup, c Channels, p Pools) KafkaWorker {
	return KafkaWorker{
		id:          id,
		kChannel:    make(chan KafkaArray, 10),
		kWriterPool: p.kWriterPool,
		Wg:          wg,
	}

}

func NewChildWorker(id int, pw *Worker) Worker {
	return Worker{
		Id:      id,
		Wg:      &sync.WaitGroup{},
		Channel: Channels{ChildWorkerReceiverChan: make(chan MsgStruct, 10000), ChildWorkerSenderChan: make(chan KafkaArray, 10000)},
		mu:      pw.mu,
	}
}

func (w *Worker) resizeChannel(oldChan *chan any, newSize int) {
	newChan := make(chan any, newSize)
	for val := range *oldChan {
		newChan <- val
	}
	*oldChan = newChan
}

func (w *Worker) Start() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	w.pool.ProcessFilePool <- w.Channel.WorkerFileChannel
	i := 0
	for {
		select {
		case file := <-w.Channel.WorkerFileChannel:
			w.ProcessFile(file)
			for {
				select {
				case val := <-w.miniWorkers[i].Channel.ChildWorkerSenderChan:
					kafkaWorkerChannel := <-w.pool.kWriterPool
					kafkaWorkerChannel <- val
					i += 1
					if i > w.miniWorkerCount-1 {
						i = 0
					}
				default:
					i += 1
					if i > w.miniWorkerCount-1 {
						i = 0
					}
				}
			}
		}
	}

}

func (w *Worker) cStart() {
	for {
		select {
		case childwork := <-w.Channel.ChildWorkerReceiverChan:
			w.workerReadFile(childwork)
		default:
			time.Sleep(Delay)
			continue
		}
	}
}

func (kw *KafkaWorker) kafkaWorkerStart() {
	kw.kWriterPool <- kw.kChannel
	for {
		select {
		case ka := <-kw.kChannel:
			kw.sendKafka(ka)
		default:
			time.Sleep(Delay)
			continue
		}
	}
}

func (kw *KafkaWorker) sendKafka(ka KafkaArray) {
	for i := 0; i < len(ka.msg); i++ {
		fmt.Println(string(ka.msg[i].Key) + " " + string(ka.msg[i].Value))
	}
	kw.kWriterPool <- kw.kChannel

}

func (w *Worker) ProcessFile(file any) {
	fD, err := file.(*os.File)
	if !err {
		fmt.Println(err)
		return
	}
	childWorkers := make([]Worker, w.miniWorkerCount)
	for i := 0; i < w.miniWorkerCount; i++ {
		childWorkers[i] = NewChildWorker(w.Id+10, w)
		w.miniWorkers = childWorkers
		go childWorkers[i].cStart()
	}
	i := 0
	reader := bufio.NewReader(fD)
	for {
		buffer := make([]byte, 4*1024)
		_, err := reader.Read(buffer)
		if err != nil && err == io.EOF {
			break
		}
		readNewline, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			continue
		}
		buffer = append(buffer, readNewline...)
		w.miniWorkers[i].Channel.ChildWorkerReceiverChan <- MsgStruct{id: w.GlobalIndex, msg: buffer, parentWorkerId: w.Id}
		w.GlobalIndex++
		i += 1
		if i > w.miniWorkerCount-1 {
			i = 0
		}
	}
	w.pool.ProcessFilePool <- w.Channel.WorkerFileChannel
}

func (w *Worker) workerReadFile(msg MsgStruct) {
	w.mu.Lock()
	lineArr := strings.Split(string(msg.msg), "\n")
	kafkaArr := KafkaArray{}
	for _, actualLine := range lineArr {
		nl := strings.Compare(string(actualLine), "")
		if nl == 0 {
			continue
		}
		kafkaArr.msg = append(kafkaArr.msg, kafka.Message{Key: []byte(strconv.Itoa(msg.parentWorkerId)), Value: []byte(actualLine), Topic: "topic1"})
	}
	kafkaArr.id = msg.id
	kafkaArr.parentWorkerId = msg.parentWorkerId
	w.Channel.ChildWorkerSenderChan <- kafkaArr
	w.mu.Unlock()
}
