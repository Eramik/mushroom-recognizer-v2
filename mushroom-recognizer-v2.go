package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

// Closes request with error if any errors occured
func anyErrors(error e, w http.ResponseWriter) bool {
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return true
	} else {
		return false
	}
}

// Generates random filename with length 32 with provided extension
// Takes salt from cfg file, takes integer seed as last insert in database
func generateFileName(fileExtension string) string {

}

// Saves file and returs its url and error
func receiveFile(r *http.Request) (string, error) {
	// Read file from request

	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("file")
	if err != nil {
		return _, err
	}

	defer file.Close()

	// Save file to filesystem

	dir, err := os.Getwd() // Get current dir
	if err != nil {
		return _, err
	}

	fileExtensionRegexp := regexp.MustCompile("(.jpg)|(.png)")
	fileExtension := fileExtensionRegexp.FindString(header.Filename)

	if fileExtension == "" {
		return _, errors.New("File extension must be .jpg or .png only")
	}

	path := dir + "\\imgs\\" + filename
	out, err := os.Create(path)
	if err != nil {
		return _, err
	}

	defer out.Close()
	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		return _, err
	}

	return path, nil
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

}

func main() {
	http.HandleFunc("/", uploadHandler)
	http.ListenAndServe(":8080", nil)
}
