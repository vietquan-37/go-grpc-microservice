package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	GatewayPort         string `mapstructure:"GATEWAY_PORT"`
	AuthServiceName     string `mapstructure:"AUTH_SERVICE_NAME"`
	ProductServiceName  string `mapstructure:"PRODUCT_SERVICE_NAME"`
	OrderServiceName    string `mapstructure:"ORDER_SERVICE_NAME"`
	ServiceName         string `mapstructure:"SERVICE_NAME"`
	ConsulAddr          string `mapstructure:"CONSUL_ADDR"`
	RequestPerTimeFrame int    `mapstructure:"REQUEST_PER_TIME_FRAME"`
	Resolve             bool   `mapstructure:"RESOLVE"`
}

func LoadConfig(path string) (config *Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
