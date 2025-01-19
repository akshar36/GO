package main

import (
	"fmt"
	"kafka-use/threadpool"
	"kafka-use/topic"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	t  topic.Topic
	tp *threadpool.Pool
)

func main() {
	startTime := time.Now()
	wg := &sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go createFile(i, wg)
	}
	wg.Wait()
	// _ = t.CreateTopic("MyNewTopic")
	tp = threadpool.NewPool(3, 10)
	tp.StartWorkers()
	go tp.Start()
	tp.SubmitWorkToBoss()
	tp.Wg.Wait()
	fmt.Println(time.Since(startTime))
}

func createFile(index int, wg *sync.WaitGroup) {
	defer wg.Done()
	stridx := strconv.Itoa(index + 1)
	file, err := os.Create("topic" + stridx + ".txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < 100000; i++ {
		_,_ = file.WriteString("topic" + stridx + ":" + strconv.Itoa(i+1)+"\n")
	}
}
