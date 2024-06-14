package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `yaml:"env" env-default="local"`
	Server      string `yaml:"http_server"`
	StoragePath string `yaml:"storage_path"`
	LogLib      string `yaml:"log_lib" env-default="devslog"`
}

func MustLoad() *Config {
	var cfg Config

	s := flag.String("a", "", "localhost:8080")
	sp := flag.String("b", "", "postgres://postgres:pwd@host:port/db")
	flag.Parse()

	if *s != "" && *sp != "" {
		cfg.Server = *s
		cfg.StoragePath = *sp
		cfg.Env = "local"
		cfg.LogLib = "devslog"
		return &cfg
	}

	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
	path, exists := os.LookupEnv("CONFIG_PATH")
	if !exists {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("can't read config file: " + err.Error())
	}

	return &cfg
}
