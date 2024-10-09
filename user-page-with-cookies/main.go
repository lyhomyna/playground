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
	tpl = template.Must(template.ParseFiles("templates/index.gohtml", "templates/bar.gohtml", "templates/info.gohtml", "templates/404.gohtml"))
}

// if isAuthenticated - go to info page, else - go to sign-in page

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/bar", barHandler)
	http.HandleFunc("/sign-in", signinHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	log.Println("server is running.")
	http.ListenAndServe(":8080", nil)
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	sessionIdCookie := getSessionCookie(req)
	if sessionIdCookie == nil {
		log.Println("SessionIdCookie is nil. Redirect to /sign-in.")
		http.Redirect(w, req, "/sign-in", http.StatusSeeOther)
	} else {
		log.Println("SessionIdCookie is present.")
		u := getUser(sessionIdCookie.Value)
		if u == nil {
			log.Println("User is nil.")
			tpl.ExecuteTemplate(w, "404", nil)
			return
		}
		handleUser(w, req)
	}
}

func handleUser(w http.ResponseWriter, req *http.Request) {
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

	log.Println("User is present. Info page.")

	w.Header().Set("content-type", "text/html")
	tpl.ExecuteTemplate(w, "info.gohtml", u)
}

func signinHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		log.Println("Register new user.")
		sessionIdCookie := &http.Cookie{
			Name:  "session-id",
			Value: uuid.NewString(),
		}

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

		http.SetCookie(w, sessionIdCookie)
		http.Redirect(w, req, "/", http.StatusSeeOther)
	} else {
		// write response
		w.Header().Set("content-type", "text/html")
		tpl.ExecuteTemplate(w, "index.gohtml", nil)
	}
}

func barHandler(w http.ResponseWriter, req *http.Request) {
	// get session id
	sessionCookie := getSessionCookie(req)
	if sessionCookie == nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
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
	http.SetCookie(w, &http.Cookie{
		Name:   "session-id",
		Value:  "",
		MaxAge: -1,
	})
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

func isAuthenicated(req *http.Request) bool {
	_, err := req.Cookie("session-id")
	if err != nil {
		return false
	}
	return true
}
