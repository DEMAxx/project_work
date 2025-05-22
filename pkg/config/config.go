package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

const AppName = "previewer"

type Config struct {
	Timeout struct {
		Read     byte `env:"TIMEOUT_READ" env-default:"5"`
		Write    byte `env:"TIMEOUT_WRITE" env-default:"5"`
		Shutdown byte `env:"TIMEOUT_SHUTDOWN" env-default:"3"`
	}
	Server struct {
		Host string `env:"SERVER_HOST" env-default:"localhost"`
		Port string `env:"SERVER_PORT" env-default:"8000"`
	}
	Capability int    `env:"CAPABILITY" env-default:"10"`
	Debug      bool   `env:"APP_DEBUG" env-default:"true"`
	Env        string `env:"APP_ENV" env-default:"local"`
	Local      bool   `env:"LOCAL"`
	LogLevel   string `env:"LOG_LEVEL" env-default:"info"`
	UploadPath string `env:"UPLOAD_PATH" env-default:"/tmp"`
}

func MustLoad(configFile string) *Config {
	cfg := Config{}

	if err := cleanenv.ReadConfig(configFile, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
