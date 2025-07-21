package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbSource           string   `mapstructure:"DB_SOURCE"`
	GrpcServerAddress  string   `mapstructure:"GRPC_SERVER_ADDRESS"`
	ServiceName        string   `mapstructure:"SERVICE_NAME"`
	ProductServiceName string   `mapstructure:"PRODUCT_SERVICE_NAME"`
	AuthServiceName    string   `mapstructure:"AUTH_SERVICE_NAME"`
	ConsulAddr         string   `mapstructure:"CONSUL_ADDR"`
	PaymentServiceName string   `mapstructure:"PAYMENT_SERVICE_NAME"`
	OrderTopic         string   `mapstructure:"ORDER_TOPIC"`
	EmailTopic         string   `mapstructure:"EMAIL_TOPIC"`
	DLQTopic           string   `mapstructure:"DLQ_TOPIC"`
	BrokerAddr         []string `mapstructure:"BROKER_ADDR"`
	GroupId            string   `mapstructure:"GROUP_ID"`
	MaxRetries         int      `mapstructure:"MAX_RETRIES"`
	WorkerCount        int      `mapstructure:"WORKER_COUNT"`
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
