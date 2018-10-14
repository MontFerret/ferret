package server

import (
	"github.com/labstack/echo"
)

type Server struct {
	engine echo.Server
}

func New() *Server {
	e := echo.New()
}

func (s *Server) Run() error {

}
