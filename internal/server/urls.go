package server

import "net/http"

func (s *Server) buildRoutes() {
	s.router.HandleFunc("/", MainHandler)
	s.router.Use(LoggingAndJson)
	s.router.NotFoundHandler = LoggingAndJson(http.HandlerFunc(My404Handler))
	s.Handler = s.router
}
