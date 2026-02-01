package handlers

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
)

type CommandResponse struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

func textToPDF(textFile, pdfFile string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	file, err := os.Open(textFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pdf.Cell(0, 10, scanner.Text())
		pdf.Ln(-1)
	}

	return pdf.OutputFileAndClose(pdfFile)
}

func ConvertFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := filepath.Base(vars["filename"])
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current directory:", dir)
	fmt.Println(filename)
	var filePath string = "./conversions/" + filename + ".pdf"
	textToPDF("./uploads/"+filename, filePath)

	cmd := exec.Command("ls", "-la")
	output, err := cmd.CombinedOutput()

	response := CommandResponse{
		Success: err == nil,
		Output:  string(output),
	}

	if err != nil {
		response.Error = err.Error()
	}

	fmt.Println(response.Output)

	var contentType string = "application/pdf"
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "inline; filename="+filename)

	http.ServeFile(w, r, filePath)
}
