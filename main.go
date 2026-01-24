package main

import (
	"fmt"
	"net/http"

	"github.com/foyko/fileconverter/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/upload", handlers.UploadHandler).Methods("POST")
	r.HandleFunc("/upload-form", handlers.UploadFormHandler).Methods("GET")
	r.HandleFunc("/files", handlers.ListFilesHandler).Methods("GET")
	r.HandleFunc("/download/{filename}", handlers.DownloadFileHandler).Methods("GET")
	r.HandleFunc("/delete/{filename}", handlers.DeleteFileHandler).Methods("GET")
	r.HandleFunc("/convert/{filename}", handlers.ConvertFileHandler).Methods("GET")
	r.HandleFunc("/view/{filename}", handlers.ViewFileHandler).Methods("GET")
	r.HandleFunc("/render/{filename}", handlers.RenderFileHandler).Methods("GET")

	port := ":80"
	fmt.Printf("Server starting on port %s\n", port)

	http.ListenAndServe(port, r)
}
