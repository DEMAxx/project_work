package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

const AppName = "previewer"
const AppPort = 8000
const ProxyHost = "http://localhost:8080"

type Env struct {
	Path string `env:"ENV_PATH" env-default:"/etc/previewer/config.toml"`
}

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

func MustLoad() *Config {
	env := Env{}
	cfg := Config{}

	if err := cleanenv.ReadEnv(&env); err != nil {
		panic("cannot read env: " + err.Error())
	}

	if err := cleanenv.ReadConfig(env.Path, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
