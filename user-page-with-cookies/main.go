package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"
    "text/template"
    "encoding/json"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "cookiesplay/models"
)
var sessionsDb = map[string]string{}
var usersDb = map[string]models.User{}
var tpl *template.Template
var sessionCookieName = "session-id"
var sessionLenInSec = 30

func init() {
    tpl = template.Must(template.ParseFiles("templates/signin.gohtml", "templates/bar.gohtml", "templates/info.gohtml", "templates/404.gohtml"))
}

func main() {
    http.HandleFunc("/", rootHandler)
    http.HandleFunc("/bar", barHandler)
    http.HandleFunc("/sign-in", signinHandler)
    http.HandleFunc("/logout", logoutHandler)

    http.HandleFunc("/log", logHandler)
    http.HandleFunc("/reg", regHandler)

    
    http.HandleFunc("/db", showAllData)
    http.Handle("/favicon.ico", http.NotFoundHandler())


    log.Println("server is running.")
    http.ListenAndServe(":8080", nil)
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
    if isAuthenicated(req) {
	sessionIdCookie, _ := req.Cookie(sessionCookieName)

	usr := getUser(sessionIdCookie.Value)
	if usr == nil {
		log.Println("User is nil.")
		tpl.ExecuteTemplate(w, "404", nil)
		return
	}
	handleUser(w, req)
    } else { log.Println("SessionIdCookie is nil. Redirect to /sign-in.")
	http.Redirect(w, req, "/sign-in", http.StatusSeeOther)
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
    w.Header().Set("content-type", "text/html")
    tpl.ExecuteTemplate(w, "signin.gohtml", nil)
}

func logHandler(w http.ResponseWriter, req *http.Request) {
    log.Println("User login.")
    if req.Method != http.MethodPost {
	http.Redirect(w, req, "/", http.StatusSeeOther)
    }

    fEmail := strings.TrimSpace(req.FormValue("email"))
    fPassword := strings.TrimSpace(req.FormValue("password"))

    // ser existance and password correctneses
    if user, ok := usersDb[fEmail]; ok {
	passwordBytes := []byte(fPassword)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), passwordBytes)
	if err != nil {
	    log.Println(fmt.Sprintf("Password for %s is incorrect.", fEmail))
	    http.Redirect(w, req, "/sign-in", http.StatusSeeOther)
	    return
	}
    } else {
	// User not exist
	log.Println(fmt.Sprintf("User %s isn't exist.", fEmail))
	http.Redirect(w, req, "/sign-in", http.StatusSeeOther)
	return
    }

    sessionCookie := createSessionCookie() 
    sessionsDb[sessionCookie.Value] = fEmail

    http.SetCookie(w, sessionCookie)
    http.Redirect(w, req, "/", http.StatusSeeOther)
}

func regHandler(w http.ResponseWriter, req *http.Request) {
    log.Println("User register.")
    if req.Method != http.MethodPost {
	http.Redirect(w, req, "/", http.StatusSeeOther)
    }

    fEmail := strings.TrimSpace(req.FormValue("email"))
    fPassword := strings.TrimSpace(req.FormValue("password"))

    fieldsHaveData := fEmail != "" && fPassword != ""
    if !fieldsHaveData {
	log.Println("Empty fields.")
	w.Header().Set("content-type", "text/html")
	tpl.ExecuteTemplate(w, "signin.gohtml", nil)
	return
    }

    // is User exist in DB?
    if _, ok := usersDb[fEmail]; ok {
	log.Println(fmt.Sprintf("User %s already exist.", fEmail))
	http.Redirect(w, req, "/sign-in", http.StatusSeeOther)
	return
    }

    pCost, err := bcrypt.GenerateFromPassword([]byte(fPassword), bcrypt.DefaultCost)
    if err != nil {
	log.Println(err)
	http.Redirect(w, req, "/sign-in", http.StatusSeeOther)
	return
    }

    sessionIdCookie := createSessionCookie() 
    newUser := models.User{
	Email:    fEmail,
	Password: string(pCost),
    }
    usersDb[newUser.Email] = newUser
    sessionsDb[sessionIdCookie.Value] = newUser.Email

    log.Println("User has been signed in. Redirect to the root.")
    http.SetCookie(w, sessionIdCookie)
    http.Redirect(w, req, "/", http.StatusSeeOther)
}

func barHandler(w http.ResponseWriter, req *http.Request) {
    // get session id
    sessionCookie := getSessionCookie(req)
    if sessionCookie == nil {
	http.Redirect(w, req, "/", http.StatusSeeOther)
	return
    }

    // get User
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

func getUser(sessionId string) *models.User {
    var u models.User
    log.Println(fmt.Sprintf("Trying to get user by sessionId: %s", sessionId))
    log.Println(sessionsDb[sessionId])
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

func createSessionCookie() *http.Cookie {
    return &http.Cookie{
	Name:  sessionCookieName,
	Value: uuid.NewString(),
	MaxAge: sessionLenInSec,
    }
}


func showAllData(w http.ResponseWriter, req *http.Request) {
    resObj := map[string]any {
	"sessions": sessionsDb,
	"users": usersDb,
    }

    res, err := json.Marshal(resObj)
    if err != nil {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("{'message': 'something went wrong while marshalling all data'}"))
	return
    }

    w.WriteHeader(http.StatusOK)
    w.Write(res)
}

