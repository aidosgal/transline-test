package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	CustomerService ServiceConfig `env-prefix:"CUSTOMER_"` 
	Postgres        PostgresConfig `env-prefix:"POSTGRES_"`
	Service         AppConfig       `env-prefix:"APP_"`    
}

type ServiceConfig struct {
	Port int    `env:"PORT" env-default:"9090"`
	URL  string `env:"URL" env-default:"localhost"`
}

type PostgresConfig struct {
	Host     string `env:"HOST" env-default:"localhost"`
	Port     int    `env:"PORT" env-default:"5432"`
	User     string `env:"USER" env-default:"postgres"`
	Password string `env:"PASSWORD" env-default:"postgres"`
	DBName   string `env:"DBNAME" env-default:"postgres"`
	SSLMode  string `env:"SSLMODE" env-default:"disable"`
}

type AppConfig struct {
	Port int    `env:"PORT" env-default:"8080"`
	Name string `env:"NAME" env-default:"shipment-service"`
}

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read environment: %v", err)
	}
	return &cfg
}
