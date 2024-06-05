package app

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Settings struct {
	Host     string `yaml:"host"`
	LogLevel string `yaml:"log_level"`
	LogPath  string `yaml:"log_path"`
	Cors     bool   `yaml:"cors"`

	HttpHost string `yaml:"http_host"`
}

type Config struct {
	Settings
}

func NewConfig(configPath string) (*Config, error) {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
