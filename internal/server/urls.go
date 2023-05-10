package server

func (s *Server) buildRoutes() {
	s.router.HandleFunc("/", mainHandler)
	s.Handler = s.router
}
