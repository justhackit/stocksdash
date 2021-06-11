package config

import (
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
)

// Configurations wraps all the config variables required by the auth service
type Configurations struct {
	ServerAddress string
	Database      DatabaseConfig
}

type DatabaseConfig struct {
	DBHost   string
	DBName   string
	DBSchema string
	DBUser   string
	DBPass   string
	DBPort   string
	DBConn   string
}

func NewConfigurationsFromFile(configFilePath string, logger hclog.Logger) *Configurations {
	viper.SetConfigFile(configFilePath)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configFilePath)
	viper.AutomaticEnv()
	var confs Configurations
	if err := viper.ReadInConfig(); err != nil {
		logger.Error("Error reading config file", "errorMsg", err)
	}
	err := viper.Unmarshal(&confs)
	if err != nil {
		logger.Error("Unable to decode into struct", "errorMsg", err)
	}
	return &confs
}
