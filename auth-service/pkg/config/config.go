package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbSource          string `mapstructure:"DB_SOURCE"`
	GrpcServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	ServiceName       string `mapstructure:"SERVICE_NAME"`
	JwtSecret         string `mapstructure:"JWT_SECRET"`
	AdminUserName     string `mapstructure:"ADMIN_USERNAME"`
	AdminPassword     string `mapstructure:"ADMIN_PASSWORD"`
	ConsulAddress     string `mapstructure:"CONSUL_ADDR"`
	ClientId          string `mapstructure:"CLIENT_ID"`
	ClientSecret      string `mapstructure:"CLIENT_SECRET"`
	GrantType         string `mapstructure:"GRANT_TYPE"`
	RedirectUri       string `mapstructure:"REDIRECT_URI"`
	EmailUsername     string `mapstructure:"EMAIL_USERNAME"`
	SMTPHost          string `mapstructure:"SMTP_HOST"`
	SMTPPort          string `mapstructure:"SMTP_PORT"`
	EmailPassword     string `mapstructure:"EMAIL_PASSWORD"`
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
