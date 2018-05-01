package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	EnvMysqlHost = "MYSQL_HOST"
	EnvMysqlUser = "MYSQL_USER"
	EnvMysqlPass = "MYSQL_PASS"
	EnvMysqlDb   = "MYSQL_DB"
)

const (
	EnvMemcacheHost = "MEMCACHE_HOST"
	EnvMemcachePort = "MEMCACHE_PORT"
)

type MysqlConfig struct {
	Host     string
	Username string
	Password string
	Database string
}

type MemcacheConfig struct {
	Host string
	Port string
}

func (m MemcacheConfig) GetConnectionString() string {
	return fmt.Sprintf("%s:%s", m.Host, m.Port)
}

func init() {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.memo")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Config file not found")
	}
}

func GetMysqlConfig() MysqlConfig {
	return MysqlConfig{
		Host:     viper.GetString(EnvMysqlHost),
		Username: viper.GetString(EnvMysqlUser),
		Password: viper.GetString(EnvMysqlPass),
		Database: viper.GetString(EnvMysqlDb),
	}
}

func GetMemcacheConfig() MemcacheConfig {
	return MemcacheConfig{
		Host: viper.GetString(EnvMemcacheHost),
		Port: viper.GetString(EnvMemcachePort),
	}
}
