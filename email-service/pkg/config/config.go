package config

import "github.com/spf13/viper"

type Config struct {
	BrokerAddr []string `mapstructure:"BROKER_ADDR"`
	Topic      string   `mapstructure:"TOPIC"`
	GroupID    string   `mapstructure:"GROUP_ID"`
	MaxRetries int      `mapstructure:"MAX_RETRIES"`

	//email
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUsername string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	URL          string `mapstructure:"URL"`
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
