package config

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Структура конфигурации
type config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Apikey   string `yaml:"apikey"`
	Schedule string `yaml:"schedule"`
}

// Загрузка конфигурации
func LoadConfig() (*config, error) {
	config := &config{}

	yamlFile, err := ioutil.ReadFile("./internal/config/config.yaml")
	if err != nil {
		logrus.Fatalln("ioutill.ReadFile: ", err)
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		logrus.Fatalln("yaml.Unmarshal: ", err)
	}

	if config.Apikey == "" {
		return nil, ErrNoApikey
	}
	if config.Database == "" {
		return nil, ErrNoDatabase
	}
	if config.Host == "" {
		return nil, ErrNoHost
	}
	if config.Password == "" {
		return nil, ErrNoPassword
	}
	if config.Port == "" {
		return nil, ErrNoPort
	}
	if config.Schedule == "" {
		return nil, ErrNoSchedule
	}
	if config.Username == "" {
		return nil, ErrNoUsername
	}

	return config, nil
}
