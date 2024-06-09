package server

import (
	"testproject/internal/handler"

	"github.com/labstack/echo/v4"
)

func (s *Server) Routes() error {
	s.e.GET("/*", func(c echo.Context) error {
		return c.File("dist/index.html")
	})
	s.e.Static("/assets/*", "dist/assets")

	s.e.GET("/api/test", handler.Test)
	s.e.GET("/api/start", handler.NewStartHandler(s).Route)
	s.e.GET("/api/stop", handler.NewStopHandler(s).Route)
	s.e.GET("/api/status", handler.Status)
	s.e.GET("/api/proxies", handler.NewGetProxiesHandler(s).Route)
	s.e.GET("/api/proxies/status", handler.NewGetProxiesStatusHandler(s).Route)
	s.e.POST("/api/proxy", handler.NewCreateProxyHandler(s).Route)
	s.e.DELETE("/api/proxy", handler.NewDeleteProxyHandler(s).Route)
	s.e.GET("/api/haproxy/logs", handler.NewGetHaproxyLogsHandler(s).Route)
	return nil
}
