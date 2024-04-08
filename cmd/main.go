package main

import (
	"log"
	"net/http"

	"github.com/botiash/sherlock/internal/handler"
)

func main() {
	http.HandleFunc("/", handler.HandleWebInterface)
	http.HandleFunc("/download/", handler.HandleFileDownload)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
