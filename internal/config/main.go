package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
)

type config struct {
	DB_FILE string `env:",notEmpty,required"`
}

func Load() *config {
	var cfg config

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalln(fmt.Errorf("load config: %w", err))
	}

	return &cfg
}
