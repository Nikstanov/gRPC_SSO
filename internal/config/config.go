package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string         `yaml:"env" env-default:"local"`
	StoragePath string         `yaml:"storage_path" env-default:"./data"`
	TokenTTL    time.Duration  `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig     `yaml:"grpc"`
	DB          DatabaseConfig `yaml:"database"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type DatabaseConfig struct {
	Port     int    `yaml:"port" env-default:"5432"`
	Dbname   string `yaml:"db" env-default:"sso"`
	User     string `yaml:"user" env-default:"admin"`
	Password string `yaml:"password" env-default:"admin"`
	DBAddr   string `yaml:"addr" env-default:"localhost"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config path does not exist: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config" + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
