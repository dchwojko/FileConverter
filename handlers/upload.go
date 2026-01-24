package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Limit the size of the incoming request body
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)

	// Parse the multipart form
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		http.Error(w, "File too large or invalid form data", http.StatusBadRequest)
		return
	}

	// Retrieve the file from form data
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the uploads directory if it doesn't exist
	if err := os.MkdirAll(UploadPath, os.ModePerm); err != nil {
		http.Error(w, "Error creating upload directory", http.StatusInternalServerError)
		return
	}

	// Create destination file
	filename := filepath.Base(fileHeader.Filename)
	dst, err := os.Create(filepath.Join(UploadPath, filename))
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/files", http.StatusSeeOther)
}
