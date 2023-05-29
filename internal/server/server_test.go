package server_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	serv "prushka/internal/server"
	"strings"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"golang.org/x/exp/slices"
)

var timeout time.Duration = time.Second * 1000

func TestIndex(t *testing.T) {

	router := mux.NewRouter()
	router.NotFoundHandler = serv.LoggingAndJson(http.HandlerFunc(serv.My404Handler))
	router.Use(serv.LoggingAndJson)
	// router.HandleFunc("/", serv.MainHandler)
	s := httptest.NewServer(router)
	defer s.Close()

	caser := func(url string) string {
		client := &http.Client{
			Timeout: timeout,
		}
		resp, _ := client.Get(url)
		if resp != nil {
			defer resp.Body.Close()
			if (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound) && slices.Contains(resp.Header["Content-Type"], "application/json") {
				bodyBytes, _ := io.ReadAll(resp.Body)
				bodyString := string(bodyBytes)
				return fmt.Sprintf("Api returns: %s", bodyString)
			}
		}

		return ""
	}

	tests := []struct {
		name  string
		caser string
	}{
		{name: "1. Main page", caser: caser(s.URL)},
		{name: "2. Random page Not found", caser: caser(s.URL + "/sasdfascs")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.caser != "" {
				if strings.Contains(tt.caser, "404") {
					t.Logf("name: %s\nresult: %s", tt.name, color.YellowString(tt.caser))
				} else {
					t.Logf("name: %s\nresult: %s", tt.name, color.GreenString(tt.caser))
				}
			} else {
				t.Logf("name: %s\nresult: %s", tt.name, color.RedString("Error"))
			}
		})
	}
}
