package main

import (
	"fmt"
	"sync"
	"time"
)

var maxCount = 10

func main() {
    pingerPonger := make(chan string)   
 
    var wd sync.WaitGroup
    wd.Add(2)

    go func () { 
	pinger(pingerPonger)
	wd.Done()
    } ()

    go func () { 
	ponger(pingerPonger)
	wd.Done()
    } ()

    wd.Wait()
}

func pinger(pingerPonger chan string) {
    pingerPonger <- "ping"
    for {
	time.Sleep(time.Second)

	msg, forMe := <- pingerPonger 
	fmt.Printf("Pinger. Message is: %s\n", msg)
	if forMe {
	    fmt.Println("Pinger says: nu blyat' it's for me ")
	} 

	if !forMe {
	    return
	}

	fmt.Println("Pinger says: ping.")
	pingerPonger <- "ping"
    }
}

func ponger(pingerPonger chan string) {
    counter := 0
    for {

	time.Sleep(time.Second)
	
	msg, forMe := <- pingerPonger
	fmt.Printf("Ponger. Message is: %s\n", msg)
	if forMe {
	    fmt.Println("Ponger says: nu blyat' it's for me")
	}

	if counter == maxCount {
	    fmt.Println("Ponger says: Haha")
	    close(pingerPonger)
	    return
	}

	fmt.Println("Ponger says: pong")

	pingerPonger <- "pong"

	counter++
    }
}
