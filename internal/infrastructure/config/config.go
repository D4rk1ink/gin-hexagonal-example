package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type app struct {
	Env  string
	Port string
}

type database struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
}

type cryptography struct {
	HashKey string
}

type jwt struct {
	Secret    string
	ExpiresIn string
}

type config struct {
	App          app
	Database     database
	Cryptography cryptography
	Jwt          jwt
}

var Config *config

func Init() error {
	configPath := "config"
	if os.Getenv("APP_PWD") != "" {
		configPath = os.Getenv("APP_PWD") + "/config"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "test" {
		viper.SetConfigName("config.test.yaml")
	} else if appEnv == "dev" {
		viper.SetConfigName("config.dev.yaml")
	} else {
		viper.SetConfigName("config.yaml")
	}

	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		return err
	}

	return nil
}
