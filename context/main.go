package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
    defer cancel()

    urls := []string {
	"https://google.com",
	"https://youtube.com",
	"https://medium.com",
    }

    var wg sync.WaitGroup

    for _, url := range urls {
	wg.Add(1)
	go func () {
	    fetch(context.WithValue(ctx, "url", url))
	    wg.Done()
	}()
    }

    wg.Wait()
}

func fetch(ctx context.Context) {
    url := ctx.Value("url").(string)
    
    client := http.DefaultClient

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
	log.Fatalf("Failure to create request. %s", err)
    }

    res, err := client.Do(req)
    if err != nil {
	log.Fatalf("Failure to make request. %s", err)
    }

    log.Printf("Response from %s. %d", url, res.StatusCode)
}
