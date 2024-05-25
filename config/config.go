package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	configPath = "/Users/seongkim/go/src/simplebank/config.yaml"
)

type Config struct {
	DBDriver      string `yaml:"dbDriver"`
	DBSource      string `yaml:"dbSource"`
	ServerAddress string `yaml:"serverAddress"`
}

func LoadConfig() (Config, error) {
	var config Config
	file, err := os.Open(configPath)
	if err != nil || file == nil {
		return config, errors.Wrapf(err, "failed to open config fiile")
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, errors.Wrapf(err, "failed to decode config file")
	}

	return config, nil
}
