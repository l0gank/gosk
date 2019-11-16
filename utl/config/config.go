package config

import (
	"github.com/caarlos0/env"
)

type ConfigurationEnvironment struct {
	DatabaseEnvironment
	ServerEnvironment
	JWTEnvironment
}

type DatabaseEnvironment struct {
	MaxIdle     int    `env:"DBMaxIdleConn"`
	MaxOpenConn int    `env:"DBMaxOpenConn"`
	DBType      string `env:"DBType"`
	DBName      string `env:"DBName"`
	DBUser      string `env:"DBUser"`
	DBPass      string `env:"DBPass"`
	DBHost      string `env:"DBHost"`
}

type ServerEnvironment struct {
	Port         string `env:"HTTPPort"`
	Debug        bool   `env:"Debug"`
	ReadTimeout  int    `env:"ReadTimeoutSeconds"`
	WriteTimeout int    `env:"WriteTimeoutSeconds"`
}

type JWTEnvironment struct {
	Secret           string `env:"Secret"`
	Duration         int    `env:"DurationMinutes"`
	SigningAlgorithm string `env:"SigningAlgorithm"`
}

var Env = ConfigurationEnvironment{}

// Load Configuration struct from environment variable value
func LoadEnv() {
	if err := env.Parse(&Env); err != nil {
		panic(err.Error())
	}
}
