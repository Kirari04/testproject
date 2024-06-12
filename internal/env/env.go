package env

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Env struct {
	Addr    string `env:"ADDR" envDefault:"0.0.0.0:8080"`
	WorkDir string `env:"WORK_DIR" envDefault:"./.data"`
}

func NewEnv() (*Env, error) {
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("failed to load .env file")
	}

	cfg := Env{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
