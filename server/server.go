package server

import (
	"github.com/dannylindquist/spend-go/sqlite"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	userService *sqlite.UserService
	router *httprouter.Router
}

func NewServer(db *sqlite.DB) *Server {
	s := &Server{
		router: httprouter.New(),
	}

	s.setup(db)

	return s
}

func (s *Server) setup(db *sqlite.DB) {

	s.userService = sqlite.NewUserService(db)

	s.setupRoutes()
}

func (s *Server) setupRoutes() {

}