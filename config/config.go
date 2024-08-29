package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
)

var AppName = "campaign-service"

type Config struct {
	AppVersion        string
	AppVersionNumber  [3]uint32
	HTTPPort          uint16        `env:"HTTP_PORT,default=5101"`
	GRPCPort          uint16        `env:"GRPC_PORT,default=5102"`
	ProfilerPort      uint16        `env:"PROFILER_PORT,default=5103"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT,default=15s"`
	LogLevel          string        `env:"LOG_LEVEL,default=warn"`
	AllowedOrigins    []string      `env:"ALLOWED_ORIGINS,delimiter=,default=*"`
	AllowedHeaders    []string      `env:"ALLOWED_HEADERS,delimiter=,default=*"`
	DatabaseURL       string        `env:"DATABASE_URL"`
	DatabaseURLReader string        `env:"DATABASE_URL_READER"`
	DatabaseName      string        `env:"DATABASE_NAME"`
	EnableProfiler    bool          `env:"ENABLE_PROFILER"`
}

func New() *Config {
	cfg := &Config{
		AppVersion:       Version,
		AppVersionNumber: VersionNumber,
	}
	if err := envconfig.Process(context.Background(), cfg); err != nil {
		panic("can't load configuration file")
	}

	return cfg
}
