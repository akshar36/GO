package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
)
type myLogger struct{}
func main() {
	resp, err := http.Get("http://google.com")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ml := myLogger{}
	io.Copy(ml, resp.Body)
	
}

func (ml myLogger) Write(p []byte) (int,error) {
	fmt.Println("INFO[] :", string(p))
	return len(p), nil
}