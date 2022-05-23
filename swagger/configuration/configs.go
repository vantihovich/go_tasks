package config

import (
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

func Load() (App, error) {
	cfg := App{}
	if err := env.Parse(&cfg); err != nil {
		return App{}, err
	}
	return cfg, nil
}
