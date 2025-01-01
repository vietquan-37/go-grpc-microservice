package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbSource          string `mapstructure:"DB_SOURCE"`
	GrpcServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	JwtSecret         string `mapstructure:"JWT_SECRET"`
	AdminUserName     string `mapstructure:"ADMIN_USERNAME"`
	AdminPassword     string `mapstructure:"ADMIN_PASSWORD"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
