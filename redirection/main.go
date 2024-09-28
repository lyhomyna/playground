package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/redirect", handleRedirect)

	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("Server is running.")
	http.ListenAndServe(":8080", nil)
}

func handleRoot(w http.ResponseWriter, req *http.Request) {
	log.Println("root request")

	// Carriage Return Line Feed - \r\n
	w.Write([]byte("Hello, this is root handler!\r\n"))
}

// 301 - moved permanently
// 303 - see other
// 307 - temporary redirect
func handleRedirect(w http.ResponseWriter, req *http.Request) {
	log.Println("redirect request")

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusMovedPermanently)
}
