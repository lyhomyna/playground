package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"qqweq/siglog/controllers"
	"qqweq/siglog/models"

	"github.com/google/uuid"
)

var tpl *template.Template
var userController *controllers.UserController
var sessionCookieName = "sessionId"
var sessions = map[string]string {}

func init() {
    userController = controllers.NewUserController()

    tpl = template.New("")
    tpl, err := tpl.ParseGlob("resources/*.html")
    if err != nil {
	log.Fatalf("Can't parse main page files. %s", err)
    }
    tpl, err = tpl.ParseGlob("resources/login/*.html")
    if err != nil {
	log.Fatalf("Can't parse login page files. %s", err)
    }
    tpl, err = tpl.ParseGlob("resources/register/*.html")
    if err != nil {
	log.Fatalf("Can't parse register page files. %s", err)
    }
}

func main() {
    http.HandleFunc("/", index)
    http.HandleFunc("/login", login)
    http.HandleFunc("/register", register)
    http.HandleFunc("/logout", logout)

    http.HandleFunc("/users", usersHandler)

    fileServer := http.FileServer(http.Dir("./resources"))
    http.Handle("/public/", http.StripPrefix("/public", fileServer))

    http.Handle("/favicon.ico", http.NotFoundHandler())
    log.Println("Server is listening on port 10443.")
    if err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil); err != nil {
      log.Println(err)
    }
}

func index(w http.ResponseWriter, req *http.Request) {
    if sessionCookie, ok := isAuthenticated(req); ok {
	username := sessions[sessionCookie.Value]
	tpl.ExecuteTemplate(w, "home.html", userController.GetUserByUsername(username)) 
    } else {
	tpl.ExecuteTemplate(w, "index.html", nil)
    }
}

func login(w http.ResponseWriter, req *http.Request) {
    if req.Method == http.MethodGet {
	if _, ok := isAuthenticated(req); ok {
	    http.Redirect(w, req, "/", http.StatusSeeOther)
	    return
	}
	tpl.ExecuteTemplate(w, "login.html", nil)
    } else if req.Method == http.MethodPost {

	var usernamePassword models.UserLog
	if err := decodeFromTo(req.Body, &usernamePassword); err != nil {
	    log.Printf("Login. %s", err)
	    w.WriteHeader(http.StatusForbidden)
	    return
	}

	// user validation
	if user := userController.GetUserByUsername(usernamePassword.Username); user != nil {
	    if err := userController.ComparePasswords(user, usernamePassword.Password); err != nil {
		log.Printf("User '%s'. %s", user.Username, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	    }
	} else {
	    log.Println("Incorrect username.")
	    w.WriteHeader(http.StatusUnauthorized)
	}
	
	log.Printf("Login: user '%s' found. Create session.", usernamePassword.Username)

	// create session
	createSession(usernamePassword.Username, w)

	log.Printf("Login: user '%s' logined.", usernamePassword.Username)
	w.WriteHeader(http.StatusOK)
    }
}

func register(w http.ResponseWriter, req *http.Request) {
    if req.Method == http.MethodGet {
	if _, ok := isAuthenticated(req); ok {
	    http.Redirect(w, req, "/", http.StatusSeeOther)
	    return
	}
	tpl.ExecuteTemplate(w, "register.html", nil)
    } else if req.Method == http.MethodPost {
	// decode
	var newUser models.User
	if err := decodeFromTo(req.Body, &newUser); err != nil {
	    log.Printf("Register. %e", err)
	    w.WriteHeader(http.StatusForbidden)
	    return
	}

	// add
	if err := userController.AddUser(&newUser); err != nil {
	    log.Println(err)
	    return
	}
	log.Printf("New user '%s' has been added.", newUser.Username)

	createSession(newUser.Username, w)

	w.WriteHeader(http.StatusCreated)
    }
}

func logout(w http.ResponseWriter, req *http.Request) {
    if _, ok := isAuthenticated(req); ok {
	http.SetCookie(w, &http.Cookie {
	    Name: sessionCookieName,
	    MaxAge: -1,
	})	

	delete(sessions, sessionCookieName)

	http.Redirect(w, req, "/", http.StatusSeeOther)
    }
}

func usersHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method == http.MethodGet { // check if username exists
	username := req.URL.Query().Get("id")
	if username == "" {
	    log.Println("Username is empty string.")
	    w.WriteHeader(http.StatusForbidden)
	    return
	}

	if user := userController.GetUserByUsername(username); user != nil {
	    log.Printf("'%s' found. Status 200.", username)
	    w.WriteHeader(http.StatusOK)
	    return
	}
	log.Printf("'%s' not found. Status 404.", username)
	w.WriteHeader(http.StatusNotFound)
    } 
}

func isAuthenticated(req *http.Request) (*http.Cookie, bool) {
    sessionCookie, err := req.Cookie(sessionCookieName)
    if err != nil {
	log.Println("Session cookies not found.")
	return nil, false
    }
    return sessionCookie, true
}

func decodeFromTo(rc io.ReadCloser, target any) error {
    decoder := json.NewDecoder(rc)
    if err := decoder.Decode(target); err != nil {
	return errors.New(fmt.Sprintf("Decode failure. %s", err))
    }
    return nil
}

func createSession(username string, w http.ResponseWriter) {
    sessionId := uuid.NewString() 

    sessions[sessionId] = username 
    log.Println("New session has been created.")

    http.SetCookie(w, &http.Cookie {
	Name: sessionCookieName,
	Value: sessionId,
    })
}
