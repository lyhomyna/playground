package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"qqweq/siglog/models"
)

var tpl *template.Template

func init() {
  tpl = template.New("")
  tpl, err := tpl.ParseGlob("templates/*.html")
  if err != nil {
    log.Fatal(err)
  }
  tpl, err = tpl.ParseGlob("templates/*.gohtml")
  if err != nil {
    log.Fatal(err)
  }
}

func main() {
  http.HandleFunc("/", index)
  http.HandleFunc("/users", usersHandler)

  fileServer := http.FileServer(http.Dir("./templates"))
  http.Handle("/public/", http.StripPrefix("/public", fileServer))

  log.Println("Server is listening on port 10443.")
  if err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil); err != nil {
    log.Println(err)
  }
}

func index(w http.ResponseWriter, req *http.Request) {
  tpl.ExecuteTemplate(w, "sign-in.html", nil)
}

var users = map[string]models.User{
    "jamesbond": {
	Username: "jamesbond",
	Password: "hashedpassword",
	Firstname: "James",
	Lastname: "Bond",
	Role: "user",
    },
}

func usersHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method == http.MethodPost {
	decoder := json.NewDecoder(req.Body)

	var username string
	if err := decoder.Decode(&username); err != nil {
	    log.Fatalf("Username decode failure. %s", err)
	}

	if _, ok := users[username]; ok {
	    log.Printf("'%s' found. Status 200.", username)
	    w.WriteHeader(http.StatusOK)
	    return
	}
	log.Printf("'%s' not found. Status 404.", username)
	w.WriteHeader(http.StatusNotFound)
    }
}


