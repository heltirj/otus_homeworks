package main

import (
	"io"
	"os"

	"github.com/heltirj/otus_homeworks/hw12_13_14_15_calendar/internal/logger"
	"github.com/heltirj/otus_homeworks/hw12_13_14_15_calendar/internal/storage"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	LogLevel    logger.LogLevel `yaml:"logLevel"`
	StorageType storage.Type    `yaml:"storageType"`
	Database    DatabaseConfig  `yaml:"database"`
	Service     ServiceConfig   `yaml:"service"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

type ServiceConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func NewConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
