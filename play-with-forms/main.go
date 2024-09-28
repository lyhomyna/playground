package main

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.gohtml"))
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func handleRoot(w http.ResponseWriter, req *http.Request) {
	form := `
    <form method="POST" enctype="multipart/form-data">
      <input type="file" name="file" />
      <input type="submit" />
    </form>
  `
	fileContent := ""

	responseData := struct {
		Form        string
		FileContent string
	}{
		Form:        form,
		FileContent: fileContent,
	}

	if req.Method == http.MethodPost {
		f, fh, err := req.FormFile("file")
		if err != nil {
			log.Println("Line 41: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		fileContent := make([]byte, fh.Size)
		if _, err := f.Read(fileContent); err != nil {
			log.Println("Line 49: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := writeFileToInternal(fh, fileContent); err != nil {
			log.Println("Line 58: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bodyContent := make([]byte, req.ContentLength)
		if _, err := req.Body.Read(bodyContent); err != nil && err != io.EOF {
			log.Println(err)
			return
		}

		responseData.FileContent = string(fileContent)
	}

	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tpl.Execute(w, responseData)
}

func writeFileToInternal(fh *multipart.FileHeader, fileContent []byte) error {
	path := filepath.Join("./user-data/", fh.Filename)
	file, err := os.Create(path)
	if err != nil {
		log.Println("Line 67: ", err)
		return err
	}
	defer file.Close()
	if _, err := file.Write(fileContent); err != nil {
		log.Println("Line 72: ", err)
		return err
	}
	return nil
}
