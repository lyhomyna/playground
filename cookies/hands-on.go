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
	http.HandleFunc("/hello-again", helloAgainHandler)
	http.HandleFunc("/delete-cookie", deleteCookie)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("server is running")
	http.ListenAndServe(":8080", nil)
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	_, err := req.Cookie("visit-counter")
	if err != nil {
		// counter cookie isn't present, so need to write it
		cookie := &http.Cookie{
			Name:  "visit-counter",
			Value: "1",
		}
		http.SetCookie(w, cookie)

		w.Header().Set("content-type", "text/html")
		io.WriteString(w, fmt.Sprintf("You have visites site for %s times.", cookie.Value))
	} else {
		// counter cookie is present - redirect to hello again page
		http.Redirect(w, req, "/hello-again", http.StatusSeeOther)
	}
}

func updateVisitCounterCookie(req *http.Request) (*http.Cookie, error) {
	cookie, err := req.Cookie("visit-counter")
	if err != nil {
		return nil, err
	}

	visitCounter, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return nil, err
	}

	visitCounter += 1
	cookie.Value = strconv.Itoa(visitCounter)

	return cookie, nil
}

func helloAgainHandler(w http.ResponseWriter, req *http.Request) {
	cookie, err := updateVisitCounterCookie(req)
	if err != nil {
		// there is not counter cookie - return to first page
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	w.Header().Set("content-type", "text/html")
	http.SetCookie(w, cookie)
	io.WriteString(w, fmt.Sprintf("You have visites site for %s times. <a href='/delete-cookie'>Reset</a>", cookie.Value))
}

func deleteCookie(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("visit-counter")
	if err != nil {
		// there is not counter cookie - return to first page
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	// delete cookie
	cookie.MaxAge = -1

	http.SetCookie(w, cookie)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}
