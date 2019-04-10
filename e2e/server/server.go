package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
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

	e.Use(func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.SetCookie(&http.Cookie{
				Name:     "x-ferret",
				Value:    "e2e",
				HttpOnly: false,
			})

			return handlerFunc(ctx)
		}
	})
	e.Static("/", settings.Dir)
	e.File("/", filepath.Join(settings.Dir, "index.html"))

	return &Server{e, settings}
}

func (s *Server) Start() error {
	return s.engine.Start(fmt.Sprintf("0.0.0.0:%d", s.settings.Port))
}

func (s *Server) Stop(ctx context.Context) error {
	return s.engine.Shutdown(ctx)
}
