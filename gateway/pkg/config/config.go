package config

import (
	"errors"
	"os"
)

type Config struct {
	GatewayPort string
	AuthUrl     string
	ProductUrl  string
	OrderUrl    string
}

func LoadConfig() (*Config, error) {

	config := &Config{
		GatewayPort: os.Getenv("GATEWAY_PORT"),
		AuthUrl:     os.Getenv("AUTH_URL"),
		ProductUrl:  os.Getenv("PRODUCT_URL"),
		OrderUrl:    os.Getenv("ORDER_URL"),
	}

	if config.GatewayPort == "" || config.AuthUrl == "" || config.ProductUrl == "" || config.OrderUrl == "" {
		return nil, errors.New("missing env variable")
	}
	return config, nil
}
