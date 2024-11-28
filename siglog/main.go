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

	_ "github.com/joho/godotenv/autoload"
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
    // returns a HTML page 
    http.HandleFunc("/login", login)
    http.HandleFunc("/register", register)
    http.HandleFunc("/logout", logout)

    // doesn't return a HTML page
    http.HandleFunc("/users", usersDataHandler)
    http.HandleFunc("/delete", deleteAcc)

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
	
	user := userController.GetUserByUsername(username)
	if user == nil {
	    // delete session and render "index.html"

	    tpl.ExecuteTemplate(w, "index.html", nil)
	    return
	}

	tpl.ExecuteTemplate(w, "home.html", user) 
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
	    log.Println("User not found.") // you're fal'shivka
	    w.WriteHeader(http.StatusUnauthorized)
	    return
	}

	sessionController.CreateSession(usernamePassword.Username, w)

	log.Printf("User '%s' logined.", usernamePassword.Username)
	w.WriteHeader(http.StatusOK)
    }
}

func register(w http.ResponseWriter, req *http.Request) {
    if req.Method == http.MethodGet {
	if _, ok := sessionController.IsAuthenticated(req); ok {
	    http.Redirect(w, req, "/", http.StatusSeeOther) // if you logined, register page isn't for you
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

	// add new user to the Database (yo, with upper letter)
	newUserId, err := userController.AddUser(&newUser)
	if  err != nil {
	    log.Println(err)
	    return
	}
	log.Printf("New user with id '%s' has been added.", newUserId)

	sessionController.CreateSession(newUser.Username, w)

	w.WriteHeader(http.StatusCreated)
    }
}

func logout(w http.ResponseWriter, req *http.Request) {
    if sessionCookie, ok := sessionController.IsAuthenticated(req); ok {
	// it's useless, but I want pretty log
	username := sessionController.GetAssosiatedUsername(sessionCookie.Value)

	sessionController.DeleteSession(sessionCookie.Value, w)
	
	log.Printf("User '%s' logged out.", username)

	http.Redirect(w, req, "/", http.StatusSeeOther)
    }
}

func usersDataHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method == http.MethodGet { // checking existence of user with some username
	username := req.URL.Query().Get("id")
	if username == "" {
	    log.Println("Can't get username from URL query.")
	    w.WriteHeader(http.StatusForbidden)
	    return
	}

	if user := userController.GetUserByUsername(username); user != nil {
	    // username is taken
	    w.WriteHeader(http.StatusOK)
	    return
	}
	// username isn't taken 
	w.WriteHeader(http.StatusNotFound)
    } 
}

func deleteAcc(w http.ResponseWriter, req *http.Request) {
    if sessionCookie, ok := sessionController.IsAuthenticated(req); ok {
	username := sessionController.GetAssosiatedUsername(sessionCookie.Value)
	userController.DeleteUser(username)

	log.Printf("User %s has been deleted.", username)

	sessionController.DeleteSession(sessionCookie.Value, w)

	http.Redirect(w, req, "/", http.StatusSeeOther)
    }
}

func decodeFromTo(rc io.ReadCloser, target any) error {
    decoder := json.NewDecoder(rc)
    if err := decoder.Decode(target); err != nil {
	return errors.New(fmt.Sprintf("Decode failure. %s", err))
    }
    return nil
}

