package main

import (
	"fmt"
	"kafka-use/threadpool"
	"sync"
	"time"
	"os"
)

var (
	tp          *threadpool.Pool
	fw          threadpool.FileWorker
	FILECOUNT   int
	WORKERCOUNT int
)
var mainfile os.File
func main() {
	FILECOUNT = 3
	WORKERCOUNT = 3
	
	startTime := time.Now()
	wg := &sync.WaitGroup{}
	for i := 0; i < FILECOUNT; i++ {
		wg.Add(1)
		fmt.Println("Creating File ", i+1)
		go fw.CreateFile(i+1, wg)
	}
	wg.Wait()
	fmt.Printf("%v files created in %v \n", FILECOUNT, time.Since(startTime))
	fmt.Printf("Creating A Worker Pool of %v Workers \n", WORKERCOUNT)
	tp = threadpool.NewPool(WORKERCOUNT, 10)
	tp.StartWorkers()
	go tp.Start()
	tp.SubmitWorkToBoss()
	tp.WorkSubmitWg.Wait()

}
