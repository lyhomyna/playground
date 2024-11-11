package main

import (
	"io"
	"net/http"
)

func main() {
    http.HandleFunc("/", index)
    http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
}

func index(w http.ResponseWriter, _ *http.Request) {
    io.WriteString(w, "HTTPS connection established.")
}
