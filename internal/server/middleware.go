package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
		ContentSecurityPolicy: "default-src 'self' style-src 'self' 'unsafe-inline';",
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
				Err(err).
				Str("URI", c.Request().URL.RequestURI()).
				Int("status", http.StatusServiceUnavailable).
				Str("remote_ip", c.RealIP()).
				Interface("request_id", c.Get("id")).
				Str("method", c.Request().Method).
				Int64("response_size", c.Response().Size).
				Msg("timeout")
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
			if v.Error == nil {
				log.Info().Str("URI", v.URI).
					Int("status", v.Status).
					Dur("latency", v.Latency).
					Str("remote_ip", v.RemoteIP).
					Interface("request_id", c.Get("id")).
					Str("method", v.Method).
					Int64("response_size", v.ResponseSize).
					Msg("request")
			}
			return nil
		},
	}))
	s.e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		he, ok := err.(*echo.HTTPError)
		if ok {
			code = he.Code
		}

		log.Error().
			Err(err).
			Str("URI", c.Request().URL.RequestURI()).
			Int("status", code).
			Str("remote_ip", c.RealIP()).
			Interface("request_id", c.Get("id")).
			Str("method", c.Request().Method).
			Int64("response_size", c.Response().Size).
			Msg("HTTPErrorHandler")

		// if is http error
		if ok {
			if strings.HasPrefix(c.Request().URL.Path, "/api") {
				c.JSON(code, map[string]string{"error": fmt.Sprint(he.Message)})
			} else {
				c.String(code, fmt.Sprint(he.Message))
			}
		}
		if !ok {
			if strings.HasPrefix(c.Request().URL.Path, "/api") {
				c.JSON(code, map[string]string{"error": "internal server error"})
			} else {
				c.String(code, "internal server error")
			}
		}

	}
	return nil
}
