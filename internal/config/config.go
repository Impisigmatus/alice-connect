package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/Impisigmatus/service_core/log"
	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel             log.Level `env:"-"`
	LogLevelStr          string    `env:"LOG_LEVEL" envDefault:"debug"`
	Port                 string    `env:"PORT" envDefault:":8000"`
	SpecPath             string    `env:"OPENAPI_SPEC" envDefault:"autogen/oapi3.yaml"`
	AliceProtocolVersion string    `env:"ALICE_PROTOCOL_VERSION" envDefault:"1.0"`
	SerialDriverAPI      SerialDriverAPIConfig
}

type SerialDriverAPIConfig struct {
	Hostname    string        `env:"SERIAL_DRIVER_API" envDefault:"http://serial-driver:8000"`
	ContentType string        `env:"SERIAL_DRIVER_API_CONTENT_TYPE" envDefault:"text/plain; charset=utf-8"`
	TimeoutInt  int           `env:"SERIAL_DRIVER_API_TIMEOUT" envDefault:"5"`
	Timeout     time.Duration `env:"-"`
	RetryCount  int           `env:"SERIAL_DRIVER_API_RETRY_COUNT" envDefault:"2"`
}

func Load() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, err
	}

	level, err := parseLogLevel(cfg.LogLevelStr)
	if err != nil {
		return Config{}, err
	}
	cfg.LogLevel = level

	cfg.SerialDriverAPI.Timeout = time.Duration(cfg.SerialDriverAPI.TimeoutInt) * time.Second
	return cfg, nil
}

func parseLogLevel(value string) (log.Level, error) {
	if strings.EqualFold(value, "warning") {
		value = "warn"
	}
	level, err := zerolog.ParseLevel(value)
	if err != nil {
		return 0, fmt.Errorf("некорректный LOG_LEVEL %q: %w", value, err)
	}
	return log.Level(level), nil
}
