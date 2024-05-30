package server

import (
	"context"
	"os"
	"testproject/internal/db"
	"testproject/internal/env"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Server struct {
	e                *echo.Echo
	BodyLimit        string
	AllowOrigins     []string
	AllowHeaders     []string
	AllowMethods     []string
	AllowCredentials bool
	env              *env.Env
	db               *gorm.DB
}

func NewServer() *Server {
	cfg, err := env.NewEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load env")
	}

	if err := os.MkdirAll(cfg.WorkDir, 0755); err != nil {
		log.Fatal().Err(err).Msg("failed to create work dir")
	}

	db, err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
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
		env:              cfg,
		db:               db,
	}
	s.Middleware()
	s.Routes()
	return s
}

func (s *Server) Start() error {
	log.Info().Msgf("Server listening on http://%s", s.ENV().Addr)
	return s.e.Start(s.ENV().Addr)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func (s *Server) DB() *gorm.DB {
	return s.db
}

func (s *Server) ENV() *env.Env {
	return s.env
}
