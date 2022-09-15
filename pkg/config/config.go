package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Port string `mapstructure:"PORT"`
	DB   DB     `mapstructure:",squash"`
	SMTP SMTP   `mapstructure:",squash"`
}

type DB struct {
	Name     string `mapstructure:"DB_NAME"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Host     string `mapstructure:"DB_HOST"`
}

type SMTP struct {
	From     string `mapstructure:"SMTP_FROM_MAIL"`
	Host     string `mapstructure:"SMTP_HOST"`
	Port     string `mapstructure:"SMTP_PORT"`
	Username string `mapstructure:"SMTP_USERNAME"`
	Password string `mapstructure:"SMTP_PASSWORD"`
}

func LoadConfig() *Config {
	viper.AddConfigPath("./pkg/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Could not read configuration:", err.Error())
	}

	var c Config
	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatalln("Could not unmarshal to Config struct:", err.Error())
	}

	return &c
}
