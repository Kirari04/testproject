package server

import (
	"net/http"
	"testproject/internal/handler"

	"github.com/labstack/echo/v4"
)

func (s *Server) Routes() error {
	s.e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	s.e.GET("/api/test", handler.Test)
	s.e.GET("/api/start", handler.Start)
	s.e.GET("/api/stop", handler.Stop)
	s.e.GET("/api/status", handler.Status)
	s.e.GET("/api/proxies", handler.NewGetProxiesHandler(s).Route)
	s.e.POST("/api/proxy", handler.NewCreateProxyHandler(s).Route)
	s.e.POST("/api/proxy/apply", handler.NewApplyHandler(s).Route)
	return nil
}
