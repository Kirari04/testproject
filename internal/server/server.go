package server

import (
	"testproject/internal/env"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type Server struct {
	e                *echo.Echo
	BodyLimit        string
	AllowOrigins     []string
	AllowHeaders     []string
	AllowMethods     []string
	AllowCredentials bool
	ENV              *env.Env
}

func NewServer() *Server {
	cfg, err := env.NewEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load env")
	}
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	s := &Server{
		e:         e,
		BodyLimit: "2M",
		AllowOrigins: []string{
			"*",
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAcceptEncoding,
			echo.HeaderAuthorization,
		},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
		AllowCredentials: false,
		ENV:              cfg,
	}
	s.Middleware()
	s.Routes()
	return s
}

func (s *Server) Start() error {
	log.Info().Msgf("Server listening on http://%s", s.ENV.Addr)
	return s.e.Start(s.ENV.Addr)
}
