package main

import (
	"fmt"
	"io"
	"os"
)

type myLogger struct{}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error :", err)
		os.Exit(1)
	}
	ml := myLogger{}
	io.Copy(ml, f)
}

func (ml myLogger) Write(p []byte) (int, error) {
	fmt.Println("Text from file :", string(p))
	return len(p), nil
}
