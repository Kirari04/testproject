package main

import (
	"testproject/cmd"

	"github.com/rs/zerolog/log"
)

func main() {
	if err := cmd.Cmd(); err != nil {
		log.Fatal().Err(err).Msg("error")
	}
}
