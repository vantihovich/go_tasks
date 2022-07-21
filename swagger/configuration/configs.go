package config

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type Database struct {
	Host     string `env:"DBHOST,required"`
	Port     int    `env:"DBPORT,required"`
	User     string `env:"DBUSER,required"`
	Password string `env:"DBPASSWORD,required"`
	Name     string `env:"DBNAME,required"`
}

type App struct {
	Database Database
}

type JWTSecret struct {
	SecretKey string `env:"JWTSECRETKEY,required"`
}

type LoginLimitParameters struct {
	InvalidLoginAttemptTTL  time.Duration `env:"INVALIDLOGINATTEMPTTTL,required"`
	MaxAllowedInvalidLogins int           `env:"MAXALLOWEDINVALIDLOGINS,required"`
}

type MailJetParameters struct {
	User        string `env:"MAILJETUSER,required"`
	Password    string `env:"MAILJETPASSWORD,required"`
	SenderEmail string `env:"MAILJETSENDEREMAIL,required"`
}

type WorldCoinIndexParameters struct {
	Key string `env:"WORLDCOININDEXKEY,required"`
	URL string `env:"WORLDCOININDEXURL,required"`
}

func LoadDB() (App, error) {
	cfg := App{}
	if err := env.Parse(&cfg); err != nil {
		return App{}, err
	}
	return cfg, nil
}

func LoadJWT() (JWTSecret, error) {
	cfg := JWTSecret{}
	if err := env.Parse(&cfg); err != nil {
		return JWTSecret{}, err
	}
	return cfg, nil
}

func LoadLogin() (LoginLimitParameters, error) {
	cfg := LoginLimitParameters{}
	if err := env.Parse(&cfg); err != nil {
		return LoginLimitParameters{}, err
	}

	return cfg, nil
}

func LoadMailJetParameters() (MailJetParameters, error) {
	cfg := MailJetParameters{}
	if err := env.Parse(&cfg); err != nil {
		return MailJetParameters{}, err
	}

	return cfg, nil
}

func LoadWCIParameter() (WorldCoinIndexParameters, error) {
	cfg := WorldCoinIndexParameters{}
	if err := env.Parse(&cfg); err != nil {
		return WorldCoinIndexParameters{}, err
	}
	return cfg, nil
}
