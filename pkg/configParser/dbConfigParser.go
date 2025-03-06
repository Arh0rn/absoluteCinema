package configParser

import (
	"github.com/go-yaml/yaml"
	"os"
)

type DbConfig struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	User    string `yaml:"user"`
	DBName  string `yaml:"dbname"`
	SSLMode string `yaml:"sslmode"`
}

// Good of current implementation:
func ParseDBConfig(filePath string) (*DbConfig, error) {
	var config DbConfig

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

//No reason to use
//func ParseConfig(filePath string) (*dbConfig, error) {
//	var config dbConfig
//
//	file, err := os.Open(filePath)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
//	bytes, err := io.ReadAll(file)
//	if err != nil {
//		log.Println(err)
//		return nil, err
//	}
//
//	err = yaml.Unmarshal(bytes, &config)
//	if err != nil {
//		return nil, err
//	}
//
//	return &config, nil
//}

//Best practice for big files
//func ParseConfig(filePath string) (*dbConfig, error) {
//	var config dbConfig
//
//	file, err := os.Open(filePath)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
//	decoder := yaml.NewDecoder(file)
//	err = decoder.Decode(&config)
//	if err != nil {
//		return nil, err
//	}
//
//	return &config, nil
//}
