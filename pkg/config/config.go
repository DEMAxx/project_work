package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

const AppName = "previewer"

type Config struct {
	Timeout struct {
		Read     byte `toml:"TIMEOUT_READ" env:"TIMEOUT_READ" env-default:"5"`
		Write    byte `toml:"TIMEOUT_WRITE" env:"TIMEOUT_WRITE" env-default:"5"`
		Shutdown byte `toml:"TIMEOUT_SHUTDOWN" env:"TIMEOUT_SHUTDOWN" env-default:"3"`
	}
	Server struct {
		Host string `toml:"SERVER_HOST" env:"SERVER_HOST" env-default:"localhost"`
		Port string `toml:"SERVER_PORT" env:"SERVER_PORT" env-default:"8000"`
	}
	Capability int    `toml:"CAPABILITY" env:"CAPABILITY" env-default:"10"`
	Debug      bool   `toml:"APP_DEBUG" env:"APP_DEBUG" env-default:"true"`
	Env        string `toml:"APP_ENV" env:"APP_ENV" env-default:"local"`
	Local      bool   `toml:"LOCAL" env:"LOCAL"`
	LogLevel   string `toml:"LOG_LEVEL" env:"LOG_LEVEL" env-default:"info"`
	UploadPath string `toml:"UPLOAD_PATH" env:"UPLOAD_PATH" env-default:"/tmp"`
}

func MustLoad(configFile string) *Config {
	cfg := Config{}

	if err := cleanenv.ReadConfig(configFile, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
