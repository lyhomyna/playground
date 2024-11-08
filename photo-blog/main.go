package main

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var sessionCookieName = "session"

var tpl *template.Template

func init() {
    tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
    http.HandleFunc("/", index)
    http.HandleFunc("/file", file)

    http.Handle("/favicon.ico", http.NotFoundHandler())
    
    log.Println("Listening on localhost:8080.")
    http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
    c := sessionCookie(w, req)
    http.SetCookie(w, c)
    tpl.ExecuteTemplate(w, "index.gohtml", c)
}

func file(w http.ResponseWriter, req *http.Request) { 
    if req.Method == http.MethodPost {
	clientFile, fh, err := req.FormFile("file")
	if err != nil {
	    log.Fatalf("I haven't got a file. %s", err)
	} 
	log.Println(fmt.Sprintf("I've got a file. '%s'", fh.Filename))
	defer clientFile.Close()

	encryptedFilename := encryptFilename(fh.Filename) 
	log.Printf("Encrypted to: '%s'", encryptedFilename)
	saveFile(clientFile, encryptedFilename)

	sessionCookie := sessionCookie(w, req)
	appendCookieValue(w, sessionCookie, encryptedFilename)
	http.Redirect(w, req, "/", http.StatusSeeOther)
    }
}

func encryptFilename(filename string) string {
    splittedFilename := strings.Split(filename, ".")
    fileName :=  splittedFilename[0]
    fileExt := ""
    if len(splittedFilename) == 2 {
	fileExt = splittedFilename[1]
    }

    // encrypt filename with SHA-1
    hash := sha1.New()
    io.WriteString(hash, fileName)

    return fmt.Sprintf("%x", hash.Sum(nil)) + "." + fileExt
}

func saveFile(clientFile multipart.File, filename string) {
    folder := filepath.Join(".", "public", "user-files")
    err := os.MkdirAll(folder, os.ModePerm)
    if err != nil {
	log.Fatal(err)
    }

    serverFilePath := filepath.Join(".", "public", "user-files", filename)
    serverFile, err := os.Create(serverFilePath)
    if err != nil {
	log.Fatal(err)
    }
    defer serverFile.Close()

    io.Copy(serverFile, clientFile)
}

func sessionCookie(w http.ResponseWriter,req *http.Request) *http.Cookie {
    c, err := req.Cookie(sessionCookieName)
    if err != nil {
	c = &http.Cookie {
	    Name: sessionCookieName,
	    Value: uuid.New().String(),
	}

	http.SetCookie(w, c)
    } 

    return c
}

func appendCookieValue(w http.ResponseWriter, c *http.Cookie, v string) {
    if !strings.Contains(c.Value, v) {
	c.Value += "|" + v 
	http.SetCookie(w, c)
    }
}

