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
	"qqweq/siglog/model/database"
	"qqweq/siglog/model/models"
)

var tpl *template.Template
var userController *controllers.UserController
var sessionController *controllers.SessionController

func init() {
    db := database.NewDatabase()
    if db == nil {
	log.Panic("Can't connect to the database.")
    }

    userController = controllers.NewUserController(db)
    sessionController = controllers.NewSessionController(db)

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
    // return a HTML page 
    http.HandleFunc("/login", login)
    http.HandleFunc("/register", register)
    http.HandleFunc("/logout", logout)

    // doesn't return a HTML page
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
    if sessionCookie, ok := sessionController.IsAuthenticated(req); ok {
	username := sessionController.GetAssosiatedUsername(sessionCookie.Value)
	tpl.ExecuteTemplate(w, "home.html", userController.GetUserByUsername(username)) 
    } else {
	tpl.ExecuteTemplate(w, "index.html", nil)
    }
}

func login(w http.ResponseWriter, req *http.Request) {
    if req.Method == http.MethodGet {
	if _, ok := sessionController.IsAuthenticated(req); ok {
	    http.Redirect(w, req, "/", http.StatusSeeOther) // if you authenticated - login page isn't for you
	    return
	}
	tpl.ExecuteTemplate(w, "login.html", nil)
    } else if req.Method == http.MethodPost {

	var usernamePassword models.UserLog
	if err := decodeFromTo(req.Body, &usernamePassword); err != nil {
	    log.Printf("Login. %s", err) // bro, your data is trash, even decoder can't decode it
	    w.WriteHeader(http.StatusForbidden)
	    return
	}

	// user validation (looks like spagetti)
	if user := userController.GetUserByUsername(usernamePassword.Username); user != nil {
	    incorrectPasswordErr := userController.ComparePasswords(user, usernamePassword.Password)
	    if incorrectPasswordErr != nil {
		log.Printf("User '%s'. %s", user.Username, incorrectPasswordErr) 
		w.WriteHeader(http.StatusUnauthorized)
		return
	    }
	} else {
	    log.Println("Incorrect username.") // you're fal'shivka
	    w.WriteHeader(http.StatusUnauthorized)
	}

	sessionController.CreateSession(usernamePassword.Username, w)

	log.Printf("User '%s' logined.", usernamePassword.Username)
	w.WriteHeader(http.StatusOK)
    }
}

func register(w http.ResponseWriter, req *http.Request) {
    if req.Method == http.MethodGet {
	if _, ok := sessionController.IsAuthenticated(req); ok {
	    http.Redirect(w, req, "/", http.StatusSeeOther)
	    return
	}
	tpl.ExecuteTemplate(w, "register.html", nil)
    } else if req.Method == http.MethodPost {

	var newUser models.User
	if err := decodeFromTo(req.Body, &newUser); err != nil {
	    log.Printf("Register. %e", err) // data is trash again. decoder is shocked
	    w.WriteHeader(http.StatusForbidden)
	    return
	}

	// add decoded user to "database" (haha)
	if err := userController.AddUser(&newUser); err != nil {
	    log.Println(err)
	    return
	}
	log.Printf("New user '%s' has been added.", newUser.Username)

	sessionController.CreateSession(newUser.Username, w)

	w.WriteHeader(http.StatusCreated)
    }
}

func logout(w http.ResponseWriter, req *http.Request) {
    if sessionCookie, ok := sessionController.IsAuthenticated(req); ok {
	// it's stupid, but I want log
	username := sessionController.GetAssosiatedUsername(sessionCookie.Value)

	sessionController.DeleteSession(w)
	
	log.Printf("User '%s' logged out.", username)
	// from this moment, I don't know who you are
	http.Redirect(w, req, "/", http.StatusSeeOther)
    }
}

func usersHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method == http.MethodGet { // username existence fact or with another words - "USERNAME ALREADY TAKEN" 
	username := req.URL.Query().Get("id")
	if username == "" {
	    log.Println("Can't get username from URL query.")
	    w.WriteHeader(http.StatusForbidden)
	    return
	}

	if user := userController.GetUserByUsername(username); user != nil {
	    // username not found
	    w.WriteHeader(http.StatusOK)
	    return
	}
	// usename found
	w.WriteHeader(http.StatusNotFound)
    } 
}

func decodeFromTo(rc io.ReadCloser, target any) error {
    decoder := json.NewDecoder(rc)
    if err := decoder.Decode(target); err != nil {
	return errors.New(fmt.Sprintf("Decode failure. %s", err))
    }
    return nil
}

