package handlers

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

// ViewFileHandler renders the file in the browser within an iframe
func ViewFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := filepath.Base(vars["filename"])

	filePath := filepath.Join(UploadPath, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, "Error reading file info", http.StatusInternalServerError)
		return
	}

	// Get file extension to determine content type
	ext := strings.ToLower(filepath.Ext(filename))

	// Determine if file is viewable in iframe
	var isViewable bool = true
	var viewType string

	switch ext {
	case ".pdf":
		viewType = "pdf"
	case ".txt", ".log", ".md", ".json", ".xml", ".csv":
		viewType = "text"
	case ".html", ".htm":
		viewType = "html"
	case ".jpg", ".jpeg", ".png", ".gif", ".svg":
		viewType = "image"
	case ".mp4", ".webm":
		viewType = "video"
	case ".mp3", ".wav":
		viewType = "audio"
	default:
		isViewable = false
	}

	if !isViewable {
		showFilePreview(w, filename, filePath)
		return
	}

	// Render the viewer page with iframe
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>View File - {{.Name}}</title>
		<style>
			* {
				margin: 0;
				padding: 0;
				box-sizing: border-box;
			}
			body {
				font-family: Arial, sans-serif;
				height: 100vh;
				display: flex;
				flex-direction: column;
			}
			.header {
				background: #333;
				color: white;
				padding: 15px 20px;
				display: flex;
				justify-content: space-between;
				align-items: center;
				box-shadow: 0 2px 4px rgba(0,0,0,0.1);
			}
			.header h2 {
				font-size: 18px;
				font-weight: normal;
			}
			.file-name {
				font-weight: bold;
			}
			.header-info {
				display: flex;
				gap: 20px;
				font-size: 14px;
				color: #ccc;
			}
			.actions {
				display: flex;
				gap: 10px;
			}
			.btn {
				padding: 8px 16px;
				border-radius: 4px;
				text-decoration: none;
				font-size: 14px;
				transition: background 0.3s;
			}
			.btn-primary {
				background: #007bff;
				color: white;
			}
			.btn-primary:hover {
				background: #0056b3;
			}
			.btn-secondary {
				background: #6c757d;
				color: white;
			}
			.btn-secondary:hover {
				background: #545b62;
			}
			.viewer-container {
				flex: 1;
				background: #f5f5f5;
				padding: 20px;
				overflow: auto;
			}
			.viewer-frame {
				width: 100%;
				height: 100%;
				border: 1px solid #ddd;
				background: white;
				border-radius: 4px;
				box-shadow: 0 2px 8px rgba(0,0,0,0.1);
			}
			iframe {
				width: 100%;
				height: 100%;
				border: none;
			}
			.image-container {
				display: flex;
				justify-content: center;
				align-items: center;
				height: 100%;
			}
			.image-container img {
				max-width: 100%;
				max-height: 100%;
				object-fit: contain;
			}
			video, audio {
				width: 100%;
				max-width: 800px;
				display: block;
				margin: 0 auto;
			}
		</style>
	</head>
	<body>
		<div class="header">
			<div>
				<h2>
					<span class="file-name">{{.Name}}</span>
				</h2>
				<div class="header-info">
					<span>Size: {{.SizeFormatted}}</span>
					<span>Modified: {{.ModTime}}</span>
				</div>
			</div>
			<div class="actions">
				<a href="/download/{{.Name}}" class="btn btn-primary">Download</a>
				<a href="/files" class="btn btn-secondary">Back to Files</a>
			</div>
		</div>

		<div class="viewer-container">
			<div class="viewer-frame">
				{{if eq .ViewType "image"}}
					<div class="image-container">
						<img src="/render/{{.Name}}" alt="{{.Name}}">
					</div>
				{{else if eq .ViewType "video"}}
					<video controls>
						<source src="/render/{{.Name}}" type="video/{{.Ext}}">
						Your browser does not support the video tag.
					</video>
				{{else if eq .ViewType "audio"}}
					<audio controls>
						<source src="/render/{{.Name}}" type="audio/{{.Ext}}">
						Your browser does not support the audio tag.
					</audio>
				{{else}}
					<iframe src="/render/{{.Name}}"></iframe>
				{{end}}
			</div>
		</div>
	</body>
	</html>
	`

	data := struct {
		Name          string
		SizeFormatted string
		ModTime       string
		ViewType      string
		Ext           string
	}{
		Name:          filename,
		SizeFormatted: FormatFileSize(fileInfo.Size()),
		ModTime:       fileInfo.ModTime().Format("2006-01-02 15:04:05"),
		ViewType:      viewType,
		Ext:           strings.TrimPrefix(ext, "."),
	}

	t, err := template.New("viewer").Parse(tmpl)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

// showFilePreview displays a page with file information for non-viewable files
func showFilePreview(w http.ResponseWriter, filename, filePath string) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, "Error reading file info", http.StatusInternalServerError)
		return
	}

	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>File Preview - {{.Name}}</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				max-width: 800px;
				margin: 50px auto;
				padding: 20px;
			}
			.file-info {
				background: #f5f5f5;
				padding: 20px;
				border-radius: 8px;
				margin-bottom: 20px;
			}
			.file-info h2 {
				margin-top: 0;
			}
			.info-row {
				margin: 10px 0;
			}
			.label {
				font-weight: bold;
				display: inline-block;
				width: 150px;
			}
			.actions {
				margin-top: 30px;
			}
			.btn {
				padding: 10px 20px;
				margin-right: 10px;
				border-radius: 5px;
				text-decoration: none;
				display: inline-block;
			}
			.btn-primary {
				background: #007bff;
				color: white;
			}
			.btn-secondary {
				background: #6c757d;
				color: white;
			}
			.warning {
				background: #fff3cd;
				border: 1px solid #ffc107;
				padding: 15px;
				border-radius: 5px;
				margin-top: 20px;
			}
		</style>
		<meta charset="UTF-8">
	</head>
	<body>
		<div class="file-info">
			<h2>File Information</h2>
			<div class="info-row">
				<span class="label">Filename:</span>
				<span>{{.Name}}</span>
			</div>
			<div class="info-row">
				<span class="label">Size:</span>
				<span>{{.SizeFormatted}}</span>
			</div>
			<div class="info-row">
				<span class="label">Modified:</span>
				<span>{{.ModTime}}</span>
			</div>
		</div>

		<div class="warning">
			<strong>⚠️ Preview not available</strong><br>
			This file type cannot be previewed in the browser. Please download it to view.
		</div>

		<div class="actions">
			<a href="/download/{{.Name}}" class="btn btn-primary">Download File</a>
			<a href="/files" class="btn btn-secondary">Back to Files</a>
		</div>
	</body>
	</html>
	`

	data := struct {
		Name          string
		SizeFormatted string
		ModTime       string
	}{
		Name:          filename,
		SizeFormatted: FormatFileSize(fileInfo.Size()),
		ModTime:       fileInfo.ModTime().Format("2006-01-02 15:04:05"),
	}

	t, err := template.New("preview").Parse(tmpl)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

// RenderFileHandler serves the actual file content for rendering in iframe
func RenderFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := filepath.Base(vars["filename"])

	filePath := filepath.Join(UploadPath, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Get file extension to determine content type
	ext := strings.ToLower(filepath.Ext(filename))

	// Set appropriate content type
	var contentType string
	switch ext {
	case ".pdf":
		contentType = "application/pdf"
	case ".txt", ".log", ".md":
		contentType = "text/plain; charset=utf-8"
	case ".html", ".htm":
		contentType = "text/html; charset=utf-8"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".svg":
		contentType = "image/svg+xml"
	case ".json":
		contentType = "application/json"
	case ".xml":
		contentType = "application/xml"
	case ".csv":
		contentType = "text/csv"
	case ".mp4":
		contentType = "video/mp4"
	case ".webm":
		contentType = "video/webm"
	case ".mp3":
		contentType = "audio/mpeg"
	case ".wav":
		contentType = "audio/wav"
	default:
		contentType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "inline; filename="+filename)

	http.ServeFile(w, r, filePath)
}
