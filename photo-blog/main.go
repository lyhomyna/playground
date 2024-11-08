package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var sessionCookieName = "session"

var tpl *template.Template

func init() {
    tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
    http.HandleFunc("/", index)
    http.Handle("/favicon.ico", http.NotFoundHandler())
    
    log.Println("Listening on localhost:8080.")
    http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
    c := sessionCookie(w, req)
    appendPictures(c)
    http.SetCookie(w, c)
    log.Println(c)
    tpl.ExecuteTemplate(w, "index.gohtml", c)
}

func sessionCookie(w http.ResponseWriter,req *http.Request) *http.Cookie {
    c, err := req.Cookie(sessionCookieName)
    if err != nil {
	c = &http.Cookie {
	    Name: sessionCookieName,
	    Value: uuid.New().String(),
	}

	http.SetCookie(w, c)
    } 

    return c
}

func appendPictures(c *http.Cookie) {
    p1 := "sunset.jpeg"
    p2 := "disneyland.jpeg"
    p3 := "beach.jpeg"

    if !strings.Contains(c.Value, p1) {
	c.Value += "|" + p1
    }

    if !strings.Contains(c.Value, p2) {
	c.Value += "|" + p2
    }

    if !strings.Contains(c.Value, p3) {
	c.Value += "|" + p3
    }
}
