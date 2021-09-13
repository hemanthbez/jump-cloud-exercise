package main

import (
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const PORT_PARAM string = "port"
const DEFAULT_PORT_NUMBER string = "8080"
const PORT_NUMBER_PARAM_MSG string = "Port number to be used"
const PASSWORD_RESPOND_TIME_LIMIT_SECS float64 = 5

type PASSWORD_DATA struct {
	id                      int
	sha512value             string
	requestTimestamp        time.Time
	processingTimeMicrosecs int64
}

func (passwordDataObj *PASSWORD_DATA) SetProcessingTimeMicrosecs(processingTimeMicrosecs int64) {
	passwordDataObj.processingTimeMicrosecs = processingTimeMicrosecs
}

type PasswordMap map[int]*PASSWORD_DATA

var passwordMap PasswordMap = make(PasswordMap)
var requestCounter int
var shutdownEnabled bool
var requestInProgress bool

func constructPwdData(requestId int, password string) *PASSWORD_DATA {
	passwordStructData := PASSWORD_DATA{
		id:               requestId,
		sha512value:      encryptPwdWithSha512(password),
		requestTimestamp: time.Now(),
	}
	return &passwordStructData
}

func encryptPwdWithSha512(password string) string {
	sha512 := sha1.New()
	sha512.Write([]byte(password))
	sha512valueStr := base64.URLEncoding.EncodeToString(sha512.Sum(nil))
	return sha512valueStr
}

func processPasswordPostRequest(w http.ResponseWriter, r *http.Request) {
	requestInProgress = true
	r.ParseForm()
	var passwordParam string = r.FormValue("password")
	if len(strings.TrimSpace(passwordParam)) == 0 {
		fmt.Fprintf(w, "ERROR: Password empty \n")
		return
	}

	currentTime := time.Now()
	requestCounter++

	passwordDataObj := constructPwdData(requestCounter, passwordParam)
	passwordMap[requestCounter] = passwordDataObj
	currentTs := time.Now().Sub(currentTime).Microseconds()
	passwordDataObj.SetProcessingTimeMicrosecs(currentTs)
	fmt.Fprintf(w, "%d \n", requestCounter)
	requestInProgress = false
	log.Printf("requestInProgress (PWD) --> %v \n", requestInProgress)
	return
}

func processPasswordGetRequest(w http.ResponseWriter, r *http.Request) {
	requestInProgress = true
	requestID, _ := strconv.Atoi(getField(r, 0))
	log.Printf("requestID: %d \n", requestID)
	if entry, ok := passwordMap[requestID]; ok {
		sha512Val := entry.sha512value
		requestedTimestamp := entry.requestTimestamp
		timeDuration := time.Now().Sub(requestedTimestamp).Seconds()

		if timeDuration >= PASSWORD_RESPOND_TIME_LIMIT_SECS {
			fmt.Fprintf(w, "%s \n", sha512Val)
		} else {
			http.Error(w, "404 Data not available", http.StatusNotFound)
			return
		}
	} else {
		http.Error(w, "400 Data not available", http.StatusBadRequest)
		return
	}
	requestInProgress = false
}

func processStats(w http.ResponseWriter, r *http.Request) {
	requestInProgress = true
	totalEntries := len(passwordMap)

	var totalTimeMs int64
	for _, element := range passwordMap {
		totalTimeMs += element.processingTimeMicrosecs
	}

	average := int(totalTimeMs) / totalEntries

	fmt.Fprintf(w, "{\"total\": %d, \"average\": %d } \n", totalEntries, average)
	requestInProgress = false
}

func processShutown(w http.ResponseWriter, r *http.Request) {
	log.Println("Shutdown Triggered....")
	if !shutdownEnabled {
		shutdownEnabled = true
	}

	for !requestInProgress {
		logMsg := "Exiting Program!"
		notifyExit(w, r, logMsg)
		os.Exit(0)
	}

}

func notifyExit(w http.ResponseWriter, r *http.Request, logMsg string) {
	fmt.Fprintf(w, "%s", logMsg)
	return
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Serve)
	portNumberPtr := flag.String(PORT_PARAM, DEFAULT_PORT_NUMBER, PORT_NUMBER_PARAM_MSG)
	flag.Parse()
	log.Printf("Starting server using the Port number: %s \n", *portNumberPtr)
	http.ListenAndServe(":"+*portNumberPtr, mux)
}
