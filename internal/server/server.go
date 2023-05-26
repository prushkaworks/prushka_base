package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"prushka/internal/db"

	"github.com/gorilla/mux"
)

type Server struct {
	http.Server
	router *mux.Router
}

func New() *Server {
	return &Server{
		router: mux.NewRouter(),
	}
}

func (s *Server) Run() {
	db.ConnectDb()
	s.buildRoutes()
	s.Addr = "127.0.0.1" + ":" + os.Getenv("PORT")
	log.Println("Server is listening...")

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := s.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		log.Printf("\tClosing...")
		close(idleConnsClosed)
	}()

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
