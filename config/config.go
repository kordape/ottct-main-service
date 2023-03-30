package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App              `yaml:"app"`
		HTTP             `yaml:"http"`
		Log              `yaml:"logger"`
		DB               `yaml:"db"`
		SecretKey        string `env-required:"true" yaml:"secret_key" env:"SECRET_KEY"`
		TwitterBearerKey string `env-required:"true" yaml:"twitter_bearer_key" env:"TWITTER_BEARER_KEY"`
		PredictorURL     string `env-required:"true" yaml:"predictor_url" env:"PREDICTOR_URL"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	DB struct {
		URL string `env-required:"true" yaml:"username" env:"DB_URL"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
