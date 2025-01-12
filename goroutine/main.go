package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	links := []string{
		"http://google.com",
		"http://amazon.com",
		"http://golang.org",
		"http://youtube.com",
		"http://stackoverflow.com",
	}
	c := make(chan string)
	for _, link := range links {
		go checkLink(link, c)
	}

	for l := range c {
		go func() {
			// time.Sleep(5*time.Second)
			checkLink(l, c)
		}()
	}
}
func checkLink(link string, c chan string) {
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link)
		c <- link
		return
	}
	fmt.Println(link)
	c <- link
}
