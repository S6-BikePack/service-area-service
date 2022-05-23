package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
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
	SSLMode  string
}

type Tracing struct {
	Host string
	Port int
}

func initDefaultValues() *Config {
	defaultConfig := &Config{}
	defaultConfig.Server.Service = "service-area-service"
	defaultConfig.Server.Port = "1234"
	defaultConfig.Server.Description = "Bikepack Service Area Service"

	defaultConfig.RabbitMQ.Host = "localhost"
	defaultConfig.RabbitMQ.Port = 5672
	defaultConfig.RabbitMQ.User = "user"
	defaultConfig.RabbitMQ.Password = "password"
	defaultConfig.RabbitMQ.Exchange = "topics"

	defaultConfig.Database.Host = "localhost"
	defaultConfig.Database.Port = 5432
	defaultConfig.Database.User = "user"
	defaultConfig.Database.Password = "password"
	defaultConfig.Database.Database = "service-area"
	defaultConfig.Database.Debug = false
	defaultConfig.Database.SSLMode = "disable"

	defaultConfig.Tracing.Host = ""
	defaultConfig.Tracing.Port = 0

	return defaultConfig
}

func UseConfig(path string) (*Config, error) {
	v := viper.New()

	defaults := initDefaultValues()

	v.SetConfigName(path)
	v.AddConfigPath(".")

	// If a config file is found, read it in. Otherwise, use defaults.
	if err := v.ReadInConfig(); err != nil {
		cfgMap := make(map[string]interface{})
		err := mapstructure.Decode(defaults, &cfgMap)
		if err != nil {
			fmt.Println("Error:", err)
		}

		cfgJsonBytes, err := json.Marshal(&cfgMap)
		if err != nil {
			fmt.Println("Error:", err)
		}

		v.SetConfigType("json")
		err = v.ReadConfig(bytes.NewReader(cfgJsonBytes))
		if err != nil {
			fmt.Println("Error:", err)
		}
	}

	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	v.AutomaticEnv()

	var config Config

	err := v.Unmarshal(&config)

	if err != nil {
		return nil, err
	}

	return &config, err
}
