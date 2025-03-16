package configParser

import (
	"github.com/go-yaml/yaml"
	"os"
)

type ConnectionConfig struct {
	Port            int    `yaml:"port"`
	Host            string `yaml:"host"`
	ShutdownTimeout int    `yaml:"shutdownTimeout"`
	TokenTTl        int    `yaml:"tokenTTL"`
}

func ParseConnectionConfig(filePath string) (*ConnectionConfig, error) {
	var config ConnectionConfig

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
