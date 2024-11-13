package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"qqweq/siglog/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template
var sessionCookieName = "sessionId"
var sessions = map[string]string {}
var users = map[string]models.User{
    "jamesbond": {
	Username: "jamesbond",
	Password: "hashedpassword",
	Firstname: "James",
	Lastname: "Bond",
	Role: "user",
    },
}

func init() {
    tpl = template.New("")
    tpl, err := tpl.ParseGlob("templates/*.html")
    if err != nil {
	log.Fatalf("Can't parse main page files. %s", err)
    }
    tpl, err = tpl.ParseGlob("templates/login/*.html")
    if err != nil {
	log.Fatalf("Can't parse login page files. %s", err)
    }
    tpl, err = tpl.ParseGlob("templates/register/*.html")
    if err != nil {
	log.Fatalf("Can't parse register page files. %s", err)
    }
}

// for debugging
func db(w http.ResponseWriter, req *http.Request) {
    d := map[string]any {
	"sessions": sessions,
	"users": users,
    }
    json.NewEncoder(w).Encode(d)
}

func main() {
    http.HandleFunc("/", index)
    http.HandleFunc("/login", login)
    http.HandleFunc("/register", register)

    http.HandleFunc("/users", usersHandler)
    http.HandleFunc("/db", db)

    fileServer := http.FileServer(http.Dir("./templates"))
    http.Handle("/public/", http.StripPrefix("/public", fileServer))

    http.Handle("/favicon.ico", http.NotFoundHandler())
    log.Println("Server is listening on port 10443.")
    if err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil); err != nil {
      log.Println(err)
    }
}

func index(w http.ResponseWriter, req *http.Request) {
    if sessionCookie, ok := isAuthenticated(req); ok {
	userId := sessions[sessionCookie.Value]
	tpl.ExecuteTemplate(w, "home.html", users[userId]) 
    } else {
	tpl.ExecuteTemplate(w, "index.html", nil)
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

func login(w http.ResponseWriter, req *http.Request) {
    if _, ok := isAuthenticated(req); ok {
	http.Redirect(w, req, "/", http.StatusSeeOther)
	return
    }

    tpl.ExecuteTemplate(w, "login.html", nil)
}

func register(w http.ResponseWriter, req *http.Request) {
    if _, ok := isAuthenticated(req); ok {
	http.Redirect(w, req, "/", http.StatusSeeOther)
	return
    }

    tpl.ExecuteTemplate(w, "register.html", nil)
}

func usersHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method == http.MethodGet { // check if username exist
	username := req.URL.Query().Get("id")
	if username == "" {
	    log.Println("Username is empty string.")
	    w.WriteHeader(http.StatusForbidden)
	    return
	}

	if _, ok := users[username]; ok {
	    log.Printf("'%s' found. Status 200.", username)
	    w.WriteHeader(http.StatusOK)
	    return
	}
	log.Printf("'%s' not found. Status 404.", username)
	w.WriteHeader(http.StatusNotFound)
    } else if (req.Method == http.MethodPost) { // add new user
	decoder := json.NewDecoder(req.Body)

	var newUser models.User
	if err := decoder.Decode(&newUser); err != nil {
	    log.Printf("Can't decode new user. %s", err)
	    w.WriteHeader(http.StatusForbidden)
	    return
	}

	encryptedPassword := encryptPassword(newUser.Password)
	if encryptedPassword == "" {
	    w.WriteHeader(http.StatusForbidden)
	    return
	}
	newUser.Password = encryptedPassword

	users[newUser.Username] = newUser
	log.Printf("New user '%s' has been added.", newUser.Username)

	sessionId := uuid.NewString() 

	sessions[sessionId] = newUser.Username
	log.Println("New session has been created.")

	http.SetCookie(w, &http.Cookie {
	    Name: sessionCookieName,
	    Value: sessionId,
	})

	log.Println("Redirect to /.")
	http.Redirect(w, req, "/", http.StatusSeeOther)
    }
}

func encryptPassword(password string) (string) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
	log.Printf("Cant encrypt password. %s", err)
	return ""
    }
    return string(bytes)
}

