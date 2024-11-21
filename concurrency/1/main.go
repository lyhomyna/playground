package main

import (
	"fmt"
	"time"
)

func main() {
    currTime := time.Now()

    arr := []int {1,3,5,10} 

    //ch := make(chan int, len(arr))
    //withChan(arr, ch)
    withoutChan(arr)

    fmt.Printf("Elapsed time: %s\n", time.Since(currTime))
}

func withChan(arr []int, ch chan int) {
     for _, v := range arr {
         go func (num int) {
             time.Sleep(time.Second * time.Duration(num))
             ch <- num * 2
         } (v)
     }

     defer close(ch)

     for i := 0; i < len(arr); i++ {
         fmt.Println(<- ch)
     }
}

func withoutChan(arr []int) {
    for _, v := range arr {
	time.Sleep(time.Second * time.Duration(v))
	fmt.Println(v * 2)
    }
}
