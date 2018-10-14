package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	"path/filepath"
)

type (
	Settings struct {
		Port uint64
		Dir  string
	}
	Server struct {
		engine   *echo.Echo
		settings Settings
	}
)

func New(settings Settings) *Server {
	e := echo.New()
	e.Debug = false
	e.HideBanner = true

	e.Static("/static", filepath.Join(settings.Dir, "static"))

	return &Server{e, settings}
}

func (s *Server) Start() error {
	return s.engine.Start(fmt.Sprintf("0.0.0.0:%d", s.settings.Port))
}

func (s *Server) Stop(ctx context.Context) error {
	return s.engine.Shutdown(ctx)
}
