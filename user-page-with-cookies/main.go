package main

import (
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/google/uuid"
)

type user struct {
	Email     string
	FirstName string
	LastName  string
}

var sessionsDb = map[string]string{}
var usersDb = map[string]user{}
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("templates/index.gohtml", "templates/bar.gohtml", "templates/info.gohtml"))
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/info", infoHandler)
	http.HandleFunc("/bar", barHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("server is running.")
	http.ListenAndServe(":8080", nil)
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	// write session id cookie
	sessionIdCookie, _ := req.Cookie("session-id")
	if sessionIdCookie == nil {
		sessionId := uuid.New().String()
		sessionIdCookie = &http.Cookie{
			Name:  "session-id",
			Value: sessionId,
		}
		http.SetCookie(w, sessionIdCookie)
	}

	// submit button pressed
	if req.Method == http.MethodPost {
		fUsername := req.FormValue("email")
		fFirstName := req.FormValue("firstname")
		fLastName := req.FormValue("lastname")

		fieldsHaveData := strings.TrimSpace(fUsername) != "" && strings.TrimSpace(fFirstName) != "" && strings.TrimSpace(fLastName) != ""

		if !fieldsHaveData {
			// write response
			w.Header().Set("content-type", "text/html")
			tpl.ExecuteTemplate(w, "index.gohtml", nil)
			return
		}

		newUser := user{
			Email: fUsername, FirstName: fFirstName,
			LastName: fLastName,
		}
		usersDb[newUser.Email] = newUser
		sessionsDb[sessionIdCookie.Value] = newUser.Email
	}

	// check if user exist
	u := getUser(sessionIdCookie.Value)
	if u != nil {
		http.Redirect(w, req, "/info", http.StatusSeeOther)
		return
	}

	// write response
	w.Header().Set("content-type", "text/html")
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func infoHandler(w http.ResponseWriter, req *http.Request) {
	// get session id
	sessionCookie := getSessionCookie(req)
	if sessionCookie == nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	u := getUser(sessionCookie.Value)
	if u == nil {
		log.Printf("User not exist for session %s", sessionCookie.Value)
		log.Println("Redirect to root")
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	w.Header().Set("content-type", "text/html")
	tpl.ExecuteTemplate(w, "info.gohtml", u)
}

func barHandler(w http.ResponseWriter, req *http.Request) {
	// get session id
	sessionCookie := getSessionCookie(req)
	if sessionCookie == nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	// get user
	u := getUser(sessionCookie.Value)
	if u == nil {
		log.Printf("User not exist for session %s", sessionCookie.Value)
		log.Println("Redirect to root")
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	// write response
	w.Header().Set("content-type", "text/html")
	tpl.ExecuteTemplate(w, "bar.gohtml", u)
}

func logoutHandler(w http.ResponseWriter, req *http.Request) {
	sessionCookie := getSessionCookie(req)
	delete(sessionsDb, sessionCookie.Value)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func getSessionCookie(req *http.Request) *http.Cookie {
	sessionIdCookie, _ := req.Cookie("session-id")
	if sessionIdCookie == nil {
		return nil // no session id
	}
	return sessionIdCookie
}

func getUser(sessionId string) *user {
	var u user
	if uid, ok := sessionsDb[sessionId]; ok {
		u = usersDb[uid]
	} else {
		return nil
	}
	return &u
}
