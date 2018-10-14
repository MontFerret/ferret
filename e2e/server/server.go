package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
)

type (
	Settings struct {
		Port uint64
	}
	Server struct {
		engine   *echo.Echo
		settings Settings
	}
)

func New(settings Settings) *Server {
	e := echo.New()

	return &Server{e, settings}
}

func (s *Server) Start() error {
	return s.engine.Start(fmt.Sprintf("0.0.0.0:%d", s.settings.Port))
}

func (s *Server) Stop(ctx context.Context) error {
	return s.engine.Shutdown(ctx)
}
