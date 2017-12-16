package main

//import "database/sql"
import "encoding/json"
import "errors"
import "fmt"
import "io"
import "net/http"
import "os"
import "regexp"
import "bytes"
import "io/ioutil"

import "github.com/rs/xid"

//import "github.com/go-sql-driver/mysql"

type Configuration struct {
	MySQLUser     string
	MySQLDatabase string
	MySQLPassword string
	PredictionURL string
	PredictionKey string
}

const CONFIGURATION_PATH = "recognizer20conf.json"
const IMGS_STORAGE_DIR = "/www/dgodovanets.shpp.me/recognizer-v2/imgs/"

// http:// will be added automatically
const REGEXP_IMG_URL = "dgodovanets.shpp.me/.+$"
const PORT = ":9090"

var configuration = Configuration{}

func loadConfiguration() error {
	file, _ := os.Open(CONFIGURATION_PATH)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	if err != nil {
		return err
	}
	return nil
}

// Closes request with error if any errors occured
func anyErrors(e error, w http.ResponseWriter) bool {
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
	filename := xid.New().String()
	return filename + fileExtension
}

// Saves file and returs its url and error
func receiveFile(r *http.Request) (string, error) {
	// Read file from request

	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("file")
	if err != nil {
		return "", err
	}

	defer file.Close()

	// Save file to filesystem

	dir, err := os.Getwd() // Get current dir
	if err != nil {
		return "", err
	}

    fileExtensionRegexp := regexp.MustCompile("(.jpg$)|(.jpeg&)|(.png$)")
	fileExtension := fileExtensionRegexp.FindString(header.Filename)

	if fileExtension == "" {
		return "", errors.New("File extension must be .jpg or .png only")
	}

	filename := generateFileName(fileExtension)

	path := dir + IMGS_STORAGE_DIR + filename
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}

	defer out.Close()
	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	// Generate url from path
	urlRegexp := regexp.MustCompile(REGEXP_IMG_URL)
	url := urlRegexp.FindString(path)
	if url == "" {
		return "", errors.New("Can't generate url from path of file")
	}

	url = "http://" + url

	return url, nil
}

func getPrediction(fileUrl string) (string, error) {
	client := &http.Client{}
	var jsonStr = []byte(`{"Url": "` + fileUrl + `"}`)
	req, err := http.NewRequest("POST", configuration.PredictionURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}

	req.Header.Add("Prediction-Key", configuration.PredictionKey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	respString := string(respData)

	return respString, nil
}

func controller(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method == "OPTIONS" {
		return
	}

	fileUrl, err := receiveFile(r)
	if anyErrors(err, w) {
		return
	}

	resp, err := getPrediction(fileUrl)
	if anyErrors(err, w) {
		return
	}

	w.Write([]byte("{ \"success\": true, \"data\": " + resp + " }"))
}

func main() {
	err := loadConfiguration()
	if err != nil {
		fmt.Println("Error while loading configuration: ", err)
		return
	}
	http.HandleFunc("/receive", controller)
	http.ListenAndServe(PORT, nil)
	fmt.Println("Running on port ", PORT)
}
