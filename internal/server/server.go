package server

import (
	"context"
	"fmt"
	"os"
	"testproject/internal/db"
	"testproject/internal/env"
	"testproject/internal/haproxy"
	"testproject/internal/util"
	"time"

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
	Haproxy          *haproxy.Haproxy
}

func NewServer() (*Server, error) {
	cfg, err := env.NewEnv()
	if err != nil {
		log.Error().Err(err).Msg("failed to load env")
		return nil, err
	}

	if err := initDirs(cfg); err != nil {
		return nil, err
	}

	db, err := db.Connect(cfg)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to db")
		return nil, err
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

	haproxy := haproxy.NewHaproxy(s)
	s.Haproxy = haproxy

	s.Middleware()
	s.Routes()

	return s, nil
}

func (s *Server) Start(tls bool) error {
	go s.HaKeepAlive()

	if tls {
		// generate self-signed certs if not exists
		_, err1 := os.Stat(fmt.Sprintf("%s/panel.crt", s.env.WorkDir))
		_, err2 := os.Stat(fmt.Sprintf("%s/panel.key", s.env.WorkDir))
		if os.IsNotExist(err1) || os.IsNotExist(err2) {
			host := "localhost"
			validFrom := ""
			validFor := 365 * 24 * time.Hour
			isCA := false
			rsaBits := 2048
			ecdsaCurve := ""
			ed25519Key := false

			cert, key, err := util.GenerateSelfSignedCert(host, validFrom, validFor, isCA, rsaBits, ecdsaCurve, ed25519Key)
			if err != nil {
				return err
			}
			if err := os.WriteFile(fmt.Sprintf("%s/panel.crt", s.env.WorkDir), []byte(cert), 0644); err != nil {
				return err
			}
			if err := os.WriteFile(fmt.Sprintf("%s/panel.key", s.env.WorkDir), []byte(key), 0644); err != nil {
				return err
			}
		}
		log.Info().Msgf("Server listening on https://%s", s.ENV().Addr)
		return s.e.StartTLS(
			s.ENV().Addr,
			fmt.Sprintf("%s/panel.crt", s.env.WorkDir),
			fmt.Sprintf("%s/panel.key", s.env.WorkDir),
		)
	} else {
		log.Info().Msgf("Server listening on http://%s", s.ENV().Addr)
		return s.e.Start(
			s.ENV().Addr,
		)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.Haproxy.StopKeepAlive()
	return s.e.Shutdown(ctx)
}

func (s *Server) DB() *gorm.DB {
	return s.db
}

func (s *Server) ENV() *env.Env {
	return s.env
}

func initDirs(cfg *env.Env) error {
	if err := os.MkdirAll(cfg.WorkDir, 0755); err != nil {
		log.Error().Err(err).Msg("failed to create work dir")
		return err
	}

	if err := os.MkdirAll(cfg.WorkDir+"/certs", 0755); err != nil {
		log.Error().Err(err).Msg("failed to create work dir certs")
		return err
	}
	if err := os.MkdirAll(cfg.WorkDir+"/haproxy", 0755); err != nil {
		log.Error().Err(err).Msg("failed to create haproxy dir")
		return err
	}
	return nil
}
