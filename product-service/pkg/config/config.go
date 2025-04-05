package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbSource        string `mapstructure:"DB_SOURCE"`
	GrpcAddr        string `mapstructure:"GRPC_SERVER_ADDRESS"`
	ServiceName     string `mapstructure:"SERVICE_NAME"`
	AuthServiceName string `mapstructure:"AUTH_SERVICE_NAME"`
	ConsulAddr      string `mapstructure:"CONSUL_ADDR"`
	Resolve         bool   `mapstructure:"RESOLVE"`
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
