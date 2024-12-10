package config

import (
	"errors"
	"os"
)

type Config struct {
	DbSource string
	GrpcAddr string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		DbSource: os.Getenv("DB_SOURCE"),
		GrpcAddr: os.Getenv("GRPC_SERVER_ADDRESS"),
	}

	if config.DbSource == "" || config.GrpcAddr == "" {
		return nil, errors.New("missing env variables")
	}

	return config, nil
}
