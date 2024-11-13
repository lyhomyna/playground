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
    tpl = template.Must(template.ParseGlob("templates/*.html"))
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
    sessionCookie, err := req.Cookie(sessionCookieName)
    if err != nil {
	log.Println("Session cookies not found.")
	tpl.ExecuteTemplate(w, "sign-in.html", nil)
	return
    }
    userId := sessions[sessionCookie.Value]
    log.Println(users[userId])
    tpl.ExecuteTemplate(w, "index.html", users[userId])
}

func usersHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method == http.MethodGet {
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
    } else if (req.Method == http.MethodPost) {
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

