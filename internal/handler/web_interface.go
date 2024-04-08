package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func HandleFileDownload(w http.ResponseWriter, r *http.Request) {
    fileName := filepath.Base(r.URL.Path)

    filePath := "./" + fileName // Adjust the file path
	log.Println("adddsef ", filePath)
    // Check if the file exists
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }

    // Set Content-Disposition header to prompt download
    w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
    http.ServeFile(w, r, filePath)
}


func HandleWebInterface(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// При GET-запросе выводим форму для ввода запроса
		tmpl := template.Must(template.ParseFiles("web/index.html"))
		tmpl.Execute(w, nil)
	}
	if r.Method == "POST" {
		if r.Method == "POST" {
			query := r.FormValue("query")
			searchType := r.FormValue("searchType")
			fileName := r.FormValue("fileName") // Get the filename from the form

			// Provide a default filename if not specified
			if fileName == "" {
				fileName = "default_result.txt"
			}

			resultChan := make(chan string)
			go Run(query, searchType, fileName, resultChan)
			<-resultChan

			// Send a link to download the file
			fmt.Fprintf(w, "<a href='/download/%s'>Download Results</a>", fileName)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
