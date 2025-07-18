package config

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DeepSeekConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	DeepSeek DeepSeekConfig `mapstructure:"deepseek"`
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetDefault("server.port", ":8080")
	viper.SetDefault("deepseek.base_url", "https://api.deepseek.com")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v, using defaults", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("Error unmarshalling config: %v", err)
	}

	return &config
}
