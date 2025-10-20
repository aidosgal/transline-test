package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	CustomerService ServiceConfig `env-prefix:"CUSTOMER_"`
	Shipment        ServiceConfig `env-prefix:"SHIPMENT_"`
	Service         AppConfig     `env-prefix:"APP_"`
	Jaeger          JaegerConfig  `env-prefix:"JAEGER_"`
}

type ServiceConfig struct {
	Port     int            `env:"PORT" env-default:"9090"`
	URL      string         `env:"URL" env-default:"localhost"`
	Postgres PostgresConfig `env-prefix:"POSTGRES_"`
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
	Name string `env:"NAME" env-default:"envoy"`
}

type JaegerConfig struct {
	URL string `env:"URL" env-default:"http://jaeger:14268/api/traces"`
}

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read environment: %v", err)
	}
	return &cfg
}

func (pc *PostgresConfig) BuildPostgresURL() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pc.Host,
		pc.Port,
		pc.User,
		pc.Password,
		pc.DBName,
		pc.SSLMode,
	)
}
