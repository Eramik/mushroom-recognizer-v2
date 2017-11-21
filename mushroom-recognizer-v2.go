package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("file")

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	defer file.Close()

	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	out, err := os.Create(dir + "\\tmp\\uploadedfile")
	if err != nil {
		//fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
		fmt.Fprintf(w, "Err: "+err.Error())
		return
	}

	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	//fmt.Fprintf(w, "File uploaded successfully : ")
	fmt.Fprintf(w, header.Filename)

}

func main() {
	http.HandleFunc("/", uploadHandler)
	http.ListenAndServe(":8080", nil)
}
