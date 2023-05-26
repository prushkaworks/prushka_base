package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	respBody, _ := json.Marshal(struct {
		Text string `json:"text"`
	}{Text: "Hello World!"})
	fmt.Fprint(w, string(respBody))
}

func My404Handler(w http.ResponseWriter, r *http.Request) {
	respBody, _ := json.Marshal(struct {
		Text string `json:"message"`
	}{Text: "404 Not found"})
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, string(respBody))
}

func LoggingAndJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		start := time.Now()
		next.ServeHTTP(w, req)
		log.Printf("%s %s%s\t%s", color.BlueString(req.Method), req.Host, req.URL.Path, color.CyanString(fmt.Sprintf("\t%fsec", time.Since(start).Seconds())))
	})
}
