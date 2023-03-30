package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App       `yaml:"app"`
		HTTP      `yaml:"http"`
		Log       `yaml:"logger"`
		DB        `yaml:"db"`
		AWS       `yaml:"aws"`
		SecretKey string `env-required:"true" yaml:"secret_key" env:"SECRET_KEY"`
	}

	// App -.
	App struct {
		Name           string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version        string `env-required:"true" yaml:"version" env:"APP_VERSION"`
		WorkerInterval int    `env-required:"true" yaml:"worker_interval" env:"WORKER_INTERVAL"`
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

	AWS struct {
		EndpointUrl       string `env-required:"true" yaml:"aws_endpoint_url" env:"AWS_ENDPOINT_URL"`
		Region            string `env-required:"true" yaml:"aws_region" env:"AWS_REGION"`
		FakeNewsQueueName string `env-required:"true" yaml:"fake_news_queue_name" env:"FAKE_NEWS_QUEUE_NAME"`
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
