package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string     `yaml:"env" env-default:"local"`
	GRPC           GRPCConfig `yaml:"grpc"`
	MigrationsPath string     `yaml:"migrations_path"`
	Db             Database   `yaml:"database"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout"`
}

type Database struct {
	Host     string `yaml:"host" env-default:"local"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"my"`
	Password string `yaml:"password" env-default:"12345"`
	Dbname   string `yaml:"dbname" env-default:"mydb"`
}

type fetchConfigPathProvider interface {
	fetchConfigPath() string
}

type defaultFetchCfgPathProvider struct{}

// Gets the path to the configuration file
func (defaultFetchCfgPathProvider) fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

var cfgPathProvider fetchConfigPathProvider = defaultFetchCfgPathProvider{}

// Parse config file, panic in case of error
func MustLoad() *Config {
	configPath := cfgPathProvider.fetchConfigPath()

	return MustLoadByPath(configPath)
}

func MustLoadByPath(configPath string) *Config {
	if configPath == "" {
		panic("config path empty")
	}

	var conf Config

	if err := cleanenv.ReadConfig(configPath, &conf); err != nil {
		panic("YAML parsing error: " + err.Error())
	}

	return &conf
}
