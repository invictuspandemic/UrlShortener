package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-required:"true"`
	Storage    string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

func LoadCFG() (string, *Config) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return fmt.Sprint("CONFIG_PATH is not defined"), nil
	}

	// check if file exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Sprintf("config file '%s' does not exist", configPath), nil
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return fmt.Sprintf("cannot read config: %s", err), nil
	}

	return "", &cfg

}
