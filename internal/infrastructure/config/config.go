package config

import "os"

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
	SecretKey string
	ExpiresIn string
}

type config struct {
	Env          string
	Database     database
	Cryptography cryptography
	Jwt          jwt
}

var Config *config

func Init() error {
	// NOTE: May import env from yaml file
	Config = &config{
		Env: os.Getenv("ENV"),
		Database: database{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		Jwt: jwt{
			SecretKey: os.Getenv("JWT_SECRET"),
			ExpiresIn: os.Getenv("JWT_EXPIRESIN"),
		},
	}

	return nil
}
