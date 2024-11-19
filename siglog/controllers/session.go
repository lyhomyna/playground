package controllers

import (
    "net/http"
    "log"
    "github.com/google/uuid"
)

type SessionController struct {}

var sessionCookieName = "sessionId"
var sessions = map[string]string {} // sessionId : username
var sessionController *SessionController

func NewSessionController() *SessionController {
    if sessionController == nil {
	sessionController = &SessionController {}
    }
    return sessionController
}

func (*SessionController) CreateSession(username string, w http.ResponseWriter) {
    sessionId := uuid.NewString() 

    sessions[sessionId] = username 
    log.Println("New session has been created.")

    http.SetCookie(w, &http.Cookie {
	Name: sessionCookieName,
	Value: sessionId,
    })
}

// TODO: bro, rename this
func (*SessionController) GetAssosiatedUsername(sessionId string) string {
    return sessions[sessionId]
}

func (*SessionController) DeleteSession(w http.ResponseWriter) {
    http.SetCookie(w, &http.Cookie {
	Name: sessionCookieName,
	MaxAge: -1,
    })	

    delete(sessions, sessionCookieName)
}

// TODO: rename it
func (*SessionController) IsAuthenticated(req *http.Request) (*http.Cookie, bool) {
    sessionCookie, err := req.Cookie(sessionCookieName)
    if err != nil {
	return nil, false
    }
    return sessionCookie, true
}

