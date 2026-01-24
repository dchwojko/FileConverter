package handlers

import (
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body>")
	fmt.Fprintf(w, "<p>Hello from Go Docker Server with File Upload!</p>")
	fmt.Fprintf(w, "<a href='http://localhost:80/upload-form' target='_self'>Upload File</a>")
	fmt.Fprintf(w, "<div class='link'><a href='/files'>View All Files</a></div>")
	fmt.Fprintf(w, "</body></html>")
}

func UploadFormHandler(w http.ResponseWriter, r *http.Request) {
	html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>File Upload</title>
        <style>
            body { font-family: Arial, sans-serif; max-width: 600px; margin: 50px auto; padding: 20px; }
            .upload-form { border: 2px dashed #ccc; padding: 30px; border-radius: 10px; }
            input[type="file"] { margin: 20px 0; }
            button { background: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; }
            button:hover { background: #0056b3; }
        </style>
    </head>
    <body>
		<a href='../'>Back</a>
        <h1>Upload a File</h1>
        <div class="upload-form">
            <form action="/upload" method="post" enctype="multipart/form-data">
                <input type="file" name="file" required>
                <br>
                <button type="submit">Upload</button>
            </form>
        </div>
    </body>
    </html>
    `
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}
