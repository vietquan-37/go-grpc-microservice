package config

import "github.com/spf13/viper"

type Config struct {
	DbSource          string `mapstructure:"DB_SOURCE"`
	GrpcServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	ProductURL        string `mapstructure:"PRODUCT_URL"`
	AuthURL           string `mapstructure:"AUTH_URL"`
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
