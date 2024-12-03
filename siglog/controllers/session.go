package controllers

import (
	"log"
	"net/http"
	"qqweq/siglog/model/database"
)

type SessionController struct {
    dao database.SiglogDao
}

var sessionCookieName = "sessionId"
var sessionController *SessionController

func NewSessionController(dao database.SiglogDao) *SessionController {
    if sessionController == nil {
	sessionController = &SessionController { dao }
    }
    return sessionController
}

func (c *SessionController) CreateSession(username string, w http.ResponseWriter) {
    sessionId, err := c.dao.CreateSession(username)
    if err != nil {
	log.Fatal(err)
    }

    log.Printf("New session '%s' has been created.", sessionId)

    http.SetCookie(w, &http.Cookie {
	Name: sessionCookieName,
	Value: sessionId,
    })
}

func (c *SessionController) GetAssosiatedUsername(sessionId string) string {
    username, err := c.dao.UsernameFromSessionId(sessionId)
    if err != nil {
	log.Fatal(err)
    }

    return username 
}

func (c *SessionController) DeleteSession(sessionId string, w http.ResponseWriter) {
    if err := c.dao.DeleteSession(sessionId); err != nil {
	log.Fatal(err)
    }

    http.SetCookie(w, &http.Cookie {
	Name: sessionCookieName,
	MaxAge: -1,
    })

    log.Printf("Session '%s' has been deleted.", sessionId)
}

// TODO: rename it 
func (*SessionController) IsAuthenticated(req *http.Request) (*http.Cookie, bool) {
    sessionCookie, err := req.Cookie(sessionCookieName)
    if err != nil {
	return nil, false
    }
    return sessionCookie, true
}

