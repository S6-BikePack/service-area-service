package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Server   Server
	RabbitMQ RabbitMQ
	Database Database
	Tracing  Tracing
}

type Server struct {
	Service     string
	Port        string
	Description string
}

type RabbitMQ struct {
	Host     string
	Port     int
	User     string
	Password string
	Exchange string
}

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Debug    bool
}

type Tracing struct {
	Host string
	Port int
}

func UseConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetDefault("server.service", "service-area-service")
	v.SetDefault("server.port", "1234")
	v.SetDefault("server.description", "bikepack project service-area-service")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "user")
	v.SetDefault("database.password", "password")

	v.SetDefault("rabbitmq.host", "localhost")
	v.SetDefault("rabbitmq.port", 5672)
	v.SetDefault("rabbitmq.user", "user")
	v.SetDefault("rabbitmq.password", "password")

	v.SetConfigName(path)
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	v.AutomaticEnv()

	var config Config

	err := v.Unmarshal(&config)

	if err != nil {
		return nil, err
	}

	return &config, err
}
