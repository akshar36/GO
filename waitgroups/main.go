package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	for i :=0; i< 5; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup){
			doSomeExpensiveTask(i, wg)
			time.Sleep(time.Second)
		}(i, &wg)
	}
	wg.Wait()
}

func doSomeExpensiveTask(i int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Task %v \n", i)
}