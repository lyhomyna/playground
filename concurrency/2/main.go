package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
    urls := []string {"https://google.com", "https://youtube.com", "https://reddit.com"}
    
    var wg sync.WaitGroup

    for _, v := range urls {
	wg.Add(1)

	go func(url string) {
	    httpGet(url)
	    wg.Done()
	} (v)
    }

    wg.Wait()
}

func httpGet(url string) {
    startTime := time.Now()

    res, err := http.Get(url)
    if err != nil {
	fmt.Printf("Failure. URL: '%s\n.", url)
	return
    }
    defer res.Body.Close()

    fmt.Printf("'%s': \tLatency %d ms.\n", url, time.Since(startTime).Milliseconds())
}
