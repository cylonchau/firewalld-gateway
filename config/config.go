package config

import (
	"github.com/spf13/viper"
)

var CONFIG *Config

type MysqlConfig struct {
	Ip       string //公有访问
	Port     string
	User     string
	Password string
	Database string
}

func (this *MysqlConfig) IsEmpty() bool {
	return this == nil
}

type Config struct { //Config对象和config.toml文件保持一致
	AppName              string
	LogLevel             string
	Address              string
	Port                 string
	Dbus_Port            string
	Async_Process        bool
	Mission_Retry_Number int
	Mysql                MysqlConfig //需要定义子类型对应的变量，如果不定义映射不成功
}

func InitConfiguration(configFile string) error {
	viper.SetDefault("LogLevel", "info")
	viper.SetDefault("Port", "7777")
	viper.SetDefault("Address", "127.0.0.1")
	viper.SetConfigType("toml")
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return err
	}
	CONFIG = config
	return nil
}
