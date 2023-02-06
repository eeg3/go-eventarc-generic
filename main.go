package main

import (
    log "github.com/sirupsen/logrus"
	"encoding/json"
	"io/ioutil"
	"net/http"
    "fmt"
	"os"
	"strings"
)

// GenericHandler receives and echos a HTTP request's headers and body.
func GenericHandler(w http.ResponseWriter, r *http.Request) {
    log.Info("Event received!")

	// Log all headers besides Authorization header.
	headerMap := make(map[string]string)
	for k, v := range r.Header {
		if k != "Authorization" {
			val := strings.Join(v, ",")
			headerMap[k] = val
            log.Debug(fmt.Sprintf("%q: %q\n", k, val))
		}
	}

	// Log body.
	bodyMap := make(map[string]interface{})
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Error parsing body: " + err.Error())
        return
	}
    log.Debug(string(bodyBytes))
    if err := json.Unmarshal(bodyBytes, &bodyMap); err != nil {
        log.Error("JSON decode error: " + err.Error())
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    // Write to HTTP response and log output in { "headers": {}, "body": {} } JSON format. 
	type result struct {
		Headers map[string]string `json:"headers"`
		Body    interface{} `json:"body"`
	}
	res := &result{
		Headers: headerMap,
		Body:    bodyMap,
	}
    output, err := json.Marshal(res);
    if err != nil {
        log.Error("JSON decode error: " + err.Error())
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }
    log.Info(string(output))
    w.WriteHeader(http.StatusOK)
    w.Write(output)
}

func main() {
    log.SetFormatter(&log.TextFormatter{
        DisableQuote: true,
        DisableTimestamp: true,
    })

	http.HandleFunc("/", GenericHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
        log.Info("Defaulting to port: ", port)
	}
    log.Info("Listening on port: ", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

