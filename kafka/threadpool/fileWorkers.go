package threadpool

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

type FileWorker struct{}

var (
	LINECOUNT   int
	WORKERCOUNT int
	mainfile *os.File
)

func (fw FileWorker) CreateFile(fileIdx int, mainWg *sync.WaitGroup) {
	defer mainWg.Done()
	LINECOUNT = 999999+1
	WORKERCOUNT = 50000
	ranges := []map[string]int{}
	for i := 0; i < WORKERCOUNT; i++ {
		ranges = append(ranges, map[string]int{"start": i * (LINECOUNT / WORKERCOUNT), "end": (i + 1) * (LINECOUNT / WORKERCOUNT)})
	}
	ranges[len(ranges)-1]["end"] = LINECOUNT
	results := make([]string, len(ranges))
	channel := make(chan struct {
		index  int
		result string
	}, len(ranges))
	wg := sync.WaitGroup{}
	for i, v := range ranges {
		wg.Add(1)
		go generateContent(v["start"], v["end"], i, &channel, &wg)
	}
	wg.Wait()
	close(channel)
	for r := range channel {
		results[r.index] = r.result
	}
	os.WriteFile("newfile"+strconv.Itoa(fileIdx)+".txt", []byte(strings.Join(results, "\n")), 0644)
}

func generateContent(start int, end int, index int, channel *chan struct {
	index  int
	result string
}, wg *sync.WaitGroup) {
	defer wg.Done()
	var output string
	for i := start; i < end; i++ {
		if i == end-1 {
			output += "topic1 " + strconv.Itoa(i)
		} else {
			output += "topic1 " + strconv.Itoa(i) + "\n"
		}
	}

	*channel <- struct {
		index  int
		result string
	}{index: index, result: output}

}
