package config

import "github.com/spf13/viper"

type Config struct {
	StripeSecretKey string `mapstructure:"STRIPE_SECRET"`
	StripeSignature string `mapstructure:"STRIPE_SIGNATURE"`
	ServiceName     string `mapstructure:"SERVICE_NAME"`
	GrpcAddress     string `mapstructure:"GRPC_ADDRESS"`
	ConsulAddress   string `mapstructure:"CONSUL_ADDRESS"`
	BrokerAddress   string `mapstructure:"BROKER_ADDRESS"`
	Topic           string `mapstructure:"TOPIC"`
	Currency        string `mapstructure:"CURRENCY"`
	SuccessURL      string `mapstructure:"SUCCESS_URL"`
	CancelURL       string `mapstructure:"CANCEL_URL"`
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
