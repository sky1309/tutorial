package tutorial_viper

import (
	"fmt"
	"log"
	"testing"

	"github.com/spf13/viper"
)

func TestViper(t *testing.T) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("viper read in config err %v", err)
	}

	logLevel := viper.GetString("log_level")
	fmt.Printf("logLevel=%s\n", logLevel)

	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")
	db := viper.GetInt("redis.db")

	fmt.Printf("redis host=%s, port=%s, db=%d", host, port, db)
}

func TestViper_UnMarshl(t *testing.T) {
	type RedisConfig struct {
		Host string `mapstructure:"host" yaml:"host"`
		Port int    `mapstructure:"port" yaml:"port"`
		DB   int    `mapstructure:"db" yaml:"db"`
	}

	type Config struct {
		LogLevel string      `mapstructure:"log_level" yaml:"log_level"`
		Redis    RedisConfig `mapstructure:"redis" yaml:"redis"`
	}
	viper.SetConfigFile("./config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("viper read in config err %v", err)
	}

	var config *Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("viper unmarshal err %s", err)
	}

	fmt.Printf("logLevel=%s\n", config.LogLevel)
	fmt.Printf("redis host=%s, port=%s, db=%d", config.Redis.Host, config.Redis.Host, config.Redis.DB)
}
