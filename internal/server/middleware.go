package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (s *Server) Middleware() error {
	s.e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:           middleware.DefaultSkipper,
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
		LogLevel:          0,
	}))
	s.e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Skipper: middleware.DefaultSkipper,
		RequestIDHandler: func(c echo.Context, id string) {
			c.Set("id", id)
		},
	}))
	s.e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		Skipper:               middleware.DefaultSkipper,
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	}))
	s.e.Use(middleware.BodyLimit(s.BodyLimit))
	s.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     s.AllowOrigins,
		AllowHeaders:     s.AllowHeaders,
		AllowMethods:     s.AllowMethods,
		AllowCredentials: s.AllowCredentials,
	}))
	s.e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "timeout",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			log.Error().
				Err(errors.New("timeout")).
				Str("URI", c.Request().URL.String()).
				Int("status", http.StatusServiceUnavailable).
				Msg("request")
		},
		Timeout: 30 * time.Second,
	}))
	s.e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogError:        true,
		LogURI:          true,
		LogStatus:       true,
		LogLatency:      true,
		LogRemoteIP:     true,
		LogMethod:       true,
		LogResponseSize: true,
		Skipper: func(c echo.Context) bool {
			return c.Request().URL.Path == "/favicon.ico"
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			var logMe *zerolog.Event
			if v.Error != nil {
				logMe = log.Error().
					Err(v.Error)
			} else {
				logMe = log.Info()
			}
			logMe.Str("URI", v.URI).
				Int("status", v.Status).
				Dur("latency", v.Latency).
				Str("remote_ip", v.RemoteIP).
				Interface("request_id", c.Get("id")).
				Str("method", v.Method).
				Int64("response_size", v.ResponseSize).
				Msg("request")
			return nil
		},
	}))
	return nil
}
