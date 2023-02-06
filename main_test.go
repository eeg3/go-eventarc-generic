package main

import (
    log "github.com/sirupsen/logrus"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenericHandler(t *testing.T) {
    log.SetOutput(ioutil.Discard)

    jsonStr := fmt.Sprintf(`{"message": "some string"}`)
    payload := strings.NewReader(jsonStr)

    req := httptest.NewRequest(http.MethodPost, "/", payload)
    req.Header.Set("Ce-Id", "1234")
    req.Header.Set("Authorization", "super-secret-value")
    req.Header.Set("Ce-Source", "//storage.googleapis.com/projects/YOUR-PROJECT")
    w := httptest.NewRecorder()
    GenericHandler(w, req)
    res := w.Result()
    defer res.Body.Close()
    data, err := ioutil.ReadAll(res.Body)
    if err != nil {
        t.Errorf("Expected error to be nil, got: %v", err)
    }
    if string(data) != `{"headers":{"Ce-Id":"1234","Ce-Source":"//storage.googleapis.com/projects/YOUR-PROJECT"},"body":{"message":"some string"}}` {
        t.Errorf("Got unexpected data: %v", string(data))
    }
}
