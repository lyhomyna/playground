package main

import (
	"fmt"
	"net/http"
)

func main() {
    http.HandleFunc("/", rootHandler)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
	fmt.Println(err)
    }
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
    w.Write([]byte("q ;) _"))
}
