package main

import "threadpoolexecutor/poolstuff"

func main(){
	Pool := poolstuff.NewWorkerPool(4,4)
	Pool.Start()
}