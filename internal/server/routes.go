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

	s.e.GET("/api/start", handler.NewStartHandler(s).Route)
	s.e.GET("/api/stop", handler.NewStopHandler(s).Route)
	s.e.GET("/api/status", handler.NewStatusHandler(s).Route)
	s.e.GET("/api/reload", handler.NewReloadHandler(s).Route)

	s.e.GET("/api/proxies", handler.NewGetProxiesHandler(s).Route)
	s.e.GET("/api/proxy", handler.NewGetProxyHandler(s).Route)
	s.e.POST("/api/proxy", handler.NewCreateProxyHandler(s).Route)
	s.e.PUT("/api/proxy", handler.NewUpdateProxyHandler(s).Route)
	s.e.DELETE("/api/proxy", handler.NewDeleteProxyHandler(s).Route)
	s.e.GET("/api/proxies/status", handler.NewGetProxiesStatusHandler(s).Route)
	s.e.GET("/api/proxy/crash", handler.NewGetCrashHandler(s).Route)

	s.e.GET("/api/haproxy/logs", handler.NewGetHaproxyLogsHandler(s).Route)

	s.e.GET("/api/settings", handler.NewGetSettingsHandler(s).Route)
	s.e.POST("/api/settings", handler.NewSetSettingsHandler(s).Route)
	s.e.POST("/api/settings/acme/cf", handler.NewCreateAcmeCf(s).Route)

	s.e.GET("/api/certificates", handler.NewGetCertificatesHandler(s).Route)
	s.e.POST("/api/certificate", handler.NewAddCertificateHandler(s).Route)
	s.e.DELETE("/api/certificate", handler.NewDeleteCertificateHandler(s).Route)
	s.e.POST("/api/certificate/request", handler.NewRequestCertificateHandler(s).Route)
	return nil
}
