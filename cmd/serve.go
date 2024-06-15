package cmd

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"testproject/internal/server"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func serve(c *cli.Context) error {
	var useTls bool
	if c.Bool("tls") {
		useTls = true
	}

	var useSocket bool
	if c.Bool("socket") {
		useSocket = true
	}
	os.Setenv("SOCKET", strconv.FormatBool(useSocket))

	s, err := server.NewServer()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	doingExit := false
	// Start server
	go func() {
		if err := s.Start(useTls); err != nil {
			if !doingExit {
				log.Error().Err(err).Msg("failed to start server")
			}
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	doingExit = true

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Info().Msg("Stopping server")
	if err := s.Stop(ctx); err != nil {
		return err
	}

	return nil
}
