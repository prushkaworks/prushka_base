package server

import (
	"net/http"
	"prushka/internal/db"
)

func (s *Server) buildRoutes() {

	// s.router.Path("/").HandlerFunc(MainHandler).Methods("GET")

	s.router.Path("/users/{id:[0-9]+}/").HandlerFunc(UserHandler).Methods("GET", "POST", "DELETE")
	s.router.Path("/users/").HandlerFunc(UserHandler).Methods("GET", "POST")
	s.router.Path("/auth/").HandlerFunc(AuthHandler).Methods("GET", "POST")

	s.router.Path("/privilege/{id:[0-9]+}/").HandlerFunc(ModelHandler(db.Privilege{})).Methods("GET", "POST", "DELETE")
	s.router.Path("/privilege/").HandlerFunc(ModelHandler(db.Privilege{})).Methods("GET", "POST")

	s.router.Path("/cards/{id:[0-9]+}/").HandlerFunc(ModelHandler(db.Card{})).Methods("GET", "POST", "DELETE")
	s.router.Path("/cards/").HandlerFunc(ModelHandler(db.Card{})).Methods("GET", "POST")

	s.router.Path("/attachment/{id:[0-9]+}/").HandlerFunc(ModelHandler(db.Attachment{})).Methods("GET", "POST", "DELETE")
	s.router.Path("/attachment/").HandlerFunc(ModelHandler(db.Attachment{})).Methods("GET", "POST")

	s.router.Path("/label/{id:[0-9]+}/").HandlerFunc(ModelHandler(db.Label{})).Methods("GET", "POST", "DELETE")
	s.router.Path("/label/").HandlerFunc(ModelHandler(db.Label{})).Methods("GET", "POST")

	s.router.Path("/column/{id:[0-9]+}/").HandlerFunc(ModelHandler(db.Column{})).Methods("GET", "POST", "DELETE")
	s.router.Path("/column/").HandlerFunc(ModelHandler(db.Column{})).Methods("GET", "POST")

	s.router.Path("/desk/{id:[0-9]+}/").HandlerFunc(ModelHandler(db.Desk{})).Methods("GET", "POST", "DELETE")
	s.router.Path("/desk/").HandlerFunc(ModelHandler(db.Desk{})).Methods("GET", "POST")

	s.router.Path("/workspace/{id:[0-9]+}/").HandlerFunc(ModelHandler(db.Workspace{})).Methods("GET", "POST", "DELETE")
	s.router.Path("/workspace/").HandlerFunc(ModelHandler(db.Label{})).Methods("GET", "POST")

	s.router.Path("/user_privilege/{id:[0-9]+}/").HandlerFunc(ModelHandler(db.UserPrivilege{})).Methods("GET", "POST", "DELETE")
	s.router.Path("/user_privilege/").HandlerFunc(ModelHandler(db.UserPrivilege{})).Methods("GET", "POST")

	s.router.Path("/cards_label/{id:[0-9]+}/").HandlerFunc(ModelHandler(db.CardsLabel{})).Methods("GET", "POST", "DELETE")
	s.router.Path("/cards_label/").HandlerFunc(ModelHandler(db.CardsLabel{})).Methods("GET", "POST")

	s.router.Use(LoggingAndJson)
	s.router.NotFoundHandler = LoggingAndJson(http.HandlerFunc(My404Handler))
	s.Handler = s.router
}
