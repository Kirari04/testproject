package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func Cmd() error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cmd := &cli.App{
		Name:  "cmd",
		Usage: "a cmd just to be a cmd",
		Action: func(*cli.Context) error {
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:   "serve",
				Usage:  "server the webserver",
				Action: serve,
			},
		},
	}

	if err := cmd.Run(os.Args); err != nil {
		return err
	}

	return nil
}
