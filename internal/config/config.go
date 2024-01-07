package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	GRPCBindConfig BindConfig    `yaml:"bind_grpc"`
	HttpBindConfig BindConfig    `yaml:"bind_http"`
	DBConfig       DBConfig      `yaml:"DB"`
	TokenTTL       time.Duration `yaml:"token_TTL"`
}

type BindConfig struct {
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}

type DBConfig struct {
	Server   string `yaml:"server"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

func GetConfig(path string) (*Config, error) {
	var cnf Config
	err := cleanenv.ReadConfig(path, &cnf)
	return &cnf, err
}
