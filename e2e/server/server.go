package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/labstack/echo"
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
	api := e.Group("/api")
	api.GET("/ts", func(ctx echo.Context) error {
		var headers string

		if len(ctx.Request().Header) > 0 {
			b, err := json.Marshal(ctx.Request().Header)

			if err != nil {
				return err
			}

			headers = string(b)
		}

		ts := time.Now().Format("2006-01-02 15:04:05")

		return ctx.HTML(http.StatusOK, fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
			<head>
				<meta charset="utf-8" />
			</head>
			<body>
				<span id="timestamp">%s</span>
				<span id="headers">%s</span>
			</body>
		</html>
	`, ts, headers))
	})
	api.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"header": ctx.Request().Header,
			"url":    ctx.Request().URL,
			"data":   "pong",
			"ts":     time.Now(),
		})
	})

	return &Server{e, settings}
}

func (s *Server) Start() error {
	return s.engine.Start(fmt.Sprintf("0.0.0.0:%d", s.settings.Port))
}

func (s *Server) Stop(ctx context.Context) error {
	return s.engine.Shutdown(ctx)
}
