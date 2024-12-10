package config

import (
	"errors"
	"os"
)

type Config struct {
	DbSource          string
	GrpcServerAddress string
	JwtSecret         string
	AdminUserName     string
	AdminPassword     string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		DbSource:          os.Getenv("DB_SOURCE"),
		GrpcServerAddress: os.Getenv("GRPC_SERVER_ADDRESS"),
		JwtSecret:         os.Getenv("JWT_SECRET"),
		AdminUserName:     os.Getenv("ADMIN_USERNAME"),
		AdminPassword:     os.Getenv("ADMIN_PASSWORD"),
	}

	if config.DbSource == "" || config.GrpcServerAddress == "" || config.JwtSecret == "" || config.AdminUserName == "" || config.AdminPassword == "" {
		return nil, errors.New("missing env variable")
	}
	return config, nil
}
