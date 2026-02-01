package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure uploads directory exists
	if err := os.MkdirAll(UploadPath, os.ModePerm); err != nil {
		http.Error(w, "Error accessing upload directory", http.StatusInternalServerError)
		return
	}

	// Read directory contents
	files, err := os.ReadDir(UploadPath)
	if err != nil {
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}

	// Gather file information
	var fileInfos []FileInfo
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		fileInfos = append(fileInfos, FileInfo{
			Name:          file.Name(),
			Size:          info.Size(),
			SizeFormatted: formatFileSize(info.Size()),
			ModTime:       info.ModTime().Format("2006-01-02 15:04:05"),
			DownloadURL:   "/download/" + file.Name(),
		})
	}

	// HTML template
	tmpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>File List</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                max-width: 1500px;
                margin: 50px auto;
                padding: 20px;
            }
            h1 {
                color: #333;
            }
            .header {
                display: flex;
                justify-content: space-between;
                align-items: center;
                margin-bottom: 30px;
            }
            .upload-btn {
                background: #28a745;
                color: white;
                padding: 10px 20px;
                border: none;
                border-radius: 5px;
                text-decoration: none;
                display: inline-block;
            }
            .upload-btn:hover {
                background: #218838;
            }
            table {
                width: 100%;
                border-collapse: collapse;
                background: white;
                box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            }
            th {
                background: #007bff;
                color: white;
                padding: 12px;
                text-align: left;
                font-weight: bold;
            }
            td {
                padding: 12px;
                border-bottom: 1px solid #ddd;
            }
            tr:hover {
                background: #f5f5f5;
            }
            .download-btn {
                background: #007bff;
                color: white;
                padding: 6px 12px;
                border-radius: 4px;
                text-decoration: none;
                font-size: 14px;
            }
			.convert-btn {
                background: #b39800ff;
                color: white;
                padding: 6px 12px;
                border-radius: 4px;
                text-decoration: none;
                font-size: 14px;
            }
            .convert-btn:hover {
                background: #dde033ff;
            }
			.view-btn {
                background: #b34800ff;
                color: white;
                padding: 6px 12px;
                border-radius: 4px;
                text-decoration: none;
                font-size: 14px;
            }
            .view-btn:hover {
                background: #e09233ff;
            }
            .delete-btn {
                background: #dc3545;
                color: white;
                padding: 6px 12px;
                border-radius: 4px;
                text-decoration: none;
                font-size: 14px;
                margin-left: 5px;
            }
            .delete-btn:hover {
                background: #c82333;
            }
            .no-files {
                text-align: center;
                padding: 40px;
                color: #666;
            }
            .file-count {
                color: #666;
                font-size: 14px;
            }
        </style>
    </head>
    <body>
        <div class="header">
            <div>
                <h1>Uploaded Files</h1>
                <p class="file-count">Total files: {{len .}}</p>
            </div>
            <a href="/upload-form" class="upload-btn">Upload New File</a>
        </div>

        {{if .}}
        <table>
            <thead>
                <tr>
                    <th>File Name</th>
                    <th>Size</th>
                    <th>Modified</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                {{range .}}
                <tr>
                    <td>{{.Name}}</td>
                    <td>{{.SizeFormatted}}</td>
                    <td>{{.ModTime}}</td>
                    <td>
                        <a href="{{.DownloadURL}}" class="download-btn">Download</a>
                        <a href="/delete/{{.Name}}" class="delete-btn" onclick="return confirm('Are you sure you want to delete this file?')">Delete</a>
						<a href="/convert/{{.Name}}" class="convert-btn">Convert to PDF</a>
						<a href="/view/{{.Name}}" class="view-btn">View</a>
                    </td>
                </tr>
                {{end}}
            </tbody>
        </table>
        {{else}}
        <div class="no-files">
            <p>No files uploaded yet.</p>
            <a href="/upload-form" class="upload-btn">Upload Your First File</a>
        </div>
        {{end}}
    </body>
    </html>
    `

	t, err := template.New("files").Parse(tmpl)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := t.Execute(w, fileInfos); err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := filepath.Base(vars["filename"])

	filePath := filepath.Join(UploadPath, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set headers for download
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")

	// Serve the file
	http.ServeFile(w, r, filePath)
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := filepath.Base(vars["filename"])

	filePath := filepath.Join(UploadPath, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		http.Error(w, "Error deleting file", http.StatusInternalServerError)
		return
	}

	log.Printf("File deleted: %s", filename)

	filePath = filepath.Join("./conversions/", filename+".pdf")
	// TODO: delete pdf file if exists
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		http.Error(w, "Error deleting file", http.StatusInternalServerError)
		return
	}

	log.Printf("File deleted: %s", filename+".pdf")

	// Redirect back to file list
	http.Redirect(w, r, "/files", http.StatusSeeOther)
}
