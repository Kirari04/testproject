package cmd

import (
	"testproject/internal/server"

	"github.com/urfave/cli/v2"
)

func serve(c *cli.Context) error {
	s := server.NewServer()
	if err := s.Start(); err != nil {
		return err
	}
	return nil
}
