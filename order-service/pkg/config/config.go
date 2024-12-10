package config

import (
	"errors"
	"os"
)

type Config struct {
	DbSource          string
	GrpcServerAddress string
	ProductURL        string
	AuthURL           string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		DbSource:          os.Getenv("DB_SOURCE"),
		GrpcServerAddress: os.Getenv("GRPC_SERVER_ADDRESS"),
		ProductURL:        os.Getenv("PRODUCT_URL"),
		AuthURL:           os.Getenv("AUTH_URL"),
	}

	if config.DbSource == "" || config.GrpcServerAddress == "" || config.ProductURL == "" || config.AuthURL == "" {
		return nil, errors.New(config.DbSource)
	}

	return config, nil

}
