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
    http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./public/user-files/"))))

    http.Handle("/favicon.ico", http.NotFoundHandler())
    
    log.Println("Listening on localhost:8080.")
    http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
    c := sessionCookie(w, req)
    filenames := getFilenames(c.Value)
    tpl.ExecuteTemplate(w, "index.gohtml", filenames)
}

func getFilenames(cookieValue string) []string {
    filenames := strings.Split(cookieValue, "|")
    if len(filenames) > 1 {
	return filenames[1:]
    }
    return nil
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
	log.Printf("Filename encrypted to: '%s'", encryptedFilename)

	saveFile(clientFile, encryptedFilename)

	sessionCookie := sessionCookie(w, req)
	updatedCookie := appendToCookieValue(sessionCookie, encryptedFilename)
	log.Printf("Updated session cookie: %+v", updatedCookie)

	http.SetCookie(w, updatedCookie)
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
	log.Fatalf("Can't create folder. %s", err)
    }

    serverFilePath := filepath.Join(".", "public", "user-files", filename)
    serverFile, err := os.Create(serverFilePath)
    if err != nil {
	log.Fatalf("Can't create file. %s", err)
    }
    defer serverFile.Close()

    
    if _, err = io.Copy(serverFile, clientFile); err != nil {
	log.Printf("Can't copy file. %s", err)
    }
    log.Println("File has been saved successfully.")
}

func sessionCookie(w http.ResponseWriter,req *http.Request) *http.Cookie {
    c, err := req.Cookie(sessionCookieName)
    if err != nil {
	log.Println("Session cookie is missing. Create.")
	c = &http.Cookie {
	    Name: sessionCookieName,
	    Value: uuid.New().String(),
	}

	http.SetCookie(w, c)
    } 
    return c
}

func appendToCookieValue(c *http.Cookie, cookiePart string) *http.Cookie {
    if !strings.Contains(c.Value, cookiePart) {
	log.Printf("Cookie part '%s' is missing. Append.", cookiePart)
	c.Value += "|" + cookiePart 
    }
    return c
}

