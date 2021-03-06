// Package config ...
package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/n-r-w/lg"
)

// Config logserver.toml
type Config struct {
	SuperAdminID         int
	RestHost             string   `toml:"REST_HOST"`
	RestPort             string   `toml:"REST_PORT"`
	GrpcHost             string   `toml:"GRPC_HOST"`
	GrpcPort             string   `toml:"GRPC_PORT"`
	DatabaseURL          string   `toml:"DATABASE_URL"`
	MaxDbSessions        int      `toml:"MAX_DB_SESSIONS"`
	MaxDbSessionIdleTime int      `toml:"MAX_DB_SESSION_IDLE_TIME"`
	MaxLogRecordsResult  int      `toml:"MAX_LOG_RECORDS_RESULT"`
	DbReadTimeout        int      `toml:"DB_READ_TIMEOUT"`
	DbWriteTimeout       int      `toml:"DB_WRITE_TIMEOUT"`
	HttpReadTimeout      int      `toml:"HTTP_READ_TIMEOUT"`
	HttpWriteTimeout     int      `toml:"HTTP_WRITE_TIMEOUT"`
	HttpShutdownTimeout  int      `toml:"HTTP_SHUTDOWN_TIMEOUT"`
	RateLimit            int      `toml:"RATE_LIMIT"`
	RateLimitBurst       int      `toml:"RATE_LIMIT_BURST"`
	TokensRead           []string `toml:"TOKENS_READ"`
	TokensWrite          []string `toml:"TOKENS_WRITE"`
}

const (
	superAdminID         = 1
	maxDbSessions        = 50
	maxDbSessionIdleTime = 50
	maxLogRecordsResult  = 100000
	defaultSessionAge    = 60 * 60 * 24 // 24 часа
)

// New Инициализация конфига значениями по умолчанию
func New(configPath string, logger lg.Logger) (*Config, error) {
	c := &Config{
		SuperAdminID:         superAdminID,
		RestHost:             "0.0.0.0",
		RestPort:             "8080",
		GrpcHost:             "0.0.0.0",
		GrpcPort:             "9000",
		DatabaseURL:          "",
		MaxDbSessions:        maxDbSessions,
		MaxDbSessionIdleTime: maxDbSessionIdleTime,
		MaxLogRecordsResult:  maxLogRecordsResult,
		DbReadTimeout:        2000,
		DbWriteTimeout:       8000,
		HttpReadTimeout:      1000,
		HttpWriteTimeout:     10000,
		HttpShutdownTimeout:  10,
		RateLimit:            10000,
		RateLimitBurst:       20000,
		TokensRead:           []string{},
		TokensWrite:          []string{},
	}

	if configPath != "" {
		if _, err := toml.DecodeFile(configPath, c); err != nil {
			return nil, err
		}
	}

	if c.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL undefined")
	}

	return c, nil
}
