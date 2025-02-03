package config

import (
	"reflect"

	"github.com/spf13/viper"
)

var CONFIG *Config

type MySQLConfig struct {
	IP                string //公有访问
	Port              string
	User              string
	Password          string
	Database          string
	MaxIdleConnection int `mapstructure:"max_idle_connection"`
	MaxOpenConnection int `mapstructure:"max_open_connection"`
}

func (this *MySQLConfig) IsEmpty() bool {
	return reflect.DeepEqual(this, MySQLConfig{})
}

type SQLiteConfig struct {
	File              string
	Database          string
	MaxIdleConnection int `mapstructure:"max_idle_connection"`
	MaxOpenConnection int `mapstructure:"max_open_connection"`
}

func (this *SQLiteConfig) IsEmpty() bool {
	return reflect.DeepEqual(this, SQLiteConfig{})
}

// Config对象和config.toml文件保持一致
type Config struct {
	AppName            string
	Address            string
	Port               string
	DbusPort           string       `mapstructure:"dbus_port"`
	AsyncProcess       bool         `mapstructure:"async_process"`
	MissionRetryNumber int          `mapstructure:"mission_retry_number"`
	DatabaseDriver     string       `mapstructure:"database_driver"`
	MySQL              MySQLConfig  //需要定义子类型对应的变量，如果不定义映射不成功
	SQLite             SQLiteConfig //需要定义子类型对应的变量，如果不定义映射不成功
	HA                 ha
}

type ha struct {
	Namespace string
}

func InitConfiguration(configFile string) error {
	viper.SetDefault("Port", "2952")
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
