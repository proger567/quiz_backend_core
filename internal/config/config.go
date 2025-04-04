package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppEnv           string `env:"APP_ENV"  envDefault:"production"`
	ListenAddr       string `env:"LISTEN_ADDR"  envDefault:":8095"`
	LogLevel         string `env:"LOG_LEVEL"  envDefault:"INFO"`
	DatabaseHost     string `env:"DB_HOST"  envDefault:"localhost"`
	DatabasePort     string `env:"DB_PORT" envDefault:"5432"`
	DatabaseUser     string `env:"DB_USER" envDefault:"postgres"`
	DatabasePassword string `env:"DB_PASSWORD" envDefault:"pgpassword"`
	DatabaseName     string `env:"DB_NAME" envDefault:"quiz"`
}

func (c *Config) Parse() error {
	flag.StringVar(&c.DatabaseHost, "dhost", "localhost", "Database address")
	flag.StringVar(&c.DatabasePort, "dport", "5432", "Database port")
	flag.StringVar(&c.DatabaseUser, "user", "postgres", "Database user")
	flag.StringVar(&c.DatabasePassword, "pass", "pgpassword", "Database user password")
	flag.StringVar(&c.DatabaseName, "db", "quiz", "Database name")
	flag.Parse()

	//settings redefinition, if env variables are used
	err := env.Parse(c)

	return err
}
