package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

const configPath = "CONFIG_PATH"

// Config is where the global configuration is stored.
type Config struct {
	Port string `mapstructure:"PORT"`
	DB   dB     `mapstructure:",squash"`
	SMTP sMTP   `mapstructure:",squash"`
	JWT  jwt    `mapstructure:",squash"`
}

type dB struct {
	Name     string `mapstructure:"DB_NAME"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Host     string `mapstructure:"DB_HOST"`
}

type sMTP struct {
	From     string `mapstructure:"SMTP_FROM_MAIL"`
	Host     string `mapstructure:"SMTP_HOST"`
	Port     string `mapstructure:"SMTP_PORT"`
	Username string `mapstructure:"SMTP_USERNAME"`
	Password string `mapstructure:"SMTP_PASSWORD"`
}

type jwt struct {
	Expiration int32  `mapstructure:"JWT_EXPIRATION_SECS"`
	Key        string `mapstructure:"JWT_KEY"`
}

// LoadConfig gets the configuration in from .env files and stores the in Config struct.
func LoadConfig() *Config {
	cPath := os.Getenv(configPath)
	if cPath == "" {
		log.Fatalln("CONFIG_PATH is not set")
	}

	viper.AddConfigPath(cPath + "/pkg/config/envs")
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
