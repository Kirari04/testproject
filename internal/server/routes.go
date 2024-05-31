package server

import (
	"testproject/internal/handler"
)

func (s *Server) Routes() error {
	s.e.Static("/", "./dist")
	s.e.GET("/api/test", handler.Test)
	s.e.GET("/api/start", handler.Start)
	s.e.GET("/api/stop", handler.Stop)
	s.e.GET("/api/status", handler.Status)
	s.e.GET("/api/proxies", handler.NewGetProxiesHandler(s).Route)
	s.e.POST("/api/proxy", handler.NewCreateProxyHandler(s).Route)
	s.e.DELETE("/api/proxy", handler.NewDeleteProxyHandler(s).Route)
	s.e.POST("/api/proxy/apply", handler.NewApplyHandler(s).Route)
	return nil
}
