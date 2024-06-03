package server

import (
	"context"
	"fmt"
	"os"
	"testproject/internal/app"
	"testproject/internal/bg"
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
}

func NewServer() *Server {
	cfg, err := env.NewEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load env")
	}

	if err := os.MkdirAll(cfg.WorkDir, 0755); err != nil {
		log.Fatal().Err(err).Msg("failed to create work dir")
	}

	if err := os.MkdirAll(cfg.WorkDir+"/certs", 0755); err != nil {
		log.Fatal().Err(err).Msg("failed to create work dir certs")
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

	app.Proxy = haproxy.NewHaproxy(s)

	// leaking goroutines
	go bg.KeepAliveProxy(s)

	return s
}

func (s *Server) Start(tls bool) error {
	if tls {
		// generate self-signed certs if not exists
		_, err1 := os.Stat(fmt.Sprintf("%s/certs/server.crt", s.env.WorkDir))
		_, err2 := os.Stat(fmt.Sprintf("%s/certs/server.key", s.env.WorkDir))
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
			if err := os.WriteFile(fmt.Sprintf("%s/certs/server.crt", s.env.WorkDir), []byte(cert), 0644); err != nil {
				return err
			}
			if err := os.WriteFile(fmt.Sprintf("%s/certs/server.key", s.env.WorkDir), []byte(key), 0644); err != nil {
				return err
			}
		}
		log.Info().Msgf("Server listening on https://%s", s.ENV().Addr)
		return s.e.StartTLS(
			s.ENV().Addr,
			fmt.Sprintf("%s/certs/server.crt", s.env.WorkDir),
			fmt.Sprintf("%s/certs/server.key", s.env.WorkDir),
		)
	} else {
		log.Info().Msgf("Server listening on http://%s", s.ENV().Addr)
		return s.e.Start(
			s.ENV().Addr,
		)
	}
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
