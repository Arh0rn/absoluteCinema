package configParser

import (
	"github.com/go-yaml/yaml"
	"os"
)

//port: 8080
//host: localhost
//shutdownTimeout: 10 #seconds
//readTimeout: 10 #seconds
//writeTimeout: 10 #seconds
//idleTimeout: 120 #seconds
//accessTokenTTL: 15 #minutes
//refreshTokenTTL: 43200 #minutes (one month)
//cacheTTL: 5 #minutes

type ConnectionConfig struct {
	Port            int    `yaml:"port"`
	Host            string `yaml:"host"`
	ShutdownTimeout int    `yaml:"shutdownTimeout"`
	ReadTimeout     int    `yaml:"readTimeout"`
	WriteTimeout    int    `yaml:"writeTimeout"`
	IdleTimeout     int    `yaml:"idleTimeout"`
	AccessTokenTTL  int    `yaml:"accessTokenTTL"`
	RefreshTokenTTL int    `yaml:"refreshTokenTTL"`
	CacheTTL        int    `yaml:"cacheTTL"`
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
