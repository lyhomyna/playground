package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", rootHandler)

	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("server is running")
	http.ListenAndServe(":8080", nil)
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("visit-counter")
	if err != nil {
		cookie = &http.Cookie{
			Name:  "visit-counter",
			Value: "0",
		}
	}

	visitCounter, err := strconv.Atoi(cookie.Value)
	if err != nil {
		log.Fatal(err)
	}

	visitCounter += 1
	cookie.Value = strconv.Itoa(visitCounter)

	http.SetCookie(w, cookie)
	io.WriteString(w, fmt.Sprintf("You have visited this site by %d time.", visitCounter))
}
