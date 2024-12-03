package config

import "github.com/spf13/viper"

type Config struct {
	GatewayPort string `mapstructure:"GATEWAY_PORT"`
	AuthUrl     string `mapstructure:"AUTH_URL"`
	ProductUrl  string `mapstructure:"PRODUCT_URL"`
	OrderUrl    string `mapstructure:"ORDER_URL"`
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
	if err != nil {
		return
	}
	return
}
