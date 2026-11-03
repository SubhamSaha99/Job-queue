package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type EnvConfig struct {
	Port string `env:"PORT" env-default:"8082"`
}

// Function to load ENVs
func LoadEnv() (*EnvConfig, error) {
	var config EnvConfig

	if err := cleanenv.ReadConfig("./.env", &config); err != nil {
		return nil, fmt.Errorf("Error loading env %w", err)
	}
	return &config, nil
}
