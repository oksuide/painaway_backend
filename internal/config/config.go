package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env              string           `mapstructure:"env"`
	HTTPServerConfig HTTPServerConfig `mapstructure:"http_server"`
	DBConfig         DBConfig         `mapstructure:"db"`
	JWTConfig        JWTConfig        `mapstructure:"jwt"`
}

type HTTPServerConfig struct {
	Address     string        `mapstructure:"address"`
	ServerPort  string        `mapstructure:"server_port"`
	Timeout     time.Duration `mapstructure:"timeout"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

type JWTConfig struct {
	SecretKey string        `mapstructure:"secret_key"`
	Duration  time.Duration `mapstructure:"duration"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	// ENV -> CONFIG
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
